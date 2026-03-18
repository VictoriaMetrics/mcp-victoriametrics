package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameRuleFilenames = "rule_filenames"

func toolRuleFilenames(_ *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of deployment alerting and recording rules filenames in VictoriaMetrics Cloud"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of filenames of alerting and recording rules",
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
	return mcp.NewTool(toolNameRuleFilenames, options...)
}

func toolRuleFilenamesHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get deployment_id parameter: %v", err)), nil
	}
	if deploymentID == "" {
		return mcp.NewToolResultError("deployment_id parameter is required for cloud mode"), nil
	}
	if cfg.IsCloudSharedInstance() {
		ctx = withCloudAccessKey(ctx, tcr)
	}
	ruleFilenames, err := cfg.VMC().ListDeploymentRuleFileNames(ctx, deploymentID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list of rule filenames: %v", err)), nil
	}
	data, err := json.Marshal(ruleFilenames)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal rule filenames: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func RegisterToolRuleFilenames(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameRuleFilenames) {
		return
	}
	if !c.IsCloud() {
		return
	}
	s.AddTool(toolRuleFilenames(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolRuleFilenamesHandler(ctx, c, request)
	})
}
