package network

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

func writeNetworkLog(t *testing.T, home, content string) {
	t.Helper()
	dir := filepath.Join(home, ".claude", "logs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "network.log"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestNetwork_NoLog(t *testing.T) {
	t.Log("When no network log file exists, the module should produce no findings")
	fakeHome(t)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Network: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings (expected 0)", len(findings))
	if len(findings) != 0 {
		t.Errorf("no log file should produce no findings, got %v", findings)
	}
}

func TestNetwork_LoggedDomainInfo(t *testing.T) {
	t.Log("When a domain is in the log but no allowlist is configured, it should emit INFO (not WARN)")
	home := fakeHome(t)
	writeNetworkLog(t, home, "GET https://api.anthropic.com/v1/messages\n")
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Network: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	if len(findings) != 1 || findings[0].Type != "NETWORK_REQUEST" {
		t.Errorf("expected NETWORK_REQUEST finding, got %v", findings)
	}
	if findings[0].Severity != "INFO" {
		t.Errorf("no allowlist configured should give INFO severity, got %s", findings[0].Severity)
	}
	t.Logf("Finding: type=%s severity=%s resource=%s", findings[0].Type, findings[0].Severity, findings[0].Resource)
}

func TestNetwork_UnallowlistedDomainWarn(t *testing.T) {
	t.Log("When a domain is NOT on the allowlist, severity should be upgraded to WARN")
	home := fakeHome(t)
	writeNetworkLog(t, home, "GET https://suspicious.example.com/data\n")
	m := &mod{}
	findings, err := m.Run(config.Config{
		Modules:   config.ModuleToggles{Network: true},
		Allowlist: config.Allowlist{Domains: []string{"api.anthropic.com"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	if len(findings) != 1 || findings[0].Severity != "WARN" {
		t.Errorf("non-allowlisted domain should be WARN, got %v", findings)
	}
	t.Logf("Finding: type=%s severity=%s resource=%s", findings[0].Type, findings[0].Severity, findings[0].Resource)
}

func TestNetwork_AllowlistedDomainInfo(t *testing.T) {
	t.Log("When a domain IS on the allowlist, it should stay INFO (not escalated)")
	home := fakeHome(t)
	writeNetworkLog(t, home, "GET https://api.anthropic.com/v1/messages\n")
	m := &mod{}
	findings, err := m.Run(config.Config{
		Modules:   config.ModuleToggles{Network: true},
		Allowlist: config.Allowlist{Domains: []string{"api.anthropic.com"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	if len(findings) != 1 || findings[0].Severity != "INFO" {
		t.Errorf("allowlisted domain should stay INFO, got %v", findings)
	}
	t.Logf("Finding: type=%s severity=%s resource=%s", findings[0].Type, findings[0].Severity, findings[0].Resource)
}
