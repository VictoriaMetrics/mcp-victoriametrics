package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	// DefaultBaseURL is the default base URL for the VMCloud API
	DefaultBaseURL = "https://api.victoriametrics.cloud"
	// AccessTokenHeader is the header name for the access token
	AccessTokenHeader = "X-VM-Cloud-Access"
	// DynamicAPIKey - use this constant as a value for API key in the context to indicate that the API key
	// should be taken from the context instead of the client configuration.
	// This allows using different API keys for different requests with the same client instance.
	DynamicAPIKey = "dynamic"
)

// VMCloudAPIClient represents a API client for VictoriaMetrics Cloud API
type VMCloudAPIClient struct {
	c         *http.Client
	apiKey    string
	baseURL   string
	parsedURL *url.URL
}

// VMCloudAPIClientOption defines a functional option to configure a VMCloudAPIClient instance.
type VMCloudAPIClientOption func(*VMCloudAPIClient)

// WithHTTPClient sets a custom HTTP client for the VMCloudAPIClient instance.
func WithHTTPClient(c *http.Client) VMCloudAPIClientOption {
	return func(client *VMCloudAPIClient) {
		client.c = c
	}
}

// WithBaseURL sets a custom base URL for the VMCloudAPIClient instance.
func WithBaseURL(baseURL string) VMCloudAPIClientOption {
	return func(client *VMCloudAPIClient) {
		client.baseURL = baseURL
	}
}

// New creates a new VMCloudAPIClient instance with the provided API key and options.
func New(apiKey string, options ...VMCloudAPIClientOption) (*VMCloudAPIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}
	if apiKey == DynamicAPIKey {
		apiKey = ""
	}
	result := &VMCloudAPIClient{
		c:       http.DefaultClient,
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
	}
	for _, option := range options {
		option(result)
	}
	var err error
	result.parsedURL, err = url.Parse(result.baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL %q: %w", result.baseURL, err)
	}
	return result, nil
}

// BaseURL returns the base URL of the VMCloudAPIClient instance.
func (a *VMCloudAPIClient) BaseURL() string {
	return a.baseURL
}

// ListCloudProviders retrieves the list of available cloud providers for deployments in VictoriaMetrics Cloud.
func (a *VMCloudAPIClient) ListCloudProviders(ctx context.Context) (CloudProviderInfoList, error) {
	return requestAPI[CloudProviderInfoList](ctx, a, http.MethodGet, nil, "/api/v1/cloud_providers")
}

// ListRegions retrieves the list of available regions for deployments in VictoriaMetrics Cloud.
func (a *VMCloudAPIClient) ListRegions(ctx context.Context) (RegionInfoList, error) {
	return requestAPI[RegionInfoList](ctx, a, http.MethodGet, nil, "/api/v1/regions")
}

// ListTiers retrieves the list of available instance tiers for deployments in VictoriaMetrics Cloud.
func (a *VMCloudAPIClient) ListTiers(ctx context.Context) (TierInfoList, error) {
	return requestAPI[TierInfoList](ctx, a, http.MethodGet, nil, "/api/v1/tiers")
}

// ListDeployments retrieves a list of deployment summaries from for the current account (API Key) in the VictoriaMetrics Cloud API.
func (a *VMCloudAPIClient) ListDeployments(ctx context.Context) (DeploymentSummaryList, error) {
	return requestAPI[DeploymentSummaryList](ctx, a, http.MethodGet, nil, "/api/v1/deployments")
}

// GetDeploymentDetails retrieves detailed information about a specific deployment using its deployment ID.
func (a *VMCloudAPIClient) GetDeploymentDetails(ctx context.Context, deploymentID string) (DeploymentInfo, error) {
	if err := checkDeploymentID(deploymentID); err != nil {
		return DeploymentInfo{}, err
	}
	return requestAPI[DeploymentInfo](ctx, a, http.MethodGet, nil, "/api/v1/deployments", deploymentID)
}

