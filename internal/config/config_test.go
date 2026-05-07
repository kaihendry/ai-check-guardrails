package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load(writeConfig(t, `{}`))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Mode != ModeMonitor {
		t.Errorf("default mode = %q, want %q", cfg.Mode, ModeMonitor)
	}
	if cfg.ScoreThreshold != 70 {
		t.Errorf("default threshold = %d, want 70", cfg.ScoreThreshold)
	}
	if !cfg.Modules.Settings {
		t.Error("settings module should be enabled by default")
	}
	if cfg.Modules.Tokens {
		t.Error("tokens module should be disabled by default")
	}
}

func TestLoad_RejectsHTTPSIEM(t *testing.T) {
	_, err := Load(writeConfig(t, `{"siem_endpoint":"http://bad.example.com"}`))
	if err == nil {
		t.Error("expected error for HTTP siem_endpoint")
	}
}

func TestLoad_RejectsRelativeScanRoot(t *testing.T) {
	_, err := Load(writeConfig(t, `{"scan_root":"relative/path"}`))
	if err == nil {
		t.Error("expected error for relative scan_root")
	}
}

func TestLoad_AcceptsHTTPSSIEM(t *testing.T) {
	_, err := Load(writeConfig(t, `{"siem_endpoint":"https://siem.example.com/events"}`))
	if err != nil {
		t.Errorf("unexpected error for HTTPS endpoint: %v", err)
	}
}

func TestLoad_TokensRequiresBaseline(t *testing.T) {
	_, err := Load(writeConfig(t, `{"modules":{"tokens":true}}`))
	if err == nil {
		t.Error("expected error: tokens enabled without baseline")
	}
}

func TestLoad_TokensWithBaseline(t *testing.T) {
	_, err := Load(writeConfig(t, `{
		"modules":{"tokens":true},
		"token_baseline":{"daily_mean":50000,"std_dev":12000}
	}`))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.json")
	if err != nil {
		t.Errorf("missing file should return defaults, got: %v", err)
	}
	if cfg.Mode != ModeMonitor {
		t.Errorf("missing file should return default mode, got %q", cfg.Mode)
	}
}
