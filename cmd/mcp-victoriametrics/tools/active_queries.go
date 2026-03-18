package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameActiveQueries = "active_queries"

func toolActiveQueries(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription(`Active queries. This tool can determine currently active queries in the VictoriaMetrics instance.
This information is obtained from the "/api/v1/status/active_queries" HTTP endpoint of VictoriaMetrics API.`),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Active queries",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	if c.IsCloud() {
		options = append(
			options,
			mcp.WithString("deployment_id",
				mcp.Required(),
				mcp.Title("Deployment ID"),
				mcp.Description("Unique identifier of the deployment in VictoriaMetrics Cloud"),
				mcp.Pattern(`^[a-zA-Z0-9\-_]+$`),
			),
		)
	}
	if c.IsCluster() || c.IsCloud() {
		options = append(options,
			mcp.WithString("tenant",
				mcp.Title("Tenant name"),
				mcp.Description("Name of the tenant for which the active queries will be displayed"),
				mcp.DefaultString("0"),
				mcp.Pattern(`^([0-9]+)(:[0-9]+)?$`),
			),
		)
	}
	return mcp.NewTool(toolNameActiveQueries, options...)
}

func toolActiveQueriesHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	req, err := CreateSelectRequest(ctx, cfg, tcr, "api", "v1", "status", "active_queries")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolActiveQueries(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameActiveQueries) {
		return
	}
	s.AddTool(toolActiveQueries(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolActiveQueriesHandler(ctx, c, request)
	})
}
