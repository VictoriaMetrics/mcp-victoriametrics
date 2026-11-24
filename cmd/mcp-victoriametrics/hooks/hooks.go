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

func RegisterMetrics(hooks *server.Hooks, ms *metrics.Set) {
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
}

func RegisterLogging(hooks *server.Hooks) {
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
		slog.Info("MCP request received", withPayload(ctx, []any{
			"request_id", id,
			"session_id", sessionID,
			"method", string(method),
		}, "message", message)...)
	})

	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		sessionID := extractSessionID(ctx)
		slog.Info("MCP request succeeded", withPayload(ctx, []any{
			"request_id", id,
			"session_id", sessionID,
			"method", string(method),
		}, "message", message, "result", result)...)
	})

	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		sessionID := extractSessionID(ctx)
		slog.Error("MCP request failed", withPayload(ctx, []any{
			"request_id", id,
			"session_id", sessionID,
			"method", string(method),
			"error", err.Error(),
		}, "message", message)...)
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

}

// extractSessionID extracts session ID from context
func extractSessionID(ctx context.Context) string {
	session := server.ClientSessionFromContext(ctx)
	if session != nil {
		return session.SessionID()
	}
	return ""
}

// withPayload appends JSON-encoded key/value payload pairs to args only when debug logging is enabled.
// Request and result bodies can be megabytes in size (e.g. resources/list), so at info level
// they are neither marshaled nor logged.
func withPayload(ctx context.Context, args []any, kvs ...any) []any {
	if !slog.Default().Enabled(ctx, slog.LevelDebug) {
		return args
	}
	for i := 0; i+1 < len(kvs); i += 2 {
		args = append(args, kvs[i], toJSON(kvs[i+1]))
	}
	return args
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
