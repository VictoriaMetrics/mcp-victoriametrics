package tools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameQuery = "query"

func toolQuery(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("Instant query executes PromQL or MetricsQL query expression at the given time. The result of Instant query is a list of time series matching the filter in query expression. Each returned series contains exactly one (timestamp, value) entry, where timestamp equals to the time query arg, while the value contains query result at the requested time. This tool uses `/api/v1/query` endpoint of VictoriaMetrics API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Instant Query",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withTargetingOptions(options, c, true, true)
	options = append(
		options,
		mcp.WithString("query",
			mcp.Required(),
			mcp.Title("MetricsQL or PromQL expression"),
			mcp.Description(`MetricsQL or PromQL expression for the query of the data`),
		),
		mcp.WithString("time",
			mcp.Title("Timestamp"),
			mcp.Description("Timestamp in millisecond precision to evaluate the query at. If omitted, time is set to now() (current timestamp). The time param can be specified in multiple allowed formats."),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithString("step",
			mcp.Title("Step"),
			mcp.Description("Optional interval for searching for raw samples in the past when executing the query (used when a sample is missing at the specified time). For example, the request /api/v1/query?query=up&step=1m looks for the last written raw sample for the metric up in the (now()-1m, now()] interval (the first millisecond is not included). If omitted, step is set to 5m (5 minutes) by default."),
			mcp.Pattern(`^([0-9]+)([a-z]+)$`),
		),
		mcp.WithString("timeout",
			mcp.Title("Timeout"),
			mcp.Description("Optional query timeout. For example, timeout=5s. Query is canceled when the timeout is reached. By default the timeout is set to the value of -search.maxQueryDuration command-line flag passed to single-node VictoriaMetrics or to vmselect component of VictoriaMetrics cluster."),
			mcp.Pattern(`^([0-9]+)([a-z]+)$`),
		),
		mcp.WithBoolean("trace",
			mcp.Title("Enable query trace"),
			mcp.Description("If true, the query will be traced and the trace will be returned in the response. This is useful for debugging and performance analysis."),
			mcp.DefaultBool(false),
		),
		mcp.WithBoolean("nocache",
			mcp.Title("Disable cache"),
			mcp.Description("If true, the query will not use the rollup cache on execution."),
			mcp.DefaultBool(false),
		),
	)
	return mcp.NewTool(toolNameQuery, options...)
}

func toolQueryHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := GetToolReqParam[string](tcr, "query", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	time, err := GetToolReqParam[string](tcr, "time", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	step, err := GetToolReqParam[string](tcr, "step", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	timeout, err := GetToolReqParam[string](tcr, "timeout", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	trace, err := GetToolReqParam[bool](tcr, "trace", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	nocache, err := GetToolReqParam[bool](tcr, "nocache", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "api", "v1", "query")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	q := req.URL.Query()
	q.Add("query", query)
	if time != "" {
		q.Add("time", time)
	}
	if step != "" {
		q.Add("step", step)
	}
	if timeout != "" {
		q.Add("timeout", timeout)
	}
	if trace {
		q.Add("trace", "1")
	}
	if nocache {
		q.Add("nocache", "1")
	}
	req.URL.RawQuery = q.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolQuery(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameQuery) {
		return
	}
	s.AddTool(toolQuery(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolQueryHandler(ctx, c, request)
	})
}
