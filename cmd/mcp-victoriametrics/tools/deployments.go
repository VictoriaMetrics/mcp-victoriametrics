package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
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
	return mcp.NewTool(toolNameDeployments, append(options, withEnvironmentParam())...)
}

func toolDeploymentsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	envName, _ := GetToolReqParam[string](tcr, "env", false)
	env, err := cfg.Environment(envName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	deployments, err := env.VMC().ListDeployments(ctx)
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
