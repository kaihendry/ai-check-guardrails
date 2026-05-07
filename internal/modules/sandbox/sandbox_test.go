package sandbox

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

func writeSettings(t *testing.T, home, content string) {
	t.Helper()
	dir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "settings.json"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSandbox_NoSettings(t *testing.T) {
	fakeHome(t)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Sandbox: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("no settings file should produce no findings, got %v", findings)
	}
}

func TestSandbox_ExplicitlyDisabled(t *testing.T) {
	home := fakeHome(t)
	writeSettings(t, home, `{"sandbox":false}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Sandbox: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].Type != "SANDBOX_VIOLATION" {
		t.Errorf("expected SANDBOX_VIOLATION, got %v", findings)
	}
	if findings[0].Severity != "HIGH" {
		t.Errorf("severity = %s, want HIGH", findings[0].Severity)
	}
}

func TestSandbox_Enabled(t *testing.T) {
	home := fakeHome(t)
	writeSettings(t, home, `{"sandbox":true}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Sandbox: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("enabled sandbox should produce no findings, got %v", findings)
	}
}
