package sandbox

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "sandbox" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Sandbox }

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	var findings []modules.Finding
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

	// Look for sandbox-related config keys.
	sandboxKeys := []string{"sandbox", "isolatedEnv", "restrictedMode", "networkAccess"}
	for _, k := range sandboxKeys {
		v, ok := raw[k]
		if !ok {
			continue
		}
		if b, isBool := v.(bool); isBool && !b {
			findings = append(findings, modules.Finding{
				Type:        "SANDBOX_VIOLATION",
				Severity:    modules.SeverityHigh,
				Module:      m.Name(),
				Resource:    settingsPath + "#" + k,
				Description: "Sandbox config key '" + k + "' is explicitly disabled.",
				Remediation: "Enable sandboxing in ~/.claude/settings.json to restrict Claude's data access.",
			})
		}
	}
	return findings, nil
}
