package config

import (
	"fmt"
	"maps"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v2"
)

var backendNameRe = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

const defaultCloudName = "cloud"

type backendsFile struct {
	Toolsets map[string]string  `yaml:"toolsets"`
	Cloud    *cloudFileSection  `yaml:"cloud"`
	Backends []backendFileEntry `yaml:"backends"`
}

type cloudFileSection struct {
	Name            string    `yaml:"name"`
	APIKey          string    `yaml:"api_key"`
	BaseURL         string    `yaml:"base_url"`
	BackendTypes    []Backend `yaml:"backend_types"`
	DefaultTenantID string    `yaml:"default_tenant_id"`
}

type backendFileEntry struct {
	Name            string            `yaml:"name"`
	Type            Backend           `yaml:"type"`
	Entrypoint      string            `yaml:"entrypoint"`
	InstanceType    string            `yaml:"instance_type"`
	BearerToken     string            `yaml:"bearer_token"`
	Headers         map[string]string `yaml:"headers"`
	DefaultTenantID string            `yaml:"default_tenant_id"`
}

func (c *Config) loadBackendsFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var f backendsFile
	if err := yaml.UnmarshalStrict(data, &f); err != nil {
		return fmt.Errorf("cannot parse file: %w", err)
	}
	if len(f.Backends) == 0 && f.Cloud == nil {
		return fmt.Errorf("at least one backend must be configured in backends list or cloud section")
	}
	if len(f.Toolsets) > 0 && len(c.toolsets) > 0 {
		return fmt.Errorf("toolsets section and %s* environment variables cannot be set at the same time", toolsetsEnvPrefix)
	}
	for name, mode := range f.Toolsets {
		if err := setToolsetMode(c.toolsets, strings.ToLower(name), mode); err != nil {
			return fmt.Errorf("toolsets.%s: %w", name, err)
		}
	}

	names := make(map[string]bool, len(f.Backends))
	legacy := make(map[string]string)
	for i, entry := range f.Backends {
		if !backendNameRe.MatchString(entry.Name) {
			return fmt.Errorf("backends[%d]: name %q must match %s", i, entry.Name, backendNameRe)
		}
		if names[entry.Name] {
			return fmt.Errorf("backends[%d]: duplicate name %q", i, entry.Name)
		}
		names[entry.Name] = true

		b, err := c.newBackend(entry)
		if err != nil {
			return fmt.Errorf("backends[%d] (%s): %w", i, entry.Name, err)
		}
		mergeLegacyToolNames(legacy, b.legacyToolNames)
		c.backends = append(c.backends, b)
	}

	if f.Cloud != nil {
		cloudBackends, err := c.newCloudBackends(*f.Cloud, names)
		if err != nil {
			return fmt.Errorf("cloud: %w", err)
		}
		for _, b := range cloudBackends {
			mergeLegacyToolNames(legacy, b.legacyToolNames)
		}
		c.backends = append(c.backends, cloudBackends...)
	}

	c.name = c.backends[0].name
	c.backend = c.backends[0].backend
	c.legacyToolNames = legacy
	return nil
}

func (c *Config) newCloudBackends(section cloudFileSection, names map[string]bool) ([]*Config, error) {
	if section.Name == "" {
		section.Name = defaultCloudName
	}
	if !backendNameRe.MatchString(section.Name) {
		return nil, fmt.Errorf("name %q must match %s", section.Name, backendNameRe)
	}
	if names[section.Name] {
		return nil, fmt.Errorf("name %q is already used by a backend", section.Name)
	}
	if section.APIKey == "" {
		return nil, fmt.Errorf("api_key is not set")
	}
	var types []Backend
	for _, t := range section.BackendTypes {
		backendType, err := parseBackend(string(t))
		if err != nil || t == "" {
			return nil, fmt.Errorf("backend_types: %q must be '%s', '%s' or '%s'", t, BackendMetrics, BackendLogs, BackendTraces)
		}
		if !slices.Contains(types, backendType) {
			types = append(types, backendType)
		}
	}
	if len(types) == 0 {
		return nil, fmt.Errorf("backend_types must contain at least one backend type")
	}

	if err := c.setCloudClient(section.APIKey, section.BaseURL); err != nil {
		return nil, err
	}
	backends := make([]*Config, 0, len(types))
	for _, t := range types {
		b, err := c.newCloudBackend(t, section.Name, section.DefaultTenantID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default_tenant_id for %s backend: %w", t, err)
		}
		backends = append(backends, b)
	}
	return backends, nil
}

func (c *Config) newBackend(entry backendFileEntry) (*Config, error) {
	if entry.Type == "" {
		return nil, fmt.Errorf("type is not set")
	}
	backendType, err := parseBackend(string(entry.Type))
	if err != nil {
		return nil, fmt.Errorf("type %w", err)
	}
	entry.Type = backendType
	if entry.Entrypoint == "" {
		return nil, fmt.Errorf("entrypoint is not set")
	}
	entryPointURL, err := url.Parse(entry.Entrypoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse entrypoint: %w", err)
	}
	if entry.Type == BackendMetrics {
		if entry.InstanceType != "single" && entry.InstanceType != "cluster" {
			return nil, fmt.Errorf("instance_type must be 'single' or 'cluster'")
		}
	} else if entry.InstanceType != "" {
		return nil, fmt.Errorf("instance_type is supported only for %s backends", BackendMetrics)
	}

	b := *c
	b.backends = nil
	b.vmc = nil
	b.apiKey = ""
	b.apiBaseURL = ""
	b.name = entry.Name
	b.backend = entry.Type
	b.entrypoint = entry.Entrypoint
	b.entryPointURL = entryPointURL
	b.instanceType = entry.InstanceType
	b.bearerToken = entry.BearerToken
	b.customHeaders = maps.Clone(entry.Headers)
	b.passthroughHeaders = slices.Clone(c.passthroughHeaders)
	b.defaultTenantID = "0"
	b.disabledTools, b.legacyToolNames = parseDisabledTools(entry.Type)
	if err := b.setDefaultTenantID(entry.DefaultTenantID); err != nil {
		return nil, fmt.Errorf("failed to parse default_tenant_id: %w", err)
	}
	return &b, nil
}
