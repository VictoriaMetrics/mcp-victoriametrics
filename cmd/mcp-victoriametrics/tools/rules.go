package tools

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameRules = "rules"

func toolRules(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of alerting and recording rules of VictoriaMetrics instance. This tool uses `/api/v1/rules` endpoint of vmalert API proxied by VictoriaMetrics API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of alerting and recording rules",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withTargetingOptions(options, c, true, true)
	options = append(
		options,
		mcp.WithString("type",
			mcp.Title("Rules type"),
			mcp.Description("Rules type to be displayed: alert or record"),
			mcp.DefaultString(""),
			mcp.Enum("alert", "record"),
		),
		mcp.WithString("filter",
			mcp.Title("Extra filter for rules"),
			mcp.Description("Extra filter for rules with possible problems: unhealthy (rules that get some errors during evaluation) or noMatch (rules that don't match any time series)"),
			mcp.DefaultString(""),
			mcp.Enum("unhealthy", "noMatch"),
		),
		mcp.WithBoolean("exclude_alerts",
			mcp.Title("Exclude alerts"),
			mcp.Description("Exclude alerts from the list"),
			mcp.DefaultBool(false),
		),
		mcp.WithArray("rule_names",
			mcp.Title("Rule names"),
			mcp.Description("Filter rules by name"),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithArray("rule_groups",
			mcp.Title("Rule groups"),
			mcp.Description("Filter rules by group names"),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithArray("rule_files",
			mcp.Title("Rule files"),
			mcp.Description("Filter rules by file names"),
			mcp.Items(map[string]any{"type": "string"}),
		),
	)
	return mcp.NewTool(toolNameRules, options...)
}

func toolRulesHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ruleType, err := GetToolReqParam[string](tcr, "type", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	filter, err := GetToolReqParam[string](tcr, "filter", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	excludeAlerts, err := GetToolReqParam[bool](tcr, "exclude_alerts", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	ruleNames, err := GetToolReqParam[[]any](tcr, "rule_names", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	ruleGroups, err := GetToolReqParam[[]any](tcr, "rule_groups", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	ruleFiles, err := GetToolReqParam[[]any](tcr, "rule_files", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "vmalert", "api", "v1", "rules")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	q := req.URL.Query()
	if ruleType != "" {
		q.Add("type", ruleType)
	}
	if filter != "" {
		q.Add("filter", filter)
	}
	if excludeAlerts {
		q.Add("exclude_alerts", "true")
	}

	for _, ruleName := range ruleNames {
		ruleNameStr, ok := ruleName.(string)
		if !ok {
			return mcp.NewToolResultError("rule_names element must be a string"), nil
		}
		q.Add("rule_name[]", ruleNameStr)
	}

	for _, ruleGroup := range ruleGroups {
		ruleGroupStr, ok := ruleGroup.(string)
		if !ok {
			return mcp.NewToolResultError("rule_groups element must be a string"), nil
		}
		q.Add("rule_group[]", ruleGroupStr)
	}

	for _, ruleFile := range ruleFiles {
		ruleFileStr, ok := ruleFile.(string)
		if !ok {
			return mcp.NewToolResultError("rule_files element must be a string"), nil
		}
		q.Add("file[]", ruleFileStr)
	}

	req.URL.RawQuery = q.Encode()

	return GetTextBodyForRequest(req, cfg), nil
}

func RegisterToolRules(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameRules) {
		return
	}
	s.AddTool(toolRules(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolRulesHandler(ctx, c, request)
	})
}
