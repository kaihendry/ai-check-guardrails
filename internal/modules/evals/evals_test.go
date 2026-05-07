package evals

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

func TestEvals_NoHooksSection(t *testing.T) {
	home := fakeHome(t)
	writeSettings(t, home, `{"model":"claude-opus-4-5"}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Evals: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != len(recommendedHooks) {
		t.Errorf("expected %d EVAL_HOOK_ABSENT findings, got %d", len(recommendedHooks), len(findings))
	}
	for _, f := range findings {
		if f.Type != "EVAL_HOOK_ABSENT" || f.Severity != "HIGH" {
			t.Errorf("unexpected finding: %+v", f)
		}
	}
}

func TestEvals_AllHooksPresent(t *testing.T) {
	home := fakeHome(t)
	writeSettings(t, home, `{"hooks":{"PreToolUse":[],"PostToolUse":[],"Stop":[]}}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Evals: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Errorf("all hooks present should produce no findings, got %v", findings)
	}
}

func TestEvals_SomeHooksMissing(t *testing.T) {
	home := fakeHome(t)
	writeSettings(t, home, `{"hooks":{"PreToolUse":[]}}`)
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Evals: true}})
	if err != nil {
		t.Fatal(err)
	}
	// PostToolUse and Stop are missing
	if len(findings) != 2 {
		t.Errorf("expected 2 missing hook findings, got %d: %v", len(findings), findings)
	}
}
