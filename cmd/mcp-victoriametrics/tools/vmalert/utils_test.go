package vmalert

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func TestCreateRequestPath(t *testing.T) {
	testCases := []struct {
		name        string
		env         map[string]string
		tenant      string
		wantPath    string
		wantAccount string
		wantProject string
	}{
		{
			name:     "metrics single",
			env:      map[string]string{"VM_INSTANCE_BACKEND": "metrics", "VM_INSTANCE_TYPE": "single"},
			wantPath: "/vmalert/api/v1/alerts",
		},
		{
			name:     "metrics cluster",
			env:      map[string]string{"VM_INSTANCE_BACKEND": "metrics", "VM_INSTANCE_TYPE": "cluster"},
			wantPath: "/select/0/prometheus/vmalert/api/v1/alerts",
		},
		{
			name:     "metrics cluster with tenant",
			env:      map[string]string{"VM_INSTANCE_BACKEND": "metrics", "VM_INSTANCE_TYPE": "cluster"},
			tenant:   "1:2",
			wantPath: "/select/1:2/prometheus/vmalert/api/v1/alerts",
		},
		{
			name:        "logs with default tenant",
			env:         map[string]string{"VM_INSTANCE_BACKEND": "logs"},
			wantPath:    "/select/vmalert/api/v1/alerts",
			wantAccount: "0",
			wantProject: "0",
		},
		{
			name:        "logs with tenant",
			env:         map[string]string{"VM_INSTANCE_BACKEND": "logs"},
			tenant:      "1:2",
			wantPath:    "/select/vmalert/api/v1/alerts",
			wantAccount: "1",
			wantProject: "2",
		},
		{
			name:        "traces with default tenant from config",
			env:         map[string]string{"VM_INSTANCE_BACKEND": "traces", "VM_DEFAULT_TENANT_ID": "3:4"},
			wantPath:    "/select/vmalert/api/v1/alerts",
			wantAccount: "3",
			wantProject: "4",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://localhost:8428")
			t.Setenv("VM_INSTANCE_TYPE", "")
			t.Setenv("VMC_API_KEY", "")
			t.Setenv("VM_DEFAULT_TENANT_ID", "")
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			cfg, err := config.InitConfig()
			if err != nil {
				t.Fatalf("failed to init config: %v", err)
			}

			tcr := mcp.CallToolRequest{}
			if tc.tenant != "" {
				tcr.Params.Arguments = map[string]any{"tenant": tc.tenant}
			}
			req, err := createRequest(context.Background(), cfg, tcr, "api", "v1", "alerts")
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			if req.URL.Path != tc.wantPath {
				t.Errorf("expected path %q, got %q", tc.wantPath, req.URL.Path)
			}
			if got := req.Header.Get("AccountID"); got != tc.wantAccount {
				t.Errorf("expected AccountID header %q, got %q", tc.wantAccount, got)
			}
			if got := req.Header.Get("ProjectID"); got != tc.wantProject {
				t.Errorf("expected ProjectID header %q, got %q", tc.wantProject, got)
			}
		})
	}
}
