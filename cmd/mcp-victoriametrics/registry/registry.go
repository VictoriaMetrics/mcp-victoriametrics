package registry

import (
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/prompts"
	metricsprompts "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/prompts/metrics"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
	cloudtools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/cloud"
	logstools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/logs"
	metricstools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/metrics"
	tracestools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/traces"
	vmalerttools "github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/vmalert"
)

// Register registers tools and prompts of the configured backend.
func Register(s *server.MCPServer, c *config.Config) {
	tools.Register(s, c)
	prompts.Register(s, c)
	cloudtools.Register(s, c)
	vmalerttools.Register(s, c)
	switch c.Backend() {
	case config.BackendLogs:
		logstools.Register(s, c)
	case config.BackendTraces:
		tracestools.Register(s, c)
	default:
		metricstools.Register(s, c)
		metricsprompts.Register(s, c)
	}
}
