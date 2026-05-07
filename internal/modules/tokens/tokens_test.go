package tokens

import (
	"testing"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
)

func TestTokens_NoBaseline(t *testing.T) {
	m := &mod{}
	findings, err := m.Run(config.Config{Modules: config.ModuleToggles{Tokens: true}})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 || findings[0].Type != "MODULE_UNAVAILABLE" {
		t.Errorf("expected MODULE_UNAVAILABLE when no baseline, got %v", findings)
	}
}

func TestTokens_WithBaselineNoUsageLog(t *testing.T) {
	// estimateDailyTokens returns 0 when no log exists, so no anomaly.
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
	if len(findings) != 0 {
		t.Errorf("zero usage should produce no anomaly findings, got %v", findings)
	}
}
