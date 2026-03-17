package tools

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestGetTextBodyForRequest(t *testing.T) {
	originalClient := http.DefaultClient
	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("test response")),
			},
		},
	}
	defer func() { http.DefaultClient = originalClient }()

	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() error = %v", err)
	}

	result := GetTextBodyForRequest(req, &config.Config{})
	if result.IsError {
		t.Fatal("expected success result")
	}
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("unexpected content type %T", result.Content[0])
	}
	if textContent.Text != "test response" {
		t.Fatalf("Text = %q", textContent.Text)
	}
}

func TestGetToolReqParam(t *testing.T) {
	tcr := mcp.CallToolRequest{}
	tcr.Params.Arguments = map[string]any{
		"string": "value",
		"number": 123.0,
		"bool":   true,
	}

	if value, err := GetToolReqParam[string](tcr, "string", true); err != nil || value != "value" {
		t.Fatalf("GetToolReqParam[string]() = %q, %v", value, err)
	}
	if value, err := GetToolReqParam[float64](tcr, "number", true); err != nil || value != 123.0 {
		t.Fatalf("GetToolReqParam[float64]() = %v, %v", value, err)
	}
	if value, err := GetToolReqParam[bool](tcr, "bool", true); err != nil || !value {
		t.Fatalf("GetToolReqParam[bool]() = %v, %v", value, err)
	}
	if _, err := GetToolReqParam[string](tcr, "missing", true); err == nil {
		t.Fatal("expected error for missing required param")
	}
}

func TestGetToolInstanceAndRouting(t *testing.T) {
	t.Setenv("VM_ENVIRONMENTS", "default,demo")
	t.Setenv("VM_DEFAULT_ENVIRONMENT", "demo")
	t.Setenv("VM_INSTANCE_DEFAULT_ENTRYPOINT", "http://default.example.com")
	t.Setenv("VM_INSTANCE_DEFAULT_TYPE", "single")
	t.Setenv("VM_INSTANCE_DEMO_ENTRYPOINT", "http://demo.example.com")
	t.Setenv("VM_INSTANCE_DEMO_TYPE", "cluster")
	t.Setenv("VM_INSTANCE_DEMO_DEFAULT_TENANT_ID", "42")

	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}

	t.Run("default env", func(t *testing.T) {
		tcr := mcp.CallToolRequest{}
		instance, err := getToolInstance(cfg, tcr)
		if err != nil {
			t.Fatalf("getToolInstance() error = %v", err)
		}
		if instance.Name() != "demo" {
			t.Fatalf("instance.Name() = %q", instance.Name())
		}
		url, err := getSelectURL(context.Background(), instance, tcr, "api", "v1", "query")
		if err != nil {
			t.Fatalf("getSelectURL() error = %v", err)
		}
		if url != "http://demo.example.com/select/42/prometheus/api/v1/query" {
			t.Fatalf("getSelectURL() = %q", url)
		}
	})

	t.Run("explicit env", func(t *testing.T) {
		tcr := mcp.CallToolRequest{}
		tcr.Params.Arguments = map[string]any{"env": "default"}
		instance, err := getToolInstance(cfg, tcr)
		if err != nil {
			t.Fatalf("getToolInstance() error = %v", err)
		}
		url, err := getSelectURL(context.Background(), instance, tcr, "api", "v1", "query")
		if err != nil {
			t.Fatalf("getSelectURL() error = %v", err)
		}
		if url != "http://default.example.com/api/v1/query" {
			t.Fatalf("getSelectURL() = %q", url)
		}
	})

	t.Run("unknown env", func(t *testing.T) {
		tcr := mcp.CallToolRequest{}
		tcr.Params.Arguments = map[string]any{"env": "missing"}
		if _, err := getToolInstance(cfg, tcr); err == nil {
			t.Fatal("expected unknown env error")
		}
	})
}

