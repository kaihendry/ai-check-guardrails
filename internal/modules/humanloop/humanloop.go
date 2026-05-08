package humanloop

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "humanloop" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.HumanLoop }

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	settingsPath := filepath.Join(home, ".claude", "settings.json")
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, nil
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil
	}

	humanLoopKeys := []string{"requireConfirmation", "humanInTheLoop", "confirmDangerousActions"}
	for _, k := range humanLoopKeys {
		if _, ok := raw[k]; ok {
			return nil, nil
		}
	}
	return []modules.Finding{
		{
			Type:        "HUMANLOOP_ABSENT",
			Severity:    modules.SeverityWarn,
			Module:      m.Name(),
			Resource:    settingsPath,
			Description: "No human-in-the-loop confirmation setting found in settings.json.",
			Remediation: "Add 'requireConfirmation: true' to ~/.claude/settings.json.",
		},
	}, nil
}
