package config

import (
	"testing"

	"github.com/VictoriaMetrics/VictoriaLogs/lib/logstorage"
)

func clearBackendEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{
		"VM_INSTANCE_BACKEND", "VM_INSTANCE_ENTRYPOINT", "VM_INSTANCE_TYPE", "VM_INSTANCE_BEARER_TOKEN", "VM_INSTANCE_HEADERS", "VM_DEFAULT_TENANT_ID", "VMC_API_KEY",
		"VL_INSTANCE_ENTRYPOINT", "VL_INSTANCE_BEARER_TOKEN", "VL_INSTANCE_HEADERS", "VL_DEFAULT_TENANT_ID",
		"VT_INSTANCE_ENTRYPOINT", "VT_INSTANCE_BEARER_TOKEN", "VT_INSTANCE_HEADERS", "VT_DEFAULT_TENANT_ID",
		"MCP_DISABLED_TOOLS", "MCP_LISTEN_ADDR", "MCP_SSE_ADDR", "MCP_SERVER_MODE",
	} {
		t.Setenv(name, "")
	}
}

func TestInitConfigBackend(t *testing.T) {
	testCases := []struct {
		name            string
		env             map[string]string
		expectedBackend Backend
		expectedURL     string
		expectedToken   string
		expectedAddr    string
		expectedTenant  logstorage.TenantID
		expectedCloud   bool
		expectError     bool
	}{
		{
			name:            "Metrics by default",
			env:             map[string]string{"VM_INSTANCE_ENTRYPOINT": "http://vm", "VM_INSTANCE_TYPE": "single"},
			expectedBackend: BackendMetrics,
			expectedURL:     "http://vm",
			expectedAddr:    "localhost:8080",
		},
		{
			name:            "Explicit logs backend",
			env:             map[string]string{"VM_INSTANCE_BACKEND": "logs", "VM_INSTANCE_ENTRYPOINT": "http://vl", "VM_INSTANCE_BEARER_TOKEN": "token", "VM_DEFAULT_TENANT_ID": "1:2"},
			expectedBackend: BackendLogs,
			expectedURL:     "http://vl",
			expectedToken:   "token",
			expectedAddr:    "localhost:8080",
			expectedTenant:  logstorage.TenantID{AccountID: 1, ProjectID: 2},
		},
		{
			name:            "Explicit traces backend is case insensitive",
			env:             map[string]string{"VM_INSTANCE_BACKEND": "Traces", "VM_INSTANCE_ENTRYPOINT": "http://vt"},
			expectedBackend: BackendTraces,
			expectedURL:     "http://vt",
			expectedAddr:    "localhost:8080",
		},
		{
			name:            "Legacy VL variables for logs backend",
			env:             map[string]string{"VM_INSTANCE_BACKEND": "logs", "VL_INSTANCE_ENTRYPOINT": "http://vl", "VL_INSTANCE_BEARER_TOKEN": "vl-token", "VL_DEFAULT_TENANT_ID": "3:4"},
			expectedBackend: BackendLogs,
			expectedURL:     "http://vl",
			expectedToken:   "vl-token",
			expectedAddr:    "localhost:8080",
			expectedTenant:  logstorage.TenantID{AccountID: 3, ProjectID: 4},
		},
		{
			name:            "Legacy VT variables for traces backend",
			env:             map[string]string{"VM_INSTANCE_BACKEND": "traces", "VT_INSTANCE_ENTRYPOINT": "http://vt"},
			expectedBackend: BackendTraces,
			expectedURL:     "http://vt",
			expectedAddr:    "localhost:8080",
		},
		{
			name:            "VM variables take precedence over legacy ones",
			env:             map[string]string{"VM_INSTANCE_BACKEND": "logs", "VM_INSTANCE_ENTRYPOINT": "http://vm", "VL_INSTANCE_ENTRYPOINT": "http://vl"},
			expectedBackend: BackendLogs,
			expectedURL:     "http://vm",
			expectedAddr:    "localhost:8080",
		},
		{
			name:        "Legacy variables without backend",
			env:         map[string]string{"VL_INSTANCE_ENTRYPOINT": "http://vl"},
			expectError: true,
		},
		{
			name:        "Unknown backend",
			env:         map[string]string{"VM_INSTANCE_BACKEND": "profiles", "VM_INSTANCE_ENTRYPOINT": "http://vm"},
			expectError: true,
		},
		{
			name:        "Logs backend without entrypoint",
			env:         map[string]string{"VM_INSTANCE_BACKEND": "logs"},
			expectError: true,
		},
		{
			name:            "Cloud API key for logs",
			env:             map[string]string{"VM_INSTANCE_BACKEND": "logs", "VMC_API_KEY": "key"},
			expectedBackend: BackendLogs,
			expectedAddr:    "localhost:8080",
			expectedCloud:   true,
		},
		{
			name:            "Cloud API key for traces",
			env:             map[string]string{"VM_INSTANCE_BACKEND": "traces", "VMC_API_KEY": "key"},
			expectedBackend: BackendTraces,
			expectedAddr:    "localhost:8080",
			expectedCloud:   true,
		},
		{
			name:        "Cloud API key and entrypoint for logs",
			env:         map[string]string{"VM_INSTANCE_BACKEND": "logs", "VM_INSTANCE_ENTRYPOINT": "http://vl", "VMC_API_KEY": "key"},
			expectError: true,
		},
		{
			name:        "Invalid logs tenant",
			env:         map[string]string{"VM_INSTANCE_BACKEND": "logs", "VM_INSTANCE_ENTRYPOINT": "http://vl", "VM_DEFAULT_TENANT_ID": "1:2:3"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearBackendEnv(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			cfg, err := InitConfig()
			if tc.expectError {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}
			if cfg.Backend() != tc.expectedBackend {
				t.Errorf("Expected backend %q, got %q", tc.expectedBackend, cfg.Backend())
			}
			if cfg.IsCloud() != tc.expectedCloud {
				t.Errorf("Expected IsCloud() %v, got %v", tc.expectedCloud, cfg.IsCloud())
			}
			if !tc.expectedCloud && cfg.EntryPointURL().String() != tc.expectedURL {
				t.Errorf("Expected entrypoint %q, got %q", tc.expectedURL, cfg.EntryPointURL().String())
			}
			if cfg.BearerToken() != tc.expectedToken {
				t.Errorf("Expected bearer token %q, got %q", tc.expectedToken, cfg.BearerToken())
			}
			if cfg.ListenAddr() != tc.expectedAddr {
				t.Errorf("Expected listen addr %q, got %q", tc.expectedAddr, cfg.ListenAddr())
			}
			if cfg.DefaultLogsTenantID() != tc.expectedTenant {
				t.Errorf("Expected logs tenant %v, got %v", tc.expectedTenant, cfg.DefaultLogsTenantID())
			}
		})
	}
}

func TestInitConfigDisabledToolsDefaults(t *testing.T) {
	clearBackendEnv(t)
	t.Setenv("VM_INSTANCE_BACKEND", "logs")
	t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://vl")

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cfg.IsToolDisabled("flags") {
		t.Error("Expected flags tool to be enabled by default for logs backend")
	}
}
