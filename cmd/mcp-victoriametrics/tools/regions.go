package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameRegions = "regions"

func toolRegions(_ *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of regions in VictoriaMetrics Cloud"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of regions",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	return mcp.NewTool(toolNameRegions, options...)
}

func toolRegionsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if cfg.IsCloudSharedInstance() {
		ctx = withCloudAccessKey(ctx, tcr)
	}
	regions, err := cfg.VMC().ListRegions(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list regions: %v", err)), nil
	}
	data, err := json.Marshal(regions)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal regions: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func RegisterToolRegions(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameRegions) {
		return
	}
	if !c.IsCloud() {
		return
	}
	s.AddTool(toolRegions(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolRegionsHandler(ctx, c, request)
	})
}
