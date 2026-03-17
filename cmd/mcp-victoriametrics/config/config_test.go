package config

import (
	"testing"
	"time"
)

func TestInitConfigLegacyInstance(t *testing.T) {
	t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
	t.Setenv("VM_INSTANCE_TYPE", "cluster")
	t.Setenv("VM_INSTANCE_BEARER_TOKEN", "secret")
	t.Setenv("VM_INSTANCE_HEADERS", "A=1,B=2")
	t.Setenv("VM_DEFAULT_TENANT_ID", "100:200")
	t.Setenv("MCP_HEARTBEAT_INTERVAL", "45s")

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}
	instance, err := cfg.ResolveInstance("")
	if err != nil {
		t.Fatalf("ResolveInstance(default) error = %v", err)
	}
	if cfg.DefaultInstanceName() != "default" {
		t.Fatalf("DefaultInstanceName() = %q", cfg.DefaultInstanceName())
	}
	if !instance.IsCluster() {
		t.Fatal("expected default instance to be cluster")
	}
	if instance.BearerToken() != "secret" {
		t.Fatalf("BearerToken() = %q", instance.BearerToken())
	}
	if got := instance.EntryPointURL().String(); got != "http://example.com" {
		t.Fatalf("EntryPointURL() = %q", got)
	}
	if got := instance.DefaultTenantID(); got != "100:200" {
		t.Fatalf("DefaultTenantID() = %q", got)
	}
	if got := instance.CustomHeaders()["A"]; got != "1" {
		t.Fatalf("CustomHeaders()[A] = %q", got)
	}
	if cfg.HeartbeatInterval() != 45*time.Second {
		t.Fatalf("HeartbeatInterval() = %v", cfg.HeartbeatInterval())
	}
	if instance.Name() != "default" {
		t.Fatalf("ResolveInstance(default).Name() = %q", instance.Name())
	}
}

func TestInitConfigLegacyCloud(t *testing.T) {
	t.Setenv("VMC_API_KEY", "test-api-key")

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}
	instance, err := cfg.ResolveInstance("")
	if err != nil {
		t.Fatalf("ResolveInstance(default) error = %v", err)
	}
	if !instance.IsCloud() {
		t.Fatal("expected default instance to be cloud")
	}
	if !cfg.HasCloudInstances() {
		t.Fatal("expected HasCloudInstances() to be true")
	}
	if !cfg.HasOnlyCloudInstances() {
		t.Fatal("expected HasOnlyCloudInstances() to be true")
	}
}

func TestInitConfigMultiInstance(t *testing.T) {
	t.Setenv("VM_ENVIRONMENTS", "demo,prod_cloud")
	t.Setenv("VM_DEFAULT_ENVIRONMENT", "prod_cloud")
	t.Setenv("VM_INSTANCE_DEMO_ENTRYPOINT", "http://demo.example.com")
	t.Setenv("VM_INSTANCE_DEMO_TYPE", "single")
	t.Setenv("VM_INSTANCE_DEMO_HEADERS", "X-Test=yes")
	t.Setenv("VM_INSTANCE_PROD_CLOUD_DEFAULT_TENANT_ID", "7")
	t.Setenv("VMC_PROD_CLOUD_API_KEY", "test-api-key")

	cfg, err := InitConfig()
	if err != nil {
		t.Fatalf("InitConfig() error = %v", err)
	}
	if !cfg.HasMultipleInstances() {
		t.Fatal("expected multiple instances")
	}
	if got := cfg.DefaultInstanceName(); got != "prod_cloud" {
		t.Fatalf("DefaultInstanceName() = %q", got)
	}
	if got := cfg.InstanceNames(); len(got) != 2 || got[0] != "demo" || got[1] != "prod_cloud" {
		t.Fatalf("InstanceNames() = %#v", got)
	}
	if !cfg.HasCloudInstances() {
		t.Fatal("expected HasCloudInstances() to be true")
	}
	if !cfg.HasClusterInstances() {
		t.Fatal("expected HasClusterInstances() to be true when cloud env exists")
	}

	demo, err := cfg.ResolveInstance("demo")
	if err != nil {
		t.Fatalf("ResolveInstance(demo) error = %v", err)
	}
	if demo.IsCloud() {
		t.Fatal("demo should not be cloud")
	}
	if got := demo.EntryPointURL().String(); got != "http://demo.example.com" {
		t.Fatalf("demo EntryPointURL() = %q", got)
	}
	if got := demo.CustomHeaders()["X-Test"]; got != "yes" {
		t.Fatalf("demo CustomHeaders()[X-Test] = %q", got)
	}

	prod, err := cfg.ResolveInstance("prod_cloud")
	if err != nil {
		t.Fatalf("ResolveInstance(prod_cloud) error = %v", err)
	}
	if !prod.IsCloud() {
		t.Fatal("prod_cloud should be cloud")
	}
	if got := prod.DefaultTenantID(); got != "7" {
		t.Fatalf("prod_cloud DefaultTenantID() = %q", got)
	}
	if _, err := cfg.ResolveInstance(""); err != nil {
		t.Fatalf("ResolveInstance(default) error = %v", err)
	}
}

