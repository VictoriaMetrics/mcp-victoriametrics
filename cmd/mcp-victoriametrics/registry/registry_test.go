package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func clearEnv(t *testing.T) {
	t.Helper()
	for _, name := range []string{
		"VM_INSTANCE_BACKEND", "VM_INSTANCE_ENTRYPOINT", "VM_INSTANCE_TYPE", "VM_INSTANCE_BEARER_TOKEN", "VM_INSTANCE_HEADERS", "VM_DEFAULT_TENANT_ID",
		"VL_INSTANCE_ENTRYPOINT", "VT_INSTANCE_ENTRYPOINT", "VMC_API_KEY", "VMC_API_BASE_URL", "VMC_BACKEND_TYPES",
		"MCP_SERVER_MODE", "MCP_CONFIG_FILE", "MCP_DISABLED_TOOLS", "MCP_TOOLSETS_ADMIN",
	} {
		t.Setenv(name, "")
	}
	os.Unsetenv("MCP_DISABLED_TOOLS")
}

func newServer(t *testing.T) *server.MCPServer {
	t.Helper()
	cfg, err := config.InitConfig()
	if err != nil {
		t.Fatalf("failed to init config: %v", err)
	}
	s := server.NewMCPServer("test", "", server.WithToolCapabilities(cfg.HasDynamicToolsets()))
	Register(s, cfg)
	return s
}

func callTool(t *testing.T, s *server.MCPServer, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()
	st := s.GetTool(name)
	if st == nil {
		t.Fatalf("tool %s is not registered", name)
	}
	var req mcp.CallToolRequest
	req.Params.Name = name
	req.Params.Arguments = args
	res, err := st.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("tool %s returned error: %v", name, err)
	}
	return res
}

func resultText(res *mcp.CallToolResult) string {
	var sb strings.Builder
	for _, c := range res.Content {
		if tc, ok := c.(mcp.TextContent); ok {
			sb.WriteString(tc.Text)
		}
	}
	return sb.String()
}

func backendEnum(t *testing.T, tool mcp.Tool) []string {
	t.Helper()
	prop, ok := tool.InputSchema.Properties[BackendParam].(map[string]any)
	if !ok {
		return nil
	}
	if !slices.Contains(tool.InputSchema.Required, BackendParam) {
		t.Errorf("tool %s: backend parameter is not required", tool.Name)
	}
	enum, _ := prop["enum"].([]string)
	return enum
}

func TestRegisterSingleBackendKeepsSchema(t *testing.T) {
	clearEnv(t)
	t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://vm")
	t.Setenv("VM_INSTANCE_TYPE", "single")

	s := newServer(t)
	for name, st := range s.ListTools() {
		if _, ok := st.Tool.InputSchema.Properties[BackendParam]; ok {
			t.Errorf("tool %s has backend parameter in single backend mode", name)
		}
	}
	if s.GetTool("metrics_query") == nil || s.GetTool("logs_query") != nil {
		t.Error("Expected only metrics tools to be registered")
	}
}

func TestRegisterMultipleStaticBackends(t *testing.T) {
	var (
		mu    sync.Mutex
		paths []string
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		_, _ = fmt.Fprint(w, `{"status":"success","data":{"resultType":"vector","result":[]}}`)
	}))
	defer srv.Close()

	clearEnv(t)
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := fmt.Sprintf(`
backends:
  - {name: prod, type: metrics, entrypoint: "%[1]s/prod", instance_type: cluster}
  - {name: staging, type: metrics, entrypoint: "%[1]s/staging", instance_type: single}
  - {name: logs, type: logs, entrypoint: "%[1]s/logs"}
`, srv.URL)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	t.Setenv("MCP_CONFIG_FILE", path)

	s := newServer(t)

	if got := backendEnum(t, s.GetTool("metrics_query").Tool); !slices.Equal(got, []string{"prod", "staging"}) {
		t.Errorf("metrics_query backend enum = %v", got)
	}
	if got := backendEnum(t, s.GetTool("alerts").Tool); !slices.Equal(got, []string{"prod", "staging", "logs"}) {
		t.Errorf("alerts backend enum = %v", got)
	}
	for _, name := range []string{"logs_query", "flags", "documentation"} {
		st := s.GetTool(name)
		if st == nil {
			t.Fatalf("tool %s is not registered", name)
		}
		if _, ok := st.Tool.InputSchema.Properties[BackendParam]; ok {
			t.Errorf("tool %s provided by a single backend has backend parameter", name)
		}
	}

	res := callTool(t, s, "metrics_query", map[string]any{"query": "up"})
	if !res.IsError || !strings.Contains(resultText(res), "available backends: prod, staging") {
		t.Errorf("Expected missing backend error, got: %s", resultText(res))
	}
	res = callTool(t, s, "metrics_query", map[string]any{"query": "up", "backend": "dev"})
	if !res.IsError || !strings.Contains(resultText(res), `unknown backend "dev"`) {
		t.Errorf("Expected unknown backend error, got: %s", resultText(res))
	}

	for _, tc := range []struct {
		backend string
		path    string
	}{
		{backend: "staging", path: "/staging/api/v1/query"},
		{backend: "prod", path: "/prod/select/0/prometheus/api/v1/query"},
	} {
		mu.Lock()
		paths = nil
		mu.Unlock()
		res = callTool(t, s, "metrics_query", map[string]any{"query": "up", "backend": tc.backend})
		if res.IsError {
			t.Fatalf("metrics_query for %s failed: %s", tc.backend, resultText(res))
		}
		mu.Lock()
		got := slices.Clone(paths)
		mu.Unlock()
		if !slices.Equal(got, []string{tc.path}) {
			t.Errorf("metrics_query for %s requested %v, want %s", tc.backend, got, tc.path)
		}
	}
}

