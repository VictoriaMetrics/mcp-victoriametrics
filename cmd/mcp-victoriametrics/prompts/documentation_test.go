package prompts

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestPromptDocumentationHandler(t *testing.T) {
	// Test cases
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "Valid query",
			query:       "metrics",
			expectError: false,
		},
		{
			name:        "Empty query",
			query:       "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock prompt request
			gpr := mcp.GetPromptRequest{}
			if tc.query != "" {
				gpr.Params.Arguments = map[string]string{
					"query": tc.query,
				}
			}

			// Call the handler
			result, err := promptDocumentationHandler(context.Background(), "VictoriaMetrics", gpr)

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
			if result == nil {
				t.Fatal("Expected non-nil result")
			}

			// Check that the result contains a message
			if len(result.Messages) == 0 {
				t.Error("Expected at least one message in result")
			}

			// Check that the message contains the query
			message := result.Messages[0]
			if message.Role != mcp.RoleUser {
				t.Errorf("Expected role %s, got: %s", mcp.RoleUser, message.Role)
			}

			// Check that the content is a TextContent with the expected text
			textContent, ok := message.Content.(mcp.TextContent)
			if !ok {
				t.Fatal("Expected TextContent, got different content type")
			}

			expectedText := fmt.Sprintf(`Please tell me about "%s" by VictoriaMetrics documentation`, tc.query)
			if textContent.Text != expectedText {
				t.Errorf("Expected text '%s', got: '%s'", expectedText, textContent.Text)
			}
		})
	}
}
