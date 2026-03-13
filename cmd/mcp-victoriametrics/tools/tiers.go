package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameTiers = "tiers"

func toolTiers(_ *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of tiers in VictoriaMetrics Cloud"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of tiers",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	return mcp.NewTool(toolNameTiers, append(options, withEnvironmentParam())...)
}

func toolTiersHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	envName, _ := GetToolReqParam[string](tcr, "env", false)
	env, err := cfg.Environment(envName)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tiers, err := env.VMC().ListTiers(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list tiers: %v", err)), nil
	}
	data, err := json.Marshal(tiers)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal tiers: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func RegisterToolTiers(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameTiers) {
		return
	}
	if !c.IsCloud() {
		return
	}
	s.AddTool(toolTiers(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolTiersHandler(ctx, c, request)
	})
}
