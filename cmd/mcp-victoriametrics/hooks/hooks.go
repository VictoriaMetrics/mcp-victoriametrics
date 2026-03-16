package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/metrics"
)

func New(ms *metrics.Set) *server.Hooks {
	hooks := &server.Hooks{}

	hooks.AddAfterInitialize(func(_ context.Context, _ any, message *mcp.InitializeRequest, _ *mcp.InitializeResult) {
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_victoriametrics_initialize_total{client_name="%s",client_version="%s"}`,
			message.Params.ClientInfo.Name,
			message.Params.ClientInfo.Version,
		)).Inc()
	})

	hooks.AddAfterListTools(func(_ context.Context, _ any, _ *mcp.ListToolsRequest, _ *mcp.ListToolsResult) {
		ms.GetOrCreateCounter(`mcp_victoriametrics_list_tools_total`).Inc()
	})

	hooks.AddAfterListResources(func(_ context.Context, _ any, _ *mcp.ListResourcesRequest, _ *mcp.ListResourcesResult) {
		ms.GetOrCreateCounter(`mcp_victoriametrics_list_resources_total`).Inc()
	})

	hooks.AddAfterListPrompts(func(_ context.Context, _ any, _ *mcp.ListPromptsRequest, _ *mcp.ListPromptsResult) {
		ms.GetOrCreateCounter(`mcp_victoriametrics_list_prompts_total`).Inc()
	})

	hooks.AddAfterCallTool(func(_ context.Context, _ any, message *mcp.CallToolRequest, result any) {
		isError := false
		if r, ok := result.(*mcp.CallToolResult); ok {
			isError = r.IsError
		}
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_victoriametrics_call_tool_total{name="%s",is_error="%t"}`,
			message.Params.Name,
			isError,
		)).Inc()
	})

	hooks.AddAfterGetPrompt(func(_ context.Context, _ any, message *mcp.GetPromptRequest, _ *mcp.GetPromptResult) {
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_victoriametrics_get_prompt_total{name="%s"}`,
			message.Params.Name,
		)).Inc()
	})

	hooks.AddAfterReadResource(func(_ context.Context, _ any, message *mcp.ReadResourceRequest, _ *mcp.ReadResourceResult) {
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_victoriametrics_read_resource_total{uri="%s"}`,
			message.Params.URI,
		)).Inc()
	})

	hooks.AddOnError(func(_ context.Context, _ any, method mcp.MCPMethod, _ any, err error) {
		ms.GetOrCreateCounter(fmt.Sprintf(
			`mcp_victoriametrics_error_total{method="%s",error="%s"}`,
			method,
			err,
		)).Inc()
	})

	return hooks
}

func NewLoggerHooks() *server.Hooks {
	hooks := &server.Hooks{}

	hooks.AddOnRegisterSession(func(_ context.Context, session server.ClientSession) {
		slog.Info("Session registered",
			"session_id", session.SessionID(),
		)
	})

	hooks.AddOnUnregisterSession(func(_ context.Context, session server.ClientSession) {
		slog.Info("Session unregistered",
			"session_id", session.SessionID(),
		)
	})

	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		sessionID := extractSessionID(ctx)
		slog.Info("MCP request received",
			"request_id", id,
			"session_id", sessionID,
			"method", string(method),
			"message", toJSON(message),
		)
	})

	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		sessionID := extractSessionID(ctx)
		slog.Info("MCP request succeeded",
			"request_id", id,
			"session_id", sessionID,
			"method", string(method),
			"message", toJSON(message),
			"result", toJSON(result),
		)
	})

	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		sessionID := extractSessionID(ctx)
		slog.Error("MCP request failed",
			"request_id", id,
			"session_id", sessionID,
			"method", string(method),
			"message", toJSON(message),
			"error", err.Error(),
		)
	})

	hooks.AddAfterInitialize(func(_ context.Context, id any, msg *mcp.InitializeRequest, _ *mcp.InitializeResult) {
		slog.Info("Client initialized",
			"request_id", id,
			"client_name", msg.Params.ClientInfo.Name,
			"client_version", msg.Params.ClientInfo.Version,
			"protocol_version", msg.Params.ProtocolVersion,
		)
	})

	hooks.AddAfterCallTool(func(_ context.Context, id any, msg *mcp.CallToolRequest, result any) {
		isError := false
		if r, ok := result.(*mcp.CallToolResult); ok {
			isError = r.IsError
		}
		slog.Info("Tool called",
			"request_id", id,
			"tool_name", msg.Params.Name,
			"is_error", isError,
		)
	})

	return hooks
}

