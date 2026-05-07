package mcp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
)

func fakeHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	return home
}

func writeMCPConfig(t *testing.T, home, content string) {
	t.Helper()
	dir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "claude_desktop_config.json"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestMCP_EmptyAllowlist(t *testing.T) {
	fakeHome(t)
	m := &mod{}
	findings, err := m.Run(config.Config{
		Modules:   config.ModuleToggles{MCP: true},
		Allowlist: config.Allowlist{MCPs: nil},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].Type != "UNCONFIGURED_ALLOWLIST" {
		t.Errorf("expected UNCONFIGURED_ALLOWLIST, got %v", findings)
	}
}

func TestMCP_ApprovedMCPNoFinding(t *testing.T) {
	home := fakeHome(t)
	writeMCPConfig(t, home, `{"mcpServers":{"filesystem":{}}}`)
	m := &mod{}
	findings, err := m.Run(config.Config{
		Modules:   config.ModuleToggles{MCP: true},
		Allowlist: config.Allowlist{MCPs: []string{"mcp://filesystem"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// Only UNCONFIGURED_ALLOWLIST is absent; approved MCP produces no finding.
	for _, f := range findings {
		if f.Type == "UNAPPROVED_MCP" {
			t.Errorf("approved MCP should not produce UNAPPROVED_MCP finding")
		}
	}
}

func TestMCP_UnapprovedMCPFlagged(t *testing.T) {
	home := fakeHome(t)
	writeMCPConfig(t, home, `{"mcpServers":{"evil-exporter":{}}}`)
	m := &mod{}
	findings, err := m.Run(config.Config{
		Modules:   config.ModuleToggles{MCP: true},
		Allowlist: config.Allowlist{MCPs: []string{"mcp://filesystem"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, f := range findings {
		if f.Type == "UNAPPROVED_MCP" && f.Resource == "mcp://evil-exporter" {
			found = true
			if f.Severity != "HIGH" {
				t.Errorf("unapproved MCP severity = %s, want HIGH", f.Severity)
			}
		}
	}
	if !found {
		t.Errorf("expected UNAPPROVED_MCP for evil-exporter, got %v", findings)
	}
}

func TestMCP_NoConfigFile(t *testing.T) {
	fakeHome(t) // empty home, no config file
	m := &mod{}
	findings, err := m.Run(config.Config{
		Modules:   config.ModuleToggles{MCP: true},
		Allowlist: config.Allowlist{MCPs: []string{"mcp://filesystem"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	// No config file = no active MCPs = only possibly UNCONFIGURED_ALLOWLIST if allowlist empty.
	for _, f := range findings {
		if f.Type == "UNAPPROVED_MCP" {
			t.Errorf("no MCP config should not produce UNAPPROVED_MCP")
		}
	}
}
