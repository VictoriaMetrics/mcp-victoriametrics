package resources

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/mark3labs/mcp-go/mcp"
)

// TestGetDocResourceContent tests the GetDocResourceContent function
func TestGetDocResourceContent(t *testing.T) {
	// Setup test data
	testURI := "docs://test.md#0"
	testContent := mcp.TextResourceContents{
		URI:      testURI,
		MIMEType: "text/markdown",
		Text:     "# Test Document\n\nThis is a test document.",
	}

	// Save original contents and restore after test
	origContents := contents
	defer func() { contents = origContents }()

	// Initialize contents map for testing
	contents = map[string]mcp.ResourceContents{
		testURI: testContent,
	}

	// Test cases
	testCases := []struct {
		name        string
		uri         string
		expectError bool
	}{
		{
			name:        "Valid URI",
			uri:         testURI,
			expectError: false,
		},
		{
			name:        "Invalid URI",
			uri:         "docs://nonexistent.md#0",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			content, err := GetDocResourceContent(tc.uri)

			// Check for errors
			if tc.expectError {
				if err == nil {
					t.Error("Expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check the content
			if content == nil {
				t.Fatal("Expected non-nil content")
			}

			// Check that the content is a TextResourceContents
			textContent, ok := content.(mcp.TextResourceContents)
			if !ok {
				t.Fatal("Expected TextResourceContents, got different content type")
			}

			// Check the URI
			if textContent.URI != tc.uri {
				t.Errorf("Expected URI '%s', got: '%s'", tc.uri, textContent.URI)
			}
		})
	}
}

// TestSearchDocResources covers queries that match fewer documents than limit.
// bleve returns at most min(matches, limit) hits, and its result slice may have
// spare capacity filled with nil pointers (limit <= 10) or none at all (limit > 10).
func TestSearchDocResources(t *testing.T) {
	origIndex, origResources := searchIndex, resources
	defer func() { searchIndex, resources = origIndex, origResources }()

	var err error
	if searchIndex, err = bleve.NewMemOnly(bleve.NewIndexMapping()); err != nil {
		t.Fatalf("error creating index: %v", err)
	}
	resources = map[string]mcp.Resource{}
	for i := range 3 {
		uri := fmt.Sprintf("docs://doc%d.md#0", i)
		resources[uri] = mcp.NewResource(uri, uri)
		if err := searchIndex.Index(uri, DocFileInfo{Content: "vmagent remote write"}); err != nil {
			t.Fatalf("error indexing %s: %v", uri, err)
		}
	}

	testCases := []struct {
		name  string
		limit int
		want  int
	}{
		{name: "limit below matches", limit: 2, want: 2},
		{name: "limit above matches, small", limit: 5, want: 3},
		{name: "limit above matches, large", limit: 30, want: 3},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rs, err := SearchDocResources("vmagent", tc.limit)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if len(rs) != tc.want {
				t.Errorf("Expected %d results, got: %d", tc.want, len(rs))
			}
		})
	}

	if _, err := SearchDocResources("nonexistentterm", 5); err == nil {
		t.Error("Expected an error for a query without matches, got nil")
	}
}

