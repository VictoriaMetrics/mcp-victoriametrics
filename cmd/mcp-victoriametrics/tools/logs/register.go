package logs

import (
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func Register(s *server.MCPServer, c *config.Config) {
	registerToolHits(s, c)
	registerToolFlags(s, c)
	registerToolQuery(s, c)
	registerToolFacets(s, c)
	registerToolStreams(s, c)
	registerToolStreamIDs(s, c)
	registerToolStatsQuery(s, c)
	registerToolFieldNames(s, c)
	registerToolFieldValues(s, c)
	registerToolStatsQueryRange(s, c)
	registerToolStreamFieldNames(s, c)
	registerToolStreamFieldValues(s, c)
}
