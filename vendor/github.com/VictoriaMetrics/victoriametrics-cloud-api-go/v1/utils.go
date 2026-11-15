package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type apiKeyContextKeyType string

const apiKeyContextKey apiKeyContextKeyType = "X-VM-Cloud-Access"

func ContextWithDynamicAPIKey(ctx context.Context, apiKey string) context.Context {
	return context.WithValue(ctx, apiKeyContextKey, apiKey)
}

func requestAPI[R any](ctx context.Context, a *VMCloudAPIClient, method string, body io.Reader, path ...string) (R, error) {
	return requestAPIWithQuery[R](ctx, a, method, body, nil, path...)
}

func requestAPIWithQuery[R any](ctx context.Context, a *VMCloudAPIClient, method string, body io.Reader, query url.Values, path ...string) (R, error) {
	var result R
	respBody, _, err := doRequest(ctx, a, method, body, query, nil, path...)
	if err != nil {
		return result, err
	}
	if len(respBody) > 0 {
		// Special case for string type - just return the response body as a string
		if stringResult, ok := any(&result).(*string); ok {
			*stringResult = string(respBody)
		} else {
			// For other types, unmarshal as JSON
			if err = json.Unmarshal(respBody, &result); err != nil {
				return result, fmt.Errorf("failed to unmarshal response body: %w", err)
			}
		}
	}
	return result, nil
}

// doRequest sends one API request and returns the response body and headers.
// A non-2xx response is returned as *APIError.
func doRequest(ctx context.Context, a *VMCloudAPIClient, method string, body io.Reader, query url.Values, headers http.Header, path ...string) ([]byte, http.Header, error) {
	u := a.parsedURL.JoinPath(path...)
	if len(query) > 0 {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}
	for name, values := range headers {
		for _, v := range values {
			req.Header.Add(name, v)
		}
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
		return nil, nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, nil, newAPIError(resp.StatusCode, respBody)
	}
	return respBody, resp.Header, nil
}

// APIError is returned for a response with a non-2xx status. Use errors.As
// to read the status, for example to tell a failed precondition (409) from a
// missing file (404).
type APIError struct {
	StatusCode int
	Message    string
	body       string
}

func newAPIError(statusCode int, body []byte) *APIError {
	message := string(body)
	var payload struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(body, &payload) == nil && payload.Error != "" {
		message = payload.Error
	}
	return &APIError{StatusCode: statusCode, Message: message, body: string(body)}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("unexpected status code: %d, body: %s", e.StatusCode, e.body)
}

// RequestOption sets something on a single request. An option that returns
// an error stops the request before it is sent.
type RequestOption func(http.Header) error

// ErrEmptyETag is returned when WithIfMatch is given an empty version, for
// example because the response to GetDeploymentRuleFile carried no ETag.
// The request is not sent, so it can't silently become an unconditional write.
var ErrEmptyETag = errors.New("if-match: etag cannot be empty")

// WithIfMatch makes a write apply only if the file is still at etag, the
// version returned by GetDeploymentRuleFile. A newer version answers 409.
// The etag is sent as given, weak (W/"...") or strong; a bare value is quoted.
// An empty etag fails the request with ErrEmptyETag.
func WithIfMatch(etag string) RequestOption {
	return func(h http.Header) error {
		if etag == "" {
			return ErrEmptyETag
		}
		h.Set("If-Match", quoteETag(etag))
		return nil
	}
}

// WithIfNoneMatchAny makes a write apply only if the file does not exist yet.
func WithIfNoneMatchAny() RequestOption {
	return func(h http.Header) error {
		h.Set("If-None-Match", "*")
		return nil
	}
}

// WithActor names the person a write is made for, for the audit log. The API
// records it only for system keys.
func WithActor(email string) RequestOption {
	return func(h http.Header) error {
		h.Set("X-VM-Actor", email)
		return nil
	}
}

func requestHeaders(options []RequestOption) (http.Header, error) {
	if len(options) == 0 {
		return nil, nil
	}
	h := http.Header{}
	for _, option := range options {
		if err := option(h); err != nil {
			return nil, err
		}
	}
	return h, nil
}

// quoteETag returns etag as an entity-tag: "*", a quoted strong tag and a
// W/-prefixed weak tag are kept; a bare value is quoted.
func quoteETag(etag string) string {
	if etag == "*" || strings.HasPrefix(etag, `"`) || strings.HasPrefix(etag, `W/"`) {
		return etag
	}
	return `"` + etag + `"`
}