// TestDocResourcesHandler tests the docResourcesHandler function
func TestDocResourcesHandler(t *testing.T) {
	// Setup test data
	testURI := "docs://test.md#0"
	testContent := mcp.TextResourceContents{
		URI:      testURI,
		MIMEType: "text/markdown",
		Text:     "# Test Document\n\nThis is a test document.",
	}

	// Save original contents and restore after test
	origContents := contents
	defer func() { contents = origContents }()

	// Initialize contents map for testing
	contents = map[string]mcp.ResourceContents{
		testURI: testContent,
	}

	// Test cases
	testCases := []struct {
		name        string
		uri         string
		expectError bool
	}{
		{
			name:        "Valid URI",
			uri:         testURI,
			expectError: false,
		},
		{
			name:        "Invalid URI",
			uri:         "docs://nonexistent.md#0",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock request
			req := mcp.ReadResourceRequest{}
			req.Params.URI = tc.uri

			// Call the handler
			result, err := docResourcesHandler(context.Background(), req)

			// Check for errors
			if tc.expectError {
				if err == nil {
					t.Error("Expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check the result
			if len(result) != 1 {
				t.Fatalf("Expected 1 result, got: %d", len(result))
			}

			// Check that the content is correct
			content := result[0]
			textContent, ok := content.(mcp.TextResourceContents)
			if !ok {
				t.Fatal("Expected TextResourceContents, got different content type")
			}

			// Check the URI
			if textContent.URI != tc.uri {
				t.Errorf("Expected URI '%s', got: '%s'", tc.uri, textContent.URI)
			}
		})
	}
}

func testDocFiles() []DocFileInfo {
	return []DocFileInfo{
		{Path: "vm/a.md", Content: "# Retention\n\nRetention filters drop old samples.", Name: "Retention"},
		{Path: "vm/b.md", Content: "# Downsampling\n\nDownsampling reduces resolution.", Name: "Downsampling"},
	}
}

func buildTestDocsIndexArchive(t *testing.T, docFiles []DocFileInfo) []byte {
	t.Helper()
	tmpDir := t.TempDir()
	indexDir := filepath.Join(tmpDir, "index")
	if err := buildDocsIndex(indexDir, docFiles); err != nil {
		t.Fatalf("buildDocsIndex: %v", err)
	}
	archivePath := filepath.Join(tmpDir, "docs-index.tar.gz")
	if err := writeDocsIndexArchive(indexDir, archivePath); err != nil {
		t.Fatalf("writeDocsIndexArchive: %v", err)
	}
	archive, err := os.ReadFile(archivePath)
	if err != nil {
		t.Fatalf("read archive: %v", err)
	}
	return archive
}

func searchTestIndex(t *testing.T, index bleve.Index, query string) []string {
	t.Helper()
	result, err := index.Search(bleve.NewSearchRequest(bleve.NewMatchQuery(query)))
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	ids := make([]string, 0, len(result.Hits))
	for _, hit := range result.Hits {
		ids = append(ids, hit.ID)
	}
	return ids
}

func TestArchivedDocsIndex(t *testing.T) {
	docFiles := testDocFiles()
	hash := docsHash(docFiles)
	archive := buildTestDocsIndexArchive(t, docFiles)
	dir := filepath.Join(t.TempDir(), "docs-index")

	index, err := openArchivedDocsIndex(archive, dir, hash)
	if err != nil {
		t.Fatalf("openArchivedDocsIndex: %v", err)
	}
	if got := searchTestIndex(t, index, "downsampling"); len(got) != 1 || got[0] != "docs://vm/b.md#0" {
		t.Fatalf("unexpected search result: %v", got)
	}
	_ = index.Close()

	marker := filepath.Join(dir, "marker")
	if err := os.WriteFile(marker, nil, 0o644); err != nil {
		t.Fatalf("write marker: %v", err)
	}
	index, err = openArchivedDocsIndex(archive, dir, hash)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	_ = index.Close()
	if _, err := os.Stat(marker); err != nil {
		t.Fatal("expected the unpacked index to be reused")
	}

	changed := testDocFiles()
	changed[1].Content += " More."
	if _, err := openArchivedDocsIndex(archive, dir, docsHash(changed)); err == nil {
		t.Fatal("expected an error for an index built from different docs")
	}
	if _, err := openArchivedDocsIndex([]byte("not an archive"), filepath.Join(t.TempDir(), "x"), hash); err == nil {
		t.Fatal("expected an error for a damaged archive")
	}
}

func TestArchivedDocsIndexReadOnly(t *testing.T) {
	docFiles := testDocFiles()
	hash := docsHash(docFiles)
	dir := filepath.Join(t.TempDir(), "docs-index")
	index, err := openArchivedDocsIndex(buildTestDocsIndexArchive(t, docFiles), dir, hash)
	if err != nil {
		t.Fatalf("openArchivedDocsIndex: %v", err)
	}
	_ = index.Close()
	if err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		mode := os.FileMode(0o444)
		if d.IsDir() {
			mode = 0o555
		}
		return os.Chmod(p, mode)
	}); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() {
		_ = filepath.WalkDir(dir, func(p string, _ fs.DirEntry, _ error) error { return os.Chmod(p, 0o755) })
	})
	index, err = openPrebuiltDocsIndex(dir, hash)
	if err != nil {
		t.Fatalf("openPrebuiltDocsIndex on read-only files: %v", err)
	}
	_ = index.Close()
}

func TestExtractDocsIndexArchiveRejectsUnsafePaths(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: "../escape", Typeflag: tar.TypeReg, Mode: 0o644}); err != nil {
		t.Fatal(err)
	}
	if err := errors.Join(tw.Close(), gz.Close()); err != nil {
		t.Fatal(err)
	}
	if err := extractDocsIndexArchive(buf.Bytes(), t.TempDir()); err == nil {
		t.Fatal("expected an error for a path outside the target directory")
	}
}

func TestOpenDocsIndexFallback(t *testing.T) {
	docFiles := testDocFiles()
	index := openDocsIndex(docFiles)
	defer func() { _ = index.Close() }()
	if got := searchTestIndex(t, index, "retention"); len(got) != 1 || got[0] != "docs://vm/a.md#0" {
		t.Fatalf("unexpected search result: %v", got)
	}
}
