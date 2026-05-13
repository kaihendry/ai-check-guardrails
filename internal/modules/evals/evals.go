// # evals
//
// Checks whether the Anthropic-recommended evaluation hooks are configured in
// ~/.claude/settings.json. These hooks enable adversarial review and private data
// protection by intercepting tool calls at key lifecycle events.
//
// # Findings
//
// | Type | Severity | Description | Remediation |
// |------|----------|-------------|-------------|
// | EVAL_HOOK_ABSENT | HIGH | An Anthropic-recommended eval hook is not configured. | Add the missing hook to the hooks section in ~/.claude/settings.json. |
//
// One finding is emitted for each missing hook. The three recommended hooks are:
//
// | Hook | Purpose |
// |------|---------|
// | PreToolUse | Inspect or block tool calls before they execute |
// | PostToolUse | Audit tool outputs after execution |
// | Stop | Review model responses before they are finalized |
//
// # Configuration
//
// | Key | Type | Default | Description |
// |-----|------|---------|-------------|
// | modules.evals | bool | false | Enable or disable this module |
package evals

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "evals" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Evals }

// Anthropic-recommended hooks for adversarial review and private data protection.
var recommendedHooks = []string{
	"PreToolUse",
	"PostToolUse",
	"Stop",
}

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

	hooksRaw, ok := raw["hooks"]
	if !ok {
		for _, h := range recommendedHooks {
			findings = append(findings, modules.Finding{
				Type:        "EVAL_HOOK_ABSENT",
				Severity:    modules.SeverityHigh,
				Module:      m.Name(),
				Resource:    settingsPath,
				Description: "Anthropic-recommended eval hook '" + h + "' is not configured.",
				Remediation: "Add a '" + h + "' hook to ~/.claude/settings.json for adversarial review.",
			})
		}
		return findings, nil
	}

	hooksMap, ok := hooksRaw.(map[string]any)
	if !ok {
		return findings, nil
	}
	for _, h := range recommendedHooks {
		if _, ok := hooksMap[h]; !ok {
			findings = append(findings, modules.Finding{
				Type:        "EVAL_HOOK_ABSENT",
				Severity:    modules.SeverityHigh,
				Module:      m.Name(),
				Resource:    settingsPath + "#hooks." + h,
				Description: "Anthropic-recommended eval hook '" + h + "' is not configured.",
				Remediation: "Add a '" + h + "' hook to the hooks section in ~/.claude/settings.json.",
			})
		}
	}
	return findings, nil
}
