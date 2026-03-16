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
	instance, err := getToolInstance(cfg, tcr)
	if err != nil {
		return nil, err
	}

	selectURL, err := getSelectURL(ctx, instance, tcr, path...)
	if err != nil {
		return nil, fmt.Errorf("failed to get select URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, selectURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	bearerToken, err := getBearerToken(ctx, instance, tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to get bearer token: %v", err)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}
	for key, value := range instance.CustomHeaders() {
		req.Header.Set(key, value)
	}

	return req, nil
}

func CreateAdminRequest(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (*http.Request, error) {
	instance, err := getToolInstance(cfg, tcr)
	if err != nil {
		return nil, err
	}

	rootURL, err := getRootURL(ctx, instance, tcr, path...)
	if err != nil {
		return nil, fmt.Errorf("failed to get select URL: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rootURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	bearerToken, err := getBearerToken(ctx, instance, tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to get bearer token: %v", err)
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}
	for key, value := range instance.CustomHeaders() {
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

func getCloudDeploymentInfo(ctx context.Context, instance *config.Instance, deploymentID string) (cloudDeploymentInfo, error) {
	key := cloudCacheKey(instance, deploymentID)
	cloudDeploymentInfoCacheMutex.RLock()
	info, ok := cloudDeploymentInfoCache[key]
	cloudDeploymentInfoCacheMutex.RUnlock()
	if ok && info.accessEndpoint != "" && info.deploymentType != "" {
		return info, nil
	}

	dd, err := instance.VMC().GetDeploymentDetails(ctx, deploymentID)
	if err != nil {
		return cloudDeploymentInfo{}, fmt.Errorf("failed to get deployment details: %v", err)
	}
	if dd.Type != vmcloud.DeploymentTypeSingleNode && dd.Type != vmcloud.DeploymentTypeCluster {
		return cloudDeploymentInfo{}, fmt.Errorf("unsupported deployment type %s for deployment %s", dd.Type, deploymentID)
	}

	info = cloudDeploymentInfo{accessEndpoint: dd.AccessEndpoint, deploymentType: dd.Type}
	cloudDeploymentInfoCacheMutex.Lock()
	cloudDeploymentInfoCache[key] = info
	cloudDeploymentInfoCacheMutex.Unlock()
	return info, nil
}

func getBearerToken(ctx context.Context, instance *config.Instance, tcr mcp.CallToolRequest) (string, error) {
	if !instance.IsCloud() {
		return instance.BearerToken(), nil
	}

	deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
	if err != nil {
		return "", fmt.Errorf("failed to get deployment_id parameter: %v", err)
	}
	if deploymentID == "" {
		return "", fmt.Errorf("deployment_id parameter is required for cloud env %q", instance.Name())
	}

	key := cloudCacheKey(instance, deploymentID)
	cloudAccessTokenCacheMutex.RLock()
	result, ok := cloudAccessTokenCache[key]
	cloudAccessTokenCacheMutex.RUnlock()
	if ok {
		return result, nil
	}

	at, err := instance.VMC().ListDeploymentAccessTokens(ctx, deploymentID)
	if err != nil {
		return "", fmt.Errorf("failed to list deployment access tokens: %v", err)
	}
	if len(at) == 0 {
		return "", fmt.Errorf("no access tokens found for deployment %s", deploymentID)
	}
	for _, t := range at {
		if t.Type == vmcloud.AccessModeWrite || t.TenantID != "" {
			continue
		}
		token, err := instance.VMC().RevealDeploymentAccessToken(ctx, deploymentID, t.ID)
		if err != nil {
			return "", fmt.Errorf("failed to reveal access token for deployment %s: %v", deploymentID, err)
		}
		cloudAccessTokenCacheMutex.Lock()
		cloudAccessTokenCache[key] = token.Secret
		cloudAccessTokenCacheMutex.Unlock()
		return token.Secret, nil
	}
	return "", fmt.Errorf("no read access tokens found for deployment %s", deploymentID)
}

func getRootURL(ctx context.Context, instance *config.Instance, tcr mcp.CallToolRequest, path ...string) (string, error) {
	entrypointURL := instance.EntryPointURL()
	if instance.IsCloud() {
		deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
		if err != nil {
			return "", fmt.Errorf("failed to get deployment_id parameter: %v", err)
		}
		if deploymentID == "" {
			return "", fmt.Errorf("deployment_id parameter is required for cloud env %q", instance.Name())
		}
		info, err := getCloudDeploymentInfo(ctx, instance, deploymentID)
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

func getSelectURL(ctx context.Context, instance *config.Instance, tcr mcp.CallToolRequest, path ...string) (string, error) {
	entrypointURL := instance.EntryPointURL()
	isSingle := instance.IsSingle()
	if instance.IsCloud() {
		deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
		if err != nil {
			return "", fmt.Errorf("failed to get deployment_id parameter: %v", err)
		}
		if deploymentID == "" {
			return "", fmt.Errorf("deployment_id parameter is required for cloud env %q", instance.Name())
		}
		info, err := getCloudDeploymentInfo(ctx, instance, deploymentID)
		if err != nil {
			return "", fmt.Errorf("failed to get cloud deployment info: %v", err)
		}
		var errParse error
		entrypointURL, errParse = url.Parse(info.accessEndpoint)
		if errParse != nil {
			return "", fmt.Errorf("failed to parse deployment entry point URL: %v", errParse)
		}
		isSingle = info.deploymentType == vmcloud.DeploymentTypeSingleNode
	}
	if isSingle {
		return entrypointURL.JoinPath(path...).String(), nil
	}

	tenant, err := GetToolReqParam[string](tcr, "tenant", false)
	if err != nil {
		return "", fmt.Errorf("failed to get tenant parameter: %v", err)
	}
	if tenant == "" {
		tenant = instance.DefaultTenantID()
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

func getToolInstance(cfg *config.Config, tcr mcp.CallToolRequest) (*config.Instance, error) {
	env, err := GetToolReqParam[string](tcr, "env", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get env parameter: %v", err)
	}
	instance, err := cfg.ResolveInstance(strings.TrimSpace(env))
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func getCloudToolInstance(cfg *config.Config, tcr mcp.CallToolRequest) (*config.Instance, error) {
	instance, err := getToolInstance(cfg, tcr)
	if err != nil {
		return nil, err
	}
	if !instance.IsCloud() {
		return nil, fmt.Errorf("env %q is not a VictoriaMetrics Cloud env", instance.Name())
	}
	return instance, nil
}

func cloudCacheKey(instance *config.Instance, deploymentID string) string {
	return instance.Name() + ":" + deploymentID
}

func maybeWithEnvironmentParam(c *config.Config) []mcp.ToolOption {
	if !c.HasMultipleInstances() {
		return nil
	}
	return []mcp.ToolOption{mcp.WithString("env",
		mcp.Title("Environment"),
		mcp.Description("Optional environment to target. If omitted, the default environment is used."),
		mcp.Pattern(`^[a-z0-9_]+$`),
	)}
}

func maybeWithDeploymentIDParam(c *config.Config) []mcp.ToolOption {
	if !c.HasCloudInstances() {
		return nil
	}
	options := []mcp.PropertyOption{
		mcp.Title("Deployment ID"),
		mcp.Description("Deployment ID in VictoriaMetrics Cloud. Required when the selected env is a cloud env."),
		mcp.Pattern(`^[a-zA-Z0-9\-_]+$`),
	}
	if c.HasOnlyCloudInstances() {
		options = append(options, mcp.Required())
	}
	return []mcp.ToolOption{mcp.WithString("deployment_id", options...)}
}

func maybeWithTenantParam(c *config.Config) []mcp.ToolOption {
	if !c.HasClusterInstances() {
		return nil
	}
	return []mcp.ToolOption{mcp.WithString("tenant",
		mcp.Title("Tenant name"),
		mcp.Description("Tenant name for cluster or cloud environments. If omitted, the selected env default is used."),
		mcp.DefaultString("0"),
		mcp.Pattern(`^([0-9]+)(:[0-9]+)?$`),
	)}
}

func ptr[T any](v T) *T {
	return &v
}
