package score

import "github.com/kaihendry/ai-check-guardrails/internal/modules"

func Calculate(findings []modules.Finding) int {
	score := 100
	for _, f := range findings {
		switch f.Severity {
		case modules.SeverityCritical:
			score -= 25
		case modules.SeverityHigh:
			score -= 15
		case modules.SeverityWarn:
			score -= 5
		}
	}
	if score < 0 {
		return 0
	}
	return score
}
