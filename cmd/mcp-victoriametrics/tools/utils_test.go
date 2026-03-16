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
