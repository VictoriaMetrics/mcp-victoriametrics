package traces

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

const toolNameDependencies = "dependencies"

var (
	toolDependencies = mcp.NewTool(toolNameDependencies,
		mcp.WithDescription("Query the service dependency graph. This tools uses `/select/jaeger/api/dependencies` endpoint of VictoriaTraces API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Service dependency graph",
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
		mcp.WithNumber("endTs",
			mcp.Title("The end of the time interval"),
			mcp.Description("The end timestamp in unix milliseconds. Current timestamp will be used if empty."),
			mcp.DefaultNumber(0),
		),
		mcp.WithNumber("lookback",
			mcp.Title("The length the time interval in milliseconds"),
			mcp.Description("the lookbehind window duration in milliseconds (i.e. start-time + lookback = endTs). Default to 1h if empty."),
			mcp.DefaultNumber(3600000),
		),
	)
)

func toolDependenciesHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	endTs, err := tools.GetToolReqParam[float64](tcr, "endTs", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	lookback, err := tools.GetToolReqParam[float64](tcr, "lookback", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "dependencies")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	if endTs > 0 || lookback > 0 {
		q := req.URL.Query()
		if endTs > 0 {
			q.Add("endTs", fmt.Sprintf("%d", uint64(endTs)))
		}
		if lookback > 0 {
			q.Add("lookback", fmt.Sprintf("%d", uint64(lookback)))
		}
		req.URL.RawQuery = q.Encode()
	}

	return tools.GetTextBodyForRequest(req, cfg), nil
}

func registerToolDependencies(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameDependencies) {
		return
	}
	s.AddTool(tools.WithCloudDeploymentID(c, toolDependencies), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolDependenciesHandler(ctx, c, request)
	})
}
