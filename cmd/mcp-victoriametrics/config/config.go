package config

import (
	"crypto/tls"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/VictoriaLogs/lib/logstorage"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/auth"
	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"
)

type Backend string

type ToolsetMode string

const (
	ToolsetOff     ToolsetMode = "off"
	ToolsetOn      ToolsetMode = "on"
	ToolsetDynamic ToolsetMode = "dynamic"
)

const ToolsetAdmin = "admin"

const toolsetsEnvPrefix = "MCP_TOOLSETS_"

const (
	BackendMetrics Backend = "metrics"
	BackendLogs    Backend = "logs"
	BackendTraces  Backend = "traces"
)

var (
	toolsDisabledByDefault = map[Backend]string{
		BackendMetrics: "export,flags,metric_relabel_debug,downsampling_filters_debug,retention_filters_debug,test_rules",
	}
	legacyEnvPrefixes = map[Backend]string{
		BackendLogs:   "VL_",
		BackendTraces: "VT_",
	}
	toolsetDefaults = map[string]ToolsetMode{
		ToolsetAdmin: ToolsetOn,
	}
	legacyToolNames = map[Backend]map[string]string{
		BackendMetrics: {
			"query":       "metrics_query",
			"query_range": "metrics_query_range",
		},
		BackendLogs: {
			"query": "logs_query",
		},
	}
)

type Config struct {
	name               string
	backend            Backend
	backends           []*Config
	serverMode         string
	listenAddr         string
	entrypoint         string
	instanceType       string
	bearerToken        string
	disabledTools      map[string]bool
	legacyToolNames    map[string]string
	toolsets           map[string]ToolsetMode
	apiKey             string
	apiBaseURL         string
	heartbeatInterval  time.Duration
	disableResources   bool
	customHeaders      map[string]string
	passthroughHeaders []string
	defaultTenantID    string
	defaultLogsTenant  logstorage.TenantID
	authToken          string

	// Logging configuration
	logFormat string
	logLevel  string

	entryPointURL *url.URL
	vmc           *vmcloud.VMCloudAPIClient
	httpClient    *http.Client
}

