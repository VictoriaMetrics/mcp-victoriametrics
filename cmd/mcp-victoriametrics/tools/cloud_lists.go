package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const (
	toolNameCloudProviders = "cloud_providers"
	toolNameDeployments    = "deployments"
	toolNameRegions        = "regions"
	toolNameTiers          = "tiers"
)

func newCloudListTool(name, description, title string, c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription(description),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           title,
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withCloudToolTargetingOptions(options, c, false)
	return mcp.NewTool(name, options...)
}

func handleCloudListTool(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, list func(context.Context, *config.Instance) (any, error), noun string) (*mcp.CallToolResult, error) {
	instance, err := getCloudToolInstance(cfg, tcr)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	items, err := list(ctx, instance)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list %s: %v", noun, err)), nil
	}
	data, err := json.Marshal(items)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal %s: %v", noun, err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func registerCloudListTool(s *server.MCPServer, c *config.Config, name, description, title, noun string, list func(context.Context, *config.Instance) (any, error)) {
	if c.IsToolDisabled(name) {
		return
	}
	s.AddTool(newCloudListTool(name, description, title, c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleCloudListTool(ctx, c, request, list, noun)
	})
}

func RegisterToolCloudProviders(s *server.MCPServer, c *config.Config) {
	registerCloudListTool(s, c, toolNameCloudProviders, "List of cloud providers in VictoriaMetrics Cloud", "List of cloud providers", "cloud providers", func(ctx context.Context, instance *config.Instance) (any, error) {
		return instance.VMC().ListCloudProviders(ctx)
	})
}

func RegisterToolDeployments(s *server.MCPServer, c *config.Config) {
	registerCloudListTool(s, c, toolNameDeployments, "List of deployments in VictoriaMetrics Cloud", "List of deployments", "deployments", func(ctx context.Context, instance *config.Instance) (any, error) {
		return instance.VMC().ListDeployments(ctx)
	})
}

func RegisterToolRegions(s *server.MCPServer, c *config.Config) {
	registerCloudListTool(s, c, toolNameRegions, "List of regions in VictoriaMetrics Cloud", "List of regions", "regions", func(ctx context.Context, instance *config.Instance) (any, error) {
		return instance.VMC().ListRegions(ctx)
	})
}

func RegisterToolTiers(s *server.MCPServer, c *config.Config) {
	registerCloudListTool(s, c, toolNameTiers, "List of tiers in VictoriaMetrics Cloud", "List of tiers", "tiers", func(ctx context.Context, instance *config.Instance) (any, error) {
		return instance.VMC().ListTiers(ctx)
	})
}
