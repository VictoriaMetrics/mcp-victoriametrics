package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/auth"
	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"
)

const (
	toolsDisabledByDefault = "export,flags,metric_relabel_debug,downsampling_filters_debug,retention_filters_debug,test_rules"
	defaultInstanceName    = "default"
)

type Instance struct {
	name            string
	entrypoint      string
	instanceType    string
	bearerToken     string
	customHeaders   map[string]string
	defaultTenantID string
	entryPointURL   *url.URL
	vmc             *vmcloud.VMCloudAPIClient
}

func (i *Instance) Name() string {
	return i.name
}

func (i *Instance) IsCluster() bool {
	return i.instanceType == "cluster"
}

func (i *Instance) IsSingle() bool {
	return i.instanceType == "single"
}

func (i *Instance) IsCloud() bool {
	return i.vmc != nil
}

func (i *Instance) VMC() *vmcloud.VMCloudAPIClient {
	return i.vmc
}

func (i *Instance) BearerToken() string {
	return i.bearerToken
}

func (i *Instance) EntryPointURL() *url.URL {
	return i.entryPointURL
}

func (i *Instance) CustomHeaders() map[string]string {
	return i.customHeaders
}

func (i *Instance) DefaultTenantID() string {
	return i.defaultTenantID
}

type Config struct {
	serverMode        string
	listenAddr        string
	disabledTools     map[string]bool
	heartbeatInterval time.Duration
	disableResources  bool
	logFormat         string
	logLevel          string
	instances         map[string]*Instance
	instanceOrder     []string
	defaultInstance   string
}

type instanceSpec struct {
	name               string
	entrypointEnv      string
	instanceTypeEnv    string
	bearerTokenEnv     string
	headersEnv         string
	defaultTenantIDEnv string
	apiKeyEnv          string
}

func InitConfig() (*Config, error) {
	disabledTools, isDisabledToolsSet := os.LookupEnv("MCP_DISABLED_TOOLS")
	if disabledTools == "" && !isDisabledToolsSet {
		disabledTools = toolsDisabledByDefault
	}

	heartbeatInterval := 30 * time.Second
	if value := os.Getenv("MCP_HEARTBEAT_INTERVAL"); value != "" {
		interval, err := time.ParseDuration(value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MCP_HEARTBEAT_INTERVAL: %w", err)
		}
		if interval < 0 {
			return nil, fmt.Errorf("MCP_HEARTBEAT_INTERVAL must be a non-negative")
		}
		heartbeatInterval = interval
	}

	disableResources := false
	if value := os.Getenv("MCP_DISABLE_RESOURCES"); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MCP_DISABLE_RESOURCES: %w", err)
		}
		disableResources = parsed
	}

	logFormat := strings.ToLower(os.Getenv("MCP_LOG_FORMAT"))
	if logFormat == "" {
		logFormat = "text"
	}
	if logFormat != "text" && logFormat != "json" {
		return nil, fmt.Errorf("MCP_LOG_FORMAT must be 'text' or 'json'")
	}

	logLevel := strings.ToLower(os.Getenv("MCP_LOG_LEVEL"))
	if logLevel == "" {
		logLevel = "info"
	}
	if logLevel != "debug" && logLevel != "info" && logLevel != "warn" && logLevel != "error" {
		return nil, fmt.Errorf("MCP_LOG_LEVEL must be 'debug', 'info', 'warn' or 'error'")
	}

	serverMode := strings.ToLower(os.Getenv("MCP_SERVER_MODE"))
	if serverMode == "" {
		serverMode = "stdio"
	}
	if serverMode != "stdio" && serverMode != "sse" && serverMode != "http" {
		return nil, fmt.Errorf("MCP_SERVER_MODE must be 'stdio', 'sse' or 'http'")
	}

	listenAddr := os.Getenv("MCP_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = os.Getenv("MCP_SSE_ADDR")
	}
	if listenAddr == "" {
		listenAddr = "localhost:8080"
	}

	instances, order, defaultInstance, err := initInstances()
	if err != nil {
		return nil, err
	}

	return &Config{
		serverMode:        serverMode,
		listenAddr:        listenAddr,
		disabledTools:     parseDisabledTools(disabledTools),
		heartbeatInterval: heartbeatInterval,
		disableResources:  disableResources,
		logFormat:         logFormat,
		logLevel:          logLevel,
		instances:         instances,
		instanceOrder:     order,
		defaultInstance:   defaultInstance,
	}, nil
}

