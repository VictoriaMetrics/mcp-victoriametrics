package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v2"

	"github.com/VictoriaMetrics/VictoriaMetrics/app/vmalert-tool/unittest"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

const toolNameTestRules = "test_rules"

var (
	toolTestRules = mcp.NewTool(toolNameTestRules,
		mcp.WithDescription("Unit test alerting and recording rules. It use **[vmalert-tool](https://docs.victoriametrics.com/victoriametrics/vmalert-tool/)** under the hood . vmalert-tool unittest is compatible with Prometheus config format for tests."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "Unit test rules",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(false),
		}),
		mcp.WithArray("rule_files",
			mcp.Required(),
			mcp.Title("List of rule files contents in vmalert/prometheus format"),
			mcp.Description(`List of rule yaml files contents in vmalert/prometheus format. Each item in the list should be a string containing the yaml content of a rule file.`),
			mcp.Items(map[string]any{
				"type": "object",
				"parameters": map[string]any{
					"filename": map[string]any{
						"type":        "string",
						"description": "Optional filename for the rule file. If not provided, a default name will be used.",
					},
					"content": map[string]any{
						"type":        "string",
						"description": "Content of the rule file in vmalert/prometheus format. It should be a valid yaml string.",
					},
				},
			}),
		),
		mcp.WithString("evaluation_interval",
			mcp.Required(),
			mcp.Title("Evaluation interval"),
			mcp.Description(`Evaluation interval for the rules  specified in "rule_files". It should be in the format "1m", "5s", etc. This is used to determine how often the rules are evaluated.`),
			mcp.Pattern(`^([0-9]+)([a-z]+)$`),
		),
		mcp.WithArray("tests",
			mcp.Required(),
			mcp.Title("List of unit tests"),
			mcp.Description(`The list of unit test files to be checked during evaluation. See "vmalert-tool" docs for details on the format of the tests.`),
			mcp.Items(map[string]any{
				"type":  "object",
				"title": "Unit test group configuration",
				"properties": map[string]any{
					"interval": map[string]any{
						"type":        "string",
						"description": "Interval between samples for input series, in the format '1m', '5s', etc. default = evaluation_interval",
						"operational": true,
					},
					"input_series": map[string]any{
						"type":        "array",
						"description": "Time series to persist into the database according to configured <interval> before running tests.",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"series": map[string]any{
									"type": "string",
									"description": `series in the following format '<metric name>{<label name>=<label value>, ...}'. 
Examples:
  - series_name{label1="value1", label2="value2"}
  - go_goroutines{job="prometheus", instance="localhost:9090"}`,
								},
								"values": map[string]any{
									"type": "string",
									"description": `Values support several special equations:
   'a+bxc' becomes 'a a+b a+(2*b) a+(3*b) … a+(c*b)'
    Read this as series starts at a, then c further samples incrementing by b.
   'a-bxc' becomes 'a a-b a-(2*b) a-(3*b) … a-(c*b)'
    Read this as series starts at a, then c further samples decrementing by b (or incrementing by negative b).
   '_' represents a missing sample from scrape
   'stale' indicates a stale sample
Examples:
    1. '-2+4x3' becomes '-2 2 6 10' - series starts at -2, then 3 further samples incrementing by 4.
    2. ' 1-2x4' becomes '1 -1 -3 -5 -7' - series starts at 1, then 4 further samples decrementing by 2.
    3. ' 1x4' becomes '1 1 1 1 1' - shorthand for '1+0x4', series starts at 1, then 4 further samples incrementing by 0.
    4. ' 1 _x3 stale' becomes '1 _ _ _ stale' - the missing sample cannot increment, so 3 missing samples are produced by `,
								},
							},
						},
					},
					"name": map[string]any{
						"type":        "string",
						"description": "Name of the test group, optional",
						"optional":    true,
					},
					"alert_rule_test": map[string]any{
						"type":        "array",
						"description": "Unit tests for alerting rules",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"eval_time": map[string]any{
									"type":        "string",
									"description": "The time elapsed from time=0s when this alerting rule should be checked. Means this rule should be firing at this point, or shouldn't be firing if 'exp_alerts' is empty.",
								},
								"groupname": map[string]any{
									"type":        "string",
									"description": "Name of the group name to be tested.",
								},
								"alertname": map[string]any{
									"type":        "string",
									"description": "Name of the alert to be tested.",
								},
								"exp_alerts": map[string]any{
									"type":        "array",
									"description": "List of the expected alerts that are firing under the given alertname at the given evaluation time. If you want to test if an alerting rule should not be firing, then you can mention only the fields above and leave 'exp_alerts' empty.",
									"items": map[string]any{
										"type":        "object",
										"description": "Expected alert configuration",
										"properties": map[string]any{
											"exp_labels": map[string]any{
												"type":        "object",
												"description": "Labels of the expected alert",
												"additionalProperties": map[string]any{
													"type": "string",
												},
											},
											"exp_annotations": map[string]any{
												"type":        "object",
												"description": "Annotations of the expected alert",
												"additionalProperties": map[string]any{
													"type": "string",
												},
											},
										},
									},
								},
							},
						},
					},
					"metricsql_expr_test": map[string]any{
						"type":        "array",
						"description": "Unit tests for Metricsql expressions",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"eval_time": map[string]any{
									"type":        "string",
									"description": "The time elapsed from time=0s when this expression should be checked.",
								},
								"expr": map[string]any{
									"type":        "string",
									"description": "Metricsql expression to evaluate",
								},
								"exp_series": map[string]any{
									"type":        "array",
									"description": "Expected samples at the given evaluation time.",
									"items": map[string]any{
										"type":        "object",
										"description": "Expected series configuration",
										"properties": map[string]any{
											"labels": map[string]any{
												"type": "string",
												"description": `Labels of the sample in usual series notation '<metric name>{<label name>=<label value>, ...}'. 
Examples:
 - series_name{label1="value1", label2="value2"}
 - go_goroutines{job="prometheus", instance="localhost:9090"}`,
											},
											"values": map[string]any{
												"type":        "string",
												"description": `The expected value of the Metricsql expression.`,
											},
										},
									},
								},
							},
						},
					},
					"external_labels": map[string]any{
						"type":        "object",
						"description": "External labels for the tests. This is not accessible for templating, use '-external.label' cmd-line flag instead.",
						"additionalProperties": map[string]any{
							"type": "string",
						},
					},
				},
			}),
		),
	)
)

