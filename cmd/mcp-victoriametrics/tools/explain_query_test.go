package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/VictoriaMetrics/metricsql"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func TestToolExplainQueryHandler(t *testing.T) {
	// Initialize functions info
	err := initFunctionsInfo()
	if err != nil {
		t.Fatalf("Failed to initialize functions info: %v", err)
	}

	// Create a mock config
	cfg := &config.Config{}

	// Test cases
	testCases := []struct {
		name         string
		query        string
		expectError  bool
		validateFunc func(t *testing.T, result map[string]any)
	}{
		{
			name:        "Simple metric query",
			query:       "http_requests_total",
			expectError: false,
			validateFunc: func(t *testing.T, result map[string]any) {
				syntaxTree, ok := result["syntax_tree"].(map[string]any)
				if !ok {
					t.Fatal("Expected syntax_tree to be a map")
				}
				if syntaxTree["type"] != "MetricExpr" {
					t.Errorf("Expected MetricExpr type, got: %s", syntaxTree["type"])
				}
			},
		},
		{
			name:        "Function query",
			query:       "rate(http_requests_total[5m])",
			expectError: false,
			validateFunc: func(t *testing.T, result map[string]any) {
				syntaxTree, ok := result["syntax_tree"].(map[string]any)
				if !ok {
					t.Fatal("Expected syntax_tree to be a map")
				}
				if syntaxTree["type"] != "FuncExpr" {
					t.Errorf("Expected FuncExpr type, got: %s", syntaxTree["type"])
				}
				if syntaxTree["name"] != "rate" {
					t.Errorf("Expected function name 'rate', got: %s", syntaxTree["name"])
				}
			},
		},
		{
			name:        "Binary operation",
			query:       "http_requests_total > 100",
			expectError: false,
			validateFunc: func(t *testing.T, result map[string]any) {
				syntaxTree, ok := result["syntax_tree"].(map[string]any)
				if !ok {
					t.Fatal("Expected syntax_tree to be a map")
				}
				if syntaxTree["type"] != "BinaryOpExpr" {
					t.Errorf("Expected BinaryOpExpr type, got: %s", syntaxTree["type"])
				}
				if syntaxTree["op"] != ">" {
					t.Errorf("Expected operator '>', got: %s", syntaxTree["op"])
				}
			},
		},
		{
			name:        "Aggregate function",
			query:       "sum(http_requests_total) by (instance)",
			expectError: false,
			validateFunc: func(t *testing.T, result map[string]any) {
				syntaxTree, ok := result["syntax_tree"].(map[string]any)
				if !ok {
					t.Fatal("Expected syntax_tree to be a map")
				}
				if syntaxTree["type"] != "AggrFuncExpr" {
					t.Errorf("Expected AggrFuncExpr type, got: %s", syntaxTree["type"])
				}
				if syntaxTree["name"] != "sum" {
					t.Errorf("Expected function name 'sum', got: %s", syntaxTree["name"])
				}
			},
		},
		{
			name:         "Invalid query",
			query:        "invalid query syntax",
			expectError:  true,
			validateFunc: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock tool request
			tcr := mcp.CallToolRequest{}
			tcr.Params.Arguments = map[string]any{
				"query": tc.query,
			}

			// Call the handler
			result, err := toolExplainQueryHandler(context.Background(), cfg, tcr)

			// Check for errors
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tc.expectError {
				if !result.IsError {
					t.Error("Expected error result, got success")
				}
				return
			}

			if result.IsError {
				t.Fatalf("Expected success, got error: %v", result.Content)
			}

			// Extract the JSON content
			if len(result.Content) == 0 {
				t.Fatal("Expected content in result, got empty content")
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatal("Expected TextContent, got different content type")
			}

			// Parse the JSON
			var parsedResult map[string]any
			err = json.Unmarshal([]byte(textContent.Text), &parsedResult)
			if err != nil {
				t.Fatalf("Failed to parse JSON: %v", err)
			}

			// Validate the result
			if tc.validateFunc != nil {
				tc.validateFunc(t, parsedResult)
			}
		})
	}
}

func TestGetQueryInfo(t *testing.T) {
	// Initialize functions info
	err := initFunctionsInfo()
	if err != nil {
		t.Fatalf("Failed to initialize functions info: %v", err)
	}

	// Initialize metrics info
	err = initMetricsInfo()
	if err != nil {
		t.Fatalf("Failed to initialize metrics info: %v", err)
	}

	// Create a mock config and request
	cfg := &config.Config{}
	tcr := mcp.CallToolRequest{}

	// Test cases
	testCases := []struct {
		name        string
		query       string
		expectError bool
		checkFunc   func(t *testing.T, info map[string]any)
	}{
		{
			name:        "Simple metric",
			query:       "http_requests_total",
			expectError: false,
			checkFunc: func(t *testing.T, info map[string]any) {
				if _, ok := info["syntax_tree"]; !ok {
					t.Error("Expected syntax_tree in result")
				}
				if _, ok := info["types_info"]; !ok {
					t.Error("Expected types_info in result")
				}
				if _, ok := info["functions_info"]; !ok {
					t.Error("Expected functions_info in result")
				}
			},
		},
		{
			name:        "Invalid query",
			query:       "invalid query syntax",
			expectError: true,
			checkFunc:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			info, err := getQueryInfo(context.Background(), cfg, tcr, tc.query)

			// Check for errors
			if tc.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check the result
			if tc.checkFunc != nil {
				tc.checkFunc(t, info)
			}
		})
	}
}

