package traces

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

const toolNameServiceOperations = "service_operations"

var (
	toolServiceOperations = mcp.NewTool(toolNameServiceOperations,
		mcp.WithDescription("Get all the span names (operations) of a service. This tool uses `/select/jaeger/api/services/{service_name}/operations` endpoint of VictoriaTraces API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of span names (operations) of a service",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
		mcp.WithString("tenant",
			mcp.Title("Tenant name (Account ID and Project ID)"),
			mcp.Description("Name of the tenant for which the data will be displayed (format AccountID:ProjectID)"),
			mcp.DefaultString("0:0"),
			mcp.Pattern(`^([0-9]+)(:[0-9]+)$`),
		),
		mcp.WithString("service_name",
			mcp.Required(),
			mcp.Title("Service name"),
			mcp.Description("Service name"),
		),
	)
)

func toolServiceOperationsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	serviceName, err := tools.GetToolReqParam[string](tcr, "service_name", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "services", serviceName, "operations")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	return tools.GetTextBodyForRequest(req, cfg), nil
}

func registerToolServiceOperations(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameServiceOperations) {
		return
	}
	s.AddTool(tools.WithCloudDeploymentID(c, toolServiceOperations), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolServiceOperationsHandler(ctx, c, request)
	})
}
