package logs

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

const toolNameFacets = "facets"

var (
	toolFacets = mcp.NewTool(toolNameFacets,
		mcp.WithDescription("The most frequent values per each log field seen in the logs returned by the given <query> on the given [<start> ... <end>] time range. This tool uses `/select/logsql/facets` endpoint of VictoriaLogs API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Most frequent values",
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
			mcp.Required(),
			mcp.Title("End timestamp"),
			mcp.Description("End timestamp in RFC3339 format. For example, 2023-10-01T00:00:00Z"),
			mcp.Pattern(`^((?:(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}(?:\.\d+)?))(Z|[\+-]\d{2}:\d{2})?)|([0-9]+)$`),
		),
		mcp.WithNumber("limit",
			mcp.Title("Limit"),
			mcp.Description("The number of values per each log field can be controlled via limit arg."),
		),
		mcp.WithNumber("max_values_per_field",
			mcp.Title("Max values per field"),
			mcp.Description("The facets tool ignores log fields, which contain too big number of unique values, since they can consume a lot of RAM to track. The limit on the number of unique values per each log field can be controlled via max_values_per_field arg."),
		),
		mcp.WithNumber("max_value_len",
			mcp.Title("Max value length"),
			mcp.Description("The facets tool ignores log fields, which contain too long values. The limit on the per-field value length can be controlled via max_value_len arg."),
		),
		mcp.WithBoolean("keep_const_fields",
			mcp.Title("Keep constant fields"),
			mcp.Description("By default the facets tool doesn’t return log fields, which contain the same constant value across all the logs matching the given query. Add keep_const_fields=true arg if you need such log fields"),
		),
	)
)

func toolFacetsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	limit, err := tools.GetToolReqParam[float64](tcr, "limit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	maxValuesPerField, err := tools.GetToolReqParam[float64](tcr, "max_values_per_field", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	maxValueLen, err := tools.GetToolReqParam[float64](tcr, "max_value_len", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	keepConstFields, err := tools.GetToolReqParam[bool](tcr, "keep_const_fields", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "facets")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("start", start)
	if end != "" {
		q.Add("end", end)
	}
	if limit != 0 {
		q.Add("limit", fmt.Sprintf("%.f", limit))
	}
	if maxValuesPerField != 0 {
		q.Add("max_values_per_field", fmt.Sprintf("%.f", maxValuesPerField))
	}
	if maxValueLen != 0 {
		q.Add("max_value_len", fmt.Sprintf("%.f", maxValueLen))
	}
	if keepConstFields {
		q.Add("keep_const_fields", "1")
	}
	req.URL.RawQuery = q.Encode()

	return tools.GetTextBodyForRequest(req, cfg), nil
}

func registerToolFacets(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameFacets) {
		return
	}
	s.AddTool(tools.WithCloudDeploymentID(c, toolFacets), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolFacetsHandler(ctx, c, request)
	})
}
