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
	if c.IsCloud() {
		options = append(
			options,
			mcp.WithString("deployment_id",
				mcp.Title("Deployment ID"),
				mcp.Description("Unique identifier of the deployment in VictoriaMetrics Cloud"),
				mcp.Pattern(`^[a-zA-Z0-9\-_]+$`),
			),
		)
	}
	if c.IsCluster() || c.IsCloud() {
		options = append(
			options,
			mcp.WithString("tenant",
				mcp.Title("Tenant name"),
				mcp.Description("Name of the tenant for which the data will be displayed"),
				mcp.DefaultString("0"),
				mcp.Pattern(`^([0-9]+)(:[0-9]+)?$`),
			),
		)
	}
	options = append(
		options,
		mcp.WithString("query",
			mcp.Required(),
			mcp.Title("MetricsQL or PromQL expression"),
			mcp.Description(`MetricsQL or PromQL expression for prettification. This is the query that will be formatted.`),
		),
	)
	return mcp.NewTool(toolNamePrettifyQuery, append(options, withEnvironmentParam())...)
}

func toolPrettifyQueryHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	envName, _ := GetToolReqParam[string](tcr, "env", false)
	env, err := cfg.Environment(envName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

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
	} else if env.IsCloud() {
		deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", false)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get deployment_id parameter: %v", err)), nil
		}
		if deploymentID == "" {
			return mcp.NewToolResultErrorFromErr("failed to prettify query: ", err), nil
		}
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
