package cloud

import (
	vmcloud "github.com/VictoriaMetrics/victoriametrics-cloud-api-go/v1"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
)

func listOptions(cfg *config.Config) []vmcloud.ListOption {
	var types []vmcloud.DeploymentType
	for _, backend := range cfg.CloudBackendTypes() {
		types = append(types, tools.CloudDeploymentTypes[backend]...)
	}
	return []vmcloud.ListOption{vmcloud.WithDeploymentType(types...)}
}

func ptr[T any](v T) *T {
	return &v
}
