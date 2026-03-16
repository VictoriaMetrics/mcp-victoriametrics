package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameRuleFilenames = "rule_filenames"

func toolRuleFilenames(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of deployment alerting and recording rules filenames in VictoriaMetrics Cloud"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of filenames of alerting and recording rules",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = append(options, maybeWithEnvironmentParam(c)...)
	options = append(
		options,
		maybeWithDeploymentIDParam(c)...,
	)
	return mcp.NewTool(toolNameRuleFilenames, options...)
}

func toolRuleFilenamesHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	instance, err := getCloudToolInstance(cfg, tcr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get deployment_id parameter: %v", err)), nil
	}
	if deploymentID == "" {
		return mcp.NewToolResultError(fmt.Sprintf("deployment_id parameter is required for cloud env %q", instance.Name())), nil
	}
	ruleFilenames, err := instance.VMC().ListDeploymentRuleFileNames(ctx, deploymentID)
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
	if !c.HasCloudInstances() {
		return
	}
	s.AddTool(toolRuleFilenames(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolRuleFilenamesHandler(ctx, c, request)
	})
}
