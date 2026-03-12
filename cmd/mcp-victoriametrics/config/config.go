package config

import (
	"fmt"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/auth"
	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"
)

const (
	toolsDisabledByDefault = "export,flags,metric_relabel_debug,downsampling_filters_debug,retention_filters_debug,test_rules"
	defaultEnvironmentName = "default"
)

type InstanceConfig struct {
	name            string
	entrypoint      string
	instanceType    string
	bearerToken     string
	customHeaders   map[string]string
	defaultTenantID string
	apiKey          string
	entryPointURL   *url.URL
	vmc             *vmcloud.VMCloudAPIClient
}

func (s *InstanceConfig) Name() string {
	return s.name
}

func (s *InstanceConfig) IsCluster() bool {
	return s.instanceType == "cluster"
}

func (s *InstanceConfig) IsSingle() bool {
	return s.instanceType == "single"
}

func (s *InstanceConfig) IsCloud() bool {
	return s.vmc != nil
}

func (s *InstanceConfig) VMC() *vmcloud.VMCloudAPIClient {
	return s.vmc
}

func (s *InstanceConfig) BearerToken() string {
	return s.bearerToken
}

func (s *InstanceConfig) EntryPointURL() *url.URL {
	return s.entryPointURL
}

func (s *InstanceConfig) CustomHeaders() map[string]string {
	return s.customHeaders
}

func (s *InstanceConfig) DefaultTenantID() string {
	return s.defaultTenantID
}

type Config struct {
	serverMode         string
	listenAddr         string
	disabledTools      map[string]bool
	heartbeatInterval  time.Duration
	disableResources   bool
	environments       map[string]*InstanceConfig
	environmentOrder   []string
	defaultEnvironment string

	// Logging configuration
	logFormat string
	logLevel  string
}

