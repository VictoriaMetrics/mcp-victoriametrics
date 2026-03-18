package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameCloudProviders = "cloud_providers"

func toolCloudProviders(_ *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of cloud providers in VictoriaMetrics Cloud"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of cloud providers",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	return mcp.NewTool(toolNameCloudProviders, options...)
}

func toolCloudProvidersHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if cfg.IsCloudSharedInstance() {
		ctx = withCloudAccessKey(ctx, tcr)
	}
	cloudProviders, err := cfg.VMC().ListCloudProviders(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list cloud providers: %v", err)), nil
	}
	data, err := json.Marshal(cloudProviders)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal cloud providers: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func RegisterToolCloudProviders(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameCloudProviders) {
		return
	}
	if !c.IsCloud() {
		return
	}
	s.AddTool(toolCloudProviders(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolCloudProvidersHandler(ctx, c, request)
	})
}
