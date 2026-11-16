package tools

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
)

func GetTextBodyForRequest(req *http.Request, cfg *config.Config, f ...func(s string) (string, error)) *mcp.CallToolResult {
	resp, err := cfg.HTTPClient().Do(req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to do request: %v", err))
	}
	defer func() { _ = resp.Body.Close() }()
	maxSize := cfg.MaxResponseSize()
	var bodyReader io.Reader = resp.Body
	if maxSize > 0 {
		// Read one byte past the limit to tell a body of exactly maxSize bytes from a larger one.
		bodyReader = io.LimitReader(resp.Body, maxSize+1)
	}
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read response body: %v", err))
	}
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("unexpected response status code %v: %s", resp.StatusCode, string(body)))
	}
	if maxSize > 0 && int64(len(body)) > maxSize {
		return mcp.NewToolResultError(fmt.Sprintf("response is larger than %d bytes (MCP_MAX_RESPONSE_SIZE); narrow the request with a more specific selector, a shorter time range or a `limit` argument where the tool has one", maxSize))
	}
	result := string(body)
	for _, fn := range f {
		if result, err = fn(result); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to process response body: %v", err))
		}
	}
	return mcp.NewToolResultText(result)
}

type ToolReqParamType interface {
	string | float64 | bool | []string | []any
}

func GetToolReqParam[T ToolReqParamType](tcr mcp.CallToolRequest, param string, required bool) (T, error) {
	var value T
	matchArg, ok := tcr.GetArguments()[param]
	if ok {
		value, ok = matchArg.(T)
		if !ok {
			return value, fmt.Errorf("%s has wrong type: %T", param, matchArg)
		}
	} else if required {
		return value, fmt.Errorf("%s param is required", param)
	}
	return value, nil
}

func ptr[T any](v T) *T {
	return &v
}
