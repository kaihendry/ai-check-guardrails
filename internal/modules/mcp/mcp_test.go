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
	t.Log("When no MCP allowlist is configured, the module should emit UNCONFIGURED_ALLOWLIST to prompt the operator to set one up")
	fakeHome(t)
	m := &mod{}
	findings, err := m.Run(config.Config{
		Modules:   config.ModuleToggles{MCP: true},
		Allowlist: config.Allowlist{MCPs: nil},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	if len(findings) != 1 || findings[0].Type != "UNCONFIGURED_ALLOWLIST" {
		t.Errorf("expected UNCONFIGURED_ALLOWLIST, got %v", findings)
	}
	t.Logf("Finding: type=%s severity=%s", findings[0].Type, findings[0].Severity)
}

func TestMCP_ApprovedMCPNoFinding(t *testing.T) {
	t.Log("An MCP server that appears on the allowlist should not generate an UNAPPROVED_MCP finding")
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
	t.Logf("Got %d findings: %+v", len(findings), findings)
	for _, f := range findings {
		if f.Type == "UNAPPROVED_MCP" {
			t.Errorf("approved MCP should not produce UNAPPROVED_MCP finding")
		}
	}
}

func TestMCP_UnapprovedMCPFlagged(t *testing.T) {
	t.Log("An MCP server not on the allowlist should be flagged as HIGH UNAPPROVED_MCP")
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
	t.Logf("Got %d findings: %+v", len(findings), findings)
	found := false
	for _, f := range findings {
		if f.Type == "UNAPPROVED_MCP" && f.Resource == "mcp://evil-exporter" {
			found = true
			t.Logf("Finding: type=%s severity=%s resource=%s", f.Type, f.Severity, f.Resource)
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
	t.Log("When no MCP config file exists, no UNAPPROVED_MCP findings should be produced (nothing to check)")
	fakeHome(t) // empty home, no config file
	m := &mod{}
	findings, err := m.Run(config.Config{
		Modules:   config.ModuleToggles{MCP: true},
		Allowlist: config.Allowlist{MCPs: []string{"mcp://filesystem"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	for _, f := range findings {
		if f.Type == "UNAPPROVED_MCP" {
			t.Errorf("no MCP config should not produce UNAPPROVED_MCP")
		}
	}
}
