package gamification

import (
	"fmt"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "gamification" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Gamification }

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	return []modules.Finding{
		{
			Type:        "SCORE_INFO",
			Severity:    modules.SeverityInfo,
			Module:      m.Name(),
			Resource:    "posture-score",
			Description: fmt.Sprintf("Security posture score will be calculated from findings in this run. Target: ≥%d/100.", cfg.ScoreThreshold),
			Remediation: "Resolve HIGH and CRITICAL findings to improve your score.",
		},
	}, nil
}
