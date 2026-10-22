package metrics

import (
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func Register(s *server.MCPServer, c *config.Config) {
	registerPromptUnusedMetrics(s, c)
	registerPromptRarelyUsedCardinalMetrics(s, c)
}
