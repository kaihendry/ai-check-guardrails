package settings

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

func writeSettingsJSON(t *testing.T, home, content string) {
	t.Helper()
	dir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "settings.json"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestSettings_MissingFile(t *testing.T) {
	fakeHome(t) // empty home, no settings.json
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Settings: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].Type != "MISSING_CONFIG" {
		t.Errorf("expected MISSING_CONFIG, got %v", findings)
	}
	if findings[0].Severity != "CRITICAL" {
		t.Errorf("expected CRITICAL severity, got %s", findings[0].Severity)
	}
}

func TestSettings_EmptyFile(t *testing.T) {
	home := fakeHome(t)
	writeSettingsJSON(t, home, `{}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Settings: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("empty settings should produce no findings, got %v", findings)
	}
}

func TestSettings_OverridesDetected(t *testing.T) {
	home := fakeHome(t)
	writeSettingsJSON(t, home, `{"model":"claude-opus-4-5","theme":"dark"}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Settings: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 2 {
		t.Errorf("expected 2 SETTINGS_OVERRIDE findings, got %d: %v", len(findings), findings)
	}
	for _, f := range findings {
		if f.Type != "SETTINGS_OVERRIDE" || f.Severity != "WARN" {
			t.Errorf("unexpected finding: %+v", f)
		}
	}
}

func TestSettings_InvalidJSON(t *testing.T) {
	home := fakeHome(t)
	writeSettingsJSON(t, home, `{not valid json`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Settings: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].Type != "SETTINGS_INVALID" {
		t.Errorf("expected SETTINGS_INVALID finding, got %v", findings)
	}
}