// CreateDeployment creates a new deployment in VictoriaMetrics Cloud based on the provided deployment configuration.
func (a *VMCloudAPIClient) CreateDeployment(ctx context.Context, deployment DeploymentCreationRequest) (DeploymentInfo, error) {
	// Validate common parameters
	err := validateCommonDeploymentParams(
		deployment.Name,
		deployment.Tier,
		deployment.MaintenanceWindow,
		deployment.StorageSize,
		deployment.StorageSizeUnit,
		deployment.Retention,
		deployment.RetentionUnit,
		deployment.DeduplicationUnit,
	)
	if err != nil {
		return DeploymentInfo{}, err
	}

	// Validate creation-specific parameters
	err = validateCreateDeploymentParams(
		deployment.Type,
		deployment.Region,
		deployment.Provider,
		deployment.StorageSize,
		deployment.StorageSizeUnit,
	)
	if err != nil {
		return DeploymentInfo{}, err
	}

	body, err := json.Marshal(deployment)
	if err != nil {
		return DeploymentInfo{}, fmt.Errorf("failed to marshal deployment create request: %w", err)
	}
	return requestAPI[DeploymentInfo](ctx, a, http.MethodPost, bytes.NewReader(body), "/api/v1/deployments")
}

// UpdateDeployment updates the configuration of an existing deployment using the provided deployment ID and update request.
func (a *VMCloudAPIClient) UpdateDeployment(ctx context.Context, deploymentID string, deployment DeploymentUpdateRequest) (DeploymentInfo, error) {
	if err := checkDeploymentID(deploymentID); err != nil {
		return DeploymentInfo{}, err
	}

	// Validate common parameters
	err := validateCommonDeploymentParams(
		deployment.Name,
		deployment.Tier,
		deployment.MaintenanceWindow,
		deployment.StorageSize,
		deployment.StorageSizeUnit,
		deployment.Retention,
		deployment.RetentionUnit,
		deployment.DeduplicationUnit,
	)
	if err != nil {
		return DeploymentInfo{}, err
	}

	body, err := json.Marshal(deployment)
	if err != nil {
		return DeploymentInfo{}, fmt.Errorf("failed to marshal deployment update request: %w", err)
	}
	return requestAPI[DeploymentInfo](ctx, a, http.MethodPut, bytes.NewReader(body), "/api/v1/deployments", deploymentID)
}

// DeleteDeployment deletes an existing deployment using its deployment ID.
func (a *VMCloudAPIClient) DeleteDeployment(ctx context.Context, deploymentID string) error {
	if err := checkDeploymentID(deploymentID); err != nil {
		return err
	}
	_, err := requestAPI[any](ctx, a, http.MethodDelete, nil, "/api/v1/deployments", deploymentID)
	if err != nil {
		return fmt.Errorf("failed to delete deployment %q: %w", deploymentID, err)
	}
	return nil
}

// ListDeploymentAccessTokens retrieves a list of access tokens for a specific deployment using its deployment ID.
func (a *VMCloudAPIClient) ListDeploymentAccessTokens(ctx context.Context, deploymentID string) (AccessTokensList, error) {
	if err := checkDeploymentID(deploymentID); err != nil {
		return nil, err
	}
	return requestAPI[AccessTokensList](ctx, a, http.MethodGet, nil, "/api/v1/deployments", deploymentID, "access_tokens")
}

// CreateDeploymentAccessToken creates a new access token for a specific deployment using its deployment ID and the provided access token creation request.
func (a *VMCloudAPIClient) CreateDeploymentAccessToken(ctx context.Context, deploymentID string, token AccessTokenCreateRequest) (AccessToken, error) {
	if err := checkDeploymentID(deploymentID); err != nil {
		return AccessToken{}, err
	}
	if token.Description == "" {
		return AccessToken{}, fmt.Errorf("access token description cannot be empty")
	}
	if token.Type != AccessModeRead && token.Type != AccessModeWrite && token.Type != AccessModeReadWrite {
		return AccessToken{}, fmt.Errorf("invalid access token type: %s", token.Type)
	}
	if isValidTenantID(token.TenantID) {
		return AccessToken{}, fmt.Errorf("invalid tenant ID format: %s, expected <accountID> or <accountID>:<projectID>", token.TenantID)
	}
	body, err := json.Marshal(token)
	if err != nil {
		return AccessToken{}, fmt.Errorf("failed to marshal access token creation request: %w", err)
	}
	return requestAPI[AccessToken](ctx, a, http.MethodPost, bytes.NewReader(body), "/api/v1/deployments", deploymentID, "access_tokens")
}

// RevealDeploymentAccessToken retrieves the details of a specific access token with full secret value for a deployment using its deployment ID and token ID.
func (a *VMCloudAPIClient) RevealDeploymentAccessToken(ctx context.Context, deploymentID, tokenID string) (AccessToken, error) {
	if err := checkDeploymentID(deploymentID); err != nil {
		return AccessToken{}, err
	}
	if tokenID == "" {
		return AccessToken{}, fmt.Errorf("token ID cannot be empty")
	}
	return requestAPI[AccessToken](ctx, a, http.MethodGet, nil, "/api/v1/deployments", deploymentID, "access_tokens", tokenID)
}

