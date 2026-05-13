// # mcp
//
// Audits active Model Context Protocol (MCP) servers against a configured allowlist.
// Unapproved MCP servers can expand Claude's capabilities in ways that have not been
// security-reviewed.
//
// # Findings
//
// | Type | Severity | Description | Remediation |
// |------|----------|-------------|-------------|
// | UNCONFIGURED_ALLOWLIST | WARN | The MCP allowlist is empty; all MCPs are treated as unapproved. | Populate allowlist.mcps in your config with approved MCP identifiers. |
// | UNAPPROVED_MCP | WARN | An active MCP is not on the allowlist, but no allowlist is configured. | Configure an allowlist or contact the security team. |
// | UNAPPROVED_MCP | HIGH | An active MCP is not on the allowlist and an allowlist is configured. | Contact the security team to request approval or remove the MCP. |
//
// Active MCPs are read from ~/.claude/claude_desktop_config.json or
// ~/Library/Application Support/Claude/claude_desktop_config.json.
// MCP identifiers take the form mcp://<server-name>.
//
// # Configuration
//
// | Key | Type | Default | Description |
// |-----|------|---------|-------------|
// | modules.mcp | bool | true | Enable or disable this module |
// | allowlist.mcps | []string | *(empty)* | Approved MCP identifiers (e.g. ["mcp://filesystem", "mcp://github"]) |
//
// Example config:
//
//	{
//	  "modules": { "mcp": true },
//	  "allowlist": {
//	    "mcps": ["mcp://filesystem", "mcp://github"]
//	  }
//	}
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
