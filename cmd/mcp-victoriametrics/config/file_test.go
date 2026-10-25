package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/VictoriaMetrics/VictoriaLogs/lib/logstorage"
)

func writeBackendsFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	return path
}

func TestInitConfigBackendsFile(t *testing.T) {
	clearBackendEnv(t)
	os.Unsetenv("MCP_DISABLED_TOOLS")
	t.Setenv("MCP_CONFIG_FILE", writeBackendsFile(t, `
backends:
  - name: prod
    type: metrics
    entrypoint: http://vmselect:8481
    instance_type: cluster
    bearer_token: prod-token
    headers:
      X-Scope: prod
    default_tenant_id: "1:2"
  - name: logs
    type: Logs
    entrypoint: http://vlogs:9428
    default_tenant_id: "3:4"
  - name: traces
    type: traces
    entrypoint: http://vtraces:10428
`))

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	backends := cfg.Backends()
	if len(backends) != 3 {
		t.Fatalf("Expected 3 backends, got %d", len(backends))
	}

	prod := backends[0]
	if prod.Name() != "prod" || prod.Backend() != BackendMetrics || !prod.IsCluster() {
		t.Errorf("Unexpected prod backend: name=%s type=%s cluster=%v", prod.Name(), prod.Backend(), prod.IsCluster())
	}
	if prod.EntryPointURL().String() != "http://vmselect:8481" || prod.BearerToken() != "prod-token" {
		t.Errorf("Unexpected prod endpoint %s or token %s", prod.EntryPointURL(), prod.BearerToken())
	}
	if prod.CustomHeaders()["X-Scope"] != "prod" || prod.DefaultTenantID() != "1:2" {
		t.Errorf("Unexpected prod headers %v or tenant %s", prod.CustomHeaders(), prod.DefaultTenantID())
	}
	if !prod.IsToolDisabled("flags") {
		t.Error("Expected flags to be disabled by default for metrics backend")
	}

	logs := backends[1]
	if logs.Backend() != BackendLogs || logs.DefaultLogsTenantID() != (logstorage.TenantID{AccountID: 3, ProjectID: 4}) {
		t.Errorf("Unexpected logs backend: type=%s tenant=%v", logs.Backend(), logs.DefaultLogsTenantID())
	}
	if logs.IsToolDisabled("flags") {
		t.Error("Expected flags to be enabled by default for logs backend")
	}

	if got := cfg.BackendTypes(); len(got) != 3 {
		t.Errorf("Expected 3 backend types, got %v", got)
	}
	if got := cfg.ProductName(); got != "VictoriaMetrics, VictoriaLogs and VictoriaTraces" {
		t.Errorf("Unexpected product name %q", got)
	}
	if cfg.IsCloud() {
		t.Error("Expected non-cloud config")
	}
}

func TestInitConfigBackendsFileLegacyToolNames(t *testing.T) {
	clearBackendEnv(t)
	t.Setenv("MCP_DISABLED_TOOLS", "query")
	t.Setenv("MCP_CONFIG_FILE", writeBackendsFile(t, `
backends:
  - {name: vm, type: metrics, entrypoint: "http://vm", instance_type: single}
  - {name: vl, type: logs, entrypoint: "http://vl"}
`))

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !cfg.Backends()[0].IsToolDisabled("metrics_query") || !cfg.Backends()[1].IsToolDisabled("logs_query") {
		t.Error("Expected legacy query name to disable query tools of both backends")
	}
	if got := cfg.LegacyToolNames()["query"]; got != "metrics_query, logs_query" {
		t.Errorf("Unexpected legacy tool names %q", got)
	}
}

