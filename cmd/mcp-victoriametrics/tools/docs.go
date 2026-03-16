package tools

import (
	"context"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics-Community/mcp-victoriametrics/cmd/mcp-victoriametrics/resources"
)

const ToolNameDocumentation = "documentation"
const defaultDocSearchLimit = 30

var (
	toolDocumentation = mcp.NewTool(ToolNameDocumentation,
		mcp.WithDescription("Search documentation resources for the given search query, returning the URIs of the resources that match the search criteria sorted by relevance. This tool can help to get context for any VictoriaMetrics related question."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Search documentation resources",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(false),
		}),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Title("Search query"),
			mcp.Description("Query for search (for example, list of keywords)"),
		),
		mcp.WithNumber("limit",
			mcp.Title("Maximum number of results"),
			mcp.Description("Maximum number of results to return"),
			mcp.DefaultNumber(defaultDocSearchLimit),
			mcp.Min(1),
		),
	)
)

func toolDocumentationHandler(_ context.Context, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := GetToolReqParam[string](tcr, "query", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	limit, err := GetToolReqParam[float64](tcr, "limit", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if limit < 1 {
		limit = defaultDocSearchLimit
	}

	rs, err := resources.SearchDocResources(query, int(limit))
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := &mcp.CallToolResult{Content: []mcp.Content{}}
	for _, resource := range rs {
		content, err := resources.GetDocResourceContent(resource.URI)
		if err != nil {
			log.Printf("error getting content for resource %s: %v", resource.URI, err)
			continue
		}
		result.Content = append(result.Content, mcp.EmbeddedResource{Type: "resource", Resource: content})
	}
	return result, nil
}

func RegisterToolDocumentation(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(ToolNameDocumentation) {
		return
	}
	s.AddTool(toolDocumentation, toolDocumentationHandler)
}