type fakeAPI struct {
	url   string
	vm    string
	vl    string
	vt    string
	mu    sync.Mutex
	paths []string
}

func newFakeAPI(t *testing.T, idPrefix string) *fakeAPI {
	t.Helper()
	api := &fakeAPI{
		vm: idPrefix + "111111-1111-4111-8111-111111111111",
		vl: idPrefix + "222222-2222-4222-8222-222222222222",
		vt: idPrefix + "333333-3333-4333-8333-333333333333",
	}
	deployments := map[string]vmcloud.DeploymentType{
		api.vm: vmcloud.DeploymentTypeSingleNode,
		api.vl: vmcloud.DeploymentTypeVLogs,
		api.vt: vmcloud.DeploymentTypeVTraces,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rest, ok := strings.CutPrefix(r.URL.Path, "/api/v1/deployments/"); ok {
			id, sub, _ := strings.Cut(rest, "/")
			typ, ok := deployments[id]
			if !ok {
				http.NotFound(w, r)
				return
			}
			var resp any
			switch sub {
			case "":
				resp = vmcloud.DeploymentInfo{ID: id, Type: typ, AccessEndpoint: api.url + "/" + id}
			case "access_tokens":
				resp = vmcloud.AccessTokensList{{ID: "read", Type: vmcloud.AccessModeRead}}
			case "access_tokens/read":
				resp = vmcloud.AccessToken{ID: "read", Type: vmcloud.AccessModeRead, Secret: "read-secret"}
			default:
				http.NotFound(w, r)
				return
			}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		api.mu.Lock()
		api.paths = append(api.paths, r.URL.Path)
		api.mu.Unlock()
		_, _ = fmt.Fprint(w, `{"status":"success","data":{"groups":[]}}`)
	}))
	t.Cleanup(srv.Close)
	api.url = srv.URL
	return api
}

func (api *fakeAPI) takePaths() []string {
	api.mu.Lock()
	defer api.mu.Unlock()
	paths := api.paths
	api.paths = nil
	return paths
}

func TestRegisterMultipleCloudBackends(t *testing.T) {
	api := newFakeAPI(t, "aa")
	clearEnv(t)
	t.Setenv("VMC_API_KEY", "key")
	t.Setenv("VMC_API_BASE_URL", api.url)
	t.Setenv("VMC_BACKEND_TYPES", "metrics,logs,traces")

	s := newServer(t)

	for _, name := range []string{"rules", "metrics_query", "logs_query", "services"} {
		st := s.GetTool(name)
		if st == nil {
			t.Fatalf("tool %s is not registered", name)
		}
		if _, ok := st.Tool.InputSchema.Properties[BackendParam]; ok {
			t.Errorf("tool %s has backend parameter in cloud mode", name)
		}
		if !slices.Contains(st.Tool.InputSchema.Required, "deployment_id") {
			t.Errorf("tool %s does not require deployment_id", name)
		}
	}

	for _, tc := range []struct {
		deploymentID string
		path         string
	}{
		{deploymentID: api.vm, path: "/" + api.vm + "/vmalert/api/v1/rules"},
		{deploymentID: api.vl, path: "/" + api.vl + "/select/vmalert/api/v1/rules"},
		{deploymentID: api.vt, path: "/" + api.vt + "/select/vmalert/api/v1/rules"},
	} {
		api.takePaths()
		res := callTool(t, s, "rules", map[string]any{"deployment_id": tc.deploymentID})
		if res.IsError {
			t.Fatalf("rules for %s failed: %s", tc.deploymentID, resultText(res))
		}
		if got := api.takePaths(); !slices.Equal(got, []string{tc.path}) {
			t.Errorf("rules for %s requested %v, want %s", tc.deploymentID, got, tc.path)
		}
	}

	res := callTool(t, s, "logs_query", map[string]any{"deployment_id": api.vm, "query": "*", "start": "1h"})
	if !res.IsError || !strings.Contains(resultText(res), "unsupported deployment type") {
		t.Errorf("Expected unsupported deployment type for logs_query on metrics deployment, got: %s", resultText(res))
	}
}

