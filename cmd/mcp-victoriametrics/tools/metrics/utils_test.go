package metrics

import (
	"context"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

// TestGetSelectURLWithDefaultTenant tests that getSelectURL uses default tenant from config
func TestGetSelectURLWithDefaultTenant(t *testing.T) {
	originalEntrypoint := os.Getenv("VM_INSTANCE_ENTRYPOINT")
	originalInstanceType := os.Getenv("VM_INSTANCE_TYPE")
	originalDefaultTenantID := os.Getenv("VM_DEFAULT_TENANT_ID")

	defer func() {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", originalEntrypoint)
		os.Setenv("VM_INSTANCE_TYPE", originalInstanceType)
		os.Setenv("VM_DEFAULT_TENANT_ID", originalDefaultTenantID)
	}()

	testCases := []struct {
		name            string
		instanceType    string
		defaultTenantID string
		requestTenant   any
		expectedURL     string
	}{
		{
			name:            "Cluster mode with default tenant from config",
			instanceType:    "cluster",
			defaultTenantID: "123",
			requestTenant:   "",
			expectedURL:     "http://example.com/select/123/prometheus/api/v1/query",
		},
		{
			name:            "Cluster mode with default tenant 0",
			instanceType:    "cluster",
			defaultTenantID: "0",
			requestTenant:   "",
			expectedURL:     "http://example.com/select/0/prometheus/api/v1/query",
		},
		{
			name:            "Cluster mode with accountID:projectID format",
			instanceType:    "cluster",
			defaultTenantID: "100:200",
			requestTenant:   "",
			expectedURL:     "http://example.com/select/100:200/prometheus/api/v1/query",
		},
		{
			name:            "Cluster mode with request tenant overrides config default",
			instanceType:    "cluster",
			defaultTenantID: "123",
			requestTenant:   "456",
			expectedURL:     "http://example.com/select/456/prometheus/api/v1/query",
		},
		{
			name:            "Cluster mode with empty config default uses 0",
			instanceType:    "cluster",
			defaultTenantID: "",
			requestTenant:   "",
			expectedURL:     "http://example.com/select/0/prometheus/api/v1/query",
		},
		{
			name:            "Cluster mode with numeric request tenant",
			instanceType:    "cluster",
			defaultTenantID: "123",
			requestTenant:   float64(456),
			expectedURL:     "http://example.com/select/456/prometheus/api/v1/query",
		},
		{
			name:            "Single mode ignores tenant",
			instanceType:    "single",
			defaultTenantID: "123",
			requestTenant:   "",
			expectedURL:     "http://example.com/api/v1/query",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables for this test case
			os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
			os.Setenv("VM_INSTANCE_TYPE", tc.instanceType)
			os.Setenv("VM_DEFAULT_TENANT_ID", tc.defaultTenantID)

			cfg, err := config.InitConfig()
			if err != nil {
				t.Fatalf("Failed to create config: %v", err)
			}

			tcr := mcp.CallToolRequest{}
			if tc.requestTenant != nil && tc.requestTenant != "" {
				tcr.Params.Arguments = map[string]any{"tenant": tc.requestTenant}
			}

			url, err := getSelectURL(context.Background(), cfg, tcr, "api", "v1", "query")
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if url != tc.expectedURL {
				t.Errorf("Expected URL %q, got %q", tc.expectedURL, url)
			}
		})
	}
}
