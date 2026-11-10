package config

import (
	"net/url"
	"os"
	"testing"
	"time"
)

func TestInitConfigMaxResponseSize(t *testing.T) {
	testCases := []struct {
		name          string
		value         string
		expectedValue int64
		expectedError string
	}{
		{
			name:          "default",
			value:         "",
			expectedValue: 16 * 1024 * 1024,
		},
		{
			name:          "custom size",
			value:         "1048576",
			expectedValue: 1048576,
		},
		{
			name:          "disabled",
			value:         "0",
			expectedValue: 0,
		},
		{
			name:          "negative",
			value:         "-1",
			expectedError: "MCP_MAX_RESPONSE_SIZE must be non-negative",
		},
		{
			name:          "not a number",
			value:         "16MiB",
			expectedError: "failed to parse MCP_MAX_RESPONSE_SIZE: strconv.ParseInt: parsing \"16MiB\": invalid syntax",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
			t.Setenv("VM_INSTANCE_TYPE", "single")
			t.Setenv("MCP_MAX_RESPONSE_SIZE", tc.value)
			cfg, err := InitConfig()
			if tc.expectedError != "" {
				if err == nil || err.Error() != tc.expectedError {
					t.Fatalf("Expected error %q, got: %v", tc.expectedError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}
			if cfg.MaxResponseSize() != tc.expectedValue {
				t.Errorf("Expected max response size %d, got: %d", tc.expectedValue, cfg.MaxResponseSize())
			}
		})
	}
}

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
	originalStateless := os.Getenv("MCP_STATELESS")
	originalDisableLocalhostProtection := os.Getenv("MCP_DISABLE_LOCALHOST_PROTECTION")

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
		os.Setenv("MCP_STATELESS", originalStateless)
		os.Setenv("MCP_DISABLE_LOCALHOST_PROTECTION", originalDisableLocalhostProtection)
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

	// Test case: Stateless defaults to false
	t.Run("Stateless default", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_STATELESS", "")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if cfg.Stateless() {
			t.Error("Expected Stateless() to default to false")
		}
	})

	// Test case: Stateless enabled
	t.Run("Stateless enabled", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_STATELESS", "true")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !cfg.Stateless() {
			t.Error("Expected Stateless() to be true when MCP_STATELESS=true")
		}

		os.Setenv("MCP_STATELESS", "")
	})

	// Test case: Stateless explicitly false
	t.Run("Stateless explicit false", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_STATELESS", "false")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if cfg.Stateless() {
			t.Error("Expected Stateless() to be false when MCP_STATELESS=false")
		}

		os.Setenv("MCP_STATELESS", "")
	})

	// Test case: Invalid stateless value
	t.Run("Invalid stateless value", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_STATELESS", "notabool")

		_, err := InitConfig()
		if err == nil {
			t.Fatal("Expected error for invalid MCP_STATELESS, got nil")
		}

		os.Setenv("MCP_STATELESS", "")
	})

	// Test case: DisableLocalhostProtection defaults to false
	t.Run("DisableLocalhostProtection default", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_DISABLE_LOCALHOST_PROTECTION", "")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if cfg.DisableLocalhostProtection() {
			t.Error("Expected DisableLocalhostProtection() to default to false")
		}
	})

	// Test case: DisableLocalhostProtection enabled
	t.Run("DisableLocalhostProtection enabled", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_DISABLE_LOCALHOST_PROTECTION", "true")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !cfg.DisableLocalhostProtection() {
			t.Error("Expected DisableLocalhostProtection() to be true when MCP_DISABLE_LOCALHOST_PROTECTION=true")
		}

		os.Setenv("MCP_DISABLE_LOCALHOST_PROTECTION", "")
	})

	// Test case: Invalid DisableLocalhostProtection value
	t.Run("Invalid DisableLocalhostProtection value", func(t *testing.T) {
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_DISABLE_LOCALHOST_PROTECTION", "notabool")

		_, err := InitConfig()
		if err == nil {
			t.Fatal("Expected error for invalid MCP_DISABLE_LOCALHOST_PROTECTION, got nil")
		}

		os.Setenv("MCP_DISABLE_LOCALHOST_PROTECTION", "")
	})
}
