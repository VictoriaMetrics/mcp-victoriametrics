package config

import (
	"os"
	"testing"
)

func TestInitConfig(t *testing.T) {
	// Save original environment variables
	envVars := []string{
		"VM_INSTANCE_ENTRYPOINT",
		"VM_INSTANCE_TYPE",
		"VM_INSTANCE_BEARER_TOKEN",
		"VM_INSTANCE_HEADERS",
		"VM_DEFAULT_TENANT_ID",
		"VM_ENVIRONMENTS",
		"VM_DEFAULT_ENVIRONMENT",
		"VMC_API_KEY",
		"MCP_SERVER_MODE",
		"MCP_SSE_ADDR",
		"MCP_HEARTBEAT_INTERVAL",
	}
	originalEnv := make(map[string]string)
	for _, v := range envVars {
		originalEnv[v] = os.Getenv(v)
	}

	// Restore environment variables after test
	defer func() {
		for k, v := range originalEnv {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()

	// Helper to clear env
	clearEnv := func() {
		for _, v := range envVars {
			os.Unsetenv(v)
		}
		// Also clear some prefixed ones we might use
		os.Unsetenv("VM_INSTANCE_DEFAULT_ENTRYPOINT")
		os.Unsetenv("VM_INSTANCE_DEFAULT_TYPE")
		os.Unsetenv("VM_INSTANCE_DEMO_ENTRYPOINT")
		os.Unsetenv("VM_INSTANCE_DEMO_TYPE")
	}

	// Test case 1: Valid configuration (standard single instance mode)
	t.Run("Valid configuration standard", func(t *testing.T) {
		clearEnv()
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		os.Setenv("VM_INSTANCE_TYPE", "single")
		os.Setenv("MCP_SERVER_MODE", "stdio")
		os.Setenv("VM_INSTANCE_BEARER_TOKEN", "test-token")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if cfg.BearerToken() != "test-token" {
			t.Errorf("Expected bearer token 'test-token', got: %s", cfg.BearerToken())
		}
		if !cfg.IsSingle() {
			t.Error("Expected IsSingle() to be true")
		}
		if !cfg.IsStdio() {
			t.Error("Expected IsStdio() to be true")
		}
		if cfg.ListenAddr() != "localhost:8080" {
			t.Errorf("Expected address 'localhost:8080', got: %s", cfg.ListenAddr())
		}
	})

	// Test case 2: Multi-instance configuration
	t.Run("Multi-instance configuration", func(t *testing.T) {
		clearEnv()
		os.Setenv("VM_ENVIRONMENTS", "default,demo")
		os.Setenv("VM_INSTANCE_DEFAULT_ENTRYPOINT", "http://default.com")
		os.Setenv("VM_INSTANCE_DEFAULT_TYPE", "single")
		os.Setenv("VM_INSTANCE_DEMO_ENTRYPOINT", "http://demo.com")
		os.Setenv("VM_INSTANCE_DEMO_TYPE", "cluster")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check default
		def, err := cfg.Environment("default")
		if err != nil {
			t.Fatal(err)
		}
		if def.EntryPointURL().String() != "http://default.com" {
			t.Errorf("Expected http://default.com, got %s", def.EntryPointURL())
		}

		// Check demo
		demo, err := cfg.Environment("demo")
		if err != nil {
			t.Fatal(err)
		}
		if demo.EntryPointURL().String() != "http://demo.com" {
			t.Errorf("Expected http://demo.com, got %s", demo.EntryPointURL())
		}
		if !demo.IsCluster() {
			t.Error("Expected demo to be cluster")
		}
	})

	// Test case 3: Conflict validation
	t.Run("Variable conflict", func(t *testing.T) {
		clearEnv()
		os.Setenv("VM_ENVIRONMENTS", "demo")
		os.Setenv("VM_INSTANCE_ENTRYPOINT", "http://default.com") // standard var

		_, err := InitConfig()
		if err == nil {
			t.Fatal("Expected error when mixing VM_ENVIRONMENTS and standard vars")
		}
	})}
