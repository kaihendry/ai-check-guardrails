package permissions

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

func TestPermissions_NoSettings(t *testing.T) {
	fakeHome(t)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Permissions: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("no settings file should produce no findings, got %v", findings)
	}
}

func TestPermissions_MissingPermissionKeys(t *testing.T) {
	home := fakeHome(t)
	writeSettings(t, home, `{"model":"claude-opus-4-5"}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Permissions: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].Type != "PERMISSION_MISCONFIGURED" {
		t.Errorf("expected PERMISSION_MISCONFIGURED, got %v", findings)
	}
}

func TestPermissions_HasPermissionKey(t *testing.T) {
	home := fakeHome(t)
	writeSettings(t, home, `{"permissions":{"allow":[],"deny":["Bash"]}}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Permissions: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("settings with permissions key should produce no findings, got %v", findings)
	}
}