func toolTestRulesHandler(_ context.Context, _ *config.Config, tcr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ruleFiles, err := GetToolReqParam[[]any](tcr, "rule_files", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	evaluationInterval, err := GetToolReqParam[string](tcr, "evaluation_interval", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tests, err := GetToolReqParam[[]any](tcr, "tests", true)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	tmpDir, err := os.MkdirTemp("", "rules-*")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create temp dir: %v", err)), nil
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	var rulePaths []string
	for i, rf := range ruleFiles {
		rfMap, ok := rf.(map[string]any)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("invalid rule file format at index %d", i)), nil
		}
		filename, ok := rfMap["filename"].(string)
		if !ok || filename == "" {
			filename = fmt.Sprintf("rules_%d.yml", i)
		}
		content, ok := rfMap["content"].(string)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("invalid content in rule file %d", i)), nil
		}

		path := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to write rule file: %v", err)), nil
		}
		rulePaths = append(rulePaths, path)
	}
	testFileContent := map[string]any{
		"rule_files":          rulePaths,
		"evaluation_interval": evaluationInterval,
		"tests":               tests,
	}
	testFileData, err := yaml.Marshal(testFileContent)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal test file content: %v", err)), nil
	}
	testFileName := "test_rules.yaml"
	testFilePath := filepath.Join(tmpDir, testFileName)
	if err := os.WriteFile(testFilePath, testFileData, 0644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to write test file: %v", err)), nil
	}

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	failed := unittest.UnitTest([]string{testFilePath}, false, nil, "", "", "")
	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = originalStdout

	result := map[string]any{
		"status":  "success",
		"details": string(out),
	}
	if failed {
		result["status"] = "failed"
	}
	data, err := json.Marshal(result)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

func RegisterToolTestRules(s *server.MCPServer, c *config.Config) {
	if c.IsToolDisabled(toolNameTestRules) {
		return
	}
	s.AddTool(toolTestRules, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolTestRulesHandler(ctx, c, request)
	})
}
