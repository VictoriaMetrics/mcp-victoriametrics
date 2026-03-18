package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type apiKeyContextKeyType string

const apiKeyContextKey apiKeyContextKeyType = "X-VM-Cloud-Access"

func ContextWithDynamicAPIKey(ctx context.Context, apiKey string) context.Context {
	return context.WithValue(ctx, apiKeyContextKey, apiKey)
}

func requestAPI[R any](ctx context.Context, a *VMCloudAPIClient, method string, body io.Reader, path ...string) (R, error) {
	var result R
	reqURL := a.parsedURL.JoinPath(path...).String()
	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}
	apiKey := a.apiKey
	if apiKey == "" {
		if apiKeyFromCtx, ok := ctx.Value(apiKeyContextKey).(string); ok && apiKeyFromCtx != "" {
			apiKey = apiKeyFromCtx
		}
	}
	req.Header.Set(AccessTokenHeader, apiKey)
	resp, err := a.c.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return result, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBodyBytes))
	}
	if len(respBodyBytes) > 0 {
		// Special case for string type - just return the response body as a string
		if stringResult, ok := any(&result).(*string); ok {
			*stringResult = string(respBodyBytes)
		} else {
			// For other types, unmarshal as JSON
			if err = json.Unmarshal(respBodyBytes, &result); err != nil {
				return result, fmt.Errorf("failed to unmarshal response body: %w", err)
			}
		}
	}
	return result, nil
}
