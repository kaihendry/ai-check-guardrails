package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                    { return "settings" }
func (m *mod) Enabled(cfg config.Config) bool  { return cfg.Modules.Settings }

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	var findings []modules.Finding
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	paths := []string{
		filepath.Join(home, ".claude", "settings.json"),
	}
	localPath := filepath.Join(home, ".claude", "settings.local.json")
	if _, err := os.Stat(localPath); err == nil {
		paths = append(paths, localPath)
	}

	foundAny := false
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		foundAny = true
		var raw map[string]any
		if err := json.Unmarshal(data, &raw); err != nil {
			findings = append(findings, modules.Finding{
				Type:        "SETTINGS_INVALID",
				Severity:    modules.SeverityHigh,
				Module:      m.Name(),
				Resource:    p,
				Description: fmt.Sprintf("settings file is not valid JSON: %v", err),
				Remediation: "Fix the JSON syntax in " + p,
			})
			continue
		}
		// Detect non-empty overrides as WARN findings.
		for k, v := range raw {
			if v != nil && v != false && v != "" {
				findings = append(findings, modules.Finding{
					Type:        "SETTINGS_OVERRIDE",
					Severity:    modules.SeverityWarn,
					Module:      m.Name(),
					Resource:    p,
					Description: fmt.Sprintf("settings key %q is set to a non-default value", k),
					Remediation: "Review whether this override is approved by your security policy.",
				})
			}
		}
	}

	if !foundAny {
		findings = append(findings, modules.Finding{
			Type:        "MISSING_CONFIG",
			Severity:    modules.SeverityCritical,
			Module:      m.Name(),
			Resource:    filepath.Join(home, ".claude", "settings.json"),
			Description: "Claude settings.json not found",
			Remediation: "Create ~/.claude/settings.json with required security keys.",
		})
	}
	return findings, nil
}