func TestCloudCacheKeyIncludesEnv(t *testing.T) {
	t.Setenv("VM_ENVIRONMENTS", "cloud_a,cloud_b")
	t.Setenv("VMC_CLOUD_A_API_KEY", "test-api-key")
	t.Setenv("VMC_CLOUD_B_API_KEY", "test-api-key")

	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}
	a, err := cfg.ResolveInstance("cloud_a")
	if err != nil {
		t.Fatalf("ResolveInstance(cloud_a) error = %v", err)
	}
	b, err := cfg.ResolveInstance("cloud_b")
	if err != nil {
		t.Fatalf("ResolveInstance(cloud_b) error = %v", err)
	}
	if cloudCacheKey(a, "dep") == cloudCacheKey(b, "dep") {
		t.Fatal("expected cloud cache keys to differ by env")
	}
}

func TestMixedEnvCloudOnlyToolsRequireExplicitTargeting(t *testing.T) {
	t.Setenv("VM_ENVIRONMENTS", "default,cloud")
	t.Setenv("VM_DEFAULT_ENVIRONMENT", "default")
	t.Setenv("VM_INSTANCE_DEFAULT_ENTRYPOINT", "http://default.example.com")
	t.Setenv("VM_INSTANCE_DEFAULT_TYPE", "single")
	t.Setenv("VMC_CLOUD_API_KEY", "test-api-key")

	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}

	t.Run("cloud providers requires env", func(t *testing.T) {
		tool := newCloudListTool(toolNameCloudProviders, "List of cloud providers in VictoriaMetrics Cloud", "List of cloud providers", cfg)
		if !toolHasRequiredProperty(tool, "env") {
			t.Fatal("cloud-only tool should require env when default env is not cloud")
		}
	})

	t.Run("access tokens requires deployment id", func(t *testing.T) {
		tool := toolAccessTokens(cfg)
		if !toolHasRequiredProperty(tool, "deployment_id") {
			t.Fatal("cloud-only tool should require deployment_id even on mixed servers")
		}
		if !toolHasRequiredProperty(tool, "env") {
			t.Fatal("cloud-only tool should require env when default env is not cloud")
		}
	})

	t.Run("rule filenames requires deployment id", func(t *testing.T) {
		tool := toolRuleFilenames(cfg)
		if !toolHasRequiredProperty(tool, "deployment_id") {
			t.Fatal("cloud-only tool should require deployment_id even on mixed servers")
		}
		if !toolHasRequiredProperty(tool, "env") {
			t.Fatal("cloud-only tool should require env when default env is not cloud")
		}
	})
}

func TestMixedDefaultCloudGenericToolsRequireExplicitEnv(t *testing.T) {
	t.Setenv("VM_ENVIRONMENTS", "cloud,local")
	t.Setenv("VM_DEFAULT_ENVIRONMENT", "cloud")
	t.Setenv("VMC_CLOUD_API_KEY", "test-api-key")
	t.Setenv("VM_INSTANCE_LOCAL_ENTRYPOINT", "http://local.example.com")
	t.Setenv("VM_INSTANCE_LOCAL_TYPE", "cluster")

	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}

	for name, tool := range map[string]mcp.Tool{
		"query":   toolQuery(cfg),
		"flags":   toolFlags(cfg),
		"tenants": toolTenants(cfg),
	} {
		t.Run(name, func(t *testing.T) {
			if !toolHasRequiredProperty(tool, "env") {
				t.Fatal("env should be required when default env is cloud on a mixed server")
			}
			envSchema := toolPropertySchema(tool, "env")
			if envSchema["description"] == "Optional environment to target. If omitted, the default environment is used." {
				t.Fatal("required env should not be described as optional")
			}
		})
	}
}

