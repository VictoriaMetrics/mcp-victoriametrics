package traces

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

const toolNameTrace = "trace"

var (
	toolTrace = mcp.NewTool(toolNameTrace,
		mcp.WithDescription("Get trace info by trace ID. This tool uses `/select/jaeger/api/traces/{trace_id}` endpoint of VictoriaTraces API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Trace information",
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
		mcp.WithString("trace_id",
			mcp.Required(),
			mcp.Title("Trace ID"),
			mcp.Description("Trace ID"),
		),
	)
)

func toolTraceHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	traceID, err := tools.GetToolReqParam[string](tcr, "trace_id", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "traces", traceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	return tools.GetTextBodyForRequest(req, cfg), nil
}

func registerToolTrace(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameTrace) {
		return
	}
	s.AddTool(tools.WithCloudDeploymentID(c, toolTrace), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolTraceHandler(ctx, c, request)
	})
}