func TestInitConfigBackendsFileErrors(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		env     map[string]string
		errPart string
	}{
		{
			name:    "Empty backends",
			content: "backends: []",
			errPart: "at least one backend",
		},
		{
			name:    "Unknown field",
			content: "backends:\n  - {name: a, type: logs, entrypoint: http://a, token: x}",
			errPart: "cannot parse file",
		},
		{
			name:    "Invalid name",
			content: "backends:\n  - {name: \"a b\", type: logs, entrypoint: http://a}",
			errPart: "must match",
		},
		{
			name:    "Duplicate name",
			content: "backends:\n  - {name: a, type: logs, entrypoint: http://a}\n  - {name: a, type: traces, entrypoint: http://b}",
			errPart: "duplicate name",
		},
		{
			name:    "Invalid type",
			content: "backends:\n  - {name: a, type: profiles, entrypoint: http://a}",
			errPart: "type must be",
		},
		{
			name:    "Missing entrypoint",
			content: "backends:\n  - {name: a, type: logs}",
			errPart: "entrypoint is not set",
		},
		{
			name:    "Missing metrics instance type",
			content: "backends:\n  - {name: a, type: metrics, entrypoint: http://a}",
			errPart: "instance_type must be",
		},
		{
			name:    "Instance type for logs",
			content: "backends:\n  - {name: a, type: logs, entrypoint: http://a, instance_type: single}",
			errPart: "instance_type is supported only",
		},
		{
			name:    "Invalid tenant",
			content: "backends:\n  - {name: a, type: logs, entrypoint: http://a, default_tenant_id: x}",
			errPart: "default_tenant_id",
		},
		{
			name:    "Conflicts with env entrypoint",
			content: "backends:\n  - {name: a, type: logs, entrypoint: http://a}",
			env:     map[string]string{"VM_INSTANCE_ENTRYPOINT": "http://vm"},
			errPart: "cannot be set at the same time",
		},
		{
			name:    "Conflicts with cloud",
			content: "backends:\n  - {name: a, type: logs, entrypoint: http://a}",
			env:     map[string]string{"VMC_API_KEY": "key"},
			errPart: "cannot be set at the same time",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearBackendEnv(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			t.Setenv("MCP_CONFIG_FILE", writeBackendsFile(t, tc.content))
			_, err := InitConfig()
			if err == nil {
				t.Fatalf("Expected error containing %q, got nil", tc.errPart)
			}
			if !strings.Contains(err.Error(), tc.errPart) {
				t.Errorf("Expected error containing %q, got: %v", tc.errPart, err)
			}
		})
	}
}

func TestInitConfigCloudBackendTypes(t *testing.T) {
	clearBackendEnv(t)
	os.Unsetenv("MCP_DISABLED_TOOLS")
	t.Setenv("VMC_API_KEY", "key")
	t.Setenv("VMC_BACKEND_TYPES", "metrics, Logs,traces,logs")

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	backends := cfg.Backends()
	if len(backends) != 3 {
		t.Fatalf("Expected 3 backends, got %d", len(backends))
	}
	for i, expected := range []Backend{BackendMetrics, BackendLogs, BackendTraces} {
		if backends[i].Backend() != expected || backends[i].Name() != string(expected) || !backends[i].IsCloud() {
			t.Errorf("Unexpected backend %d: name=%s type=%s cloud=%v", i, backends[i].Name(), backends[i].Backend(), backends[i].IsCloud())
		}
	}
	if !backends[0].IsToolDisabled("flags") || backends[1].IsToolDisabled("flags") {
		t.Error("Expected flags to be disabled only for metrics backend")
	}

	t.Setenv("VMC_BACKEND_TYPES", "logs")
	cfg, err = InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(cfg.Backends()) != 1 || cfg.Backend() != BackendLogs {
		t.Errorf("Expected single logs backend, got %d backends of type %s", len(cfg.Backends()), cfg.Backend())
	}
}

