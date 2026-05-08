package hooks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
)

func makeRepo(t *testing.T, withHooks ...string) string {
	t.Helper()
	dir := t.TempDir()
	hooksDir := filepath.Join(dir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		t.Fatal(err)
	}
	for _, h := range withHooks {
		if err := os.WriteFile(filepath.Join(hooksDir, h), []byte("#!/bin/sh\n"), 0755); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func cfg(scanRoot string, hooks ...string) config.Config {
	return config.Config{
		Modules:   config.ModuleToggles{Hooks: true},
		ScanRoot:  scanRoot,
		Allowlist: config.Allowlist{PreCommitHooks: hooks},
	}
}

func TestHooks_MissingRequiredHook(t *testing.T) {
	t.Log("A repo without any pre-commit hooks should produce a HIGH MISSING_PRECOMMIT_HOOK finding for each required hook")
	repoDir := makeRepo(t) // no hooks installed
	m := &mod{}
	findings, err := m.Run(cfg(repoDir, "gitleaks"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	if len(findings) != 1 || findings[0].Type != "MISSING_PRECOMMIT_HOOK" {
		t.Errorf("expected MISSING_PRECOMMIT_HOOK, got %v", findings)
	}
	if findings[0].Severity != "HIGH" {
		t.Errorf("severity = %s, want HIGH", findings[0].Severity)
	}
	t.Logf("Finding: type=%s severity=%s resource=%s", findings[0].Type, findings[0].Severity, findings[0].Resource)
}

func TestHooks_HookPresent(t *testing.T) {
	t.Log("A repo with the required hook installed should produce no findings")
	repoDir := makeRepo(t, "gitleaks")
	m := &mod{}
	findings, err := m.Run(cfg(repoDir, "gitleaks"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings (expected 0)", len(findings))
	if len(findings) != 0 {
		t.Errorf("hook present should produce no findings, got %v", findings)
	}
}

func TestHooks_MultipleHooksSomeMissing(t *testing.T) {
	t.Log("When only some required hooks are installed, one finding per missing hook should be emitted")
	repoDir := makeRepo(t, "gitleaks") // detect-secrets missing
	m := &mod{}
	findings, err := m.Run(cfg(repoDir, "gitleaks", "detect-secrets"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	if len(findings) != 1 || findings[0].Type != "MISSING_PRECOMMIT_HOOK" {
		t.Errorf("expected 1 missing hook finding, got %v", findings)
	}
	t.Logf("Missing hook: %s", findings[0].Resource)
}

func TestHooks_NoRepos(t *testing.T) {
	t.Log("A scan root with no git repos should produce no findings")
	emptyDir := t.TempDir()
	m := &mod{}
	findings, err := m.Run(cfg(emptyDir, "gitleaks"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings (expected 0)", len(findings))
	if len(findings) != 0 {
		t.Errorf("no repos should produce no findings, got %v", findings)
	}
}
