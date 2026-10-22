package logs

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

const toolNameHits = "hits"

var (
	toolHits = mcp.NewTool(toolNameHits,
		mcp.WithDescription("The number of matching log entries for the given <query> on the given [<start> ... <end>] time range grouped by <step> buckets. The returned results are sorted by time. This tool uses `/select/logsql/hits` endpoint of VictoriaLogs API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Log entries hits",
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
			mcp.Description(`LogsQL expression for the query of the logs data`),
		),
		mcp.WithString("start",
			mcp.Required(),
			mcp.Title("Start timestamp"),
			mcp.Description("Start timestamp in RFC3339 format. For example, 2023-10-01T00:00:00Z"),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithString("end",
			mcp.Title("End timestamp"),
			mcp.Description("End timestamp in RFC3339 format. For example, 2023-10-01T00:00:00Z"),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithString("step",
			mcp.Title("Step"),
			mcp.Description("The step is used to group the log entries by time. For example, 60s means that the log entries counts will be grouped by 1 minute bucket. Default is 1d."),
			mcp.Pattern(`^([0-9]+)([a-z]+)$`),
			mcp.DefaultString("1d"),
		),
		mcp.WithArray("field",
			mcp.Title("Field"),
			mcp.Description("Additionally, any number of field=<field_name> args can be passed to /select/logsql/hits for grouping hits buckets by the mentioned <field_name> fields. The grouped fields are put inside 'fields' object of response."),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithNumber("fields_limit",
			mcp.Title("fields_limit"),
			mcp.Description("Optional fields_limit=N query arg can be passed for limiting the number of unique 'fields' groups to return to N. If more than N unique 'fields' groups is found, then top N 'fields' groups with the maximum number of 'total' hits are returned. The remaining hits are returned in 'fields': {} group."),
		),
	)
)

func toolHitsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	step, err := tools.GetToolReqParam[string](tcr, "step", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Use []any because JSON unmarshals arrays as []interface{}, not []string.
	// Then convert to []string with validation.
	fieldsAny, err := tools.GetToolReqParam[[]any](tcr, "field", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	var fields []string
	for i, f := range fieldsAny {
		s, ok := f.(string)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("field[%d] must be a string, got %T", i, f)), nil
		}
		fields = append(fields, s)
	}

	fieldsLimit, err := tools.GetToolReqParam[float64](tcr, "fields_limit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "hits")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("start", start)
	if end != "" {
		q.Add("end", end)
	}
	if step != "" {
		q.Add("step", step)
	}
	if len(fields) > 0 {
		for _, field := range fields {
			if field != "" {
				q.Add("field", field)
			}
		}
		if fieldsLimit > 0 {
			q.Add("fields_limit", fmt.Sprintf("%.f", fieldsLimit))
		}
	}
	req.URL.RawQuery = q.Encode()

	return tools.GetTextBodyForRequest(req, cfg), nil
}

func registerToolHits(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameHits) {
		return
	}
	s.AddTool(tools.WithCloudDeploymentID(c, toolHits), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolHitsHandler(ctx, c, request)
	})
}
