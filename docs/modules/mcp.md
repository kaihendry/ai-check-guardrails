# mcp

Audits active Model Context Protocol (MCP) servers against a configured allowlist. Unapproved MCP servers can expand Claude's capabilities in ways that have not been security-reviewed.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `UNCONFIGURED_ALLOWLIST` | WARN | The MCP allowlist is empty; all MCPs are treated as unapproved. | Populate `allowlist.mcps` in your config with approved MCP identifiers. |
| `UNAPPROVED_MCP` | WARN | An active MCP is not on the allowlist, but no allowlist is configured. | Configure an allowlist or contact the security team. |
| `UNAPPROVED_MCP` | HIGH | An active MCP is not on the allowlist and an allowlist is configured. | Contact the security team to request approval or remove the MCP. |

Active MCPs are read from `~/.claude/claude_desktop_config.json` or `~/Library/Application Support/Claude/claude_desktop_config.json`. MCP identifiers take the form `mcp://<server-name>`.

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.mcp` | bool | `true` | Enable or disable this module |
| `allowlist.mcps` | `[]string` | *(empty)* | Approved MCP identifiers (e.g. `["mcp://filesystem", "mcp://github"]`) |

**Example config**:

```json
{
  "modules": { "mcp": true },
  "allowlist": {
    "mcps": ["mcp://filesystem", "mcp://github"]
  }
}
```
