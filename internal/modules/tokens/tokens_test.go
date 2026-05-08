package tokens

import (
	"testing"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
)

func TestTokens_NoBaseline(t *testing.T) {
	t.Log("Without a token baseline configured, the module cannot perform anomaly detection and should emit MODULE_UNAVAILABLE")
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Tokens: true}})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got %d findings: %+v", len(findings), findings)
	if len(findings) != 1 || findings[0].Type != "MODULE_UNAVAILABLE" {
		t.Errorf("expected MODULE_UNAVAILABLE when no baseline, got %v", findings)
	}
	t.Logf("Finding: type=%s severity=%s", findings[0].Type, findings[0].Severity)
}

func TestTokens_WithBaselineNoUsageLog(t *testing.T) {
	t.Log("With a baseline configured but no usage log, daily token estimate is 0 — below threshold, so no anomaly findings")
	m := &mod{}
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
	t.Logf("Got %d findings (expected 0)", len(findings))
	if len(findings) != 0 {
		t.Errorf("zero usage should produce no anomaly findings, got %v", findings)
	}
}
