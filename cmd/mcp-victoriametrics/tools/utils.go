package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func CreateSelectRequest(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (*http.Request, error) {
	environment, err := getToolEnvironment(cfg, tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve env: %v", err)
	}

	selectURL, err := getSelectURL(ctx, environment, tcr, path...)
	if err != nil {
		return nil, fmt.Errorf("failed to get select URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, selectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	bearerToken, err := getBearerToken(ctx, environment, tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to get bearer token: %v", err)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}

	// Add custom headers from configuration
	for key, value := range environment.CustomHeaders() {
		req.Header.Set(key, value)
	}

	return req, nil
}

func CreateAdminRequest(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (*http.Request, error) {
	environment, err := getToolEnvironment(cfg, tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve env: %v", err)
	}

	selectURL, err := getRootURL(ctx, environment, tcr, path...)
	if err != nil {
		return nil, fmt.Errorf("failed to get select URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, selectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	bearerToken, err := getBearerToken(ctx, environment, tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to get bearer token: %v", err)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}

	// Add custom headers from configuration
	for key, value := range environment.CustomHeaders() {
		req.Header.Set(key, value)
	}

	return req, nil
}

type cloudDeploymentInfo struct {
	accessEndpoint string
	deploymentType vmcloud.DeploymentType
}

var (
	cloudAccessTokenCacheMutex    = &sync.RWMutex{}
	cloudAccessTokenCache         = make(map[string]string)
	cloudDeploymentInfoCacheMutex = &sync.RWMutex{}
	cloudDeploymentInfoCache      = make(map[string]cloudDeploymentInfo)
)

func getCloudDeploymentInfo(ctx context.Context, environment *config.InstanceConfig, deploymentID string) (cloudDeploymentInfo, error) {
	cloudDeploymentInfoCacheMutex.RLock()
	info, ok := cloudDeploymentInfoCache[deploymentID]
	cloudDeploymentInfoCacheMutex.RUnlock()
	if ok && info.accessEndpoint != "" && info.deploymentType != "" {
		return info, nil
	}

	dd, err := environment.VMC().GetDeploymentDetails(ctx, deploymentID)
	if err != nil {
		return cloudDeploymentInfo{}, fmt.Errorf("failed to get deployment details: %v", err)
	}

	if dd.Type != vmcloud.DeploymentTypeSingleNode && dd.Type != vmcloud.DeploymentTypeCluster {
		return cloudDeploymentInfo{}, fmt.Errorf("unsupported deployment type %s for deployment %s", dd.Type, deploymentID)
	}

	info = cloudDeploymentInfo{
		accessEndpoint: dd.AccessEndpoint,
		deploymentType: dd.Type,
	}

	cloudDeploymentInfoCacheMutex.Lock()
	cloudDeploymentInfoCache[deploymentID] = info
	cloudDeploymentInfoCacheMutex.Unlock()

	return info, nil
}

func getBearerToken(ctx context.Context, environment *config.InstanceConfig, tcr mcp.CallToolRequest) (string, error) {
	if !environment.IsCloud() {
		return environment.BearerToken(), nil
	}

	deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
	if err != nil {
		return "", fmt.Errorf("failed to get deployment_id parameter: %v", err)
	}
	if deploymentID == "" {
		return "", fmt.Errorf("deployment_id parameter is required for cloud mode")
	}
	cloudAccessTokenCacheMutex.RLock()
	result, ok := cloudAccessTokenCache[deploymentID]
	if ok {
		cloudAccessTokenCacheMutex.RUnlock()
		return result, nil
	}
	cloudAccessTokenCacheMutex.RUnlock()

	at, err := environment.VMC().ListDeploymentAccessTokens(ctx, deploymentID)
	if err != nil {
		return "", fmt.Errorf("failed to list deployment access tokens: %v", err)
	}
	if len(at) == 0 {
		return "", fmt.Errorf("no access tokens found for deployment %s", deploymentID)
	}
	for _, t := range at {
		if t.Type == vmcloud.AccessModeWrite {
			continue // Skip write only tokens
		}
		if t.TenantID != "" {
			continue // Skip tokens with specific tenant ID
		}
		token, err := environment.VMC().RevealDeploymentAccessToken(ctx, deploymentID, t.ID)
		if err != nil {
			return "", fmt.Errorf("failed to reveal access token for deployment %s: %v", deploymentID, err)
		}
		cloudAccessTokenCacheMutex.Lock()
		cloudAccessTokenCache[deploymentID] = token.Secret
		cloudAccessTokenCacheMutex.Unlock()
		result = token.Secret
		return result, nil
	}
	return result, fmt.Errorf("no read access tokens found for deployment %s", deploymentID)
}

func getRootURL(ctx context.Context, environment *config.InstanceConfig, tcr mcp.CallToolRequest, path ...string) (string, error) {
	entrypointURL := environment.EntryPointURL()
	if environment.IsCloud() {
		deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
		if err != nil {
			return "", fmt.Errorf("failed to get deployment_id parameter: %v", err)
		}
		if deploymentID == "" {
			return "", fmt.Errorf("deployment_id parameter is required for cloud mode")
		}
		info, err := getCloudDeploymentInfo(ctx, environment, deploymentID)
		if err != nil {
			return "", fmt.Errorf("failed to get cloud deployment info: %v", err)
		}
		entrypointURL, err = url.Parse(info.accessEndpoint)
		if err != nil {
			return "", fmt.Errorf("failed to parse deployment entry point URL: %v", err)
		}
	}
	return entrypointURL.JoinPath(path...).String(), nil
}

func getSelectURL(ctx context.Context, environment *config.InstanceConfig, tcr mcp.CallToolRequest, path ...string) (string, error) {
	var err error
	deploymentID := ""
	entrypointURL := environment.EntryPointURL()
	isSingle := environment.IsSingle()

	// Cloud mode
	if environment.IsCloud() {
		deploymentID, err = GetToolReqParam[string](tcr, "deployment_id", true)
		if err != nil {
			return "", fmt.Errorf("failed to get deployment_id parameter: %v", err)
		}
		if deploymentID == "" {
			return "", fmt.Errorf("deployment_id parameter is required for cloud mode")
		}
		info, err := getCloudDeploymentInfo(ctx, environment, deploymentID)
		if err != nil {
			return "", fmt.Errorf("failed to get cloud deployment info: %v", err)
		}
		entrypointURL, err = url.Parse(info.accessEndpoint)
		if err != nil {
			return "", fmt.Errorf("failed to parse deployment entry point URL: %v", err)
		}
		isSingle = info.deploymentType == vmcloud.DeploymentTypeSingleNode
	}

	// Single node
	if isSingle {
		return entrypointURL.JoinPath(path...).String(), nil
	}

	// Cluster mode
	tenant, err := GetToolReqParam[string](tcr, "tenant", false)
	if err != nil {
		return "", fmt.Errorf("failed to get tenant parameter: %v", err)
	}
	if tenant == "" {
		tenant = environment.DefaultTenantID()
	}
	args := []string{"select", tenant, "prometheus"}
	return entrypointURL.JoinPath(append(args, path...)...).String(), nil
}

func GetTextBodyForRequest(req *http.Request, _ *config.Config, f ...func(s string) (string, error)) *mcp.CallToolResult {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to do request: %v", err))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read response body: %v", err))
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("unexpected response status code %v: %s", resp.StatusCode, string(body)))
	}
	result := string(body)
	for _, fn := range f {
		if result, err = fn(result); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to process response body: %v", err))
		}
	}
	return mcp.NewToolResultText(result)
}

