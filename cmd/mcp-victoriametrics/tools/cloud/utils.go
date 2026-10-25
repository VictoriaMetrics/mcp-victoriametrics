package cloud

import (
	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

func listOptions(cfg *config.Config) []vmcloud.ListOption {
	return []vmcloud.ListOption{vmcloud.WithDeploymentType(tools.CloudDeploymentTypes[cfg.Backend()]...)}
}

func ptr[T any](v T) *T {
	return &v
}
