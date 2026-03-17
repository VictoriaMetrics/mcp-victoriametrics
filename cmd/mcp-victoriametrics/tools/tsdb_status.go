package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameTSDBStatus = "tsdb_status"

func toolTSDBStatus(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription(`The following tool returns various cardinality statistics about the VictoriaMetrics TSDB:
- Metric names with the highest number of series.
- Labels with the highest number of series.
- Values with the highest number of series for the selected label (aka focusLabel).
- label=name pairs with the highest number of series.
- Labels with the highest number of unique values.

This tool returns TSDB stats from "/api/v1/status/tsdb" endpoint of VictoriaMetrics API (in the way similar to Prometheus).
`),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "TSDB Stats (information about cardinality of the data)",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withTargetingOptions(options, c, true, true)
	options = append(
		options,
		mcp.WithNumber("topN",
			mcp.Title("Top N"),
			mcp.Description("is the number of top entries to return in the response"),
			mcp.DefaultNumber(10),
			mcp.Min(1),
		),
		mcp.WithString("focusLabel",
			mcp.Title("Focus label"),
			mcp.Description("Returns label values with the highest number of time series for the given label name in the seriesCountByFocusLabelValue list."),
			mcp.DefaultString(""),
		),
		mcp.WithString("date",
			mcp.Title("Date"),
			mcp.Description("The date for collecting the stats. By default, the stats is collected for the current day. Pass date=1970-01-01 in order to collect global stats across all the days."),
			mcp.DefaultString(""),
			mcp.Pattern(`^\d{4}-\d{2}-\d{2}$`),
		),
		mcp.WithString("match",
			mcp.Title("Match series selector"),
			mcp.Description("Arbitrary time series selector for series to take into account during stats calculation. By default all the series are taken into account."),
			mcp.DefaultString(""),
		),
		mcp.WithString("extraLabel",
			mcp.Title("Extra label"),
			mcp.Description(`Optional extra_label=<label_name>=<label_value> query arg, which can be used for enforcing additional label filters for queries. For example, /api/v1/query_range?extra_label=user_id=123&extra_label=group_id=456&query=<query> would automatically add {user_id="123",group_id="456"} label filters to the given <query>. This functionality can be used for limiting the scope of time series visible to the given tenant.`),
			mcp.DefaultString(""),
		),
	)
	return mcp.NewTool(toolNameTSDBStatus, options...)
}

func toolTSDBStatusHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topN, err := GetToolReqParam[float64](tcr, "topN", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if topN < 1 {
		topN = 10
	}

	focusLabel, err := GetToolReqParam[string](tcr, "focusLabel", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	date, err := GetToolReqParam[string](tcr, "date", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	match, err := GetToolReqParam[string](tcr, "match", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	extraLabel, err := GetToolReqParam[string](tcr, "extraLabel", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "api", "v1", "status", "tsdb")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	query := req.URL.Query()
	query.Set("topN", fmt.Sprintf("%d", int(topN)))
	if focusLabel != "" {
		query.Set("focusLabel", focusLabel)
	}
	if date != "" {
		query.Set("date", date)
	}
	if match != "" {
		query.Set("match[]", match)
	}
	if extraLabel != "" {
		query.Set("extra_label", extraLabel)
	}
	req.URL.RawQuery = query.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolTSDBStatus(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameTSDBStatus) {
		return
	}
	s.AddTool(toolTSDBStatus(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolTSDBStatusHandler(ctx, c, request)
	})
}
