package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameRuleFile = "rule_file"

func toolRuleFile(_ *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("Get content of deployment alerting and recording rules file in VictoriaMetrics Cloud"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Get rules file content",
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
	options = append(
		options,
		mcp.WithString("filename",
			mcp.Required(),
			mcp.Title("Rules filename"),
			mcp.Description("Name of the rules file to retrieve. This should be one of the files listed by the `rule_filenames` tool."),
			mcp.Min(1),
			mcp.Max(255),
		),
	)
	return mcp.NewTool(toolNameRuleFile, options...)
}

func toolRuleFileHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get deployment_id parameter: %v", err)), nil
	}
	if deploymentID == "" {
		return mcp.NewToolResultError("deployment_id parameter is required for cloud mode"), nil
	}

	filename, err := GetToolReqParam[string](tcr, "filename", true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get rules_filename parameter: %v", err)), nil
	}

	if cfg.IsCloudSharedInstance() {
		ctx = withCloudAccessKey(ctx, tcr)
	}
	ruleFilenames, err := cfg.VMC().GetDeploymentRuleFileContent(ctx, deploymentID, filename)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list of rule filenames: %v", err)), nil
	}
	data, err := json.Marshal(ruleFilenames)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal rule filenames: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func RegisterToolRuleFile(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameRuleFile) {
		return
	}
	if !c.IsCloud() {
		return
	}
	s.AddTool(toolRuleFile(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolRuleFileHandler(ctx, c, request)
	})
}
