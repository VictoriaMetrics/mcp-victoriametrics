package registry

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/prompts"
	metricsprompts "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/prompts/metrics"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
	cloudtools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/cloud"
	logstools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/logs"
	metricstools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/metrics"
	tracestools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/traces"
	vmalerttools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/vmalert"
)

const (
	BackendParam      = "backend"
	EnableToolsetTool = "enable_toolset"
)

var toolsetTools = map[string][]string{
	config.ToolsetAdmin: {
		"tsdb_status", "top_queries", "active_queries", "metric_statistics", "tenants", "flags", "export",
		"tiers", "regions", "cloud_providers", "access_tokens", "rule_file", "rule_filenames",
	},
}

// Register registers tools and prompts of the configured backends.
func Register(s *server.MCPServer, c *config.Config) {
	tools.Register(s, c)
	prompts.Register(s, c)

	byName := make(map[string][]backendTool)
	collect := func(b *config.Config, register func(s *server.MCPServer, c *config.Config)) {
		tmp := server.NewMCPServer(b.Name(), "", server.WithToolCapabilities(false))
		register(tmp, b)
		for name, st := range tmp.ListTools() {
			byName[name] = append(byName[name], backendTool{backend: b, tool: st.Tool, handler: st.Handler})
		}
	}
	collect(c, cloudtools.Register)
	for _, b := range c.Backends() {
		collect(b, registerBackendTools)
	}

	dynamic := make(map[string][]server.ServerTool)
	for _, name := range slices.Sorted(maps.Keys(byName)) {
		entries := byName[name]
		st := server.ServerTool{Tool: entries[0].tool, Handler: entries[0].handler}
		if len(entries) > 1 {
			st.Tool, st.Handler = mergeTools(entries)
		}
		toolset := toolsetOf(name)
		switch c.ToolsetMode(toolset) {
		case config.ToolsetOff:
		case config.ToolsetDynamic:
			dynamic[toolset] = append(dynamic[toolset], st)
		default:
			s.AddTools(st)
		}
	}
	registerEnableToolset(s, dynamic)

	if c.ToolsetMode(config.ToolsetAdmin) != config.ToolsetOff {
		for _, b := range c.Backends() {
			if b.Backend() == config.BackendMetrics {
				metricsprompts.Register(s, b)
				break
			}
		}
	}
}

func toolsetOf(name string) string {
	for toolset, names := range toolsetTools {
		if slices.Contains(names, name) {
			return toolset
		}
	}
	return ""
}

func registerBackendTools(s *server.MCPServer, c *config.Config) {
	vmalerttools.Register(s, c)
	switch c.Backend() {
	case config.BackendLogs:
		logstools.Register(s, c)
	case config.BackendTraces:
		tracestools.Register(s, c)
	default:
		metricstools.Register(s, c)
	}
}

func registerEnableToolset(s *server.MCPServer, dynamic map[string][]server.ServerTool) {
	if len(dynamic) == 0 {
		return
	}
	toolsets := slices.Sorted(maps.Keys(dynamic))
	descriptions := make([]string, 0, len(toolsets))
	for _, toolset := range toolsets {
		descriptions = append(descriptions, fmt.Sprintf("%s (%s)", toolset, strings.Join(toolNames(dynamic[toolset]), ", ")))
	}
	tool := mcp.NewTool(EnableToolsetTool,
		mcp.WithDescription("Enable a toolset to make its tools available. Toolsets: "+strings.Join(descriptions, "; ")),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Enable toolset",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(false),
		}),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Title("Toolset name"),
			mcp.Enum(toolsets...),
		),
	)
	s.AddTool(tool, func(ctx context.Context, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := tools.GetToolReqParam[string](tcr, "name", true)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		enabled, ok := dynamic[name]
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("unknown toolset %q; available toolsets: %s", name, strings.Join(toolsets, ", "))), nil
		}
		session := server.ClientSessionFromContext(ctx)
		if _, ok := session.(server.SessionWithTools); ok {
			if err := s.AddSessionTools(session.SessionID(), enabled...); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("failed to enable toolset %q: %v", name, err)), nil
			}
		} else {
			s.AddTools(enabled...)
		}
		return mcp.NewToolResultText(fmt.Sprintf("Toolset %q is enabled, available tools: %s", name, strings.Join(toolNames(enabled), ", "))), nil
	})
}

