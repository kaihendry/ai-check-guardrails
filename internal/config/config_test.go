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
	t.Log("Loading an empty config should apply all defaults: mode=monitor, threshold=70, settings enabled, tokens disabled")
	cfg, err := Load(writeConfig(t, `{}`))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("mode=%q scoreThreshold=%d modules.settings=%v modules.tokens=%v",
		cfg.Mode, cfg.ScoreThreshold, cfg.Modules.Settings, cfg.Modules.Tokens)
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
	t.Log("A plaintext HTTP siem_endpoint must be rejected to prevent credential leakage in transit")
	_, err := Load(writeConfig(t, `{"siem_endpoint":"http://bad.example.com"}`))
	t.Logf("Load error: %v", err)
	if err == nil {
		t.Error("expected error for HTTP siem_endpoint")
	}
}

func TestLoad_RejectsRelativeScanRoot(t *testing.T) {
	t.Log("A relative scan_root path must be rejected; only absolute paths are unambiguous across working directories")
	_, err := Load(writeConfig(t, `{"scan_root":"relative/path"}`))
	t.Logf("Load error: %v", err)
	if err == nil {
		t.Error("expected error for relative scan_root")
	}
}

func TestLoad_AcceptsHTTPSSIEM(t *testing.T) {
	t.Log("An HTTPS siem_endpoint is valid and should load without error")
	_, err := Load(writeConfig(t, `{"siem_endpoint":"https://siem.example.com/events"}`))
	t.Logf("Load error: %v", err)
	if err != nil {
		t.Errorf("unexpected error for HTTPS endpoint: %v", err)
	}
}

func TestLoad_TokensWithoutBaseline(t *testing.T) {
	t.Log("Enabling the tokens module without a baseline is valid — it reports usage totals without anomaly detection")
	_, err := Load(writeConfig(t, `{"modules":{"tokens":true}}`))
	t.Logf("Load error: %v", err)
	if err != nil {
		t.Errorf("expected no error for tokens without baseline, got: %v", err)
	}
}

func TestLoad_TokensWithBaseline(t *testing.T) {
	t.Log("Enabling tokens with a valid baseline (daily_mean + std_dev) should load successfully")
	_, err := Load(writeConfig(t, `{
		"modules":{"tokens":true},
		"token_baseline":{"daily_mean":50000,"std_dev":12000}
	}`))
	t.Logf("Load error: %v", err)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	t.Log("A missing config file should silently return defaults (not an error — first-run friendly)")
	cfg, err := Load("/nonexistent/path/config.json")
	t.Logf("Load error: %v, mode=%q", err, cfg.Mode)
	if err != nil {
		t.Errorf("missing file should return defaults, got: %v", err)
	}
	if cfg.Mode != ModeMonitor {
		t.Errorf("missing file should return default mode, got %q", cfg.Mode)
	}
}
