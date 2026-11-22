package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddlewareLogsSessionID(t *testing.T) {
	testCases := []struct {
		name          string
		target        string
		requestHeader string
		responseID    string
		want          string
	}{
		{name: "streamable HTTP request header", target: "/mcp", requestHeader: "session-1", want: "session-1"},
		{name: "SSE message query parameter", target: "/message?sessionId=session-2", want: "session-2"},
		{name: "initialize response header", target: "/mcp", responseID: "session-3", want: "session-3"},
		{name: "no session", target: "/mcp", want: ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			prev := slog.Default()
			slog.SetDefault(slog.New(slog.NewJSONHandler(&buf, nil)))
			t.Cleanup(func() { slog.SetDefault(prev) })

			handler := (&Logger{}).Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if tc.responseID != "" {
					w.Header().Set("Mcp-Session-Id", tc.responseID)
				}
				w.WriteHeader(http.StatusOK)
			}))
			req := httptest.NewRequest(http.MethodPost, tc.target, nil)
			if tc.requestHeader != "" {
				req.Header.Set("Mcp-Session-Id", tc.requestHeader)
			}
			handler.ServeHTTP(httptest.NewRecorder(), req)

			var completed map[string]any
			for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
				var entry map[string]any
				if err := json.Unmarshal([]byte(line), &entry); err != nil {
					t.Fatalf("cannot parse log line %q: %v", line, err)
				}
				if entry["msg"] == "HTTP request completed" {
					completed = entry
				}
			}
			if completed == nil {
				t.Fatalf("no completion log line in %q", buf.String())
			}
			if got := completed["session_id"]; got != tc.want {
				t.Fatalf("session_id = %q, want %q", got, tc.want)
			}
		})
	}
}
