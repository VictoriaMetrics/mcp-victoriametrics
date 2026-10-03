package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func TestToolAlertsHandlerLimitOffset(t *testing.T) {
	// vmalert response with a single firing alert.
	const body = `{"status":"success","data":{"alerts":[` +
		`{"id":"1","state":"firing","labels":{"alertgroup":"g"}}]}}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	t.Setenv("VM_INSTANCE_ENTRYPOINT", srv.URL)
	t.Setenv("VM_INSTANCE_TYPE", "single")
	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	testCases := []struct {
		name      string
		limit     float64
		offset    float64
		wantCount int
	}{
		{name: "limit larger than result count", limit: 50, offset: 0, wantCount: 1},
		{name: "offset past end", limit: 5, offset: 10, wantCount: 0},
		{name: "limit within range", limit: 1, offset: 0, wantCount: 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tcr := mcp.CallToolRequest{}
			tcr.Params.Arguments = map[string]any{
				"state":  "firing",
				"limit":  tc.limit,
				"offset": tc.offset,
			}

			result, err := toolAlertsHandler(context.Background(), cfg, tcr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.IsError {
				t.Fatalf("expected success, got error: %v", result.Content)
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatal("expected TextContent")
			}
			var parsed map[string]any
			if err := json.Unmarshal([]byte(textContent.Text), &parsed); err != nil {
				t.Fatalf("failed to parse JSON: %v", err)
			}
			data := parsed["data"].(map[string]any)
			alerts := data["alerts"].([]any)
			if len(alerts) != tc.wantCount {
				t.Errorf("expected %d alerts, got %d", tc.wantCount, len(alerts))
			}
		})
	}
}
