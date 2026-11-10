package hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// captureLogs sets a JSON slog default logger with the given level for the duration of the test
// and returns the buffer the logs are written to.
func captureLogs(t *testing.T, level slog.Level) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: level})))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

func parseLogLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 log line, got %d: %q", len(lines), buf.String())
	}
	entry := map[string]any{}
	if err := json.Unmarshal([]byte(lines[0]), &entry); err != nil {
		t.Fatalf("failed to parse log line %q: %v", lines[0], err)
	}
	return entry
}

func TestLoggerHooksPayloads(t *testing.T) {
	ctx := context.Background()
	message := &mcp.ListResourcesRequest{}
	// large result, similar to resources/list with the embedded documentation index
	result := &mcp.ListResourcesResult{
		Resources: []mcp.Resource{mcp.NewResource("docs://"+strings.Repeat("a", 1<<20), "doc")},
	}

	hooks := NewLoggerHooks()
	calls := []struct {
		name         string
		call         func()
		level        string
		payloadKeys  []string
		requiredKeys []string
	}{
		{
			name:         "request received",
			call:         func() { hooks.OnBeforeAny[0](ctx, 1, mcp.MethodResourcesList, message) },
			level:        "INFO",
			payloadKeys:  []string{"message"},
			requiredKeys: []string{"request_id", "session_id", "method"},
		},
		{
			name:         "request succeeded",
			call:         func() { hooks.OnSuccess[0](ctx, 1, mcp.MethodResourcesList, message, result) },
			level:        "INFO",
			payloadKeys:  []string{"message", "result"},
			requiredKeys: []string{"request_id", "session_id", "method"},
		},
		{
			name:         "request failed",
			call:         func() { hooks.OnError[0](ctx, 1, mcp.MethodResourcesList, message, errors.New("boom")) },
			level:        "ERROR",
			payloadKeys:  []string{"message"},
			requiredKeys: []string{"request_id", "session_id", "method", "error"},
		},
	}

	for _, c := range calls {
		t.Run(c.name+" at info level", func(t *testing.T) {
			buf := captureLogs(t, slog.LevelInfo)
			c.call()
			entry := parseLogLine(t, buf)
			if entry["level"] != c.level {
				t.Fatalf("expected level %s, got %v", c.level, entry["level"])
			}
			for _, k := range c.requiredKeys {
				if _, ok := entry[k]; !ok {
					t.Fatalf("expected key %q in log entry %v", k, entry)
				}
			}
			for _, k := range c.payloadKeys {
				if _, ok := entry[k]; ok {
					t.Fatalf("unexpected payload key %q at info level", k)
				}
			}
			if buf.Len() > 1024 {
				t.Fatalf("expected small log line at info level, got %d bytes", buf.Len())
			}
		})

		t.Run(c.name+" at debug level", func(t *testing.T) {
			buf := captureLogs(t, slog.LevelDebug)
			c.call()
			entry := parseLogLine(t, buf)
			if entry["level"] != c.level {
				t.Fatalf("expected level %s, got %v", c.level, entry["level"])
			}
			for _, k := range append(c.requiredKeys, c.payloadKeys...) {
				if k == "request_id" {
					continue
				}
				v, ok := entry[k].(string)
				if !ok {
					t.Fatalf("expected string key %q in log entry", k)
				}
				if k == "result" && !strings.Contains(v, `"resources"`) {
					t.Fatalf("expected marshaled result, got %q", v[:min(len(v), 100)])
				}
			}
		})
	}
}

func TestWithPayload(t *testing.T) {
	ctx := context.Background()
	base := []any{"method", "tools/list"}

	captureLogs(t, slog.LevelInfo)
	if got := withPayload(ctx, base, "message", map[string]int{"a": 1}); len(got) != len(base) {
		t.Fatalf("expected args unchanged at info level, got %v", got)
	}

	captureLogs(t, slog.LevelDebug)
	got := withPayload(ctx, base, "message", map[string]int{"a": 1}, "result", nil)
	want := []any{"method", "tools/list", "message", `{"a":1}`, "result", ""}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}
