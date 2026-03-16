package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameAccessTokens = "access_tokens"

func toolAccessTokens(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of deployment access tokens in VictoriaMetrics Cloud"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of deployment access tokens",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withTargetingOptions(options, c, true, false)
	return mcp.NewTool(toolNameAccessTokens, options...)
}

func toolAccessTokensHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance, err := getCloudToolInstance(cfg, tcr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	deploymentID, err := requireCloudDeploymentID(instance, tcr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	accessTokens, err := instance.VMC().ListDeploymentAccessTokens(ctx, deploymentID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list access_tokens: %v", err)), nil
	}
	data, err := json.Marshal(accessTokens)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal access_tokens: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func RegisterToolAccessTokens(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameAccessTokens) {
		return
	}
	s.AddTool(toolAccessTokens(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolAccessTokensHandler(ctx, c, request)
	})
}
