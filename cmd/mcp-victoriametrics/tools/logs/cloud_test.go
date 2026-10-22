package logs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

func newCloudConfig(t *testing.T) *config.Config {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp any
		switch r.URL.Path {
		case "/api/v1/deployments/0b6f5f1e-6a3c-4f1e-9a51-6c1d2f0a7b01":
			resp = vmcloud.DeploymentInfo{ID: "0b6f5f1e-6a3c-4f1e-9a51-6c1d2f0a7b01", Type: vmcloud.DeploymentTypeVLogs, AccessEndpoint: "https://vl.example.com"}
		case "/api/v1/deployments/7c2d9e4a-1b8f-4c3e-8d6a-2e5f9b0c4d12":
			resp = vmcloud.DeploymentInfo{ID: "7c2d9e4a-1b8f-4c3e-8d6a-2e5f9b0c4d12", Type: vmcloud.DeploymentTypeSingleNode, AccessEndpoint: "https://vm.example.com"}
		case "/api/v1/deployments/0b6f5f1e-6a3c-4f1e-9a51-6c1d2f0a7b01/access_tokens":
			resp = vmcloud.AccessTokensList{{ID: "write", Type: vmcloud.AccessModeWrite}, {ID: "read", Type: vmcloud.AccessModeRead}}
		case "/api/v1/deployments/0b6f5f1e-6a3c-4f1e-9a51-6c1d2f0a7b01/access_tokens/read":
			resp = vmcloud.AccessToken{ID: "read", Type: vmcloud.AccessModeRead, Secret: "read-secret"}
		default:
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	t.Setenv("VM_INSTANCE_BACKEND", "logs")
	t.Setenv("VM_INSTANCE_ENTRYPOINT", "")
	t.Setenv("VL_INSTANCE_ENTRYPOINT", "")
	t.Setenv("VMC_API_KEY", "test-key")
	t.Setenv("VMC_API_BASE_URL", srv.URL)
	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("Failed to init config: %v", err)
	}
	return cfg
}

func TestCreateSelectRequestCloud(t *testing.T) {
	cfg := newCloudConfig(t)

	tcr := mcp.CallToolRequest{}
	tcr.Params.Arguments = map[string]any{"deployment_id": "0b6f5f1e-6a3c-4f1e-9a51-6c1d2f0a7b01"}
	req, err := CreateSelectRequest(context.Background(), cfg, tcr, "query")
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if got, want := req.URL.String(), "https://vl.example.com/select/logsql/query"; got != want {
		t.Errorf("Expected URL %q, got %q", want, got)
	}
	if got, want := req.Header.Get("Authorization"), "Bearer read-secret"; got != want {
		t.Errorf("Expected Authorization %q, got %q", want, got)
	}

	tcr.Params.Arguments = map[string]any{"deployment_id": "7c2d9e4a-1b8f-4c3e-8d6a-2e5f9b0c4d12"}
	if _, err := CreateSelectRequest(context.Background(), cfg, tcr, "query"); err == nil {
		t.Error("Expected error for VictoriaMetrics deployment, got nil")
	}

	tcr.Params.Arguments = map[string]any{}
	if _, err := CreateSelectRequest(context.Background(), cfg, tcr, "query"); err == nil {
		t.Error("Expected error for missing deployment_id, got nil")
	}
}

func TestWithCloudDeploymentID(t *testing.T) {
	cfg := newCloudConfig(t)

	tool := tools.WithCloudDeploymentID(cfg, toolQuery)
	if _, ok := tool.InputSchema.Properties["deployment_id"]; !ok {
		t.Error("Expected deployment_id property in cloud mode")
	}
	if !slices.Contains(tool.InputSchema.Required, "deployment_id") {
		t.Error("Expected deployment_id to be required in cloud mode")
	}
	if _, ok := toolQuery.InputSchema.Properties["deployment_id"]; ok {
		t.Error("Expected original tool definition to stay unchanged")
	}
}
