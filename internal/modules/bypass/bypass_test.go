package bypass

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

func TestBypass_CleanHistory(t *testing.T) {
	home := fakeHome(t)
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte("git status\ngit log\n"), 0600); err != nil {
		t.Fatal(err)
	}
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Bypass: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("clean history should produce no findings, got %v", findings)
	}
}

func TestBypass_DangerousFlagInHistory(t *testing.T) {
	home := fakeHome(t)
	content := "git status\nclaude --dangerously-skip-permissions\ngit log\n"
	if err := os.WriteFile(filepath.Join(home, ".zsh_history"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Bypass: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Errorf("expected 1 POLICY_BYPASS finding, got %d: %v", len(findings), findings)
	}
	if findings[0].Type != "POLICY_BYPASS" || findings[0].Severity != "CRITICAL" {
		t.Errorf("unexpected finding: %+v", findings[0])
	}
}

func TestBypass_NoHistoryFile(t *testing.T) {
	fakeHome(t) // empty home
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Bypass: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("missing history file should produce no findings, got %v", findings)
	}
}