func initInstances() (map[string]*Instance, []string, string, error) {
	specs, defaultName, err := getInstanceSpecs()
	if err != nil {
		return nil, nil, "", err
	}

	instances := make(map[string]*Instance, len(specs))
	order := make([]string, 0, len(specs))
	for _, spec := range specs {
		instance, err := newInstance(spec)
		if err != nil {
			return nil, nil, "", err
		}
		instances[spec.name] = instance
		order = append(order, spec.name)
	}
	return instances, order, defaultName, nil
}

func getInstanceSpecs() ([]instanceSpec, string, error) {
	if value := os.Getenv("VM_ENVIRONMENTS"); value != "" {
		if err := validateLegacyVarsUnused(); err != nil {
			return nil, "", err
		}

		names, err := parseInstanceNames(value)
		if err != nil {
			return nil, "", err
		}

		defaultName := strings.TrimSpace(strings.ToLower(os.Getenv("VM_DEFAULT_ENVIRONMENT")))
		if defaultName == "" {
			defaultName = names[0]
		}
		if !contains(names, defaultName) {
			return nil, "", fmt.Errorf("VM_DEFAULT_ENVIRONMENT %q is not listed in VM_ENVIRONMENTS", defaultName)
		}

		specs := make([]instanceSpec, 0, len(names))
		for _, name := range names {
			specs = append(specs, instanceSpecForEnv(name))
		}
		return specs, defaultName, nil
	}

	return []instanceSpec{legacyInstanceSpec()}, defaultInstanceName, nil
}

func legacyInstanceSpec() instanceSpec {
	return instanceSpec{
		name:               defaultInstanceName,
		entrypointEnv:      "VM_INSTANCE_ENTRYPOINT",
		instanceTypeEnv:    "VM_INSTANCE_TYPE",
		bearerTokenEnv:     "VM_INSTANCE_BEARER_TOKEN",
		headersEnv:         "VM_INSTANCE_HEADERS",
		defaultTenantIDEnv: "VM_DEFAULT_TENANT_ID",
		apiKeyEnv:          "VMC_API_KEY",
	}
}

func instanceSpecForEnv(name string) instanceSpec {
	prefix := instancePrefix(name)
	return instanceSpec{
		name:               name,
		entrypointEnv:      prefix + "ENTRYPOINT",
		instanceTypeEnv:    prefix + "TYPE",
		bearerTokenEnv:     prefix + "BEARER_TOKEN",
		headersEnv:         prefix + "HEADERS",
		defaultTenantIDEnv: prefix + "DEFAULT_TENANT_ID",
		apiKeyEnv:          "VMC_" + strings.ToUpper(name) + "_API_KEY",
	}
}

func newInstance(spec instanceSpec) (*Instance, error) {
	entrypoint := os.Getenv(spec.entrypointEnv)
	instanceType := os.Getenv(spec.instanceTypeEnv)
	bearerToken := os.Getenv(spec.bearerTokenEnv)
	headers := parseHeaders(os.Getenv(spec.headersEnv))
	defaultTenantID := os.Getenv(spec.defaultTenantIDEnv)
	apiKey := os.Getenv(spec.apiKeyEnv)

	if entrypoint == "" && apiKey == "" {
		return nil, fmt.Errorf("%s or %s is not set", spec.entrypointEnv, spec.apiKeyEnv)
	}
	if entrypoint != "" && apiKey != "" {
		return nil, fmt.Errorf("env %q: %s and %s cannot be set at the same time", spec.name, spec.entrypointEnv, spec.apiKeyEnv)
	}
	if entrypoint != "" && instanceType == "" {
		return nil, fmt.Errorf("%s is not set", spec.instanceTypeEnv)
	}
	if entrypoint != "" && instanceType != "single" && instanceType != "cluster" {
		return nil, fmt.Errorf("%s must be 'single' or 'cluster'", spec.instanceTypeEnv)
	}

	resolvedTenantID := "0"
	if defaultTenantID != "" {
		tenantID, err := auth.NewToken(strings.ToLower(defaultTenantID))
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s %q: %w", spec.defaultTenantIDEnv, defaultTenantID, err)
		}
		resolvedTenantID = tenantID.String()
	}

	instance := &Instance{
		name:            spec.name,
		entrypoint:      entrypoint,
		instanceType:    instanceType,
		bearerToken:     bearerToken,
		customHeaders:   headers,
		defaultTenantID: resolvedTenantID,
	}

	var err error
	if apiKey != "" {
		instance.vmc, err = vmcloud.New(apiKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create VMCloud API client from %s: %w", spec.apiKeyEnv, err)
		}
		return instance, nil
	}

	instance.entryPointURL, err = url.Parse(entrypoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL from %s: %w", spec.entrypointEnv, err)
	}
	return instance, nil
}

