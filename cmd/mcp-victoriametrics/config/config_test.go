package config

import (
	"net/url"
	"os"
	"testing"
	"time"
)

func TestInitConfig(t *testing.T) {
	// Save original environment variables
	originalEntrypoint := os.Getenv("VM_INSTANCE_ENTRYPOINT")
	originalInstanceType := os.Getenv("VM_INSTANCE_TYPE")
	originalServerMode := os.Getenv("MCP_SERVER_MODE")
	originalSSEAddr := os.Getenv("MCP_SSE_ADDR")
	originalBearerToken := os.Getenv("VM_INSTANCE_BEARER_TOKEN")
	originalHeartbeatInterval := os.Getenv("MCP_HEARTBEAT_INTERVAL")
	originalHeaders := os.Getenv("VM_INSTANCE_HEADERS")
	originalPassthroughHeaders := os.Getenv("MCP_PASSTHROUGH_HEADERS")

	// Restore environment variables after test
	defer func() {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", originalEntrypoint)
		os.Setenv("VM_INSTANCE_TYPE", originalInstanceType)
		os.Setenv("MCP_SERVER_MODE", originalServerMode)
		os.Setenv("MCP_SSE_ADDR", originalSSEAddr)
		os.Setenv("VM_INSTANCE_BEARER_TOKEN", originalBearerToken)
		os.Setenv("MCP_HEARTBEAT_INTERVAL", originalHeartbeatInterval)
		os.Setenv("VM_INSTANCE_HEADERS", originalHeaders)
		os.Setenv("MCP_PASSTHROUGH_HEADERS", originalPassthroughHeaders)
	}()

	// Test case 1: Valid configuration
	t.Run("Valid configuration", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_SERVER_MODE", "stdio")
		os.Setenv("MCP_SSE_ADDR", "localhost:8080")
		os.Setenv("VM_INSTANCE_BEARER_TOKEN", "test-token")

		// Initialize config
		cfg, err := InitConfig()

		// Check for errors
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check config values
		if cfg.BearerToken() != "test-token" {
			t.Errorf("Expected bearer token 'test-token', got: %s", cfg.BearerToken())
		}
		if !cfg.IsSingle() {
			t.Error("Expected IsSingle() to be true")
		}
		if cfg.IsCluster() {
			t.Error("Expected IsCluster() to be false")
		}
		if !cfg.IsStdio() {
			t.Error("Expected IsStdio() to be true")
		}
		if cfg.IsSSE() {
			t.Error("Expected IsSSE() to be false")
		}
		if cfg.ListenAddr() != "localhost:8080" {
			t.Errorf("Expected SSE address 'localhost:8080', got: %s", cfg.ListenAddr())
		}
		expectedURL, _ := url.Parse("http://example.com")
		if cfg.EntryPointURL().String() != expectedURL.String() {
			t.Errorf("Expected entrypoint URL 'http://example.com', got: %s", cfg.EntryPointURL().String())
		}
		if !cfg.IsSingle() {
			t.Error("Expected IsSingle() to be true")
		}
		if cfg.IsCluster() {
			t.Error("Expected IsCluster() to be false")
		}
	})

	// Test case 2: Missing entrypoint
	t.Run("Missing entrypoint", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "")
		os.Setenv("VM_INSTANCE_TYPE", "single")

		// Initialize config
		_, err := InitConfig()

		// Check for errors
		if err == nil {
			t.Fatal("Expected error for missing entrypoint, got nil")
		}
	})

	// Test case 3: Missing instance type
	t.Run("Missing instance type", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "")

		// Initialize config
		_, err := InitConfig()

		// Check for errors
		if err == nil {
			t.Fatal("Expected error for missing instance type, got nil")
		}
	})

	// Test case 4: Invalid instance type
	t.Run("Invalid instance type", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "invalid")

		// Initialize config
		_, err := InitConfig()

		// Check for errors
		if err == nil {
			t.Fatal("Expected error for invalid instance type, got nil")
		}
	})

	// Test case 5: Invalid server mode
	t.Run("Invalid server mode", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_SERVER_MODE", "invalid")

		// Initialize config
		_, err := InitConfig()

		// Check for errors
		if err == nil {
			t.Fatal("Expected error for invalid server mode, got nil")
		}
	})

	// Test case 6: Default values
	t.Run("Default values", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_SERVER_MODE", "")
		os.Setenv("MCP_SSE_ADDR", "")

		// Initialize config
		cfg, err := InitConfig()

		// Check for errors
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check default values
		if !cfg.IsStdio() {
			t.Error("Expected default server mode to be stdio")
		}
		if cfg.ListenAddr() != "localhost:8080" {
			t.Errorf("Expected default SSE address 'localhost:8080', got: %s", cfg.ListenAddr())
		}
		if !cfg.IsSingle() {
			t.Error("Expected IsSingle() to be true")
		}
		if cfg.IsCluster() {
			t.Error("Expected IsCluster() to be false")
		}
	})

	// Test case 7: Cluster
	t.Run("Missing entrypoint", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "cluster")

		// Initialize config
		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check values
		if cfg.IsSingle() {
			t.Error("Expected IsSingle() to be true")
		}
		if !cfg.IsCluster() {
			t.Error("Expected IsCluster() to be false")
		}
	})
	// Test case 8: Heartbeat interval
	t.Run("Correct heartbeat interval", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_HEARTBEAT_INTERVAL", "30s")
		// Initialize config
		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		// Check values
		if cfg.HeartbeatInterval() != 30*time.Second {
			t.Errorf("Expected heartbeat interval to be 30 seconds, got: %d", cfg.HeartbeatInterval())
		}
	})
	// Test case 9: Invalid heartbeat interval
	t.Run("Incorrect heartbeat interval", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_HEARTBEAT_INTERVAL", "123")
		// Initialize config
		_, err := InitConfig()
		if err != nil && err.Error() != "failed to parse MCP_HEARTBEAT_INTERVAL: time: missing unit in duration \"123\"" {
			t.Errorf("Expected error 'invalid heartbeat interval: 123', got: %v", err)
		}

		os.Setenv("MCP_HEARTBEAT_INTERVAL", originalHeartbeatInterval)
	})

	// Test case 10: Custom headers parsing
	t.Run("Custom headers parsing", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_HEADERS", "CF-Access-Client-Id=test-client-id,CF-Access-Client-Secret=test-client-secret,Custom-Header=test-value")

		// Initialize config
		cfg, err := InitConfig()

		// Check for errors
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check custom headers
		headers := cfg.CustomHeaders()
		expectedHeaders := map[string]string{
			"CF-Access-Client-Id":     "test-client-id",
			"CF-Access-Client-Secret": "test-client-secret",
			"Custom-Header":           "test-value",
		}

		if len(headers) != len(expectedHeaders) {
			t.Errorf("Expected %d headers, got %d", len(expectedHeaders), len(headers))
		}

		for key, expectedValue := range expectedHeaders {
			if actualValue, exists := headers[key]; !exists {
				t.Errorf("Expected header %s to exist", key)
			} else if actualValue != expectedValue {
				t.Errorf("Expected header %s to have value %s, got %s", key, expectedValue, actualValue)
			}
		}
	})

	// Test case 11: Empty custom headers
	t.Run("Empty custom headers", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_HEADERS", "")

		// Initialize config
		cfg, err := InitConfig()

		// Check for errors
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check custom headers
		headers := cfg.CustomHeaders()
		if len(headers) != 0 {
			t.Errorf("Expected 0 headers, got %d", len(headers))
		}
	})

	// Test case 12: Invalid header format (should be ignored)
	t.Run("Invalid header format", func(t *testing.T) {
		// Set environment variables
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_HEADERS", "invalid-header,valid-header=value,another-invalid")

		// Initialize config
		cfg, err := InitConfig()

		// Check for errors
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check custom headers (only valid ones should be parsed)
		headers := cfg.CustomHeaders()
		expectedHeaders := map[string]string{
			"valid-header": "value",
		}

		if len(headers) != len(expectedHeaders) {
			t.Errorf("Expected %d headers, got %d", len(expectedHeaders), len(headers))
		}

		for key, expectedValue := range expectedHeaders {
			if actualValue, exists := headers[key]; !exists {
				t.Errorf("Expected header %s to exist", key)
			} else if actualValue != expectedValue {
				t.Errorf("Expected header %s to have value %s, got %s", key, expectedValue, actualValue)
			}
		}
	})

	t.Run("Default tenant ID in cluster mode", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "cluster")
		os.Setenv("VM_DEFAULT_TENANT_ID", "100:200")

		cfg, err := InitConfig()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if cfg.DefaultTenantID() != "100:200" {
			t.Errorf("Expected default tenant ID '100:200', got: %s", cfg.DefaultTenantID())
		}
	})

	t.Run("Empty default tenant ID (used default 0)", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "cluster")
		os.Setenv("VM_DEFAULT_TENANT_ID", "")

		cfg, err := InitConfig()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if cfg.DefaultTenantID() == "" {
			t.Errorf("Expected no default tenant ID, got: %s", cfg.DefaultTenantID())
		}
	})

	// Test case: Passthrough headers parsing
	t.Run("Passthrough headers parsing", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_PASSTHROUGH_HEADERS", "Authorization,X-Custom-Token,X-Request-ID")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		headers := cfg.PassthroughHeaders()
		expected := []string{"Authorization", "X-Custom-Token", "X-Request-ID"}

		if len(headers) != len(expected) {
			t.Fatalf("Expected %d passthrough headers, got %d", len(expected), len(headers))
		}
		for i, h := range headers {
			if h != expected[i] {
				t.Errorf("Expected header %q at index %d, got %q", expected[i], i, h)
			}
		}
	})

	// Test case: Empty passthrough headers
	t.Run("Empty passthrough headers", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_PASSTHROUGH_HEADERS", "")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(cfg.PassthroughHeaders()) != 0 {
			t.Errorf("Expected 0 passthrough headers, got %d", len(cfg.PassthroughHeaders()))
		}
	})

	// Test case: Passthrough headers with whitespace and empty entries
	t.Run("Passthrough headers whitespace trimming", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_PASSTHROUGH_HEADERS", " Authorization , , X-Custom-Token , ")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		headers := cfg.PassthroughHeaders()
		expected := []string{"Authorization", "X-Custom-Token"}

		if len(headers) != len(expected) {
			t.Fatalf("Expected %d passthrough headers, got %d", len(expected), len(headers))
		}
		for i, h := range headers {
			if h != expected[i] {
				t.Errorf("Expected header %q at index %d, got %q", expected[i], i, h)
			}
		}
	})
}
