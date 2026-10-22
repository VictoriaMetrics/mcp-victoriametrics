package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/registry"
)

const (
	toolsStartMarker = "<!-- tools:start -->"
	toolsEndMarker   = "<!-- tools:end -->"
)

type mode struct {
	name string
	env  []string
}

type backend struct {
	id      string
	product string
	modes   []mode
}

var backends = []backend{
	{
		id:      "metrics",
		product: "VictoriaMetrics",
		modes: []mode{
			{name: "single", env: []string{"VM_INSTANCE_ENTRYPOINT=http://localhost:8428", "VM_INSTANCE_TYPE=single"}},
			{name: "cluster", env: []string{"VM_INSTANCE_ENTRYPOINT=http://localhost:8481", "VM_INSTANCE_TYPE=cluster"}},
		},
	},
	{
		id:      "logs",
		product: "VictoriaLogs",
		modes:   []mode{{name: "single", env: []string{"VM_INSTANCE_ENTRYPOINT=http://localhost:9428"}}},
	},
	{
		id:      "traces",
		product: "VictoriaTraces",
		modes:   []mode{{name: "single", env: []string{"VM_INSTANCE_ENTRYPOINT=http://localhost:10428"}}},
	},
}

type toolInfo struct {
	name           string
	title          string
	enabledDefault bool
	modes          []string
	cloudOnly      bool
}

type backendInfo struct {
	backend backend
	tools   []*toolInfo
	prompts []mcp.Prompt
}

// collectBackend registers tools and prompts of a backend in every mode and in cloud mode
// with all tools enabled, and takes the tools enabled by default from the config.
func collectBackend(b backend) (*backendInfo, error) {
	tools := map[string]*toolInfo{}
	var prompts []mcp.Prompt
	collect := func(modeName string, cloud bool, env []string) error {
		s, err := newServer(b, append(env, "MCP_DISABLED_TOOLS="))
		if err != nil {
			return err
		}
		for name, st := range s.ListTools() {
			ti, ok := tools[name]
			if !ok {
				ti = &toolInfo{name: name, title: strings.TrimSpace(st.Tool.Annotations.Title), cloudOnly: true}
				if ti.title == "" {
					ti.title = strings.TrimSpace(strings.SplitN(st.Tool.Description, ".", 2)[0])
				}
				tools[name] = ti
			}
			if !cloud {
				ti.cloudOnly = false
				ti.modes = append(ti.modes, modeName)
			}
		}
		for name, sp := range s.ListPrompts() {
			if !slices.ContainsFunc(prompts, func(p mcp.Prompt) bool { return p.Name == name }) {
				prompts = append(prompts, sp.Prompt)
			}
		}
		return nil
	}
	for _, m := range b.modes {
		if err := collect(m.name, false, m.env); err != nil {
			return nil, fmt.Errorf("%s/%s: %w", b.id, m.name, err)
		}
	}
	if err := collect("", true, []string{"VMC_API_KEY=readmegen"}); err != nil {
		return nil, fmt.Errorf("%s/cloud: %w", b.id, err)
	}

	cfg, err := newConfig(b, b.modes[len(b.modes)-1].env)
	if err != nil {
		return nil, err
	}
	info := &backendInfo{backend: b, prompts: prompts}
	for _, t := range tools {
		t.enabledDefault = !cfg.IsToolDisabled(t.name)
		info.tools = append(info.tools, t)
	}
	slices.SortFunc(info.tools, func(a, b *toolInfo) int { return strings.Compare(a.name, b.name) })
	slices.SortFunc(info.prompts, func(a, b mcp.Prompt) int { return strings.Compare(a.Name, b.Name) })
	return info, nil
}

func newServer(b backend, env []string) (*server.MCPServer, error) {
	cfg, err := newConfig(b, env)
	if err != nil {
		return nil, err
	}
	s := server.NewMCPServer("readmegen", "0", server.WithToolCapabilities(false), server.WithPromptCapabilities(false))
	registry.Register(s, cfg)
	return s, nil
}

// newConfig initializes the server config of the backend from the given environment only.
func newConfig(b backend, env []string) (*config.Config, error) {
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		for _, prefix := range []string{"VM_", "VMC_", "VL_", "VT_", "MCP_"} {
			if strings.HasPrefix(name, prefix) {
				_ = os.Unsetenv(name)
			}
		}
	}
	for _, kv := range append(slices.Clone(env), "VM_INSTANCE_BACKEND="+b.id) {
		name, value, _ := strings.Cut(kv, "=")
		_ = os.Setenv(name, value)
	}
	cfg, err := config.InitConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot init %s config: %w", b.id, err)
	}
	return cfg, nil
}

