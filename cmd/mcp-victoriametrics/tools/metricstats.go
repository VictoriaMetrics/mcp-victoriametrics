package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameMetricStats = "metric_statistics"

func toolMetricStats(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("Statistics of metrics usage in queries of the VictoriaMetrics instance. This tool helps to identify unused (never queried) or rarely used metrics or conversely actively queried metrics. This tool uses `/api/v1/status/metric_names_stats` endpoint of VictoriaMetrics API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Metric statistics",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withTargetingOptions(options, c, true, true)
	options = append(
		options,
		mcp.WithString("match_pattern",
			mcp.Title("A regex pattern to match metric names"),
			mcp.Description("A regex pattern to match metric names for showing usage statistics. For example, ?match_pattern=vm_ will match any metric names with vm_ pattern, like vm_http_requests, max_vm_memory_available."),
			mcp.DefaultString(""),
		),
		mcp.WithNumber("limit",
			mcp.Title("Maximum number of metric names"),
			mcp.Description("Integer value to limit the number of metric names in response. By default, API returns 1000 records."),
			mcp.DefaultNumber(1000),
		),
		mcp.WithNumber("le",
			mcp.Title("Less than or equal"),
			mcp.Description("less than or equal, is an integer threshold for filtering metric names by their usage count in queries. For example, with ?le=1 API returns metric names that were queried <=1 times."),
		),
	)
	return mcp.NewTool(toolNameMetricStats, options...)
}

func toolMetricStatsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	matchPattern, err := GetToolReqParam[string](tcr, "match_pattern", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit, err := GetToolReqParam[float64](tcr, "limit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	le, err := GetToolReqParam[float64](tcr, "le", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "api", "v1", "status", "metric_names_stats")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	q := req.URL.Query()
	if matchPattern != "" {
		q.Add("match_pattern", matchPattern)
	}
	if limit != 0 {
		q.Add("limit", fmt.Sprintf("%.f", limit))
	}
	if le != 0 {
		q.Add("le", fmt.Sprintf("%.f", le))
	}
	req.URL.RawQuery = q.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolMetricStats(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameMetricStats) {
		return
	}
	s.AddTool(toolMetricStats(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolMetricStatsHandler(ctx, c, request)
	})
}
