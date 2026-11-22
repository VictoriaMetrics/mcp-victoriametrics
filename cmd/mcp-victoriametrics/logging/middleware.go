package logging

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

// responseWriter wraps http.ResponseWriter to capture status code and size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// Flush implements http.Flusher interface for SSE support
func (rw *responseWriter) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// requestSessionID returns the MCP session id sent by the client: the Mcp-Session-Id header
// (streamable HTTP) or the sessionId query parameter (SSE message endpoint).
// server.ClientSessionFromContext can't be used here: mcp-go attaches the session
// to its own request context after this middleware has run, so it is always nil.
func requestSessionID(r *http.Request) string {
	if id := r.Header.Get(server.HeaderKeySessionID); id != "" {
		return id
	}
	return r.URL.Query().Get("sessionId")
}

// Middleware creates HTTP logging middleware
func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip noisy endpoints
		if strings.HasPrefix(r.URL.Path, "/health") || r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// session id can be empty if not provided, e.g. for the initialize request,
		// which gets its session id in the response header
		sessionID := requestSessionID(r)

		// Log request start
		slog.Info("HTTP request started",
			"session_id", sessionID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		// Wrap response writer
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call next handler
		next.ServeHTTP(wrapped, r)
		if sessionID == "" {
			sessionID = wrapped.Header().Get(server.HeaderKeySessionID)
		}

		// Log request completion
		duration := time.Since(start)
		slog.Info("HTTP request completed",
			"session_id", sessionID,
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"size", wrapped.size,
			"duration_ms", duration.Milliseconds(),
		)
	})
}