// DeleteDeploymentAccessToken deletes a specific access token for a deployment using the deployment ID and token ID.
func (a *VMCloudAPIClient) DeleteDeploymentAccessToken(ctx context.Context, deploymentID, tokenID string) error {
	if err := checkDeploymentID(deploymentID); err != nil {
		return err
	}
	if tokenID == "" {
		return fmt.Errorf("token ID cannot be empty")
	}
	_, err := requestAPI[any](ctx, a, http.MethodDelete, nil, "/api/v1/deployments", deploymentID, "access_tokens", tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete access token %q for deployment %q: %w", tokenID, deploymentID, err)
	}
	return nil
}

// ListDeploymentRuleFileNames retrieves the list of slerting/recording rules file names associated with a specific deployment by deployment ID.
func (a *VMCloudAPIClient) ListDeploymentRuleFileNames(ctx context.Context, deploymentID string) ([]string, error) {
	if err := checkDeploymentID(deploymentID); err != nil {
		return nil, err
	}
	return requestAPI[[]string](ctx, a, http.MethodGet, nil, "/api/v1/deployments", deploymentID, "rule-sets", "files")
}

// GetDeploymentRuleFileContent retrieves the content of a specific alerting/recording rules file for a deployment by deployment ID and file name.
func (a *VMCloudAPIClient) GetDeploymentRuleFileContent(ctx context.Context, deploymentID, ruleFileName string) (string, error) {
	if err := checkDeploymentID(deploymentID); err != nil {
		return "", err
	}
	if ruleFileName == "" {
		return "", fmt.Errorf("rule file name cannot be empty")
	}
	return requestAPI[string](ctx, a, http.MethodGet, nil, "/api/v1/deployments", deploymentID, "rule-sets", "files", ruleFileName)
}

// UpdateDeploymentRuleFileContent updates the content of an existing alerting/recording rules file for a deployment by deployment ID and file name.
func (a *VMCloudAPIClient) UpdateDeploymentRuleFileContent(ctx context.Context, deploymentID, ruleFileName, content string) error {
	if err := checkDeploymentID(deploymentID); err != nil {
		return err
	}
	if ruleFileName == "" {
		return fmt.Errorf("rule file name cannot be empty")
	}
	body := bytes.NewBufferString(content)
	_, err := requestAPI[any](ctx, a, http.MethodPost, body, "/api/v1/deployments", deploymentID, "rule-sets", "files", ruleFileName)
	if err != nil {
		return fmt.Errorf("failed to update rule file %q for deployment %q: %w", ruleFileName, deploymentID, err)
	}
	return nil
}

// CreateDeploymentRuleFileContent creates a new alerting/recording rules file for a deployment by deployment ID and file name.
func (a *VMCloudAPIClient) CreateDeploymentRuleFileContent(ctx context.Context, deploymentID, ruleFileName, content string) error {
	if err := checkDeploymentID(deploymentID); err != nil {
		return err
	}
	if ruleFileName == "" {
		return fmt.Errorf("rule file name cannot be empty")
	}
	body := bytes.NewBufferString(content)
	_, err := requestAPI[any](ctx, a, http.MethodPost, body, "/api/v1/deployments", deploymentID, "rule-sets", "files", ruleFileName)
	if err != nil {
		return fmt.Errorf("failed to update rule file %q for deployment %q: %w", ruleFileName, deploymentID, err)
	}
	return nil
}

// DeleteDeploymentRuleFile deletes an existing alerting/recording rules file for a deployment by deployment ID and file name.
func (a *VMCloudAPIClient) DeleteDeploymentRuleFile(ctx context.Context, deploymentID, ruleFileName string) error {
	if err := checkDeploymentID(deploymentID); err != nil {
		return err
	}
	if ruleFileName == "" {
		return fmt.Errorf("rule file name cannot be empty")
	}
	_, err := requestAPI[any](ctx, a, http.MethodDelete, nil, "/api/v1/deployments", deploymentID, "rule-sets", "files", ruleFileName)
	if err != nil {
		return fmt.Errorf("failed to delete rule file %q for deployment %q: %w", ruleFileName, deploymentID, err)
	}
	return nil
}
