package tools

import (
	"context"
	"fmt"
	"maps"
	"net/url"
	"slices"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

var CloudDeploymentTypes = map[config.Backend][]vmcloud.DeploymentType{
	config.BackendMetrics: {vmcloud.DeploymentTypeSingleNode, vmcloud.DeploymentTypeCluster},
	config.BackendLogs:    {vmcloud.DeploymentTypeVLogs},
	config.BackendTraces:  {vmcloud.DeploymentTypeVTraces},
}

type CloudDeploymentInfo struct {
	AccessEndpoint string
	DeploymentType vmcloud.DeploymentType
}

var (
	cloudAccessTokenCacheMutex    = &sync.RWMutex{}
	cloudAccessTokenCache         = make(map[string]string)
	cloudDeploymentInfoCacheMutex = &sync.RWMutex{}
	cloudDeploymentInfoCache      = make(map[string]CloudDeploymentInfo)
)

func GetCloudDeploymentID(tcr mcp.CallToolRequest) (string, error) {
	deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
	if err != nil {
		return "", fmt.Errorf("failed to get deployment_id parameter: %v", err)
	}
	if deploymentID == "" {
		return "", fmt.Errorf("deployment_id parameter is required for cloud mode")
	}
	return deploymentID, nil
}

func GetCloudDeploymentInfo(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (CloudDeploymentInfo, error) {
	deploymentID, err := GetCloudDeploymentID(tcr)
	if err != nil {
		return CloudDeploymentInfo{}, err
	}
	apiKey := ""
	if cfg.IsCloudSharedInstance() {
		ctx = WithCloudAccessKey(ctx, tcr)
		apiKey = tcr.Header.Get(vmcloud.AccessTokenHeader)
	}
	cacheKey := fmt.Sprintf("%s/%s", apiKey, deploymentID)

	cloudDeploymentInfoCacheMutex.RLock()
	info, ok := cloudDeploymentInfoCache[cacheKey]
	cloudDeploymentInfoCacheMutex.RUnlock()
	if ok && info.AccessEndpoint != "" && info.DeploymentType != "" {
		return info, nil
	}

	dd, err := cfg.VMC().GetDeploymentDetails(ctx, deploymentID)
	if err != nil {
		return CloudDeploymentInfo{}, fmt.Errorf("failed to get deployment details: %v", err)
	}

	if !slices.Contains(CloudDeploymentTypes[cfg.Backend()], dd.Type) {
		return CloudDeploymentInfo{}, fmt.Errorf("unsupported deployment type %s for deployment %s", dd.Type, deploymentID)
	}

	info = CloudDeploymentInfo{
		AccessEndpoint: dd.AccessEndpoint,
		DeploymentType: dd.Type,
	}

	cloudDeploymentInfoCacheMutex.Lock()
	cloudDeploymentInfoCache[cacheKey] = info
	cloudDeploymentInfoCacheMutex.Unlock()

	return info, nil
}

func GetEntryPointURL(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*url.URL, error) {
	if !cfg.IsCloud() {
		return cfg.EntryPointURL(), nil
	}
	info, err := GetCloudDeploymentInfo(ctx, cfg, tcr)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud deployment info: %v", err)
	}
	entrypointURL, err := url.Parse(info.AccessEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deployment entry point URL: %v", err)
	}
	return entrypointURL, nil
}

func GetBearerToken(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (string, error) {
	if !cfg.IsCloud() {
		return cfg.BearerToken(), nil
	}

	deploymentID, err := GetCloudDeploymentID(tcr)
	if err != nil {
		return "", err
	}
	apiKey := ""
	if cfg.IsCloudSharedInstance() {
		ctx = WithCloudAccessKey(ctx, tcr)
		apiKey = tcr.Header.Get(vmcloud.AccessTokenHeader)
	}
	cacheKey := fmt.Sprintf("%s/%s", apiKey, deploymentID)
	cloudAccessTokenCacheMutex.RLock()
	result, ok := cloudAccessTokenCache[cacheKey]
	if ok {
		cloudAccessTokenCacheMutex.RUnlock()
		return result, nil
	}
	cloudAccessTokenCacheMutex.RUnlock()

	at, err := cfg.VMC().ListDeploymentAccessTokens(ctx, deploymentID)
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
		token, err := cfg.VMC().RevealDeploymentAccessToken(ctx, deploymentID, t.ID)
		if err != nil {
			return "", fmt.Errorf("failed to reveal access token for deployment %s: %v", deploymentID, err)
		}
		cloudAccessTokenCacheMutex.Lock()
		cloudAccessTokenCache[cacheKey] = token.Secret
		cloudAccessTokenCacheMutex.Unlock()
		result = token.Secret
		return result, nil
	}
	return result, fmt.Errorf("no read access tokens found for deployment %s", deploymentID)
}

func WithCloudAccessKey(ctx context.Context, tcr mcp.CallToolRequest) context.Context {
	accessKey := tcr.Header.Get(vmcloud.AccessTokenHeader)
	if accessKey != "" {
		ctx = vmcloud.ContextWithDynamicAPIKey(ctx, accessKey)
	}
	return ctx
}

func WithCloudDeploymentID(cfg *config.Config, tool mcp.Tool) mcp.Tool {
	if !cfg.IsCloud() {
		return tool
	}
	tool.InputSchema.Properties = maps.Clone(tool.InputSchema.Properties)
	tool.InputSchema.Required = slices.Clone(tool.InputSchema.Required)
	mcp.WithString("deployment_id",
		mcp.Required(),
		mcp.Title("Deployment ID"),
		mcp.Description("Unique identifier of the deployment in VictoriaMetrics Cloud"),
		mcp.Pattern(`^[a-zA-Z0-9\-_]+$`),
	)(&tool)
	return tool
}
