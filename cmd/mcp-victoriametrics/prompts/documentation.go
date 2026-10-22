package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func promptDocumentation(product string) mcp.Prompt {
	return mcp.NewPrompt("documentation",
		mcp.WithPromptDescription(fmt.Sprintf("Search by %s documentation or give context for answering questions", product)),
		mcp.WithArgument("query",
			mcp.RequiredArgument(),
			mcp.ArgumentDescription(fmt.Sprintf("What do you want to search in the %s documentation?", product)),
		),
	)
}

func promptDocumentationHandler(_ context.Context, product string, gpr mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	query, err := GetPromptReqParam(gpr, "query", true)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %v", err)
	}
	return mcp.NewGetPromptResult(
		"",
		[]mcp.PromptMessage{
			{
				Role:    mcp.RoleUser,
				Content: mcp.NewTextContent(fmt.Sprintf(`Please tell me about "%v" by %s documentation`, query, product)),
			},
		},
	), nil
}

func registerPromptDocumentation(s *server.MCPServer, c *config.Config) {
	product := c.Backend().ProductName()
	s.AddPrompt(promptDocumentation(product), func(ctx context.Context, gpr mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return promptDocumentationHandler(ctx, product, gpr)
	})
}