func TestMixedEnvTenantSchemaDoesNotAdvertiseGlobalDefault(t *testing.T) {
	t.Setenv("VM_ENVIRONMENTS", "demo,prod")
	t.Setenv("VM_DEFAULT_ENVIRONMENT", "demo")
	t.Setenv("VM_INSTANCE_DEMO_ENTRYPOINT", "http://demo.example.com")
	t.Setenv("VM_INSTANCE_DEMO_TYPE", "cluster")
	t.Setenv("VM_INSTANCE_DEMO_DEFAULT_TENANT_ID", "42")
	t.Setenv("VM_INSTANCE_PROD_ENTRYPOINT", "http://prod.example.com")
	t.Setenv("VM_INSTANCE_PROD_TYPE", "cluster")
	t.Setenv("VM_INSTANCE_PROD_DEFAULT_TENANT_ID", "7")

	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}

	tool := toolQuery(cfg)
	tenantSchema := toolPropertySchema(tool, "tenant")
	if _, ok := tenantSchema["default"]; ok {
		t.Fatal("tenant schema should not advertise a single default when tenant fallback depends on selected env")
	}
}

func TestTargetingParamsTrimAndRejectWhitespace(t *testing.T) {
	t.Run("deployment id rejects whitespace and trims valid input", func(t *testing.T) {
		t.Setenv("VMC_API_KEY", "test-api-key")
		cfg, err := config.InitConfig()
		if err != nil {
			t.Fatalf("InitConfig() error = %v", err)
		}
		instance, err := cfg.ResolveInstance("")
		if err != nil {
			t.Fatalf("ResolveInstance(default) error = %v", err)
		}

		blank := mcp.CallToolRequest{}
		blank.Params.Arguments = map[string]any{"deployment_id": "   "}
		if _, err := requireCloudDeploymentID(instance, blank); err == nil {
			t.Fatal("expected whitespace-only deployment_id to be rejected")
		}

		trimmed := mcp.CallToolRequest{}
		trimmed.Params.Arguments = map[string]any{"deployment_id": "  dep-1  "}
		value, err := requireCloudDeploymentID(instance, trimmed)
		if err != nil {
			t.Fatalf("requireCloudDeploymentID() error = %v", err)
		}
		if value != "dep-1" {
			t.Fatalf("requireCloudDeploymentID() = %q, want %q", value, "dep-1")
		}
	})

	t.Run("tenant rejects whitespace and trims valid input", func(t *testing.T) {
		t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://cluster.example.com")
		t.Setenv("VM_INSTANCE_TYPE", "cluster")
		t.Setenv("VM_DEFAULT_TENANT_ID", "7")
		cfg, err := config.InitConfig()
		if err != nil {
			t.Fatalf("InitConfig() error = %v", err)
		}
		instance, err := cfg.ResolveInstance("")
		if err != nil {
			t.Fatalf("ResolveInstance(default) error = %v", err)
		}

		blank := mcp.CallToolRequest{}
		blank.Params.Arguments = map[string]any{"tenant": "   "}
		if _, err := getSelectURL(context.Background(), instance, blank, "api", "v1", "query"); err == nil {
			t.Fatal("expected whitespace-only tenant to be rejected")
		}

		trimmed := mcp.CallToolRequest{}
		trimmed.Params.Arguments = map[string]any{"tenant": " 42 "}
		url, err := getSelectURL(context.Background(), instance, trimmed, "api", "v1", "query")
		if err != nil {
			t.Fatalf("getSelectURL() error = %v", err)
		}
		if url != "http://cluster.example.com/select/42/prometheus/api/v1/query" {
			t.Fatalf("getSelectURL() = %q", url)
		}
	})
}

func toolHasRequiredProperty(tool mcp.Tool, property string) bool {
	for _, required := range tool.InputSchema.Required {
		if required == property {
			return true
		}
	}
	return false
}

func toolPropertySchema(tool mcp.Tool, property string) map[string]any {
	schema, ok := tool.InputSchema.Properties[property].(map[string]any)
	if !ok {
		return nil
	}
	return schema
}