func toolNames(tools []server.ServerTool) []string {
	names := make([]string, 0, len(tools))
	for _, t := range tools {
		names = append(names, t.Tool.Name)
	}
	return names
}

func ptr[T any](v T) *T {
	return &v
}

type backendTool struct {
	backend *config.Config
	tool    mcp.Tool
	handler server.ToolHandlerFunc
}

func mergeTools(entries []backendTool) (mcp.Tool, server.ToolHandlerFunc) {
	tool := entries[0].tool
	tool.InputSchema.Properties = maps.Clone(tool.InputSchema.Properties)
	if tool.InputSchema.Properties == nil {
		tool.InputSchema.Properties = make(map[string]any)
	}
	required := slices.Clone(tool.InputSchema.Required)
	for _, e := range entries[1:] {
		for name, prop := range e.tool.InputSchema.Properties {
			if _, ok := tool.InputSchema.Properties[name]; !ok {
				tool.InputSchema.Properties[name] = prop
			}
		}
		required = slices.DeleteFunc(required, func(name string) bool {
			return !slices.Contains(e.tool.InputSchema.Required, name)
		})
	}

	tool.InputSchema.Required = required

	var staticEntries, cloudEntries []backendTool
	for _, e := range entries {
		if e.backend.IsCloud() {
			cloudEntries = append(cloudEntries, e)
		} else {
			staticEntries = append(staticEntries, e)
		}
	}
	if len(staticEntries) == 0 {
		return tool, cloudDispatchHandler(cloudEntries)
	}

	names := make([]string, 0, len(staticEntries)+1)
	descriptions := make([]string, 0, len(staticEntries)+1)
	for _, e := range staticEntries {
		names = append(names, e.backend.Name())
		descriptions = append(descriptions, fmt.Sprintf("%s (%s)", e.backend.Name(), e.backend.Backend().ProductName()))
	}
	cloudName := ""
	if len(cloudEntries) > 0 {
		cloudName = cloudEntries[0].backend.Name()
		products := make([]string, 0, len(cloudEntries))
		for _, e := range cloudEntries {
			products = append(products, e.backend.Backend().ProductName())
		}
		names = append(names, cloudName)
		descriptions = append(descriptions, fmt.Sprintf("%s (VictoriaMetrics Cloud deployments of %s, requires %s parameter)", cloudName, strings.Join(products, ", "), tools.DeploymentIDParam))
		if prop, ok := tool.InputSchema.Properties[tools.DeploymentIDParam].(map[string]any); ok {
			prop = maps.Clone(prop)
			description, _ := prop["description"].(string)
			prop["description"] = fmt.Sprintf("%s. Required when %s is %q", description, BackendParam, cloudName)
			tool.InputSchema.Properties[tools.DeploymentIDParam] = prop
		}
	}
	mcp.WithString(BackendParam,
		mcp.Required(),
		mcp.Title("Backend"),
		mcp.Description("Name of the backend to send the request to: "+strings.Join(descriptions, ", ")),
		mcp.Enum(names...),
	)(&tool)

	cloudHandler := cloudDispatchHandler(cloudEntries)
	handler := func(ctx context.Context, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := tools.GetToolReqParam[string](tcr, BackendParam, true)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("%v; available backends: %s", err, strings.Join(names, ", "))), nil
		}
		if cloudName != "" && name == cloudName {
			return cloudHandler(ctx, tcr)
		}
		for _, e := range staticEntries {
			if e.backend.Name() == name {
				return e.handler(ctx, tcr)
			}
		}
		return mcp.NewToolResultError(fmt.Sprintf("unknown backend %q; available backends: %s", name, strings.Join(names, ", "))), nil
	}
	return tool, handler
}

func cloudDispatchHandler(entries []backendTool) server.ToolHandlerFunc {
	return func(ctx context.Context, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		info, err := tools.LookupCloudDeploymentInfo(ctx, entries[0].backend, tcr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get cloud deployment info: %v", err)), nil
		}
		for _, e := range entries {
			if slices.Contains(tools.CloudDeploymentTypes[e.backend.Backend()], info.DeploymentType) {
				return e.handler(ctx, tcr)
			}
		}
		return mcp.NewToolResultError(fmt.Sprintf("unsupported deployment type %s for deployment %s", info.DeploymentType, info.ID)), nil
	}
}
