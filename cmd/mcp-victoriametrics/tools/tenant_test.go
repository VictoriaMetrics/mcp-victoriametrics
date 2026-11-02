package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestGetToolReqTenantParam(t *testing.T) {
	tests := []struct {
		name          string
		args          map[string]any
		expectedValue string
		expectError   bool
	}{
		{name: "string account", args: map[string]any{"tenant": "100"}, expectedValue: "100"},
		{name: "string account and project", args: map[string]any{"tenant": "100:200"}, expectedValue: "100:200"},
		{name: "numeric tenant", args: map[string]any{"tenant": float64(123)}, expectedValue: "123"},
		{name: "numeric zero tenant", args: map[string]any{"tenant": float64(0)}, expectedValue: "0"},
		{name: "large numeric tenant", args: map[string]any{"tenant": float64(421598)}, expectedValue: "421598"},
		{name: "missing tenant", expectedValue: ""},
		{name: "empty string tenant", args: map[string]any{"tenant": ""}, expectedValue: ""},
		{name: "fractional numeric tenant", args: map[string]any{"tenant": 1.5}, expectError: true},
		{name: "negative numeric tenant", args: map[string]any{"tenant": float64(-1)}, expectError: true},
		{name: "numeric tenant out of range", args: map[string]any{"tenant": float64(1 << 33)}, expectError: true},
		{name: "malformed string tenant", args: map[string]any{"tenant": "abc"}, expectError: true},
		{name: "string tenant with too many parts", args: map[string]any{"tenant": "1:2:3"}, expectError: true},
		{name: "invalid tenant type", args: map[string]any{"tenant": true}, expectError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tcr := mcp.CallToolRequest{}
			tcr.Params.Arguments = test.args

			value, err := GetToolReqTenantParam(tcr)
			if test.expectError {
				if err == nil {
					t.Fatalf("Expected an error, got value %q", value)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}
			if value != test.expectedValue {
				t.Errorf("Expected %q, got %q", test.expectedValue, value)
			}
		})
	}
}