func InitConfig() (*Config, error) {
	disabledTools, isDisabledToolsSet := os.LookupEnv("MCP_DISABLED_TOOLS")
	if disabledTools == "" && !isDisabledToolsSet {
		disabledTools = toolsDisabledByDefault
	}
	disabledToolsMap := make(map[string]bool)
	if disabledTools != "" {
		for _, tool := range strings.Split(disabledTools, ",") {
			tool = strings.Trim(tool, " ,")
			if tool != "" {
				disabledToolsMap[tool] = true
			}
		}
	}

	heartbeatInterval := 30 * time.Second
	heartbeatIntervalStr := os.Getenv("MCP_HEARTBEAT_INTERVAL")
	if heartbeatIntervalStr != "" {
		interval, err := time.ParseDuration(heartbeatIntervalStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MCP_HEARTBEAT_INTERVAL: %w", err)
		}
		if interval < 0 {
			return nil, fmt.Errorf("MCP_HEARTBEAT_INTERVAL must be a non-negative")
		}
		heartbeatInterval = interval
	}

	disableResources := false
	disableResourcesStr := os.Getenv("MCP_DISABLE_RESOURCES")
	if disableResourcesStr != "" {
		var err error
		disableResources, err = strconv.ParseBool(disableResourcesStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse MCP_DISABLE_RESOURCES: %w", err)
		}
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

	result := &Config{
		serverMode:        strings.ToLower(os.Getenv("MCP_SERVER_MODE")),
		listenAddr:        os.Getenv("MCP_LISTEN_ADDR"),
		disabledTools:     disabledToolsMap,
		heartbeatInterval: heartbeatInterval,
		disableResources:  disableResources,
		logFormat:         logFormat,
		logLevel:          logLevel,
	}

	if result.serverMode != "" && result.serverMode != "stdio" && result.serverMode != "sse" && result.serverMode != "http" {
		return nil, fmt.Errorf("MCP_SERVER_MODE must be 'stdio', 'sse' or 'http'")
	}
	if result.serverMode == "" {
		result.serverMode = "stdio"
	}
	if result.listenAddr == "" {
		result.listenAddr = os.Getenv("MCP_SSE_ADDR")
	}
	if result.listenAddr == "" {
		result.listenAddr = "localhost:8080"
	}

	environments, environmentOrder, defaultEnvironment, err := initEnvironmentConfigs()
	if err != nil {
		return nil, err
	}
	result.environments = environments
	result.environmentOrder = environmentOrder
	result.defaultEnvironment = defaultEnvironment

	return result, nil
}

func initEnvironmentConfigs() (map[string]*InstanceConfig, []string, string, error) {
	if envNamesValue := os.Getenv("VM_ENVIRONMENTS"); envNamesValue != "" {
		if err := validateNoStandardInstanceConfig(); err != nil {
			return nil, nil, "", err
		}

		envNames, err := parseEnvironmentNames(envNamesValue)
		if err != nil {
			return nil, nil, "", err
		}

		defaultEnvironment := strings.TrimSpace(strings.ToLower(os.Getenv("VM_DEFAULT_ENVIRONMENT")))
		if defaultEnvironment == "" {
			defaultEnvironment = envNames[0]
		}
		if !slices.Contains(envNames, defaultEnvironment) {
			return nil, nil, "", fmt.Errorf("VM_DEFAULT_ENVIRONMENT %q is not listed in VM_ENVIRONMENTS", defaultEnvironment)
		}

		environments := make(map[string]*InstanceConfig, len(envNames))
		for _, envName := range envNames {
			prefix := environmentVarPrefix(envName)
			instance, err := newInstanceConfig(
				envName,
				os.Getenv(prefix+"ENTRYPOINT"),
				os.Getenv(prefix+"TYPE"),
				os.Getenv(prefix+"BEARER_TOKEN"),
				parseHeaders(os.Getenv(prefix+"HEADERS")),
				os.Getenv(prefix+"DEFAULT_TENANT_ID"),
				os.Getenv("VMC_"+strings.ToUpper(envName)+"_API_KEY"),
			)
			if err != nil {
				return nil, nil, "", err
			}
			environments[envName] = instance
		}

		return environments, envNames, defaultEnvironment, nil
	}

	instance, err := newInstanceConfig(
		defaultEnvironmentName,
		os.Getenv("VM_INSTANCE_ENTRYPOINT"),
		os.Getenv("VM_INSTANCE_TYPE"),
		os.Getenv("VM_INSTANCE_BEARER_TOKEN"),
		parseHeaders(os.Getenv("VM_INSTANCE_HEADERS")),
		os.Getenv("VM_DEFAULT_TENANT_ID"),
		os.Getenv("VMC_API_KEY"),
	)
	if err != nil {
		return nil, nil, "", err
	}

	return map[string]*InstanceConfig{defaultEnvironmentName: instance}, []string{defaultEnvironmentName}, defaultEnvironmentName, nil
}

func validateNoStandardInstanceConfig() error {
	for _, envVar := range []string{
		"VM_INSTANCE_ENTRYPOINT",
		"VM_INSTANCE_TYPE",
		"VM_INSTANCE_BEARER_TOKEN",
		"VM_INSTANCE_HEADERS",
		"VM_DEFAULT_TENANT_ID",
		"VMC_API_KEY",
	} {
		if os.Getenv(envVar) != "" {
			return fmt.Errorf("%s cannot be combined with VM_ENVIRONMENTS; use per-environment variables instead", envVar)
		}
	}
	return nil
}

func parseEnvironmentNames(value string) ([]string, error) {
	names := make([]string, 0)
	seenNames := make(map[string]struct{})
	seenPrefixes := make(map[string]string)

	for _, rawName := range strings.Split(value, ",") {
		name := strings.TrimSpace(strings.ToLower(rawName))
		if name == "" {
			continue
		}
		if !isValidEnvironmentName(name) {
			return nil, fmt.Errorf("VM_ENVIRONMENTS contains invalid env name %q; only letters, numbers, dashes, and underscores are allowed", rawName)
		}
		if _, ok := seenNames[name]; ok {
			return nil, fmt.Errorf("VM_ENVIRONMENTS contains duplicate env name %q", name)
		}

		prefix := environmentVarPrefix(name)
		if existingName, ok := seenPrefixes[prefix]; ok {
			return nil, fmt.Errorf("VM_ENVIRONMENTS names %q and %q map to the same environment variable prefix %q", existingName, name, prefix)
		}

		seenNames[name] = struct{}{}
		seenPrefixes[prefix] = name
		names = append(names, name)
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("VM_ENVIRONMENTS is set but does not contain any env names")
	}

	return names, nil
}

func isValidEnvironmentName(value string) bool {
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return false
	}
	return true
}

func environmentVarPrefix(name string) string {
	var b strings.Builder
	b.WriteString("VM_INSTANCE_")
	for _, r := range strings.ToUpper(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		b.WriteByte('_')
	}
	b.WriteByte('_')
	return b.String()
}

