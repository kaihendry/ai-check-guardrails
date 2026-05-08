package banner

import (
	"fmt"
	"os"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "banner" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Banner }

// Run emits no findings; banner display is triggered by the runner via Display.
func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	return nil, nil
}

// Display writes a score-based security banner to stderr.
// Called by the runner after all findings are collected and scored.
func Display(score int, cfg config.Config) {
	if !cfg.Modules.Banner {
		return
	}
	threshold := cfg.ScoreThreshold
	if threshold == 0 {
		threshold = 70
	}
	if score >= 100 {
		return
	}
	if score >= threshold {
		fmt.Fprintf(os.Stderr, "\n[WARN] Security score: %d/100 — review findings to improve your posture.\n\n", score)
		return
	}
	url := cfg.BannerURL
	if url == "" {
		url = "your security training portal"
	}
	fmt.Fprintf(os.Stderr, "\n[ALERT] Security score: %d/100 — below threshold (%d).\n         Action required: %s\n\n", score, threshold, url)
}
