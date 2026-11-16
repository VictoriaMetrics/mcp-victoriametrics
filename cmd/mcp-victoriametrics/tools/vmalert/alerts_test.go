package vmalert

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func TestToolAlertsHandlerLimitOffset(t *testing.T) {
	// vmalert response with two firing alerts and one pending alert.
	const body = `{"status":"success","data":{"alerts":[` +
		`{"id":"1","state":"firing","labels":{"alertgroup":"g"}},` +
		`{"id":"2","state":"firing","labels":{"alertgroup":"g"}},` +
		`{"id":"3","state":"pending","labels":{"alertgroup":"g"}}]}}`

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
		name    string
		state   string
		limit   float64
		offset  float64
		wantIDs []string
	}{
		{name: "limit larger than result count", state: "firing", limit: 50, wantIDs: []string{"1", "2"}},
		{name: "offset past end", state: "firing", limit: 5, offset: 10, wantIDs: []string{}},
		{name: "limit within range", state: "firing", limit: 1, wantIDs: []string{"1"}},
		{name: "offset within range", state: "firing", limit: 1, offset: 1, wantIDs: []string{"2"}},
		{name: "negative offset", state: "firing", limit: 1, offset: -1, wantIDs: []string{"1"}},
		{name: "offset without limit", state: "firing", offset: 1, wantIDs: []string{"2"}},
		{name: "default state with limit", limit: 2, wantIDs: []string{"1", "2"}},
		{name: "default state with offset", offset: 2, wantIDs: []string{"3"}},
		{name: "default state without pagination", wantIDs: []string{"1", "2", "3"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tcr := mcp.CallToolRequest{}
			tcr.Params.Arguments = map[string]any{
				"limit":  tc.limit,
				"offset": tc.offset,
			}
			if tc.state != "" {
				tcr.Params.Arguments.(map[string]any)["state"] = tc.state
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
			gotIDs := make([]string, 0, len(alerts))
			for _, a := range alerts {
				gotIDs = append(gotIDs, a.(map[string]any)["id"].(string))
			}
			if !slices.Equal(gotIDs, tc.wantIDs) {
				t.Errorf("expected alerts %v, got %v", tc.wantIDs, gotIDs)
			}
		})
	}
}
