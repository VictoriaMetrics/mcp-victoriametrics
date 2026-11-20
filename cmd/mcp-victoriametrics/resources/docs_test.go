package resources

import (
	"context"
	"fmt"
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
