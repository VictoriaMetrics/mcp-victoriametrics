package auth

import (
	"crypto/subtle"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

// Middleware returns an HTTP middleware that validates the Authorization: Bearer <token>
// header on incoming requests. It skips authentication for /health/*, /metrics,
// and the root / (SPA static files). If the token is empty, the middleware is a no-op.
func Middleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for non-sensitive endpoints
			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				slog.Warn("Auth failed: missing Authorization header",
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
				)
				writeUnauthorized(w, "missing Authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				slog.Warn("Auth failed: invalid Authorization header format",
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
				)
				writeUnauthorized(w, "invalid Authorization header format")
				return
			}

			provided := parts[1]
			expected := token

			if subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) != 1 {
				slog.Warn("Auth failed: invalid token",
					"method", r.Method,
					"path", r.URL.Path,
					"remote_addr", r.RemoteAddr,
				)
				writeUnauthorized(w, "invalid token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isPublicPath(path string) bool {
	if path == "/" {
		return true
	}
	if strings.HasPrefix(path, "/health") {
		return true
	}
	if path == "/metrics" {
		return true
	}
	// Static assets served by SPA handler (files with extensions)
	if strings.Contains(path, ".") && !strings.HasPrefix(path, "/mcp") {
		return true
	}
	return false
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	}); err != nil {
		slog.Error("Failed to write unauthorized response", "error", err)
	}
}