func TestInitConfigBackendsFileCloud(t *testing.T) {
	clearBackendEnv(t)
	os.Unsetenv("MCP_DISABLED_TOOLS")
	t.Setenv("MCP_CONFIG_FILE", writeBackendsFile(t, `
cloud:
  api_key: key
  base_url: https://api.example.com
  backend_types: [metrics, Logs, logs]
  default_tenant_id: "1:2"
backends:
  - {name: onprem, type: traces, entrypoint: "http://vt"}
`))

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	backends := cfg.Backends()
	if len(backends) != 3 {
		t.Fatalf("Expected 3 backends, got %d", len(backends))
	}
	if backends[0].Name() != "onprem" || backends[0].IsCloud() {
		t.Errorf("Expected static onprem backend first, got name=%s cloud=%v", backends[0].Name(), backends[0].IsCloud())
	}
	for i, expected := range []Backend{BackendMetrics, BackendLogs} {
		b := backends[i+1]
		if b.Name() != "cloud" || b.Backend() != expected || !b.IsCloud() {
			t.Errorf("Unexpected cloud backend %d: name=%s type=%s cloud=%v", i, b.Name(), b.Backend(), b.IsCloud())
		}
	}
	if backends[1].DefaultTenantID() != "1:2" {
		t.Errorf("Unexpected cloud metrics tenant %s", backends[1].DefaultTenantID())
	}
	if !cfg.IsCloud() || cfg.VMC().BaseURL() != "https://api.example.com" {
		t.Error("Expected top-level config to have cloud client with configured base URL")
	}
	if got := cfg.CloudBackendTypes(); len(got) != 2 || got[0] != BackendMetrics || got[1] != BackendLogs {
		t.Errorf("Unexpected cloud backend types %v", got)
	}
	if !backends[1].IsToolDisabled("flags") || backends[2].IsToolDisabled("flags") {
		t.Error("Expected flags to be disabled only for cloud metrics backend")
	}
}

func TestInitConfigBackendsFileCloudOnly(t *testing.T) {
	clearBackendEnv(t)
	t.Setenv("MCP_CONFIG_FILE", writeBackendsFile(t, `
cloud:
  name: vmcloud
  api_key: key
  backend_types: [traces]
`))

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(cfg.Backends()) != 1 || cfg.Backends()[0].Name() != "vmcloud" || cfg.Backend() != BackendTraces || !cfg.Backends()[0].IsCloud() {
		t.Errorf("Unexpected backends: %d, type %s", len(cfg.Backends()), cfg.Backend())
	}
}

func TestInitConfigBackendsFileCloudErrors(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		env     map[string]string
		errPart string
	}{
		{name: "Missing api key", content: "cloud: {backend_types: [metrics]}", errPart: "api_key is not set"},
		{name: "Missing backend types", content: "cloud: {api_key: key}", errPart: "at least one backend type"},
		{name: "Invalid backend type", content: "cloud: {api_key: key, backend_types: [metrics, profiles]}", errPart: "backend_types"},
		{name: "Empty backend type", content: "cloud: {api_key: key, backend_types: [\"\"]}", errPart: "backend_types"},
		{name: "Invalid name", content: "cloud: {name: \"my cloud\", api_key: key, backend_types: [metrics]}", errPart: "must match"},
		{name: "Name used by backend", content: "cloud: {api_key: key, backend_types: [metrics]}\nbackends:\n  - {name: cloud, type: logs, entrypoint: http://a}", errPart: "already used"},
		{name: "Unknown field", content: "cloud: {api_key: key, backend_types: [metrics], token: x}", errPart: "cannot parse file"},
		{name: "Conflicts with env backend types", content: "cloud: {api_key: key, backend_types: [metrics]}", env: map[string]string{"VMC_BACKEND_TYPES": "logs"}, errPart: "cannot be set at the same time"},
		{name: "Conflicts with env base url", content: "cloud: {api_key: key, backend_types: [metrics]}", env: map[string]string{"VMC_API_BASE_URL": "https://api.example.com"}, errPart: "cannot be set at the same time"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearBackendEnv(t)
			t.Setenv("VMC_API_BASE_URL", "")
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			t.Setenv("MCP_CONFIG_FILE", writeBackendsFile(t, tc.content))
			_, err := InitConfig()
			if err == nil {
				t.Fatalf("Expected error containing %q, got nil", tc.errPart)
			}
			if !strings.Contains(err.Error(), tc.errPart) {
				t.Errorf("Expected error containing %q, got: %v", tc.errPart, err)
			}
		})
	}
}

