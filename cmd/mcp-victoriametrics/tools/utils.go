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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	// Add custom headers from configuration
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	// Add custom headers from configuration
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

	deploymentID, err := requireCloudDeploymentID(instance, tcr)
	if err != nil {
		return "", err
	}

	key := cloudCacheKey(instance, deploymentID)
	cloudAccessTokenCacheMutex.RLock()
	result, ok := cloudAccessTokenCache[key]
	if ok {
		cloudAccessTokenCacheMutex.RUnlock()
		return result, nil
	}
	cloudAccessTokenCacheMutex.RUnlock()

	at, err := instance.VMC().ListDeploymentAccessTokens(ctx, deploymentID)
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
		token, err := instance.VMC().RevealDeploymentAccessToken(ctx, deploymentID, t.ID)
		if err != nil {
			return "", fmt.Errorf("failed to reveal access token for deployment %s: %v", deploymentID, err)
		}
		cloudAccessTokenCacheMutex.Lock()
		cloudAccessTokenCache[key] = token.Secret
		cloudAccessTokenCacheMutex.Unlock()
		result = token.Secret
		return result, nil
	}
	return result, fmt.Errorf("no read access tokens found for deployment %s", deploymentID)
}

func getRootURL(ctx context.Context, instance *config.Instance, tcr mcp.CallToolRequest, path ...string) (string, error) {
	entrypointURL := instance.EntryPointURL()
	if instance.IsCloud() {
		deploymentID, err := requireCloudDeploymentID(instance, tcr)
		if err != nil {
			return "", err
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
	var err error
	deploymentID := ""
	entrypointURL := instance.EntryPointURL()
	isSingle := instance.IsSingle()

	// Cloud mode
	if instance.IsCloud() {
		deploymentID, err = requireCloudDeploymentID(instance, tcr)
		if err != nil {
			return "", err
		}
		info, err := getCloudDeploymentInfo(ctx, instance, deploymentID)
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
	tenant, present, err := getTrimmedStringParam(tcr, "tenant")
	if err != nil {
		return "", fmt.Errorf("failed to get tenant parameter: %v", err)
	}
	if present && tenant == "" {
		return "", fmt.Errorf("tenant parameter must not be empty")
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

func getClusterAdminToolInstance(cfg *config.Config, tcr mcp.CallToolRequest) (*config.Instance, error) {
	instance, err := getToolInstance(cfg, tcr)
	if err != nil {
		return nil, err
	}
	if !instance.IsCluster() && !instance.IsCloud() {
		return nil, fmt.Errorf("env %q does not support admin tenant operations", instance.Name())
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

func requireCloudDeploymentID(instance *config.Instance, tcr mcp.CallToolRequest) (string, error) {
	deploymentID, present, err := getTrimmedStringParam(tcr, "deployment_id")
	if err != nil {
		return "", fmt.Errorf("failed to get deployment_id parameter: %v", err)
	}
	if !present || deploymentID == "" {
		return "", fmt.Errorf("deployment_id parameter is required for cloud env %q", instance.Name())
	}
	return deploymentID, nil
}

func withTargetingOptions(options []mcp.ToolOption, c *config.Config, includeDeploymentID, includeTenant bool) []mcp.ToolOption {
	return withTargetingOptionsFor(options, c, includeDeploymentID, includeTenant, defaultCloudToolRequiresEnv(c, includeDeploymentID))
}

func withCloudToolTargetingOptions(options []mcp.ToolOption, c *config.Config, includeDeploymentID bool) []mcp.ToolOption {
	options = withTargetingOptionsFor(options, c, false, false, cloudToolRequiresEnv(c))
	if includeDeploymentID {
		options = appendDeploymentIDOption(options, true)
	}
	return options
}

func withClusterAdminTargetingOptions(options []mcp.ToolOption, c *config.Config, includeDeploymentID bool) []mcp.ToolOption {
	return withTargetingOptionsFor(options, c, includeDeploymentID, false, clusterToolRequiresEnv(c))
}

func withTargetingOptionsFor(options []mcp.ToolOption, c *config.Config, includeDeploymentID, includeTenant, requireEnv bool) []mcp.ToolOption {
	if c.HasMultipleInstances() {
		envDescription := "Optional environment to target. If omitted, the default environment is used."
		if requireEnv {
			envDescription = "Environment to target. This is required for the current server configuration."
		}
		envOptions := []mcp.PropertyOption{
			mcp.Title("Environment"),
			mcp.Description(envDescription),
			mcp.Pattern(`^[a-z0-9_]+$`),
		}
		if requireEnv {
			envOptions = append(envOptions, mcp.Required())
		}
		options = append(options, mcp.WithString("env", envOptions...))
	}
	if includeDeploymentID && c.HasCloudInstances() {
		options = appendDeploymentIDOption(options, c.HasOnlyCloudInstances())
	}
	if includeTenant && c.HasClusterInstances() {
		options = append(options, mcp.WithString("tenant",
			mcp.Title("Tenant name"),
			mcp.Description("Tenant name for cluster or cloud environments. If omitted, the selected env default is used."),
			mcp.Pattern(`^([0-9]+)(:[0-9]+)?$`),
		))
	}
	return options
}

func appendDeploymentIDOption(options []mcp.ToolOption, required bool) []mcp.ToolOption {
	propertyOptions := []mcp.PropertyOption{
		mcp.Title("Deployment ID"),
		mcp.Description("Deployment ID in VictoriaMetrics Cloud. Required when the selected env is a cloud env."),
		mcp.Pattern(`^[a-zA-Z0-9\-_]+$`),
	}
	if required {
		propertyOptions = append(propertyOptions, mcp.Required())
	}
	return append(options, mcp.WithString("deployment_id", propertyOptions...))
}

func defaultCloudToolRequiresEnv(c *config.Config, includeDeploymentID bool) bool {
	if !includeDeploymentID || !c.HasMultipleInstances() || c.HasOnlyCloudInstances() {
		return false
	}
	instance, err := c.ResolveInstance("")
	return err == nil && instance.IsCloud()
}

func cloudToolRequiresEnv(c *config.Config) bool {
	if !c.HasMultipleInstances() {
		return false
	}
	instance, err := c.ResolveInstance("")
	return err != nil || !instance.IsCloud()
}

func clusterToolRequiresEnv(c *config.Config) bool {
	if !c.HasMultipleInstances() {
		return false
	}
	instance, err := c.ResolveInstance("")
	return err != nil || !instance.IsCluster()
}

func getTrimmedStringParam(tcr mcp.CallToolRequest, param string) (string, bool, error) {
	matchArg, ok := tcr.GetArguments()[param]
	if !ok {
		return "", false, nil
	}
	value, ok := matchArg.(string)
	if !ok {
		return "", true, fmt.Errorf("%s has wrong type: %T", param, matchArg)
	}
	return strings.TrimSpace(value), true, nil
}

func ptr[T any](v T) *T {
	return &v
}
