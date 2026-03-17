package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameMetricsMetadata = "metrics_metadata"

func toolMetricsMetadata(c *config.Config) mcp.Tool {
	options := []mcp.ToolOption{
		mcp.WithDescription("Returns stored metrics details (metadata) such as name, type, help (description), and unit, and can be used to generate natural language queries. This tool uses `/api/v1/metadata` endpoint of VictoriaMetrics API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of metrics metadata (name, type, help and unit)",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	}
	options = withTargetingOptions(options, c, true, true)
	options = append(
		options,
		mcp.WithString("search",
			mcp.Title("Search keyword"),
			mcp.Description("A keyword for search in description and metric name"),
			mcp.DefaultString(""),
		),
		mcp.WithString("type",
			mcp.Title("Metric type name"),
			mcp.Description("The metric type by which the list of metrics will be filtered"),
			mcp.DefaultString(""),
		),
		mcp.WithString("unit",
			mcp.Title("Unit name"),
			mcp.Description("The unit by which the list of metrics will be filtered"),
			mcp.DefaultString(""),
		),
		mcp.WithString("metric",
			mcp.Title("Metric name"),
			mcp.Description("A metric name to filter metadata for. All metric metadata is retrieved if left empty"),
			mcp.DefaultString(""),
		),
		mcp.WithNumber("limit",
			mcp.Title("Maximum number of metrics"),
			mcp.Description("Maximum number of metrics to return"),
			mcp.DefaultNumber(0),
			mcp.Min(0),
		),
	)
	return mcp.NewTool(toolNameMetricsMetadata, options...)
}

func toolMetricsMetadataHandler(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	metric, err := GetToolReqParam[string](tcr, "metric", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit, err := GetToolReqParam[float64](tcr, "limit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	search, err := GetToolReqParam[string](tcr, "search", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	metricType, err := GetToolReqParam[string](tcr, "type", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	unit, err := GetToolReqParam[string](tcr, "unit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	req, err := CreateSelectRequest(ctx, cfg, tcr, "api", "v1", "metadata")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create select request: %v", err)), nil
	}

	q := req.URL.Query()
	if metric != "" {
		q.Add("metric", metric)
	}
	// Don't apply limit to API call if we need to filter client-side
	if limit != 0 && search == "" && metricType == "" && unit == "" {
		q.Add("limit", fmt.Sprintf("%.f", limit))
	}
	req.URL.RawQuery = q.Encode()

	// Get the response from the API
	result := GetTextBodyForRequest(req, cfg)

	// If no client-side filtering is needed, return the result as-is
	if search == "" && metricType == "" && unit == "" {
		return result, nil
	}

	// Parse the response to apply client-side filtering
	if len(result.Content) == 0 {
		return result, nil
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		return result, nil
	}

	var apiResponse struct {
		Status string `json:"status"`
		Data   map[string][]struct {
			Type string `json:"type"`
			Help string `json:"help"`
			Unit string `json:"unit"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(textContent.Text), &apiResponse); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to parse API response: %v", err)), nil
	}

	// Filter the data based on search, type, and unit parameters
	filteredData := make(map[string][]struct {
		Type string `json:"type"`
		Help string `json:"help"`
		Unit string `json:"unit"`
	})

	searchLower := strings.ToLower(search)
	typeLower := strings.ToLower(metricType)
	unitLower := strings.ToLower(unit)

	count := 0
	for metricName, metadataList := range apiResponse.Data {
		for _, metadata := range metadataList {
			// Apply filters
			if search != "" {
				metricNameLower := strings.ToLower(metricName)
				helpLower := strings.ToLower(metadata.Help)
				if !strings.Contains(metricNameLower, searchLower) && !strings.Contains(helpLower, searchLower) {
					continue
				}
			}

			if metricType != "" && strings.ToLower(metadata.Type) != typeLower {
				continue
			}

			if unit != "" && strings.ToLower(metadata.Unit) != unitLower {
				continue
			}

			// Add to filtered results
			if _, exists := filteredData[metricName]; !exists {
				filteredData[metricName] = []struct {
					Type string `json:"type"`
					Help string `json:"help"`
					Unit string `json:"unit"`
				}{}
			}
			filteredData[metricName] = append(filteredData[metricName], metadata)
			count++

			// Apply limit if specified
			if limit != 0 && count >= int(limit) {
				break
			}
		}

		if limit != 0 && count >= int(limit) {
			break
		}
	}

	// Create a filtered response
	filteredResponse := struct {
		Status string `json:"status"`
		Data   map[string][]struct {
			Type string `json:"type"`
			Help string `json:"help"`
			Unit string `json:"unit"`
		} `json:"data"`
	}{
		Status: apiResponse.Status,
		Data:   filteredData,
	}

	filteredJSON, err := json.Marshal(filteredResponse)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal filtered response: %v", err)), nil
	}

	return mcp.NewToolResultText(string(filteredJSON)), nil
}

func RegisterToolMetricsMetadata(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameMetricsMetadata) {
		return
	}
	s.AddTool(toolMetricsMetadata(c), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolMetricsMetadataHandler(ctx, c, request)
	})
}
