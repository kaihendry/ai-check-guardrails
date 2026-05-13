package tokens

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
)

func testdataDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "testdata", "tokens")
}

// makeJSONLFile creates a temp JSONL file containing one assistant message record per entry.
// Each entry is (inputTokens, outputTokens, cacheCreation, cacheRead, timestamp).
func makeProjectDir(t *testing.T, records []struct {
	input, output, cacheCreate, cacheRead int
	ts                                    time.Time
}) string {
	t.Helper()
	projectDir := t.TempDir()
	subDir := filepath.Join(projectDir, "project")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(filepath.Join(subDir, "session.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, r := range records {
		line := fmt.Sprintf(
			`{"timestamp":%q,"message":{"usage":{"input_tokens":%d,"output_tokens":%d,"cache_creation_input_tokens":%d,"cache_read_input_tokens":%d}}}`,
			r.ts.UTC().Format(time.RFC3339), r.input, r.output, r.cacheCreate, r.cacheRead,
		)
		if _, err := fmt.Fprintln(f, line); err != nil {
			t.Fatal(err)
		}
	}
	return projectDir
}

func TestTokens_NoBaseline(t *testing.T) {
	t.Log("Without a token baseline, module should still emit TOKEN_USAGE INFO with 0 tokens (empty dir)")
	m := &mod{projectsDirOverride: t.TempDir()}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Tokens: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	if len(findings) != 1 || findings[0].Type != "TOKEN_USAGE" {
		t.Errorf("expected TOKEN_USAGE INFO when no baseline, got %v", findings)
	}
	t.Logf("Finding: type=%s severity=%s desc=%s", findings[0].Type, findings[0].Severity, findings[0].Description)
}

func TestTokens_WithBaselineNoUsageLog(t *testing.T) {
	t.Log("With a baseline configured but no transcripts, only TOKEN_USAGE INFO should appear (no anomaly)")
	m := &mod{projectsDirOverride: t.TempDir()} // empty dir → zero tokens
	findings, err := m.Run(config.Config{
		Modules: config.ModuleToggles{Tokens: true},
		TokenBaseline: &config.TokenBaseline{
			DailyMean:  50000,
			StdDev:     12000,
			Multiplier: 3.0,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings", len(findings))
	if len(findings) != 1 || findings[0].Type != "TOKEN_USAGE" {
		t.Errorf("expected only TOKEN_USAGE when zero usage, got %v", findings)
	}
}

func TestTokens_AboveThreshold(t *testing.T) {
	t.Log("High recent usage should emit TOKEN_USAGE INFO + TOKEN_ANOMALY WARN")
	now := time.Now()
	dir := makeProjectDir(t, []struct {
		input, output, cacheCreate, cacheRead int
		ts                                    time.Time
	}{
		{50000, 30000, 60000, 80000, now.Add(-1 * time.Hour)},
		{40000, 20000, 50000, 70000, now.Add(-2 * time.Hour)},
	})

	m := &mod{projectsDirOverride: dir}
	findings, err := m.Run(config.Config{
		Modules: config.ModuleToggles{Tokens: true},
		TokenBaseline: &config.TokenBaseline{
			DailyMean:  1000,
			StdDev:     100,
			Multiplier: 3.0,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings", len(findings))
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings (TOKEN_USAGE + TOKEN_ANOMALY), got %d: %v", len(findings), findings)
	}
	if findings[0].Type != "TOKEN_USAGE" {
		t.Errorf("expected first finding TOKEN_USAGE, got %s", findings[0].Type)
	}
	if findings[1].Type != "TOKEN_ANOMALY" || findings[1].Severity != "WARN" {
		t.Errorf("expected second finding TOKEN_ANOMALY WARN, got %s %s", findings[1].Type, findings[1].Severity)
	}
	t.Logf("Usage: %s", findings[0].Description)
}

func TestTokens_BelowThreshold(t *testing.T) {
	t.Log("Low recent usage should emit TOKEN_USAGE INFO only, no TOKEN_ANOMALY")
	now := time.Now()
	dir := makeProjectDir(t, []struct {
		input, output, cacheCreate, cacheRead int
		ts                                    time.Time
	}{
		{100, 150, 100, 150, now.Add(-30 * time.Minute)},
	})

	m := &mod{projectsDirOverride: dir}
	findings, err := m.Run(config.Config{
		Modules: config.ModuleToggles{Tokens: true},
		TokenBaseline: &config.TokenBaseline{
			DailyMean:  50000,
			StdDev:     12000,
			Multiplier: 3.0,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings", len(findings))
	if len(findings) != 1 || findings[0].Type != "TOKEN_USAGE" {
		t.Errorf("expected only TOKEN_USAGE for low usage, got %v", findings)
	}
	t.Logf("Usage: %s", findings[0].Description)
}

func TestTokens_NoTranscriptsDir(t *testing.T) {
	t.Log("Non-existent projects dir should return zero stats, nil")
	stats, err := collectStats("/nonexistent/does/not/exist/"+t.Name(), 24)
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got: %v", err)
	}
	if stats.TotalTokens != 0 {
		t.Errorf("expected 0 tokens for missing dir, got %d", stats.TotalTokens)
	}
}

func TestTokens_MalformedLines(t *testing.T) {
	t.Log("Malformed JSONL lines should be skipped; valid lines counted")
	projectDir := t.TempDir()
	subDir := filepath.Join(projectDir, "myproject")
	if err := os.Mkdir(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	content := fmt.Sprintf(`not valid json
{"message":{"usage":{"input_tokens":100,"output_tokens":50}},"timestamp":%q}
{broken
{"message":{"usage":{"input_tokens":200,"output_tokens":100}},"timestamp":%q}
`,
		now.Add(-1*time.Hour).UTC().Format(time.RFC3339),
		now.Add(-2*time.Hour).UTC().Format(time.RFC3339),
	)
	if err := os.WriteFile(filepath.Join(subDir, "session.jsonl"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	stats, err := collectStats(projectDir, 24)
	if err != nil {
		t.Fatal(err)
	}
	// 2 valid records: (100+50) + (200+100) = 450
	if stats.TotalTokens != 450 {
		t.Errorf("expected 450 tokens from valid lines, got %d", stats.TotalTokens)
	}
}

func TestTokens_LookbackExcludesOldSessions(t *testing.T) {
	t.Log("Sessions older than lookback window should be excluded")
	now := time.Now()
	dir := makeProjectDir(t, []struct {
		input, output, cacheCreate, cacheRead int
		ts                                    time.Time
	}{
		{200000, 100000, 150000, 200000, now.Add(-48 * time.Hour)},
	})

	stats, err := collectStats(dir, 1)
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalTokens != 0 {
		t.Errorf("expected 0 tokens with 1h lookback for 48h-old session, got %d", stats.TotalTokens)
	}
}

func TestTokens_LookbackIncludesRecentSessions(t *testing.T) {
	t.Log("Sessions within lookback window should be included")
	now := time.Now()
	dir := makeProjectDir(t, []struct {
		input, output, cacheCreate, cacheRead int
		ts                                    time.Time
	}{
		{1000, 500, 800, 1200, now.Add(-30 * time.Minute)},
	})

	stats, err := collectStats(dir, 24)
	if err != nil {
		t.Fatal(err)
	}
	if stats.TotalTokens == 0 {
		t.Error("expected non-zero tokens with 24h lookback for 30-min-old session")
	}
	t.Logf("Recent session tokens: %d", stats.TotalTokens)
}
