package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestToolTestRulesHandlerRejectsPathTraversal(t *testing.T) {
	// Point the temp dir to a nested dir, so escaped files would land inside baseDir.
	baseDir := t.TempDir()
	tmpRoot := filepath.Join(baseDir, "tmp")
	if err := os.Mkdir(tmpRoot, 0755); err != nil {
		t.Fatalf("failed to create temp root: %v", err)
	}
	t.Setenv("TMPDIR", tmpRoot)

	absTarget := filepath.Join(baseDir, "abs.yml")
	filenames := []string{
		"../../escaped.yml",
		"../escaped.yml",
		"sub/rules.yml",
		`..\escaped.yml`,
		absTarget,
		".",
		"..",
	}
	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			var tcr mcp.CallToolRequest
			tcr.Params.Name = toolNameTestRules
			tcr.Params.Arguments = map[string]any{
				"rule_files": []any{
					map[string]any{"filename": filename, "content": "pwned"},
				},
				"evaluation_interval": "1m",
				"tests":               []any{},
			}
			result, err := toolTestRulesHandler(context.Background(), nil, tcr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !result.IsError {
				t.Fatalf("expected error result for filename %q", filename)
			}
		})
	}

	// Nothing must be left outside the (already removed) per-call temp dirs.
	for _, dir := range []string{baseDir, tmpRoot} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("failed to read %q: %v", dir, err)
		}
		for _, e := range entries {
			if e.Name() != "tmp" {
				t.Fatalf("unexpected file %q left in %q", e.Name(), dir)
			}
		}
	}
}