func parseDisabledTools(value string) map[string]bool {
	disabled := make(map[string]bool)
	if value == "" {
		return disabled
	}
	for _, tool := range strings.Split(value, ",") {
		tool = strings.Trim(tool, " ,")
		if tool != "" {
			disabled[tool] = true
		}
	}
	return disabled
}

func parseHeaders(value string) map[string]string {
	headers := make(map[string]string)
	if value == "" {
		return headers
	}
	for _, header := range strings.Split(value, ",") {
		header = strings.TrimSpace(header)
		parts := strings.SplitN(header, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		headerValue := strings.TrimSpace(parts[1])
		if key != "" && headerValue != "" {
			headers[key] = headerValue
		}
	}
	return headers
}

func validateLegacyVarsUnused() error {
	for _, envVar := range []string{
		"VM_INSTANCE_ENTRYPOINT",
		"VM_INSTANCE_TYPE",
		"VM_INSTANCE_BEARER_TOKEN",
		"VM_INSTANCE_HEADERS",
		"VM_DEFAULT_TENANT_ID",
		"VMC_API_KEY",
	} {
		if os.Getenv(envVar) != "" {
			return fmt.Errorf("%s cannot be combined with VM_ENVIRONMENTS", envVar)
		}
	}
	return nil
}

func parseInstanceNames(value string) ([]string, error) {
	seen := make(map[string]struct{})
	names := make([]string, 0)
	for _, raw := range strings.Split(value, ",") {
		name := strings.TrimSpace(strings.ToLower(raw))
		if name == "" {
			continue
		}
		if !isValidInstanceName(name) {
			return nil, fmt.Errorf("VM_ENVIRONMENTS contains invalid env name %q; use lowercase letters, numbers, and underscores only", raw)
		}
		if _, ok := seen[name]; ok {
			return nil, fmt.Errorf("VM_ENVIRONMENTS contains duplicate env name %q", name)
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	if len(names) == 0 {
		return nil, fmt.Errorf("VM_ENVIRONMENTS is set but does not contain any env names")
	}
	return names, nil
}

func isValidInstanceName(value string) bool {
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			continue
		}
		return false
	}
	return true
}

func instancePrefix(name string) string {
	return "VM_INSTANCE_" + strings.ToUpper(name) + "_"
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func (c *Config) ResolveInstance(name string) (*Instance, error) {
	if len(c.instances) == 0 {
		return nil, fmt.Errorf("no VictoriaMetrics instances configured")
	}
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		name = c.defaultInstance
	}
	instance, ok := c.instances[name]
	if !ok {
		return nil, fmt.Errorf("unknown env %q; available envs: %s", name, strings.Join(c.instanceOrder, ", "))
	}
	return instance, nil
}

func (c *Config) DefaultInstanceName() string {
	return c.defaultInstance
}

func (c *Config) InstanceNames() []string {
	return append([]string(nil), c.instanceOrder...)
}

func (c *Config) HasMultipleInstances() bool {
	return len(c.instanceOrder) > 1
}

func (c *Config) HasCloudInstances() bool {
	for _, name := range c.instanceOrder {
		if c.instances[name].IsCloud() {
			return true
		}
	}
	return false
}

func (c *Config) HasClusterInstances() bool {
	for _, name := range c.instanceOrder {
		instance := c.instances[name]
		if instance.IsCluster() || instance.IsCloud() {
			return true
		}
	}
	return false
}

func (c *Config) HasOnlyCloudInstances() bool {
	if len(c.instanceOrder) == 0 {
		return false
	}
	for _, name := range c.instanceOrder {
		if !c.instances[name].IsCloud() {
			return false
		}
	}
	return true
}

func (c *Config) IsStdio() bool {
	return c.serverMode == "stdio"
}

func (c *Config) IsSSE() bool {
	return c.serverMode == "sse"
}

func (c *Config) ServerMode() string {
	return c.serverMode
}

func (c *Config) ListenAddr() string {
	return c.listenAddr
}

func (c *Config) IsToolDisabled(toolName string) bool {
	if c.disabledTools == nil {
		return false
	}
	disabled, ok := c.disabledTools[toolName]
	return ok && disabled
}

func (c *Config) IsResourcesDisabled() bool {
	return c.disableResources
}

func (c *Config) HeartbeatInterval() time.Duration {
	return c.heartbeatInterval
}

func (c *Config) LogFormat() string {
	return c.logFormat
}

func (c *Config) LogLevel() string {
	return c.logLevel
}
