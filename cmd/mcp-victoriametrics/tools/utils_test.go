package tools

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

// TestGetTextBodyForRequest tests the GetTextBodyForRequest function
func TestGetTextBodyForRequest(t *testing.T) {
	// Create a mock config
	cfg := &config.Config{}

	// Save the original HTTP client
	originalClient := http.DefaultClient

	// Create a mock HTTP client
	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("test response")),
			},
		},
	}
	defer func() { http.DefaultClient = originalClient }()

	// Create a test request
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Call the function
	result := GetTextBodyForRequest(req, cfg)

	// Check the result
	if result.IsError {
		t.Error("Expected no error, got an error result")
	}

	// Extract the text content from the result
	if len(result.Content) == 0 {
		t.Fatal("Expected content in result, got empty content")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatal("Expected TextContent, got different content type")
	}

	if textContent.Text != "test response" {
		t.Errorf("Expected 'test response', got: %s", textContent.Text)
	}
}

// mockTransport is a mock implementation of http.RoundTripper
type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return m.response, m.err
}

// TestGetTextBodyForRequestError tests the error handling in GetTextBodyForRequest
func TestGetTextBodyForRequestError(t *testing.T) {
	// Create a mock config
	cfg := &config.Config{}

	// Save the original HTTP client
	originalClient := http.DefaultClient

	// Create a mock HTTP client that returns an error
	http.DefaultClient = &http.Client{
		Transport: &mockTransport{
			response: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString("error message")),
			},
		},
	}
	defer func() { http.DefaultClient = originalClient }()

	// Create a test request
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Call the function
	result := GetTextBodyForRequest(req, cfg)

	// Check the result
	if !result.IsError {
		t.Error("Expected an error result, got success")
	}
}

// TestGetToolReqParam tests the GetToolReqParam function
func TestGetToolReqParam(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		args          map[string]any
		param         string
		required      bool
		expectedValue string
		expectError   bool
	}{
		{
			name:          "Valid string parameter",
			args:          map[string]any{"test": "value"},
			param:         "test",
			required:      true,
			expectedValue: "value",
			expectError:   false,
		},
		{
			name:          "Missing required parameter",
			args:          map[string]any{},
			param:         "test",
			required:      true,
			expectedValue: "",
			expectError:   true,
		},
		{
			name:          "Missing optional parameter",
			args:          map[string]any{},
			param:         "test",
			required:      false,
			expectedValue: "",
			expectError:   false,
		},
		{
			name:          "Wrong type parameter",
			args:          map[string]any{"test": 123},
			param:         "test",
			required:      true,
			expectedValue: "",
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock tool request
			tcr := mcp.CallToolRequest{}
			tcr.Params.Arguments = tc.args

			// Call the function
			value, err := GetToolReqParam[string](tcr, tc.param, tc.required)

			// Check the result
			if tc.expectError && err == nil {
				t.Error("Expected an error, got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if value != tc.expectedValue {
				t.Errorf("Expected '%s', got: '%s'", tc.expectedValue, value)
			}
		})
	}
}

// TestGetToolReqParamFloat tests the GetToolReqParam function with float64 type
func TestGetToolReqParamFloat(t *testing.T) {
	// Create a mock tool request
	tcr := mcp.CallToolRequest{}
	tcr.Params.Arguments = map[string]any{
		"float": 123.45,
	}

	// Call the function
	value, err := GetToolReqParam[float64](tcr, "float", true)

	// Check the result
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if value != 123.45 {
		t.Errorf("Expected 123.45, got: %f", value)
	}
}

// TestGetToolReqParamBool tests the GetToolReqParam function with bool type
func TestGetToolReqParamBool(t *testing.T) {
	// Create a mock tool request
	tcr := mcp.CallToolRequest{}
	tcr.Params.Arguments = map[string]any{
		"bool": true,
	}

	// Call the function
	value, err := GetToolReqParam[bool](tcr, "bool", true)

	// Check the result
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !value {
		t.Error("Expected true, got false")
	}
}

// TestGetToolReqParamStringSlice tests the GetToolReqParam function with []string type
func TestGetToolReqParamStringSlice(t *testing.T) {
	// Create a mock tool request
	tcr := mcp.CallToolRequest{}
	tcr.Params.Arguments = map[string]any{
		"slice": []string{"a", "b", "c"},
	}

	// Call the function
	value, err := GetToolReqParam[[]string](tcr, "slice", true)

	// Check the result
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(value) != 3 || value[0] != "a" || value[1] != "b" || value[2] != "c" {
		t.Errorf("Expected [a b c], got: %v", value)
	}
}