func TestInitConfigServerModeAndListenDefaults(t *testing.T) {
	t.Run("defaults to stdio and default listen addr", func(t *testing.T) {
		t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		t.Setenv("VM_INSTANCE_TYPE", "single")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("InitConfig() error = %v", err)
		}
		if !cfg.IsStdio() {
			t.Fatal("expected default server mode to be stdio")
		}
		if got := cfg.ListenAddr(); got != "localhost:8080" {
			t.Fatalf("ListenAddr() = %q", got)
		}
	})

	t.Run("uses MCP_SSE_ADDR as listen fallback", func(t *testing.T) {
		t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		t.Setenv("VM_INSTANCE_TYPE", "single")
		t.Setenv("MCP_SSE_ADDR", "127.0.0.1:18080")

		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("InitConfig() error = %v", err)
		}
		if got := cfg.ListenAddr(); got != "127.0.0.1:18080" {
			t.Fatalf("ListenAddr() = %q", got)
		}
	})
}

func TestInitConfigValidationErrors(t *testing.T) {
	t.Run("missing config", func(t *testing.T) {
		if _, err := InitConfig(); err == nil {
			t.Fatal("expected error for missing VM config")
		}
	})

	t.Run("mixed legacy and multi", func(t *testing.T) {
		t.Setenv("VM_ENVIRONMENTS", "demo")
		t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		if _, err := InitConfig(); err == nil {
			t.Fatal("expected error when mixing legacy and multi config")
		}
	})

	t.Run("invalid env name", func(t *testing.T) {
		t.Setenv("VM_ENVIRONMENTS", "demo,prod-east")
		t.Setenv("VM_INSTANCE_DEMO_ENTRYPOINT", "http://demo.example.com")
		t.Setenv("VM_INSTANCE_DEMO_TYPE", "single")
		if _, err := InitConfig(); err == nil {
			t.Fatal("expected invalid env name error")
		}
	})

	t.Run("invalid server mode", func(t *testing.T) {
		t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		t.Setenv("VM_INSTANCE_TYPE", "single")
		t.Setenv("MCP_SERVER_MODE", "invalid")
		if _, err := InitConfig(); err == nil {
			t.Fatal("expected invalid server mode error")
		}
	})

	t.Run("missing legacy instance type", func(t *testing.T) {
		t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		if _, err := InitConfig(); err == nil {
			t.Fatal("expected missing legacy instance type error")
		}
	})

	t.Run("invalid legacy instance type", func(t *testing.T) {
		t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		t.Setenv("VM_INSTANCE_TYPE", "invalid")
		if _, err := InitConfig(); err == nil {
			t.Fatal("expected invalid legacy instance type error")
		}
	})

	t.Run("unknown default env", func(t *testing.T) {
		t.Setenv("VM_ENVIRONMENTS", "demo")
		t.Setenv("VM_DEFAULT_ENVIRONMENT", "prod")
		t.Setenv("VM_INSTANCE_DEMO_ENTRYPOINT", "http://demo.example.com")
		t.Setenv("VM_INSTANCE_DEMO_TYPE", "single")
		if _, err := InitConfig(); err == nil {
			t.Fatal("expected unknown default env error")
		}
	})

	t.Run("missing per env type", func(t *testing.T) {
		t.Setenv("VM_ENVIRONMENTS", "demo")
		t.Setenv("VM_INSTANCE_DEMO_ENTRYPOINT", "http://demo.example.com")
		if _, err := InitConfig(); err == nil {
			t.Fatal("expected missing per-env type error")
		}
	})

	t.Run("unknown resolved env", func(t *testing.T) {
		t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		t.Setenv("VM_INSTANCE_TYPE", "single")
		cfg, err := InitConfig()
		if err != nil {
			t.Fatalf("InitConfig() error = %v", err)
		}
		if _, err := cfg.ResolveInstance("missing"); err == nil {
			t.Fatal("expected unknown env resolution error")
		}
	})

	t.Run("invalid heartbeat interval", func(t *testing.T) {
		t.Setenv("VM_INSTANCE_ENTRYPOINT", "http://example.com")
		t.Setenv("VM_INSTANCE_TYPE", "single")
		t.Setenv("MCP_HEARTBEAT_INTERVAL", "123")
		if _, err := InitConfig(); err == nil {
			t.Fatal("expected heartbeat interval error")
		}
	})
}