func TestInitConfigCloudBackendTypesErrors(t *testing.T) {
	testCases := []struct {
		name    string
		env     map[string]string
		errPart string
	}{
		{
			name:    "Without cloud",
			env:     map[string]string{"VM_INSTANCE_ENTRYPOINT": "http://vm", "VM_INSTANCE_TYPE": "single", "VMC_BACKEND_TYPES": "logs"},
			errPart: "only with VMC_API_KEY",
		},
		{
			name:    "With VM_INSTANCE_BACKEND",
			env:     map[string]string{"VMC_API_KEY": "key", "VM_INSTANCE_BACKEND": "logs", "VMC_BACKEND_TYPES": "logs"},
			errPart: "cannot be set at the same time",
		},
		{
			name:    "Invalid type",
			env:     map[string]string{"VMC_API_KEY": "key", "VMC_BACKEND_TYPES": "metrics,profiles"},
			errPart: "must be",
		},
		{
			name:    "Empty list",
			env:     map[string]string{"VMC_API_KEY": "key", "VMC_BACKEND_TYPES": " , "},
			errPart: "at least one backend type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearBackendEnv(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			_, err := InitConfig()
			if err == nil {
				t.Fatalf("Expected error containing %q, got nil", tc.errPart)
			}
			if !strings.Contains(err.Error(), tc.errPart) {
				t.Errorf("Expected error containing %q, got: %v", tc.errPart, err)
			}
		})
	}
}

func TestInitConfigToolsets(t *testing.T) {
	clearBackendEnv(t)
	t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://vm")
	t.Setenv("VM_INSTANCE_TYPE", "single")

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cfg.ToolsetMode(ToolsetAdmin) != ToolsetOn || cfg.HasDynamicToolsets() {
		t.Errorf("Expected admin toolset to be on by default, got %s", cfg.ToolsetMode(ToolsetAdmin))
	}

	t.Setenv("MCP_TOOLSETS_ADMIN", "Dynamic")
	cfg, err = InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cfg.ToolsetMode(ToolsetAdmin) != ToolsetDynamic || !cfg.HasDynamicToolsets() {
		t.Errorf("Expected dynamic admin toolset, got %s", cfg.ToolsetMode(ToolsetAdmin))
	}

	t.Setenv("MCP_TOOLSETS_ADMIN", "")
	t.Setenv("VM_INSTANCE_ENTRYPOINT", "")
	t.Setenv("VM_INSTANCE_TYPE", "")
	t.Setenv("MCP_CONFIG_FILE", writeBackendsFile(t, "toolsets: {admin: off}\nbackends:\n  - {name: a, type: logs, entrypoint: http://a}"))
	cfg, err = InitConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if cfg.ToolsetMode(ToolsetAdmin) != ToolsetOff {
		t.Errorf("Expected admin toolset off from file, got %s", cfg.ToolsetMode(ToolsetAdmin))
	}
}

func TestInitConfigToolsetsErrors(t *testing.T) {
	testCases := []struct {
		name    string
		env     map[string]string
		file    string
		errPart string
	}{
		{name: "Unknown env toolset", env: map[string]string{"MCP_TOOLSETS_DEBUG": "on"}, errPart: `unknown toolset "debug"`},
		{name: "Invalid env mode", env: map[string]string{"MCP_TOOLSETS_ADMIN": "lazy"}, errPart: "toolset mode must be"},
		{name: "Unknown file toolset", file: "toolsets: {debug: on}", errPart: `unknown toolset "debug"`},
		{name: "Invalid file mode", file: "toolsets: {admin: lazy}", errPart: "toolset mode must be"},
		{name: "File and env", env: map[string]string{"MCP_TOOLSETS_ADMIN": "on"}, file: "toolsets: {admin: off}", errPart: "cannot be set at the same time"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clearBackendEnv(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			if tc.file != "" {
				t.Setenv("MCP_CONFIG_FILE", writeBackendsFile(t, tc.file+"\nbackends:\n  - {name: a, type: logs, entrypoint: http://a}"))
			} else {
				t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://vm")
				t.Setenv("VM_INSTANCE_TYPE", "single")
			}
			_, err := InitConfig()
			if err == nil {
				t.Fatalf("Expected error containing %q, got nil", tc.errPart)
			}
			if !strings.Contains(err.Error(), tc.errPart) {
				t.Errorf("Expected error containing %q, got: %v", tc.errPart, err)
			}
		})
	}
}
