package logs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

const toolNameFlags = "flags"

var (
	toolFlags = mcp.NewTool(toolNameFlags,
		mcp.WithDescription("List of non-default flags (parameters) of the VictoriaLogs instance. This tools uses `/flags` endpoint of VictoriaLogs API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of non-default flags (parameters)",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	)
)

func toolFlagsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if cfg.IsCloud() {
		deploymentID, err := tools.GetCloudDeploymentID(tcr)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if cfg.IsCloudSharedInstance() {
			ctx = tools.WithCloudAccessKey(ctx, tcr)
		}
		dd, err := cfg.VMC().GetDeploymentDetails(ctx, deploymentID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to get deployment details: %v", err)), nil
		}
		if dd.Type != vmcloud.DeploymentTypeVLogs {
			return mcp.NewToolResultError(fmt.Sprintf("unsupported deployment type %s for deployment %s", dd.Type, deploymentID)), nil
		}
		data, err := json.Marshal(map[string]any{"victorialogs": dd.VMSingleSettings})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to marshal deployment details: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	}

	req, err := CreateAdminRequest(ctx, cfg, tcr, "flags")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	return tools.GetTextBodyForRequest(req, cfg), nil
}

func registerToolFlags(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameFlags) {
		return
	}
	s.AddTool(tools.WithCloudDeploymentID(c, toolFlags), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolFlagsHandler(ctx, c, request)
	})
}
