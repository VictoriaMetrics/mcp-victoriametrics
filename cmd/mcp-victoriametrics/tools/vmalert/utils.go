package vmalert

import (
	"context"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/logs"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/metrics"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/tools/traces"
)

// createRequest creates a request to the vmalert API proxied by the configured backend:
// /vmalert/... for VictoriaMetrics (/select/<tenant>/prometheus/vmalert/... in cluster mode)
// and /select/vmalert/... for VictoriaLogs and VictoriaTraces.
func createRequest(ctx context.Context, cfg *config.Config, tcr mcp.CallToolRequest, path ...string) (*http.Request, error) {
	switch cfg.Backend() {
	case config.BackendLogs:
		req, err := logs.CreateAdminRequest(ctx, cfg, tcr, append([]string{"select", "vmalert"}, path...)...)
		if err != nil {
			return nil, err
		}
		return req, logs.SetTenantHeaders(req, cfg, tcr)
	case config.BackendTraces:
		req, err := traces.CreateAdminRequest(ctx, cfg, tcr, append([]string{"select", "vmalert"}, path...)...)
		if err != nil {
			return nil, err
		}
		return req, traces.SetTenantHeaders(req, cfg, tcr)
	default:
		return metrics.CreateSelectRequest(ctx, cfg, tcr, append([]string{"vmalert"}, path...)...)
	}
}

// tenantParam returns the tenant parameter in the format of the configured backend:
// a path part for VictoriaMetrics cluster and AccountID/ProjectID headers for VictoriaLogs and VictoriaTraces.
func tenantParam(c *config.Config, description string) (mcp.ToolOption, bool) {
	switch c.Backend() {
	case config.BackendLogs, config.BackendTraces:
		return mcp.WithString("tenant",
			mcp.Title("Tenant name (Account ID and Project ID)"),
			mcp.Description(description+" (format AccountID:ProjectID)"),
			mcp.DefaultString("0:0"),
			mcp.Pattern(`^([0-9]+)(:[0-9]+)$`),
		), true
	default:
		if c.IsCluster() || c.IsCloud() {
			return tools.WithTenantParam(description), true
		}
		return nil, false
	}
}

func ptr[T any](v T) *T {
	return &v
}
