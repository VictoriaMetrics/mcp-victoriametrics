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

func toolAccessTokens(_ *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of deployment access tokens in VictoriaMetrics Cloud"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of deployment access tokens",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = append(
		options,
		mcp.WithString("deployment_id",
			mcp.Required(),
			mcp.Title("Deployment ID"),
			mcp.Description("Unique identifier of the deployment in VictoriaMetrics Cloud"),
			mcp.Pattern(`^[a-zA-Z0-9\-_]+$`),
		),
	)
	return mcp.NewTool(toolNameAccessTokens, append(options, withEnvironmentParam())...)
}

func toolAccessTokensHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	envName, _ := GetToolReqParam[string](tcr, "env", false)
	env, err := cfg.Environment(envName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get deployment_id parameter: %v", err)), nil
	}
	if deploymentID == "" {
		return mcp.NewToolResultError("deployment_id parameter is required for cloud mode"), nil
	}
	accessTokens, err := env.VMC().ListDeploymentAccessTokens(ctx, deploymentID)
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
	if !c.IsCloud() {
		return
	}
	s.AddTool(toolAccessTokens(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolAccessTokensHandler(ctx, c, request)
	})
}