func (ti *toolInfo) description(b backend) string {
	if len(ti.modes) == 0 || len(ti.modes) == len(b.modes) {
		return ti.title
	}
	return fmt.Sprintf("%s (%s mode only)", ti.title, strings.Join(ti.modes, ", "))
}

type table struct {
	Header []string
	Rows   [][]string
}

type backendSection struct {
	Product string
	ID      string
	Default bool
	Tools   table
	Prompts table
}

// renderTools renders a section per backend with its tools and prompts,
// and a single section with additional tools for VictoriaMetrics Cloud.
func renderTools(infos []*backendInfo) (string, error) {
	var sections []backendSection
	cloudRows := map[string][]string{}
	var cloudNames []string
	for i, info := range infos {
		section := backendSection{
			Product: info.backend.product,
			ID:      info.backend.id,
			Default: i == 0,
			Tools:   table{Header: []string{"Tool", "Description", "Enabled by default"}},
			Prompts: table{Header: []string{"Prompt", "Description"}},
		}
		for _, t := range info.tools {
			if !t.cloudOnly {
				section.Tools.Rows = append(section.Tools.Rows, []string{"`" + t.name + "`", t.description(info.backend), enabledMark(t.enabledDefault)})
				continue
			}
			row, ok := cloudRows[t.name]
			if !ok {
				row = []string{"`" + t.name + "`", t.title, "", enabledMark(t.enabledDefault)}
				cloudNames = append(cloudNames, t.name)
			}
			row[2] = strings.TrimPrefix(row[2]+", `"+info.backend.id+"`", ", ")
			cloudRows[t.name] = row
		}
		for _, p := range info.prompts {
			section.Prompts.Rows = append(section.Prompts.Rows, []string{"`" + p.Name + "`", p.Description})
		}
		sections = append(sections, section)
	}

	var cloud *table
	if len(cloudNames) > 0 {
		slices.Sort(cloudNames)
		allBackends := make([]string, 0, len(infos))
		for _, info := range infos {
			allBackends = append(allBackends, "`"+info.backend.id+"`")
		}
		cloud = &table{Header: []string{"Tool", "Description", "Backends", "Enabled by default"}}
		sameBackends := true
		for _, name := range cloudNames {
			cloud.Rows = append(cloud.Rows, cloudRows[name])
			sameBackends = sameBackends && cloudRows[name][2] == strings.Join(allBackends, ", ")
		}
		// the Backends column carries no information when every tool supports every backend
		if sameBackends {
			cloud.Header = slices.Delete(cloud.Header, 2, 3)
			for i := range cloud.Rows {
				cloud.Rows[i] = slices.Delete(cloud.Rows[i], 2, 3)
			}
		}
	}
	return renderTemplate("tools.md.tmpl", map[string]any{
		"Backends": sections,
		"Multiple": len(sections) > 1,
		"Cloud":    cloud,
	})
}

func enabledMark(enabled bool) string {
	if enabled {
		return "✅"
	}
	return "❌"
}

func (t table) String() string {
	header, rows := t.Header, t.Rows
	centered := func(i int) bool { return header[i] == "Enabled by default" }
	widths := make([]int, len(header))
	for _, row := range append([][]string{header}, rows...) {
		for i, cell := range row {
			widths[i] = max(widths[i], displayWidth(cell))
		}
	}
	var b strings.Builder
	writeRow := func(row []string) {
		for i, cell := range row {
			pad := widths[i] - displayWidth(cell)
			left := 0
			if centered(i) {
				left = pad / 2
			}
			fmt.Fprintf(&b, "| %s%s%s ", strings.Repeat(" ", left), cell, strings.Repeat(" ", pad-left))
		}
		b.WriteString("|\n")
	}
	writeRow(header)
	for i := range header {
		if centered(i) {
			b.WriteString("|:" + strings.Repeat("-", widths[i]) + ":")
			continue
		}
		b.WriteString("|" + strings.Repeat("-", widths[i]+2))
	}
	b.WriteString("|\n")
	for _, row := range rows {
		writeRow(row)
	}
	return strings.TrimSuffix(b.String(), "\n")
}

// displayWidth returns the number of columns the text takes in a monospace font,
// counting emoji such as ✅ and ❌ as two columns.
func displayWidth(text string) int {
	width := 0
	for _, r := range text {
		if r == '✅' || r == '❌' || r >= 0x1F300 {
			width += 2
			continue
		}
		width++
	}
	return width
}

func replaceBetween(content, start, end, replacement string) (string, error) {
	i := strings.Index(content, start)
	j := strings.Index(content, end)
	if i < 0 || j < i {
		return "", fmt.Errorf("markers %s and %s are not found", start, end)
	}
	return content[:i] + replacement + content[j+len(end):], nil
}
