package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/VictoriaMetrics/metricsql"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNamePrettifyQuery = "prettify_query"

func toolPrettifyQuery(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("Prettify (format) MetricsQL query. This tool uses `/prettify-query` endpoint of VictoriaMetrics API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Prettify Query",
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
			mcp.Description(`MetricsQL or PromQL expression for prettification. This is the query that will be formatted.`),
		),
	)
	return mcp.NewTool(toolNamePrettifyQuery, options...)
}

func toolPrettifyQueryHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := GetToolReqParam[string](tcr, "query", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	prettifiedQuery, err := metricsql.Prettify(query)
	if err != nil {
		result := map[string]string{
			"status": "success",
			"query":  prettifiedQuery,
		}
		data, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "prettify-query")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolPrettifyQuery(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNamePrettifyQuery) {
		return
	}
	s.AddTool(toolPrettifyQuery(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolPrettifyQueryHandler(ctx, c, request)
	})
}