func TestRegisterStaticAndCloudBackends(t *testing.T) {
	api := newFakeAPI(t, "bb")
	clearEnv(t)
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := fmt.Sprintf(`
cloud:
  api_key: key
  base_url: %[1]s
  backend_types: [metrics, logs]
backends:
  - {name: onprem, type: logs, entrypoint: "%[1]s/onprem"}
`, api.url)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	t.Setenv("MCP_CONFIG_FILE", path)

	s := newServer(t)

	for _, name := range []string{"rules", "logs_query"} {
		tool := s.GetTool(name).Tool
		if got := backendEnum(t, tool); !slices.Equal(got, []string{"onprem", "cloud"}) {
			t.Errorf("%s backend enum = %v", name, got)
		}
		if slices.Contains(tool.InputSchema.Required, "deployment_id") {
			t.Errorf("%s requires deployment_id for static backends", name)
		}
		prop, _ := tool.InputSchema.Properties["deployment_id"].(map[string]any)
		if description, _ := prop["description"].(string); !strings.Contains(description, `Required when backend is "cloud"`) {
			t.Errorf("%s deployment_id description = %q", name, description)
		}
	}
	metricsQuery := s.GetTool("metrics_query").Tool
	if _, ok := metricsQuery.InputSchema.Properties[BackendParam]; ok {
		t.Error("metrics_query provided only by cloud has backend parameter")
	}
	if !slices.Contains(metricsQuery.InputSchema.Required, "deployment_id") {
		t.Error("metrics_query provided only by cloud does not require deployment_id")
	}
	if s.GetTool("deployments") == nil {
		t.Error("cloud tools are not registered")
	}

	for _, tc := range []struct {
		args map[string]any
		path string
	}{
		{args: map[string]any{"backend": "onprem"}, path: "/onprem/select/vmalert/api/v1/rules"},
		{args: map[string]any{"backend": "cloud", "deployment_id": api.vm}, path: "/" + api.vm + "/vmalert/api/v1/rules"},
		{args: map[string]any{"backend": "cloud", "deployment_id": api.vl}, path: "/" + api.vl + "/select/vmalert/api/v1/rules"},
	} {
		api.takePaths()
		res := callTool(t, s, "rules", tc.args)
		if res.IsError {
			t.Fatalf("rules with %v failed: %s", tc.args, resultText(res))
		}
		if got := api.takePaths(); !slices.Equal(got, []string{tc.path}) {
			t.Errorf("rules with %v requested %v, want %s", tc.args, got, tc.path)
		}
	}

	res := callTool(t, s, "rules", map[string]any{"backend": "cloud"})
	if !res.IsError || !strings.Contains(resultText(res), "deployment_id") {
		t.Errorf("Expected missing deployment_id error, got: %s", resultText(res))
	}
	res = callTool(t, s, "rules", map[string]any{"backend": "cloud", "deployment_id": api.vt})
	if !res.IsError || !strings.Contains(resultText(res), "unsupported deployment type") {
		t.Errorf("Expected unsupported deployment type for traces deployment, got: %s", resultText(res))
	}
}

type fakeSession struct {
	mu    sync.Mutex
	tools map[string]server.ServerTool
	ch    chan mcp.JSONRPCNotification
}

func (fs *fakeSession) Initialize()       {}
func (fs *fakeSession) Initialized() bool { return true }
func (fs *fakeSession) SessionID() string { return "fake-session" }
func (fs *fakeSession) NotificationChannel() chan<- mcp.JSONRPCNotification {
	return fs.ch
}
func (fs *fakeSession) GetSessionTools() map[string]server.ServerTool {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return maps.Clone(fs.tools)
}
func (fs *fakeSession) SetSessionTools(tools map[string]server.ServerTool) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.tools = maps.Clone(tools)
}

