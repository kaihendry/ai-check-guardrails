package gamification

import (
	"testing"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func TestGamification_EmitsSCOREINFO(t *testing.T) {
	m := &mod{}
	findings, err := m.Run(config.Config{
		Modules:        config.ModuleToggles{Gamification: true},
		ScoreThreshold: 70,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Type != "SCORE_INFO" {
		t.Errorf("expected SCORE_INFO, got %s", findings[0].Type)
	}
	if findings[0].Severity != modules.SeverityInfo {
		t.Errorf("expected INFO severity, got %s", findings[0].Severity)
	}
}
