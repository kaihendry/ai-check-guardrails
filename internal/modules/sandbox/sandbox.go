// # sandbox
//
// Checks whether sandbox or isolation settings are explicitly disabled in
// ~/.claude/settings.json. Disabling sandboxing removes restrictions on Claude's data
// access and network behaviour.
//
// # Findings
//
// | Type | Severity | Description | Remediation |
// |------|----------|-------------|-------------|
// | SANDBOX_VIOLATION | HIGH | A sandbox-related config key is explicitly set to false. | Enable sandboxing in ~/.claude/settings.json to restrict Claude's data access. |
//
// One finding is emitted per disabled key. Checked keys: sandbox, isolatedEnv,
// restrictedMode, networkAccess. A finding is only emitted when the key is present
// and explicitly set to false — absent keys do not trigger findings.
//
// # Configuration
//
// | Key | Type | Default | Description |
// |-----|------|---------|-------------|
// | modules.sandbox | bool | true | Enable or disable this module |
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
