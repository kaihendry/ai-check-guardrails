package permissions

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "permissions" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Permissions }

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

	// Check for permission-related keys.
	permKeys := []string{"permissions", "allowedTools", "blockedTools", "allow", "deny"}
	found := false
	for _, k := range permKeys {
		if _, ok := raw[k]; ok {
			found = true
			break
		}
	}
	if !found {
		findings = append(findings, modules.Finding{
			Type:        "PERMISSION_MISCONFIGURED",
			Severity:    modules.SeverityHigh,
			Module:      m.Name(),
			Resource:    settingsPath,
			Description: "No permission/allowedTools configuration found in settings.json.",
			Remediation: "Add permission restrictions to ~/.claude/settings.json to limit tool access.",
		})
	}
	return findings, nil
}
