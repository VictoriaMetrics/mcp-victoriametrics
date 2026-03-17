package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameAlerts = "alerts"

func toolAlerts(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("List of firing and pending alerts of the VictoriaMetrics instance. This tool uses `/api/v1/alerts` endpoint of vmalert API, proxied by VictoriaMetrics API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of alerts",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withTargetingOptions(options, c, true, true)
	options = append(
		options,
		mcp.WithString("state",
			mcp.Title("Filter by alert state"),
			mcp.Description("Filter alerts by their state. Possible values: 'firing', 'pending', 'all'. Default is 'all'."),
			mcp.DefaultString("all"),
			mcp.Enum("firing", "pending", "all"),
		),
		mcp.WithString("group",
			mcp.Title("Filter by alert group"),
			mcp.Description("Filter alerts by their group name. If not specified, all groups are included."),
			mcp.DefaultString(""),
		),
		mcp.WithNumber("limit",
			mcp.Title("Limit the number of alerts"),
			mcp.Description("Maximum number of alerts to return. If not specified, all alerts are returned."),
			mcp.DefaultNumber(0),
		),
		mcp.WithNumber("offset",
			mcp.Title("Offset for pagination"),
			mcp.Description("Number of alerts to skip before starting to collect the result set. Default is 0."),
			mcp.DefaultNumber(0),
		),
	)
	return mcp.NewTool(toolNameAlerts, options...)
}

func toolAlertsHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	state, err := GetToolReqParam[string](tcr, "state", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	state = strings.ToLower(state)
	if state == "" {
		state = "all" // Default value if not specified
	}
	if state != "firing" && state != "pending" && state != "all" {
		return mcp.NewToolResultError("invalid state parameter, must be 'firing', 'pending', or 'all'"), nil
	}
	group, err := GetToolReqParam[string](tcr, "group", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit, err := GetToolReqParam[float64](tcr, "limit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	offset, err := GetToolReqParam[float64](tcr, "offset", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "vmalert", "api", "v1", "alerts")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}

	processResult := func(s string) (string, error) {
		if state == "all" && group == "" {
			return s, nil // No filtering needed
		}

		result := make(map[string]any)
		if err = json.Unmarshal([]byte(s), &result); err != nil {
			return "", fmt.Errorf("failed to unmarshal response: %v", err)
		}

		_data, ok := result["data"]
		if !ok {
			return "", fmt.Errorf("unexpected response format, 'data' field not found")
		}
		data, ok := _data.(map[string]any)
		if !ok {
			return "", fmt.Errorf("unexpected response format, 'data' field not found or not an object")
		}
		_alerts, ok := data["alerts"]
		if !ok {
			return "", fmt.Errorf("unexpected response format, 'alerts' field not found")
		}
		alerts, ok := _alerts.([]any)
		if !ok {
			return "", fmt.Errorf("unexpected response format, 'alerts' field not found or not an array")
		}
		filteredAlerts := make([]map[string]any, 0, len(alerts))
		for _, _alert := range alerts {
			alert, ok := _alert.(map[string]any)
			if !ok {
				return "", fmt.Errorf("unexpected response format, 'alerts' field has invalid element")
			}
			if state != "all" && alert["state"] != state {
				continue // Skip alerts that do not match the state filter
			}
			_labels, ok := alert["labels"]
			if !ok {
				continue // Skip alerts without labels
			}
			labels, ok := _labels.(map[string]any)
			if !ok {
				continue // Skip alerts where labels are not an object
			}
			if group != "" && labels["alertgroup"] != group {
				continue // Skip alerts that do not match the group filter
			}
			filteredAlerts = append(filteredAlerts, alert)
		}
		slices.SortFunc(filteredAlerts, func(a, b map[string]any) int {
			aID, _ := a["id"].(string)
			bID, _ := b["id"].(string)
			return strings.Compare(aID, bID)
		})
		if limit > 0 {
			filteredAlerts = filteredAlerts[int(offset):int(offset+limit)]
		}
		data["alerts"] = filteredAlerts // Update the alerts field with filtered alerts
		result["data"] = data           // Update the data field with the modified alerts

		b, err := json.Marshal(result)
		if err != nil {
			return "", fmt.Errorf("failed to marshal filtered response: %v", err)
		}

		return string(b), nil
	}

	return GetTextBodyForRequest(req, cfg, processResult), nil
}

func RegisterToolAlerts(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameAlerts) {
		return
	}
	s.AddTool(toolAlerts(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolAlertsHandler(ctx, c, request)
	})
}
