package mcp

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

func (m *mod) Name() string                   { return "mcp" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.MCP }

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	var findings []modules.Finding

	if len(cfg.Allowlist.MCPs) == 0 {
		findings = append(findings, modules.Finding{
			Type:        "UNCONFIGURED_ALLOWLIST",
			Severity:    modules.SeverityWarn,
			Module:      m.Name(),
			Resource:    "allowlist.mcps",
			Description: "MCP allowlist is empty; all MCPs will be treated as unapproved.",
			Remediation: "Populate allowlist.mcps in your config with approved MCP identifiers.",
		})
	}

	active, err := readActiveMCPs()
	if err != nil {
		return findings, nil
	}

	allowed := make(map[string]bool, len(cfg.Allowlist.MCPs))
	for _, a := range cfg.Allowlist.MCPs {
		allowed[a] = true
	}

	for _, id := range active {
		if !allowed[id] {
			sev := modules.SeverityHigh
			if len(cfg.Allowlist.MCPs) == 0 {
				sev = modules.SeverityWarn
			}
			findings = append(findings, modules.Finding{
				Type:        "UNAPPROVED_MCP",
				Severity:    sev,
				Module:      m.Name(),
				Resource:    id,
				Description: fmt.Sprintf("MCP %q is not on the approved allowlist.", id),
				Remediation: "Contact the security team to request allowlist approval or remove the MCP.",
			})
		}
	}
	return findings, nil
}

type claudeDesktopConfig struct {
	MCPServers map[string]any `json:"mcpServers"`
}

func readActiveMCPs() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	candidates := []string{
		filepath.Join(home, ".claude", "claude_desktop_config.json"),
		filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json"),
	}
	for _, p := range candidates {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var cfg claudeDesktopConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			continue
		}
		var ids []string
		for k := range cfg.MCPServers {
			ids = append(ids, "mcp://"+k)
		}
		return ids, nil
	}
	return nil, nil
}
