package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameMetricRelabelDebug = "metric_relabel_debug"

func toolMetricRelabelDebug(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription(`Metric relabel debug tool can help with step-by-step debugging of Prometheus-compatible relabeling rules. It can be used to check how the relabeling rules are applied to the given metric. 
The tool use "/metric-relabel-debug" endpoint of the VictoriaMetrics API. `),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Metric relabel debugger",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withTargetingOptions(options, c, true, false)
	options = append(
		options,
		mcp.WithString("relabel_configs",
			mcp.Required(),
			mcp.Title("Relabel config"),
			mcp.Description("Prometheus-compatible relabeling rules"),
		),
		mcp.WithString("metric",
			mcp.Required(),
			mcp.Title("Metrics"),
			mcp.Description(`Set of metrics to be relabeled. The metrics should be in the format of {<label_name>="<label_value>",...}.`),
			mcp.Pattern(`^\{\s*(([a-zA-Z-_]+\s*\=\s*\".*\"))?(\s*,\s*([a-zA-Z-_]+\s*\=\s*\".*\"))*\s*\}$`),
		),
	)
	return mcp.NewTool(toolNameMetricRelabelDebug, options...)
}

func toolMetricRelabelDebugHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	relabelConfigs, err := GetToolReqParam[string](tcr, "relabel_configs", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	metric, err := GetToolReqParam[string](tcr, "metric", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "metric-relabel-debug")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create select request: %v", err)), nil
	}

	query := req.URL.Query()
	query.Set("relabel_configs", relabelConfigs)
	query.Set("metric", metric)
	query.Set("format", "json")
	req.URL.RawQuery = query.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolMetricRelabelDebug(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameMetricRelabelDebug) {
		return
	}
	s.AddTool(toolMetricRelabelDebug(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolMetricRelabelDebugHandler(ctx, c, request)
	})
}
