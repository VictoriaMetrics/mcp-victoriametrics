package logs

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

const toolNameQuery = "query"

var (
	toolQuery = mcp.NewTool(toolNameQuery,
		mcp.WithDescription("Executes LogsQL query expression to search log entries. This tool uses `/select/logsql/query` endpoint of VictoriaLogs API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Query logs",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
		mcp.WithString("tenant",
			mcp.Title("Tenant name (Account ID and Project ID)"),
			mcp.Description("Name of the tenant for which the data will be displayed (format AccountID:ProjectID)"),
			mcp.DefaultString("0:0"),
			mcp.Pattern(`^([0-9]+)(:[0-9]+)$`),
		),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Title("LogsQL expression"),
			mcp.Description(`LogsQL expression for the query of the logs data`),
		),
		mcp.WithString("start",
			mcp.Required(),
			mcp.Title("Start timestamp"),
			mcp.Description("Start timestamp in RFC3339 format. For example, 2023-10-01T00:00:00Z"),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithString("end",
			mcp.Title("End timestamp"),
			mcp.Description("End timestamp in RFC3339 format. For example, 2023-10-01T00:00:00Z"),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithNumber("limit",
			mcp.Title("Limit"),
			mcp.Description("The maximum number of log entries, which can be returned in the response"),
			mcp.DefaultNumber(1000),
		),
		mcp.WithString("timeout",
			mcp.Title("Timeout"),
			mcp.Description("Optional query timeout. For example, timeout=5s. Query is canceled when the timeout is reached. "),
			mcp.Pattern(`^([0-9]+)([a-z]+)$`),
		),
	)
)

func toolQueryHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := tools.GetToolReqParam[string](tcr, "query", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start, err := tools.GetToolReqParam[string](tcr, "start", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	end, err := tools.GetToolReqParam[string](tcr, "end", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit, err := tools.GetToolReqParam[float64](tcr, "limit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if limit <= 0 {
		limit = 1000
	}

	timeout, err := tools.GetToolReqParam[string](tcr, "timeout", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "query")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	q := req.URL.Query()
	q.Add("query", query)
	if start != "" {
		q.Add("start", start)
	}
	if end != "" {
		q.Add("end", end)
	}
	if limit != 0 {
		q.Add("limit", fmt.Sprintf("%.f", limit))
	}
	if timeout != "" {
		q.Add("timeout", timeout)
	}
	req.URL.RawQuery = q.Encode()

	return tools.GetTextBodyForRequest(req, cfg), nil
}

func registerToolQuery(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameQuery) {
		return
	}
	s.AddTool(tools.WithCloudDeploymentID(c, toolQuery), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolQueryHandler(ctx, c, request)
	})
}
