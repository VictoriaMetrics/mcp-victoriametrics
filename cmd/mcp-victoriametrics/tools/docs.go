package tools

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/resources"
)

const ToolNameDocumentation = "documentation"
const defaultDocSearchLimit = 30

func toolDocumentation(c *config.Config) mcp.Tool {
	return mcp.NewTool(ToolNameDocumentation,
		mcp.WithDescription(fmt.Sprintf("Search the embedded %[1]s documentation and blog posts, returning the matching documents sorted by relevance. Use only for questions about %[1]s, not for other software.", c.Backend().ProductName())),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           fmt.Sprintf("Search %s documentation", c.Backend().ProductName()),
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
}

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

func registerToolDocumentation(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(ToolNameDocumentation) {
		return
	}
	s.AddTool(toolDocumentation(c), toolDocumentationHandler)
}
