package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameDownsamplingFiltersDebug = "downsampling_filters_debug"

func toolDownsamplingFiltersDebug(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription(`Downsampling filters debug tool is used to debug flag "downsampling.period" with some series and see what downsampling strategy will be applied for which series in Enterprise version of VictoriaMetrics.
This tool use "/downsampling-filters-debug" API endpoint of VictoriaMetrics API.`),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Downsampling filters debugger ",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withTargetingOptions(options, c, true, false)
	options = append(
		options,
		mcp.WithString("flags",
			mcp.Required(),
			mcp.Title("Value of `downsampling.period` flag"),
			mcp.Description("Commandline flags values for `downsampling.period`. For example: `-downsampling.period={env=\"dev\"}:7d:5m,{env=\"dev\"}:30d:30m -downsampling.period=30d:1m -downsampling.period=60d:5m`"),
		),
		mcp.WithString("metrics",
			mcp.Required(),
			mcp.Title("Metrics"),
			mcp.Description(`Set of metrics to be debugged. The metrics should be in the format of <metric_name>{<label_name>="<label_value>",...}.`),
			mcp.Pattern(`^([a-zA-Z_]*\{\s*(([a-zA-Z-_]+\s*\=\s*\".*\"))?(\s*,\s*([a-zA-Z-_]+\s*\=\s*\".*\"))*\s*\}\n)+$`),
		),
	)
	return mcp.NewTool(toolNameDownsamplingFiltersDebug, options...)
}

func toolDownsamplingFiltersDebugHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	flags, err := GetToolReqParam[string](tcr, "flags", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	metrics, err := GetToolReqParam[string](tcr, "metrics", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "downsampling-filters-debug")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	query := req.URL.Query()
	query.Set("flags", flags)
	query.Set("metrics", metrics)
	req.URL.RawQuery = query.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolDownsamplingFiltersDebug(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameDownsamplingFiltersDebug) {
		return
	}
	s.AddTool(toolDownsamplingFiltersDebug(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolDownsamplingFiltersDebugHandler(ctx, c, request)
	})
}
