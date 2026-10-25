package resources

import (
	"strings"
	"testing"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func TestIsBackendDoc(t *testing.T) {
	testCases := []struct {
		path    string
		metrics bool
		logs    bool
		traces  bool
	}{
		{"vm/content/victoriametrics/README.md", true, false, false},
		{"vm/content/victorialogs/README.md", false, true, false},
		{"vm/content/victorialogs-grafana-datasource/README.md", false, true, false},
		{"vm/content/victoriatraces/README.md", false, false, true},
		{"vm/content/helm/victoria-logs-single/README.md", false, true, false},
		{"vm/content/helm/victoria-traces-cluster/README.md", false, false, true},
		{"vm/content/helm/victoria-metrics-single/README.md", true, false, false},
		{"vmsite/content/blog/victorialogs-release/index.md", true, true, false},
		{"vmsite/content/blog/distributed-tracing-intro/index.md", true, false, true},
		{"vmsite/content/blog/categories/index.md", true, false, false},
	}
	for _, tc := range testCases {
		for backend, expected := range map[config.Backend]bool{
			config.BackendMetrics: tc.metrics,
			config.BackendLogs:    tc.logs,
			config.BackendTraces:  tc.traces,
		} {
			if got := isBackendDoc(backend, tc.path); got != expected {
				t.Errorf("isBackendDoc(%q, %q) = %v, want %v", backend, tc.path, got, expected)
			}
		}
	}
}

func TestListDocFilesByBackend(t *testing.T) {
	for backend, productDir := range map[config.Backend]string{
		config.BackendLogs:   "vm/content/victorialogs/",
		config.BackendTraces: "vm/content/victoriatraces/",
	} {
		docs, err := ListDocFiles(backend)
		if err != nil {
			t.Fatalf("ListDocFiles(%q) failed: %v", backend, err)
		}
		found := false
		for _, doc := range docs {
			if strings.HasPrefix(doc.Path, "vm/content/victoriametrics/") {
				t.Fatalf("ListDocFiles(%q) returned VictoriaMetrics doc %q", backend, doc.Path)
			}
			found = found || strings.HasPrefix(doc.Path, productDir)
		}
		if !found {
			t.Errorf("ListDocFiles(%q) returned no docs from %q", backend, productDir)
		}
	}
}