type ToolReqParamType interface {
	string | float64 | bool | []string | []any
}

func GetToolReqParam[T ToolReqParamType](tcr mcp.CallToolRequest, param string, required bool) (T, error) {
	var value T
	matchArg, ok := tcr.GetArguments()[param]
	if ok {
		value, ok = matchArg.(T)
		if !ok {
			return value, fmt.Errorf("%s has wrong type: %T", param, matchArg)
		}
	} else if required {
		return value, fmt.Errorf("%s param is required", param)
	}
	return value, nil
}

func GetToolReqEnv(tcr mcp.CallToolRequest) (string, error) {
	env, err := GetToolReqParam[string](tcr, "env", false)
	if err != nil {
		return "", fmt.Errorf("failed to get env: %v", err)
	}
	if env != "" {
		return strings.ToLower(strings.TrimSpace(env)), nil
	}

	environment, err := GetToolReqParam[string](tcr, "environment", false)
	if err != nil {
		return "", fmt.Errorf("failed to get environment: %v", err)
	}
	return strings.ToLower(strings.TrimSpace(environment)), nil
}

func getToolEnvironment(cfg *config.Config, tcr mcp.CallToolRequest) (*config.InstanceConfig, error) {
	envName, err := GetToolReqEnv(tcr)
	if err != nil {
		return nil, err
	}
	return cfg.Environment(envName)
}

func withEnvironmentParam() mcp.ToolOption {
	return mcp.WithString("env",
		mcp.Title("Environment"),
		mcp.Description("Optional VictoriaMetrics environment to target. If omitted, the server default environment is used."),
		mcp.Pattern(`^[A-Za-z0-9_-]+$`),
	)
}

func ptr[T any](v T) *T {
	return &v
}
