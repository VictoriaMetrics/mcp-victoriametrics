package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

var (
	promptRarelyUsedCardinalMetrics = mcp.NewPrompt("rarely_used_metrics_with_high_cardinality",
		mcp.WithPromptDescription("List of rarely used metrics with high cardinality in the VictoriaMetrics instance API"),
		mcp.WithArgument("tenant",
			mcp.ArgumentDescription("Name of the tenant for which the list of rarely used metrics with high cardinality will be displayed"),
		),
	)
)

func promptRarelyUsedCardinalMetricsHandler(_ context.Context, gpr mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	tenant, err := GetPromptReqParam(gpr, "tenant", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %v", err)
	}
	return mcp.NewGetPromptResult(
		"",
		[]mcp.PromptMessage{
			{
				Role:    mcp.RoleUser,
				Content: mcp.NewTextContent(fmt.Sprintf("Do i have metrics with high cardinality that are never or rarely queried in tenant %v?", tenant)),
			},
		},
	), nil
}

func RegisterPromptRarelyUsedCardinalMetrics(s *server.MCPServer, _ *config.Config) {
	s.AddPrompt(promptRarelyUsedCardinalMetrics, promptRarelyUsedCardinalMetricsHandler)
}
