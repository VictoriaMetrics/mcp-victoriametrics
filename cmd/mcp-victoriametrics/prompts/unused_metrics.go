package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

var (
	promptUnusedMetrics = mcp.NewPrompt("unused_metrics",
		mcp.WithPromptDescription("List of unused (never queried) metrics in the VictoriaMetrics instance API"),
		mcp.WithArgument("tenant",
			mcp.ArgumentDescription("Name of the tenant for which the list of unused (never queried) metrics will be displayed"),
		),
	)
)

func promptUnusedMetricsHandler(_ context.Context, gpr mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	tenant, err := GetPromptReqParam(gpr, "tenant", false)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %v", err)
	}
	return mcp.NewGetPromptResult(
		"",
		[]mcp.PromptMessage{
			{
				Role:    mcp.RoleUser,
				Content: mcp.NewTextContent(fmt.Sprintf("Please show me the list of metrics that are never queried in tenant %v of VictoriaMetrics instance and create relabel config to stop push these metrics", tenant)),
			},
		},
	), nil
}

func RegisterPromptUnusedMetrics(s *server.MCPServer, _ *config.Config) {
	s.AddPrompt(promptUnusedMetrics, promptUnusedMetricsHandler)
}
