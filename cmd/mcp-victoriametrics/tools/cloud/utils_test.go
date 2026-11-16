package cloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func TestListToolsFilterByDeploymentType(t *testing.T) {
	var queries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.Path+"?"+r.URL.RawQuery)
		_, _ = w.Write([]byte("[]"))
	}))
	defer srv.Close()

	for backend, want := range map[string]string{
		"metrics": "type=single_node&type=cluster",
		"logs":    "type=vlogs_single",
		"traces":  "type=vtraces_single",
	} {
		t.Run(backend, func(t *testing.T) {
			queries = nil
			t.Setenv("VM_INSTANCE_BACKEND", backend)
			t.Setenv("VM_INSTANCE_ENTRYPOINT", "")
			t.Setenv("VMC_API_KEY", "test-key")
			t.Setenv("VMC_API_BASE_URL", srv.URL)
			cfg, err := config.InitConfig()
			if err != nil {
				t.Fatalf("Failed to init config: %v", err)
			}

			for _, handler := range []func(context.Context, *config.Config, mcp.CallToolRequest) (*mcp.CallToolResult, error){
				toolDeploymentsHandler,
				toolTiersHandler,
			} {
				result, err := handler(context.Background(), cfg, mcp.CallToolRequest{})
				if err != nil || result.IsError {
					t.Fatalf("Unexpected handler failure: err=%v, result=%v", err, result)
				}
			}

			wantQueries := []string{"/api/v1/deployments?" + want, "/api/v1/tiers?" + want}
			if len(queries) != len(wantQueries) {
				t.Fatalf("Captured queries %q, want %q", queries, wantQueries)
			}
			for i := range wantQueries {
				if queries[i] != wantQueries[i] {
					t.Errorf("Query %d = %q, want %q", i, queries[i], wantQueries[i])
				}
			}
		})
	}
}
