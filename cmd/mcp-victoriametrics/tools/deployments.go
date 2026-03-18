package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameDeployments = "deployments"

func toolDeployments(_ *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of deployments in VictoriaMetrics Cloud"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of deployments",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	return mcp.NewTool(toolNameDeployments, options...)
}

func toolDeploymentsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if cfg.IsCloudSharedInstance() {
		ctx = withCloudAccessKey(ctx, tcr)
	}
	deployments, err := cfg.VMC().ListDeployments(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list deployments: %v", err)), nil
	}
	data, err := json.Marshal(deployments)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal deployments: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func RegisterToolDeployments(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameDeployments) {
		return
	}
	if !c.IsCloud() {
		return
	}
	s.AddTool(toolDeployments(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolDeploymentsHandler(ctx, c, request)
	})
}
