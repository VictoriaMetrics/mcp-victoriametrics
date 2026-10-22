package tools

import (
	"fmt"
	"math"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/auth"
)

// WithTenantParam adds the tenant parameter to the tool input schema.
func WithTenantParam(description string) mcp.ToolOption {
	return mcp.WithString("tenant",
		mcp.Title("Tenant name"),
		mcp.Description(description),
		mcp.DefaultString("0"),
		mcp.Pattern(`^([0-9]+)(:[0-9]+)?$`),
	)
}

// GetToolReqTenantString returns the raw tenant parameter. Some models pass the tenant
// as a number despite the string schema, so whole non-negative numbers are accepted
// as well. An absent tenant results in an empty string.
func GetToolReqTenantString(tcr mcp.CallToolRequest) (string, error) {
	switch v := tcr.GetArguments()["tenant"].(type) {
	case nil:
		return "", nil
	case string:
		return v, nil
	case float64:
		if v < 0 || v > math.MaxUint32 || v != math.Trunc(v) {
			return "", fmt.Errorf("tenant must be a non-negative integer, got %v", v)
		}
		return strconv.FormatUint(uint64(v), 10), nil
	default:
		return "", fmt.Errorf("tenant has wrong type: %T", v)
	}
}

// GetToolReqTenantParam returns the tenant parameter validated in the `accountID[:projectID]` form.
// An absent tenant results in an empty string.
func GetToolReqTenantParam(tcr mcp.CallToolRequest) (string, error) {
	tenant, err := GetToolReqTenantString(tcr)
	if err != nil || tenant == "" {
		return "", err
	}
	token, err := auth.NewToken(tenant)
	if err != nil {
		return "", fmt.Errorf("invalid tenant %q, expected accountID or accountID:projectID: %v", tenant, err)
	}
	return token.String(), nil
}