func InitConfig() (*Config, error) {
	if os.Getenv("MCP_CONFIG_FILE") != "" {
		for _, name := range []string{"VM_INSTANCE_BACKEND", "VM_INSTANCE_ENTRYPOINT", "VL_INSTANCE_ENTRYPOINT", "VT_INSTANCE_ENTRYPOINT", "VMC_API_KEY", "VMC_API_BASE_URL", "VMC_BACKEND_TYPES"} {
			if os.Getenv(name) != "" {
				return nil, fmt.Errorf("MCP_CONFIG_FILE and %s cannot be set at the same time", name)
			}
		}
	}
	backend, err := parseBackend(os.Getenv("VM_INSTANCE_BACKEND"))
	if err != nil {
		return nil, fmt.Errorf("VM_INSTANCE_BACKEND %w", err)
	}
	var cloudBackends []Backend
	if cloudBackendTypes := os.Getenv("VMC_BACKEND_TYPES"); cloudBackendTypes != "" {
		if os.Getenv("VMC_API_KEY") == "" {
			return nil, fmt.Errorf("VMC_BACKEND_TYPES can be set only with VMC_API_KEY")
		}
		if os.Getenv("VM_INSTANCE_BACKEND") != "" {
			return nil, fmt.Errorf("VM_INSTANCE_BACKEND and VMC_BACKEND_TYPES cannot be set at the same time")
		}
		for _, t := range strings.Split(cloudBackendTypes, ",") {
			if t = strings.TrimSpace(t); t == "" {
				continue
			}
			b, err := parseBackend(t)
			if err != nil {
				return nil, fmt.Errorf("VMC_BACKEND_TYPES: %q %w", t, err)
			}
			if !slices.Contains(cloudBackends, b) {
				cloudBackends = append(cloudBackends, b)
			}
		}
		if len(cloudBackends) == 0 {
			return nil, fmt.Errorf("VMC_BACKEND_TYPES must contain at least one backend type")
		}
		backend = cloudBackends[0]
	}
	getEnv := func(name string) string {
		if v := os.Getenv("VM_" + name); v != "" {
			return v
		}
		if prefix, ok := legacyEnvPrefixes[backend]; ok {
			return os.Getenv(prefix + name)
		}
		return ""
	}

	disabledToolsMap, usedLegacyToolNames := parseDisabledTools(backend)
	customHeadersMap := parseHeaders(getEnv("INSTANCE_HEADERS"))

	var passthroughHeaders []string
	passthroughHeadersStr := os.Getenv("MCP_PASSTHROUGH_HEADERS")
	if passthroughHeadersStr != "" {
		for _, h := range strings.Split(passthroughHeadersStr, ",") {
			h = strings.TrimSpace(h)
			if h != "" {
				passthroughHeaders = append(passthroughHeaders, h)
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

	httpClient := http.DefaultClient
	insecureSkipVerifyStr := os.Getenv("VM_TLS_INSECURE_SKIP_VERIFY")
	if insecureSkipVerifyStr != "" {
		insecureSkipVerify, err := strconv.ParseBool(insecureSkipVerifyStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse VM_TLS_INSECURE_SKIP_VERIFY: %w", err)
		}
		if insecureSkipVerify {
			transport := http.DefaultTransport.(*http.Transport).Clone()
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			httpClient = &http.Client{Transport: transport}
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
		name:               string(backend),
		backend:            backend,
		serverMode:         strings.ToLower(os.Getenv("MCP_SERVER_MODE")),
		listenAddr:         os.Getenv("MCP_LISTEN_ADDR"),
		entrypoint:         getEnv("INSTANCE_ENTRYPOINT"),
		instanceType:       os.Getenv("VM_INSTANCE_TYPE"),
		bearerToken:        getEnv("INSTANCE_BEARER_TOKEN"),
		disabledTools:      disabledToolsMap,
		legacyToolNames:    usedLegacyToolNames,
		apiKey:             os.Getenv("VMC_API_KEY"),
		apiBaseURL:         os.Getenv("VMC_API_BASE_URL"),
		heartbeatInterval:  heartbeatInterval,
		disableResources:   disableResources,
		customHeaders:      customHeadersMap,
		passthroughHeaders: passthroughHeaders,
		authToken:          os.Getenv("MCP_AUTH_TOKEN"),
		logFormat:          logFormat,
		logLevel:           logLevel,
		defaultTenantID:    "0",
		httpClient:         httpClient,
	}
	// Left for backward compatibility
	if result.listenAddr == "" {
		result.listenAddr = os.Getenv("MCP_SSE_ADDR")
	}
	result.toolsets, err = parseToolsetsEnv()
	if err != nil {
		return nil, err
	}

	if result.serverMode != "" && result.serverMode != "stdio" && result.serverMode != "sse" && result.serverMode != "http" {
		return nil, fmt.Errorf("MCP_SERVER_MODE must be 'stdio', 'sse' or 'http'")
	}
	if result.serverMode == "" {
		result.serverMode = "stdio"
	}
	if result.listenAddr == "" {
		result.listenAddr = "localhost:8080"
	}

	if configFile := os.Getenv("MCP_CONFIG_FILE"); configFile != "" {
		if err := result.loadBackendsFile(configFile); err != nil {
			return nil, fmt.Errorf("failed to load MCP_CONFIG_FILE %q: %w", configFile, err)
		}
		return result, nil
	}

	if result.entrypoint == "" && result.apiKey == "" {
		return nil, fmt.Errorf("VM_INSTANCE_ENTRYPOINT or VMC_API_KEY is not set")
	}
	if result.entrypoint != "" && result.apiKey != "" {
		return nil, fmt.Errorf("VM_INSTANCE_ENTRYPOINT and VMC_API_KEY cannot be set at the same time")
	}
	if backend == BackendMetrics {
		if result.entrypoint != "" && result.instanceType == "" {
			return nil, fmt.Errorf("VM_INSTANCE_TYPE is not set")
		}
		if result.entrypoint != "" && result.instanceType != "cluster" && result.instanceType != "single" {
			return nil, fmt.Errorf("VM_INSTANCE_TYPE must be 'single' or 'cluster'")
		}
	}
	if result.apiKey == "" {
		result.entryPointURL, err = url.Parse(result.entrypoint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL from VM_INSTANCE_ENTRYPOINT: %w", err)
		}
	}
	if result.apiKey != "" {
		if err := result.setCloudClient(result.apiKey, result.apiBaseURL); err != nil {
			return nil, err
		}
	}

	if err := result.setDefaultTenantID(getEnv("DEFAULT_TENANT_ID")); err != nil {
		return nil, fmt.Errorf("failed to parse VM_DEFAULT_TENANT_ID: %w", err)
	}

	if len(cloudBackends) > 1 {
		legacy := make(map[string]string)
		for _, cb := range cloudBackends {
			b, err := result.newCloudBackend(cb, string(cb), os.Getenv("VM_DEFAULT_TENANT_ID"))
			if err != nil {
				return nil, fmt.Errorf("failed to parse VM_DEFAULT_TENANT_ID for %s backend: %w", cb, err)
			}
			mergeLegacyToolNames(legacy, b.legacyToolNames)
			result.backends = append(result.backends, b)
		}
		result.legacyToolNames = legacy
	}

	return result, nil
}

func (c *Config) setCloudClient(apiKey, baseURL string) error {
	options := []vmcloud.VMCloudAPIClientOption{vmcloud.WithHTTPClient(c.httpClient)}
	if baseURL != "" {
		options = append(options, vmcloud.WithBaseURL(baseURL))
	}
	vmc, err := vmcloud.New(apiKey, options...)
	if err != nil {
		return fmt.Errorf("failed to create VMCloud API client: %w", err)
	}
	c.apiKey = apiKey
	c.apiBaseURL = baseURL
	c.vmc = vmc
	return nil
}

func (c *Config) newCloudBackend(backend Backend, name, defaultTenantID string) (*Config, error) {
	b := *c
	b.backends = nil
	b.name = name
	b.backend = backend
	b.entrypoint = ""
	b.entryPointURL = nil
	b.instanceType = ""
	b.defaultTenantID = "0"
	b.defaultLogsTenant = logstorage.TenantID{}
	b.disabledTools, b.legacyToolNames = parseDisabledTools(backend)
	if err := b.setDefaultTenantID(defaultTenantID); err != nil {
		return nil, err
	}
	return &b, nil
}

func parseToolsetsEnv() (map[string]ToolsetMode, error) {
	toolsets := make(map[string]ToolsetMode)
	for _, kv := range os.Environ() {
		key, value, _ := strings.Cut(kv, "=")
		name, ok := strings.CutPrefix(key, toolsetsEnvPrefix)
		if !ok || value == "" {
			continue
		}
		if err := setToolsetMode(toolsets, strings.ToLower(name), value); err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}
	}
	return toolsets, nil
}

func setToolsetMode(toolsets map[string]ToolsetMode, name, value string) error {
	if _, ok := toolsetDefaults[name]; !ok {
		return fmt.Errorf("unknown toolset %q, supported toolsets: %s", name, strings.Join(Toolsets(), ", "))
	}
	mode := ToolsetMode(strings.ToLower(strings.TrimSpace(value)))
	switch mode {
	case ToolsetOff, ToolsetOn, ToolsetDynamic:
	default:
		return fmt.Errorf("toolset mode must be '%s', '%s' or '%s'", ToolsetOff, ToolsetOn, ToolsetDynamic)
	}
	toolsets[name] = mode
	return nil
}

func Toolsets() []string {
	return slices.Sorted(maps.Keys(toolsetDefaults))
}

func (c *Config) ToolsetMode(name string) ToolsetMode {
	if mode, ok := c.toolsets[name]; ok {
		return mode
	}
	return toolsetDefaults[name]
}

func (c *Config) HasDynamicToolsets() bool {
	return slices.ContainsFunc(Toolsets(), func(name string) bool {
		return c.ToolsetMode(name) == ToolsetDynamic
	})
}

func parseBackend(value string) (Backend, error) {
	backend := Backend(strings.ToLower(strings.TrimSpace(value)))
	switch backend {
	case "":
		return BackendMetrics, nil
	case BackendMetrics, BackendLogs, BackendTraces:
		return backend, nil
	default:
		return "", fmt.Errorf("must be '%s', '%s' or '%s'", BackendMetrics, BackendLogs, BackendTraces)
	}
}

func mergeLegacyToolNames(dst, src map[string]string) {
	for oldName, newName := range src {
		if prev, ok := dst[oldName]; ok && !strings.Contains(prev, newName) {
			newName = prev + ", " + newName
		}
		dst[oldName] = newName
	}
}

func parseDisabledTools(backend Backend) (map[string]bool, map[string]string) {
	disabledTools, isDisabledToolsSet := os.LookupEnv("MCP_DISABLED_TOOLS")
	if disabledTools == "" && !isDisabledToolsSet {
		disabledTools = toolsDisabledByDefault[backend]
	}
	disabledToolsMap := make(map[string]bool)
	usedLegacyToolNames := make(map[string]string)
	for _, tool := range strings.Split(disabledTools, ",") {
		tool = strings.Trim(tool, " ,")
		if tool == "" {
			continue
		}
		if newName, ok := legacyToolNames[backend][tool]; ok {
			usedLegacyToolNames[tool] = newName
			tool = newName
		}
		disabledToolsMap[tool] = true
	}
	return disabledToolsMap, usedLegacyToolNames
}

func parseHeaders(headers string) map[string]string {
	result := make(map[string]string)
	for _, header := range strings.Split(headers, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(header), "=")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if ok && key != "" && value != "" {
			result[key] = value
		}
	}
	return result
}

func (c *Config) setDefaultTenantID(tenant string) error {
	tenant = strings.ToLower(tenant)
	if tenant == "" {
		return nil
	}
	if c.backend == BackendMetrics {
		tenantID, err := auth.NewToken(tenant)
		if err != nil {
			return fmt.Errorf("cannot parse tenant %q: %w", tenant, err)
		}
		c.defaultTenantID = tenantID.String()
		return nil
	}
	tenantID, err := logstorage.ParseTenantID(tenant)
	if err != nil {
		return fmt.Errorf("cannot parse tenant %q: %w", tenant, err)
	}
	c.defaultLogsTenant = tenantID
	return nil
}

func (c *Config) Backend() Backend {
	return c.backend
}

func (c *Config) Name() string {
	return c.name
}

func (c *Config) Backends() []*Config {
	if len(c.backends) == 0 {
		return []*Config{c}
	}
	return c.backends
}

func (c *Config) BackendTypes() []Backend {
	var types []Backend
	for _, b := range c.Backends() {
		if !slices.Contains(types, b.backend) {
			types = append(types, b.backend)
		}
	}
	return types
}

func (c *Config) CloudBackendTypes() []Backend {
	var types []Backend
	for _, b := range c.Backends() {
		if b.IsCloud() && !slices.Contains(types, b.backend) {
			types = append(types, b.backend)
		}
	}
	return types
}

func (c *Config) ProductName() string {
	types := c.BackendTypes()
	names := make([]string, 0, len(types))
	for _, t := range types {
		names = append(names, t.ProductName())
	}
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

func (c *Config) InstanceType() string {
	return c.instanceType
}

func (c *Config) IsCluster() bool {
	return c.instanceType == "cluster"
}

func (c *Config) IsSingle() bool {
	return c.instanceType == "single"
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

func (c *Config) IsCloud() bool {
	return c.vmc != nil
}

func (c *Config) IsCloudSharedInstance() bool {
	return c.apiKey == vmcloud.DynamicAPIKey
}

func (c *Config) VMC() *vmcloud.VMCloudAPIClient {
	return c.vmc
}

func (c *Config) ListenAddr() string {
	return c.listenAddr
}

func (c *Config) BearerToken() string {
	return c.bearerToken
}

func (c *Config) EntryPointURL() *url.URL {
	return c.entryPointURL
}

func (c *Config) IsToolDisabled(toolName string) bool {
	if c.disabledTools == nil {
		return false
	}
	disabled, ok := c.disabledTools[toolName]
	return ok && disabled
}

func (c *Config) LegacyToolNames() map[string]string {
	return c.legacyToolNames
}

func (c *Config) IsResourcesDisabled() bool {
	return c.disableResources
}

func (c *Config) HeartbeatInterval() time.Duration {
	return c.heartbeatInterval
}

func (c *Config) CustomHeaders() map[string]string {
	return c.customHeaders
}

func (c *Config) PassthroughHeaders() []string {
	return c.passthroughHeaders
}

func (c *Config) LogFormat() string {
	return c.logFormat
}

func (c *Config) LogLevel() string {
	return c.logLevel
}

func (c *Config) DefaultTenantID() string {
	return c.defaultTenantID
}

func (c *Config) DefaultLogsTenantID() logstorage.TenantID {
	return c.defaultLogsTenant
}

func (c *Config) AuthToken() string {
	return c.authToken
}

func (c *Config) HTTPClient() *http.Client {
	if c == nil || c.httpClient == nil {
		return http.DefaultClient
	}
	return c.httpClient
}

func (b Backend) ProductName() string {
	switch b {
	case BackendLogs:
		return "VictoriaLogs"
	case BackendTraces:
		return "VictoriaTraces"
	default:
		return "VictoriaMetrics"
	}
}
