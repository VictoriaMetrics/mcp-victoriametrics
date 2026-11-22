package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestRejectStandaloneGET(t *testing.T) {
	mcpServer := server.NewMCPServer("test", "1.0.0")
	handler := rejectStandaloneGET(server.NewStreamableHTTPServer(mcpServer, server.WithStateLess(true)))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/mcp", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if allow := rec.Header().Get("Allow"); allow != "POST, DELETE, HEAD" {
		t.Fatalf("Allow = %q, want %q", allow, "POST, DELETE, HEAD")
	}

	// The methods listed in Allow keep working.
	for _, method := range []string{http.MethodDelete, http.MethodHead} {
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(method, "/mcp", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status = %d, want %d", method, rec.Code, http.StatusOK)
		}
	}
	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}`
	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("POST status = %d, want %d: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}
