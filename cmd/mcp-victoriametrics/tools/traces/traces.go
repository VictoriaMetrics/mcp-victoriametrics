package traces

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

const toolNameTraces = "traces"

var (
	toolTraces = mcp.NewTool(toolNameTraces,
		mcp.WithDescription("Query traces. This tool uses `/select/jaeger/api/traces` endpoint of VictoriaTraces API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Query traces",
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
		mcp.WithString("service",
			mcp.Required(),
			mcp.Title("Service name"),
			mcp.Description("Service name"),
		),
		mcp.WithString("operation",
			mcp.Title("Span name (operation)"),
			mcp.Description("The span name (also known as the operation name in Jaeger)"),
			mcp.DefaultString(""),
		),
		mcp.WithNumber("start",
			mcp.Title("Start timestamp"),
			mcp.Description("Start timestamp in unix milliseconds."),
		),
		mcp.WithNumber("end",
			mcp.Title("End timestamp"),
			mcp.Description("End timestamp in unix milliseconds."),
		),
		mcp.WithString("minDuration",
			mcp.Title("Minimum duration"),
			mcp.Description("The minimum duration of the span, with units: ns, us, ms, s, m, or h."),
			mcp.Pattern(`^([0-9]+)(ns|us|ms|s|m|h)$`),
		),
		mcp.WithString("maxDuration",
			mcp.Title("Maximum duration"),
			mcp.Description("The maximum duration of the span, with units: ns, us, ms, s, m, or h."),
			mcp.Pattern(`^([0-9]+)(ns|us|ms|s|m|h)$`),
		),
		mcp.WithNumber("limit",
			mcp.Required(),
			mcp.Title("Trace limit"),
			mcp.Description("The maximum number of traces in query results, default 20."),
			mcp.DefaultNumber(20),
		),
	)
)

func toolTracesHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	service, err := tools.GetToolReqParam[string](tcr, "service", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	operation, err := tools.GetToolReqParam[string](tcr, "operation", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start, err := tools.GetToolReqParam[float64](tcr, "start", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	end, err := tools.GetToolReqParam[float64](tcr, "end", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	minDuration, err := tools.GetToolReqParam[string](tcr, "minDuration", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	maxDuration, err := tools.GetToolReqParam[string](tcr, "maxDuration", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit, err := tools.GetToolReqParam[float64](tcr, "limit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if limit == 0 {
		limit = 20
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "traces")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	q := req.URL.Query()
	q.Add("service", service)
	if operation != "" {
		q.Add("operation", operation)
	}
	if start > 0 {
		q.Add("start", fmt.Sprintf("%d", uint64(start*1000)))
	}
	if end > 0 {
		q.Add("end", fmt.Sprintf("%d", uint64(end*1000)))
	}
	if minDuration != "" {
		q.Add("minDuration", minDuration)
	}
	if maxDuration != "" {
		q.Add("maxDuration", maxDuration)
	}
	q.Add("limit", fmt.Sprintf("%d", uint64(limit)))
	req.URL.RawQuery = q.Encode()

	return tools.GetTextBodyForRequest(req, cfg), nil
}

func registerToolTraces(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameTraces) {
		return
	}
	s.AddTool(tools.WithCloudDeploymentID(c, toolTraces), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolTracesHandler(ctx, c, request)
	})
}