func Merge(hooksList ...*server.Hooks) *server.Hooks {
	combined := &server.Hooks{}
	for _, h := range hooksList {
		if h == nil {
			continue
		}
		combined.OnRegisterSession = append(combined.OnRegisterSession, h.OnRegisterSession...)
		combined.OnUnregisterSession = append(combined.OnUnregisterSession, h.OnUnregisterSession...)
		combined.OnBeforeAny = append(combined.OnBeforeAny, h.OnBeforeAny...)
		combined.OnSuccess = append(combined.OnSuccess, h.OnSuccess...)
		combined.OnError = append(combined.OnError, h.OnError...)
		combined.OnRequestInitialization = append(combined.OnRequestInitialization, h.OnRequestInitialization...)
		combined.OnBeforeInitialize = append(combined.OnBeforeInitialize, h.OnBeforeInitialize...)
		combined.OnAfterInitialize = append(combined.OnAfterInitialize, h.OnAfterInitialize...)
		combined.OnBeforePing = append(combined.OnBeforePing, h.OnBeforePing...)
		combined.OnAfterPing = append(combined.OnAfterPing, h.OnAfterPing...)
		combined.OnBeforeSetLevel = append(combined.OnBeforeSetLevel, h.OnBeforeSetLevel...)
		combined.OnAfterSetLevel = append(combined.OnAfterSetLevel, h.OnAfterSetLevel...)
		combined.OnBeforeListResources = append(combined.OnBeforeListResources, h.OnBeforeListResources...)
		combined.OnAfterListResources = append(combined.OnAfterListResources, h.OnAfterListResources...)
		combined.OnBeforeListResourceTemplates = append(combined.OnBeforeListResourceTemplates, h.OnBeforeListResourceTemplates...)
		combined.OnAfterListResourceTemplates = append(combined.OnAfterListResourceTemplates, h.OnAfterListResourceTemplates...)
		combined.OnBeforeReadResource = append(combined.OnBeforeReadResource, h.OnBeforeReadResource...)
		combined.OnAfterReadResource = append(combined.OnAfterReadResource, h.OnAfterReadResource...)
		combined.OnBeforeListPrompts = append(combined.OnBeforeListPrompts, h.OnBeforeListPrompts...)
		combined.OnAfterListPrompts = append(combined.OnAfterListPrompts, h.OnAfterListPrompts...)
		combined.OnBeforeGetPrompt = append(combined.OnBeforeGetPrompt, h.OnBeforeGetPrompt...)
		combined.OnAfterGetPrompt = append(combined.OnAfterGetPrompt, h.OnAfterGetPrompt...)
		combined.OnBeforeListTools = append(combined.OnBeforeListTools, h.OnBeforeListTools...)
		combined.OnAfterListTools = append(combined.OnAfterListTools, h.OnAfterListTools...)
		combined.OnBeforeCallTool = append(combined.OnBeforeCallTool, h.OnBeforeCallTool...)
		combined.OnAfterCallTool = append(combined.OnAfterCallTool, h.OnAfterCallTool...)
	}
	return combined
}

// extractSessionID extracts session ID from context
func extractSessionID(ctx context.Context) string {
	session := server.ClientSessionFromContext(ctx)
	if session != nil {
		return session.SessionID()
	}
	return ""
}

// toJSON converts any value to JSON string for logging
func toJSON(v any) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
