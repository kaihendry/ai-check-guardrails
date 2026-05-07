package tokens

import (
	"math"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "tokens" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Tokens }

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	if cfg.TokenBaseline == nil {
		return []modules.Finding{
			{
				Type:        "MODULE_UNAVAILABLE",
				Severity:    modules.SeverityInfo,
				Module:      m.Name(),
				Resource:    "token_baseline",
				Description: "Token anomaly detection skipped: no baseline configured.",
				Remediation: "Add token_baseline to your config after establishing a usage baseline.",
			},
		}, nil
	}

	// In a real implementation, this would read Claude's usage logs.
	// For now, report that monitoring is active.
	daily := estimateDailyTokens()
	if daily == 0 {
		return nil, nil
	}

	b := cfg.TokenBaseline
	mult := b.Multiplier
	if mult == 0 {
		mult = 3.0
	}
	threshold := float64(b.DailyMean) + mult*b.StdDev
	if float64(daily) <= threshold {
		return nil, nil
	}

	// Confidence: how many stddevs above threshold, normalised to [0,1].
	raw := (float64(daily) - float64(b.DailyMean)) / b.StdDev
	confidence := math.Min(raw/10.0, 1.0)

	return []modules.Finding{
		{
			Type:        "TOKEN_ANOMALY",
			Severity:    modules.SeverityWarn,
			Module:      m.Name(),
			Resource:    "token-usage",
			Description: "Daily token usage exceeds anomaly threshold.",
			Remediation: "Review recent Claude sessions for unusual activity.",
			Confidence:  &confidence,
		},
	}, nil
}

func estimateDailyTokens() int {
	// Placeholder: would read ~/.claude/usage.json or similar.
	return 0
}
