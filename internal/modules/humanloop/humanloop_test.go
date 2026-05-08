package humanloop

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

func TestHumanLoop_KeyAbsent(t *testing.T) {
	t.Log("Settings without requireConfirmation key should produce a WARN HUMANLOOP_ABSENT finding")
	home := fakeHome(t)
	writeSettings(t, home, `{"model":"claude-opus-4-5"}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{HumanLoop: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	if len(findings) != 1 || findings[0].Type != "HUMANLOOP_ABSENT" {
		t.Errorf("expected HUMANLOOP_ABSENT, got %v", findings)
	}
	t.Logf("Finding: type=%s severity=%s", findings[0].Type, findings[0].Severity)
}

func TestHumanLoop_KeyPresent(t *testing.T) {
	t.Log("Settings with requireConfirmation:true indicate human oversight is active — no findings expected")
	home := fakeHome(t)
	writeSettings(t, home, `{"requireConfirmation":true}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{HumanLoop: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings (expected 0)", len(findings))
	if len(findings) != 0 {
		t.Errorf("requireConfirmation present should produce no findings, got %v", findings)
	}
}

func TestHumanLoop_NoSettingsFile(t *testing.T) {
	t.Log("When no settings file exists, the module should skip gracefully with no findings")
	fakeHome(t)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{HumanLoop: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings (expected 0)", len(findings))
	if len(findings) != 0 {
		t.Errorf("no settings file should produce no findings (module skips gracefully), got %v", findings)
	}
}