func newInstanceConfig(name, entrypoint, instanceType, bearerToken string, customHeaders map[string]string, defaultTenantID, apiKey string) (*InstanceConfig, error) {
	if entrypoint == "" && apiKey == "" {
		if name == defaultEnvironmentName {
			return nil, fmt.Errorf("VM_INSTANCE_ENTRYPOINT or VMC_API_KEY is not set")
		}
		return nil, fmt.Errorf("%sENTRYPOINT or VMC_%s_API_KEY is not set", environmentVarPrefix(name), strings.ToUpper(name))
	}
	if entrypoint != "" && apiKey != "" {
		return nil, fmt.Errorf("env %q: ENTRYPOINT and API_KEY cannot be set at the same time", name)
	}
	if entrypoint != "" && instanceType == "" {
		return nil, fmt.Errorf("env %q: INSTANCE_TYPE is not set", name)
	}
	if entrypoint != "" && instanceType != "cluster" && instanceType != "single" {
		return nil, fmt.Errorf("env %q: INSTANCE_TYPE must be 'single' or 'cluster'", name)
	}

	var entryPointURL *url.URL
	var vmc *vmcloud.VMCloudAPIClient
	var err error

	if apiKey == "" {
		entryPointURL, err = url.Parse(entrypoint)
		if err != nil {
			return nil, fmt.Errorf("env %q: failed to parse URL from ENTRYPOINT: %w", name, err)
		}
	} else {
		vmc, err = vmcloud.New(apiKey)
		if err != nil {
			return nil, fmt.Errorf("env %q: failed to create VMCloud API client: %w", name, err)
		}
	}

	resolvedTenantID := "0"
	if defaultTenantID != "" {
		tenantID, err := auth.NewToken(defaultTenantID)
		if err != nil {
			return nil, fmt.Errorf("env %q: failed to parse DEFAULT_TENANT_ID %q: %w", name, defaultTenantID, err)
		}
		resolvedTenantID = tenantID.String()
	}

	return &InstanceConfig{
		name:            name,
		entrypoint:      entrypoint,
		instanceType:    instanceType,
		bearerToken:     bearerToken,
		customHeaders:   customHeaders,
		defaultTenantID: resolvedTenantID,
		apiKey:          apiKey,
		entryPointURL:   entryPointURL,
		vmc:             vmc,
	}, nil
}

func parseHeaders(value string) map[string]string {
	headers := make(map[string]string)
	if value == "" {
		return headers
	}

	for _, header := range strings.Split(value, ",") {
		header = strings.TrimSpace(header)
		if header == "" {
			continue
		}

		parts := strings.SplitN(header, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		headerValue := strings.TrimSpace(parts[1])
		if key == "" || headerValue == "" {
			continue
		}
		headers[key] = headerValue
	}

	return headers
}

func (c *Config) DefaultEnvironment() string {
	return c.defaultEnvironment
}

func (c *Config) EnvironmentNames() []string {
	return slices.Clone(c.environmentOrder)
}

func (c *Config) Environment(name string) (*InstanceConfig, error) {
	if len(c.environments) == 0 {
		return nil, fmt.Errorf("no VictoriaMetrics environments configured")
	}

	resolvedName := strings.TrimSpace(strings.ToLower(name))
	if resolvedName == "" {
		resolvedName = c.defaultEnvironment
	}

	env, ok := c.environments[resolvedName]
	if !ok {
		return nil, fmt.Errorf("unknown VictoriaMetrics env %q; available envs: %s", resolvedName, strings.Join(c.environmentOrder, ", "))
	}
	return env, nil
}

// Helper methods for accessing default environment configuration
func (c *Config) IsCluster() bool {
	s, _ := c.Environment("")
	if s != nil {
		return s.IsCluster()
	}
	return false
}

func (c *Config) IsSingle() bool {
	s, _ := c.Environment("")
	if s != nil {
		return s.IsSingle()
	}
	return false
}

func (c *Config) IsCloud() bool {
	s, _ := c.Environment("")
	if s != nil {
		return s.IsCloud()
	}
	return false
}

func (c *Config) VMC() *vmcloud.VMCloudAPIClient {
	s, _ := c.Environment("")
	if s != nil {
		return s.vmc
	}
	return nil
}

func (c *Config) BearerToken() string {
	s, _ := c.Environment("")
	if s != nil {
		return s.BearerToken()
	}
	return ""
}

func (c *Config) EntryPointURL() *url.URL {
	s, _ := c.Environment("")
	if s != nil {
		return s.EntryPointURL()
	}
	return nil
}

func (c *Config) CustomHeaders() map[string]string {
	s, _ := c.Environment("")
	if s != nil {
		return s.CustomHeaders()
	}
	return nil
}

func (c *Config) DefaultTenantID() string {
	s, _ := c.Environment("")
	if s != nil {
		return s.DefaultTenantID()
	}
	return "0"
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
