package tools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameTenants = "tenants"

func toolTenants(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of tenants of the VictoriaMetrics instance.  This tool uses `/admin/tenants` endpoint of VictoriaMetrics API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of tenants",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withClusterAdminTargetingOptions(options, c, true)
	return mcp.NewTool(toolNameTenants, options...)
}

func toolTenantsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if _, err := getClusterAdminToolInstance(cfg, tcr); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := CreateAdminRequest(ctx, cfg, tcr, "admin", "tenants")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolTenants(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameTenants) {
		return
	}
	if !c.HasClusterInstances() {
		return
	}
	s.AddTool(toolTenants(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolTenantsHandler(ctx, c, request)
	})
}
