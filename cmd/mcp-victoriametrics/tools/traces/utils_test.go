package traces

import (
	"context"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

// TestGetToolReqTenant tests the GetToolReqTenant function
func TestGetToolReqTenant(t *testing.T) {
	testCases := []struct {
		name              string
		tenant            any
		expectedAccountID string
		expectedProjectID string
		expectError       bool
	}{
		{
			name:              "Empty tenant returns default 0:0",
			tenant:            "",
			expectedAccountID: "",
			expectedProjectID: "",
			expectError:       false,
		},
		{
			name:              "Single number tenant",
			tenant:            "123",
			expectedAccountID: "123",
			expectedProjectID: "0",
			expectError:       false,
		},
		{
			name:              "Numeric tenant",
			tenant:            float64(123),
			expectedAccountID: "123",
			expectedProjectID: "0",
			expectError:       false,
		},
		{
			name:        "Fractional numeric tenant",
			tenant:      1.5,
			expectError: true,
		},
		{
			name:              "Full tenant format",
			tenant:            "123:456",
			expectedAccountID: "123",
			expectedProjectID: "456",
			expectError:       false,
		},
		{
			name:              "Zero tenant",
			tenant:            "0:0",
			expectedAccountID: "0",
			expectedProjectID: "0",
			expectError:       false,
		},
		{
			name:              "Account only with colon",
			tenant:            "123:",
			expectedAccountID: "123",
			expectedProjectID: "0",
			expectError:       false,
		},
		{
			name:              "Project only with colon",
			tenant:            ":456",
			expectedAccountID: "0",
			expectedProjectID: "456",
			expectError:       false,
		},
		{
			name:        "Invalid tenant format - too many colons",
			tenant:      "1:2:3",
			expectError: true,
		},
		{
			name:        "Invalid tenant format - non-numeric",
			tenant:      "abc:def",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tcr := mcp.CallToolRequest{}
			tcr.Params.Arguments = map[string]any{
				"tenant": tc.tenant,
			}

			accountID, projectID, err := GetToolReqTenant(tcr)

			if tc.expectError {
				if err == nil {
					t.Error("Expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
				return
			}

			if accountID != tc.expectedAccountID {
				t.Errorf("Expected accountID %q, got %q", tc.expectedAccountID, accountID)
			}
			if projectID != tc.expectedProjectID {
				t.Errorf("Expected projectID %q, got %q", tc.expectedProjectID, projectID)
			}
		})
	}
}

// TestCreateSelectRequest_DefaultTenant tests that CreateSelectRequest uses default tenant from config
func TestCreateSelectRequest_DefaultTenant(t *testing.T) {
	t.Setenv("VM_INSTANCE_BACKEND", "traces")

	// Save original environment variables
	originalEntrypoint := os.Getenv("VT_INSTANCE_ENTRYPOINT")
	originalDefaultTenantID := os.Getenv("VT_DEFAULT_TENANT_ID")

	// Restore environment variables after test
	defer func() {
		os.Setenv("VT_INSTANCE_ENTRYPOINT", originalEntrypoint)
		os.Setenv("VT_DEFAULT_TENANT_ID", originalDefaultTenantID)
	}()

	testCases := []struct {
		name              string
		defaultTenantID   string
		requestTenant     string
		expectedAccountID string
		expectedProjectID string
	}{
		{
			name:              "Empty request tenant uses default from config",
			defaultTenantID:   "100:200",
			requestTenant:     "",
			expectedAccountID: "100",
			expectedProjectID: "200",
		},
		{
			name:              "Request tenant overrides default",
			defaultTenantID:   "100:200",
			requestTenant:     "300:400",
			expectedAccountID: "300",
			expectedProjectID: "400",
		},
		{
			name:              "Empty default tenant uses 0:0",
			defaultTenantID:   "",
			requestTenant:     "",
			expectedAccountID: "0",
			expectedProjectID: "0",
		},
		{
			name:              "Partial request tenant (account only) uses default project",
			defaultTenantID:   "100:200",
			requestTenant:     "500",
			expectedAccountID: "500",
			expectedProjectID: "0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv("VT_INSTANCE_ENTRYPOINT", "http://example.com")
			os.Setenv("VT_DEFAULT_TENANT_ID", tc.defaultTenantID)

			// Initialize config
			cfg, err := config.InitConfig()
			if err != nil {
				t.Fatalf("Failed to init config: %v", err)
			}

			// Create tool request
			tcr := mcp.CallToolRequest{}
			tcr.Params.Arguments = map[string]any{
				"tenant": tc.requestTenant,
			}

			// Call CreateSelectRequest
			req, err := CreateSelectRequest(context.Background(), cfg, tcr, "query")
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Check headers
			accountID := req.Header.Get("AccountID")
			projectID := req.Header.Get("ProjectID")

			if accountID != tc.expectedAccountID {
				t.Errorf("Expected AccountID %q, got %q", tc.expectedAccountID, accountID)
			}
			if projectID != tc.expectedProjectID {
				t.Errorf("Expected ProjectID %q, got %q", tc.expectedProjectID, projectID)
			}
		})
	}
}