func TestGetSyntaxTree(t *testing.T) {
	// Parse a simple expression
	expr, err := metricsql.Parse("http_requests_total")
	if err != nil {
		t.Fatalf("Failed to parse expression: %v", err)
	}

	// Call getSyntaxTree
	types := make(map[string]struct{})
	functions := make(map[string]struct{})
	metrics := make(map[string]struct{})
	tree := getSyntaxTree(expr, types, functions, metrics)

	// Check the result
	if tree["type"] != "MetricExpr" {
		t.Errorf("Expected MetricExpr type, got: %s", tree["type"])
	}

	// Check that types were collected
	if _, ok := types["MetricExpr"]; !ok {
		t.Error("Expected MetricExpr in types")
	}
}

func TestGetTypesDescriptions(t *testing.T) {
	// Create a map of types
	types := map[string]struct{}{
		"MetricExpr":  {},
		"NumberExpr":  {},
		"UnknownType": {},
	}

	// Call getTypesDescriptions
	descriptions := getTypesDescriptions(types)

	// Check the result
	if len(descriptions) != 3 {
		t.Errorf("Expected 3 descriptions, got: %d", len(descriptions))
	}

	// Check known types
	if desc, ok := descriptions["MetricExpr"]; !ok {
		t.Error("Expected MetricExpr in descriptions")
	} else if desc.Name != "Metric" {
		t.Errorf("Expected name 'Metric', got: %s", desc.Name)
	}

	if desc, ok := descriptions["NumberExpr"]; !ok {
		t.Error("Expected NumberExpr in descriptions")
	} else if desc.Name != "Number" {
		t.Errorf("Expected name 'Number', got: %s", desc.Name)
	}

	// Check unknown type
	if desc, ok := descriptions["UnknownType"]; !ok {
		t.Error("Expected UnknownType in descriptions")
	} else if desc.Name != "UnknownType" {
		t.Errorf("Expected name 'UnknownType', got: %s", desc.Name)
	}
}

func TestGetFunctionsInfo(t *testing.T) {
	// Initialize functions info
	err := initFunctionsInfo()
	if err != nil {
		t.Fatalf("Failed to initialize functions info: %v", err)
	}

	// Create a map of functions
	functions := map[string]struct{}{
		"rate":             {},
		"sum":              {},
		"unknown_function": {},
	}

	// Call getFunctionsInfo
	info := getFunctionsInfo(functions)

	// Check the result
	if len(info) != 3 {
		t.Errorf("Expected 3 function infos, got: %d", len(info))
	}

	// Check known functions
	if funcInfo, ok := info["rate"]; !ok {
		t.Error("Expected rate in function info")
	} else if funcInfo.Name != "rate" {
		t.Errorf("Expected name 'rate', got: %s", funcInfo.Name)
	}

	if funcInfo, ok := info["sum"]; !ok {
		t.Error("Expected sum in function info")
	} else if funcInfo.Name != "sum" {
		t.Errorf("Expected name 'sum', got: %s", funcInfo.Name)
	}

	// Check unknown function
	if funcInfo, ok := info["unknown_function"]; !ok {
		t.Error("Expected unknown_function in function info")
	} else if funcInfo.Name != "unknown_function" {
		t.Errorf("Expected name 'unknown_function', got: %s", funcInfo.Name)
	} else if funcInfo.Description != "Unknown function unknown_function" {
		t.Errorf("Expected description 'Unknown function unknown_function', got: %s", funcInfo.Description)
	}
}

func TestInitFunctionsInfo(t *testing.T) {
	// Call initFunctionsInfo
	err := initFunctionsInfo()

	// Check for errors
	if err != nil {
		t.Fatalf("Failed to initialize functions info: %v", err)
	}

	// Check that functionsInfo is not empty
	if len(functionsInfo) == 0 {
		t.Error("Expected non-empty functionsInfo")
	}

	// Check for some common functions
	commonFunctions := []string{"rate", "sum", "avg", "histogram_quantile"}
	for _, fn := range commonFunctions {
		if _, ok := functionsInfo[fn]; !ok {
			t.Errorf("Expected function %s in functionsInfo", fn)
		}
	}
}
