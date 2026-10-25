package metrics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

func CreateSelectRequest(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (*http.Request, error) {
	selectURL, err := getSelectURL(ctx, cfg, tcr, path...)
	if err != nil {
		return nil, fmt.Errorf("failed to get select URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, selectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	bearerToken, err := tools.GetBearerToken(ctx, cfg, tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to get bearer token: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	// Add custom headers from configuration
	for key, value := range cfg.CustomHeaders() {
		req.Header.Set(key, value)
	}

	// Apply passthrough headers from the incoming MCP request
	for _, name := range cfg.PassthroughHeaders() {
		if value := tcr.Header.Get(name); value != "" {
			req.Header.Set(name, value)
		}
	}

	return req, nil
}

func CreateAdminRequest(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (*http.Request, error) {
	selectURL, err := getRootURL(ctx, cfg, tcr, path...)
	if err != nil {
		return nil, fmt.Errorf("failed to get select URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, selectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	bearerToken, err := tools.GetBearerToken(ctx, cfg, tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to get bearer token: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	// Add custom headers from configuration
	for key, value := range cfg.CustomHeaders() {
		req.Header.Set(key, value)
	}

	// Apply passthrough headers from the incoming MCP request
	for _, name := range cfg.PassthroughHeaders() {
		if value := tcr.Header.Get(name); value != "" {
			req.Header.Set(name, value)
		}
	}

	return req, nil
}

func getRootURL(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (string, error) {
	entrypointURL, err := tools.GetEntryPointURL(ctx, cfg, tcr)
	if err != nil {
		return "", err
	}
	return entrypointURL.JoinPath(path...).String(), nil
}

func getSelectURL(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (string, error) {
	entrypointURL, err := tools.GetEntryPointURL(ctx, cfg, tcr)
	if err != nil {
		return "", err
	}

	isSingle := cfg.IsSingle()
	if cfg.IsCloud() {
		info, err := tools.GetCloudDeploymentInfo(ctx, cfg, tcr)
		if err != nil {
			return "", fmt.Errorf("failed to get cloud deployment info: %v", err)
		}
		isSingle = info.DeploymentType == vmcloud.DeploymentTypeSingleNode
	}

	// Single node
	if isSingle {
		return entrypointURL.JoinPath(path...).String(), nil
	}

	// Cluster mode
	tenant, err := tools.GetToolReqTenantParam(tcr)
	if err != nil {
		return "", fmt.Errorf("failed to get tenant parameter: %v", err)
	}
	if tenant == "" {
		tenant = cfg.DefaultTenantID()
	}
	args := []string{"select", tenant, "prometheus"}
	return entrypointURL.JoinPath(append(args, path...)...).String(), nil
}

func ptr[T any](v T) *T {
	return &v
}
