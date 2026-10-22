package traces

import (
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func Register(s *server.MCPServer, c *config.Config) {
	registerToolTrace(s, c)
	registerToolTraces(s, c)
	registerToolServices(s, c)
	registerToolDependencies(s, c)
	registerToolServiceOperations(s, c)
}
