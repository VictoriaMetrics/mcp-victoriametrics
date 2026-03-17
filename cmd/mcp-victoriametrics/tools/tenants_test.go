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

func TestMixedEnvTenantsToolRequiresEnvWhenDefaultIsSingle(t *testing.T) {
	t.Setenv("VM_ENVIRONMENTS", "default,cluster")
	t.Setenv("VM_DEFAULT_ENVIRONMENT", "default")
	t.Setenv("VM_INSTANCE_DEFAULT_ENTRYPOINT", "http://default.example.com")
	t.Setenv("VM_INSTANCE_DEFAULT_TYPE", "single")
	t.Setenv("VM_INSTANCE_CLUSTER_ENTRYPOINT", "http://cluster.example.com")
	t.Setenv("VM_INSTANCE_CLUSTER_TYPE", "cluster")

	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}

	tool := toolTenants(cfg)
	if !toolHasRequiredProperty(tool, "env") {
		t.Fatal("tenants tool should require env when default env cannot serve /admin/tenants")
	}
}

func TestTenantsHandlerRejectsDefaultSingleWithoutEnvBeforeRequest(t *testing.T) {
	t.Setenv("VM_ENVIRONMENTS", "default,cluster")
	t.Setenv("VM_DEFAULT_ENVIRONMENT", "default")
	t.Setenv("VM_INSTANCE_DEFAULT_ENTRYPOINT", "http://default.example.com")
	t.Setenv("VM_INSTANCE_DEFAULT_TYPE", "single")
	t.Setenv("VM_INSTANCE_CLUSTER_ENTRYPOINT", "http://cluster.example.com")
	t.Setenv("VM_INSTANCE_CLUSTER_TYPE", "cluster")

	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}

	transport := &trackingTransport{
		response: &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewBufferString("ok"))},
	}
	originalClient := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: transport}
	defer func() { http.DefaultClient = originalClient }()

	result, err := toolTenantsHandler(context.Background(), cfg, mcp.CallToolRequest{})
	if err != nil {
		t.Fatalf("toolTenantsHandler() error = %v", err)
	}
	if !result.IsError {
		t.Fatal("expected tenants handler to fail before issuing request")
	}
	if transport.called {
		t.Fatal("expected tenants handler to fail locally without making an HTTP request")
	}
}

type trackingTransport struct {
	called   bool
	response *http.Response
	err      error
}

func (t *trackingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	t.called = true
	return t.response, t.err
}
