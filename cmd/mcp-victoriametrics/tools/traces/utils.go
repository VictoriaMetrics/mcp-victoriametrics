package traces

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/VictoriaMetrics/VictoriaLogs/lib/logstorage"
	"github.com/mark3labs/mcp-go/mcp"

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
	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}

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

	if err := SetTenantHeaders(req, cfg, tcr); err != nil {
		return nil, err
	}

	return req, nil
}

// SetTenantHeaders sets AccountID and ProjectID headers from the tenant parameter of the request,
// falling back to the default tenant from the config.
func SetTenantHeaders(req *http.Request, cfg *config.Config, tcr mcp.CallToolRequest) error {
	accountID, projectID, err := GetToolReqTenant(tcr)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %v", err)
	}
	defaultTenantID := cfg.DefaultLogsTenantID()
	if accountID == "" {
		accountID = strconv.FormatUint(uint64(defaultTenantID.AccountID), 10)
	}
	if projectID == "" {
		projectID = strconv.FormatUint(uint64(defaultTenantID.ProjectID), 10)
	}
	req.Header.Set("AccountID", accountID)
	req.Header.Set("ProjectID", projectID)
	return nil
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
	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}

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
	return entrypointURL.JoinPath("select", "jaeger", "api").JoinPath(path...).String(), nil
}

func GetToolReqTenant(tcr mcp.CallToolRequest) (string, string, error) {
	tenant, err := tools.GetToolReqTenantString(tcr)
	if err != nil {
		return "", "", fmt.Errorf("failed to get tenant: %v", err)
	}
	if tenant == "" {
		return "", "", nil
	}
	tenantID, err := logstorage.ParseTenantID(tenant)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse tenant %q: %v", tenant, err)
	}
	accountID := strconv.FormatUint(uint64(tenantID.AccountID), 10)
	projectID := strconv.FormatUint(uint64(tenantID.ProjectID), 10)
	return accountID, projectID, nil
}

func ptr[T any](v T) *T {
	return &v
}