func writeAllBackendsConfig(t *testing.T, toolsets string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := toolsets + `
cloud: {api_key: key, backend_types: [metrics]}
backends:
  - {name: vm, type: metrics, entrypoint: "http://vm", instance_type: cluster}
  - {name: vl, type: logs, entrypoint: "http://vl"}
  - {name: vt, type: traces, entrypoint: "http://vt"}
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	t.Setenv("MCP_CONFIG_FILE", path)
}

func TestToolsetToolsExist(t *testing.T) {
	clearEnv(t)
	t.Setenv("MCP_DISABLED_TOOLS", "")
	writeAllBackendsConfig(t, "")

	s := newServer(t)
	for toolset, names := range toolsetTools {
		for _, name := range names {
			if s.GetTool(name) == nil {
				t.Errorf("tool %s of toolset %s is not registered by any backend", name, toolset)
			}
		}
	}
}

func TestToolsetAdminModes(t *testing.T) {
	adminTools := []string{"tsdb_status", "top_queries", "flags", "tiers", "access_tokens"}
	regularTools := []string{"metrics_query", "logs_query", "services", "deployments", "alerts", "documentation"}

	for _, tc := range []struct {
		name           string
		env            string
		file           string
		expectAdmin    bool
		expectEnable   bool
		expectPrompts  bool
		expectListened bool
	}{
		{name: "Default", expectAdmin: true, expectPrompts: true},
		{name: "Env off", env: "off"},
		{name: "File off", file: "toolsets: {admin: off}"},
		{name: "Env dynamic", env: "dynamic", expectEnable: true, expectPrompts: true},
		{name: "File on", file: "toolsets: {Admin: ON}", expectAdmin: true, expectPrompts: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			clearEnv(t)
			t.Setenv("MCP_DISABLED_TOOLS", "")
			if tc.env != "" {
				t.Setenv("MCP_TOOLSETS_ADMIN", tc.env)
			}
			writeAllBackendsConfig(t, tc.file)

			s := newServer(t)
			for _, name := range adminTools {
				if got := s.GetTool(name) != nil; got != tc.expectAdmin {
					t.Errorf("admin tool %s registered = %v, want %v", name, got, tc.expectAdmin)
				}
			}
			for _, name := range regularTools {
				if s.GetTool(name) == nil {
					t.Errorf("regular tool %s is not registered", name)
				}
			}
			if got := s.GetTool(EnableToolsetTool) != nil; got != tc.expectEnable {
				t.Errorf("%s registered = %v, want %v", EnableToolsetTool, got, tc.expectEnable)
			}
			prompts := s.ListPrompts()
			if _, got := prompts["unused_metrics"]; got != tc.expectPrompts {
				t.Errorf("unused_metrics prompt registered = %v, want %v", got, tc.expectPrompts)
			}
		})
	}
}

func TestEnableToolset(t *testing.T) {
	clearEnv(t)
	t.Setenv("MCP_DISABLED_TOOLS", "")
	t.Setenv("MCP_TOOLSETS_ADMIN", "dynamic")
	writeAllBackendsConfig(t, "")

	s := newServer(t)
	enable := s.GetTool(EnableToolsetTool)
	if got := backendEnumValues(enable.Tool, "name"); !slices.Equal(got, []string{"admin"}) {
		t.Errorf("enable_toolset name enum = %v", got)
	}

	res := callTool(t, s, EnableToolsetTool, map[string]any{"name": "debug"})
	if !res.IsError || !strings.Contains(resultText(res), `unknown toolset "debug"`) {
		t.Errorf("Expected unknown toolset error, got: %s", resultText(res))
	}

	session := &fakeSession{ch: make(chan mcp.JSONRPCNotification, 10)}
	if err := s.RegisterSession(context.Background(), session); err != nil {
		t.Fatalf("failed to register session: %v", err)
	}
	var req mcp.CallToolRequest
	req.Params.Name = EnableToolsetTool
	req.Params.Arguments = map[string]any{"name": "admin"}
	res, err := enable.Handler(s.WithContext(context.Background(), session), req)
	if err != nil || res.IsError {
		t.Fatalf("enable_toolset in session failed: %v %s", err, resultText(res))
	}
	sessionTools := session.GetSessionTools()
	for _, name := range []string{"tsdb_status", "flags", "tiers"} {
		if _, ok := sessionTools[name]; !ok {
			t.Errorf("tool %s is not added to session", name)
		}
	}
	if fb, ok := sessionTools["flags"]; ok {
		if got := backendEnumValues(fb.Tool, BackendParam); !slices.Equal(got, []string{"vm", "vl", "cloud"}) {
			t.Errorf("session flags backend enum = %v", got)
		}
	}
	if s.GetTool("tsdb_status") != nil {
		t.Error("tool enabled in session is registered globally")
	}
	select {
	case n := <-session.ch:
		if n.Method != "notifications/tools/list_changed" {
			t.Errorf("unexpected notification %s", n.Method)
		}
	default:
		t.Error("expected tools/list_changed notification")
	}

	res = callTool(t, s, EnableToolsetTool, map[string]any{"name": "admin"})
	if res.IsError || !strings.Contains(resultText(res), "tsdb_status") {
		t.Fatalf("enable_toolset without session failed: %s", resultText(res))
	}
	if s.GetTool("tsdb_status") == nil {
		t.Error("tool enabled without session is not registered globally")
	}
}

func backendEnumValues(tool mcp.Tool, param string) []string {
	prop, _ := tool.InputSchema.Properties[param].(map[string]any)
	enum, _ := prop["enum"].([]string)
	return enum
}
