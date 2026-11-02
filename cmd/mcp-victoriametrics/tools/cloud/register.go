package cloud

import (
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func Register(s *server.MCPServer, c *config.Config) {
	registerToolTiers(s, c)
	registerToolRegions(s, c)
	registerToolDeployments(s, c)
	registerToolAccessTokens(s, c)
	registerToolCloudProviders(s, c)
	registerToolRuleFile(s, c)
	registerToolRuleFilenames(s, c)
}
