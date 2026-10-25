package logs

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

const toolNameStreamFieldNames = "stream_field_names"

var (
	toolStreamFieldNames = mcp.NewTool(toolNameStreamFieldNames,
		mcp.WithDescription("Get log stream field names from results of the given <query> on the given [<start> ... <end>] time range. The response also contains the number of log results per every field name.. This tool uses `/select/logsql/stream_field_names` endpoint of VictoriaLogs API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of stream fields",
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
		mcp.WithString("query",
			mcp.Required(),
			mcp.Title("LogsQL expression"),
			mcp.Description(`LogsQL expression for the query of the logs streams`),
		),
		mcp.WithString("start",
			mcp.Required(),
			mcp.Title("Start timestamp"),
			mcp.Description("Start timestamp in RFC3339"),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithString("end",
			mcp.Title("End timestamp"),
			mcp.Description("End timestamp in RFC3339. If <end> is missing, then it equals to the maximum timestamp across logs stored in VictoriaLogs."),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
	)
)

func toolStreamFieldNamesHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := tools.GetToolReqParam[string](tcr, "query", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	start, err := tools.GetToolReqParam[string](tcr, "start", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	end, err := tools.GetToolReqParam[string](tcr, "end", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "stream_field_names")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("start", start)
	if end != "" {
		q.Add("end", end)
	}
	req.URL.RawQuery = q.Encode()

	return tools.GetTextBodyForRequest(req, cfg), nil
}

func registerToolStreamFieldNames(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameStreamFieldNames) {
		return
	}
	s.AddTool(tools.WithCloudDeploymentID(c, toolStreamFieldNames), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolStreamFieldNamesHandler(ctx, c, request)
	})
}
