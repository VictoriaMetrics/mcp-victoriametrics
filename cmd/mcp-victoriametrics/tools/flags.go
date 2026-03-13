package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameFlags = "flags"

func toolFlags(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of non-default flags (parameters) of the VictoriaMetrics instance. This tools uses `/flags` endpoint of VictoriaMetrics API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of non-default flags (parameters)",
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
	return mcp.NewTool(toolNameFlags, append(options, withEnvironmentParam())...)
}

func toolFlagsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	envName, _ := GetToolReqParam[string](tcr, "env", false)
	env, err := cfg.Environment(envName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	if env.IsCloud() {
		deploymentID, err := GetToolReqParam[string](tcr, "deployment_id", true)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get deployment_id parameter: %v", err)), nil
		}
		if deploymentID == "" {
			return mcp.NewToolResultError("deployment_id parameter is required for cloud mode"), nil
		}
		dd, err := env.VMC().GetDeploymentDetails(ctx, deploymentID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get deployment details: %v", err)), nil
		}
		result := map[string]any{}
		switch dd.Type {
		case vmcloud.DeploymentTypeSingleNode:
			result["vmsingle"] = dd.VMSingleSettings
		case vmcloud.DeploymentTypeCluster:
			result["vmselect"] = dd.VMSelectSettings
			result["vmstorage"] = dd.VMStorageSettings
			result["vminsert"] = dd.VMInsertSettings
		}
		data, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal deployment details: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	}

	req, err := CreateAdminRequest(ctx, cfg, tcr, "flags")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolFlags(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameFlags) {
		return
	}
	s.AddTool(toolFlags(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolFlagsHandler(ctx, c, request)
	})
}
