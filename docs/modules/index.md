# Modules

Each module performs a focused audit check. Modules are run in sequence; each produces zero or more [findings](../reference/module-interface.md).

| Module | Purpose | Default State |
|--------|---------|---------------|
| [banner](banner.md) | Displays a security score summary banner on stderr after all findings are collected | Enabled |
| [bypass](bypass.md) | Detects permission-bypass flags in shell history | Enabled |
| [envkeys](envkeys.md) | Detects API keys set as plaintext environment variables | Enabled |
| [evals](evals.md) | Checks for Anthropic-recommended eval hooks in Claude settings | Disabled |
| [gamification](gamification.md) | Emits an INFO finding anchoring the security posture score in output | Enabled |
| [hooks](hooks.md) | Checks that required pre-commit hooks are installed in Git repositories | Enabled |
| [humanloop](humanloop.md) | Checks for human-in-the-loop confirmation settings in Claude settings | Disabled |
| [mcp](mcp.md) | Audits active MCP servers against a configured allowlist | Enabled |
| [network](network.md) | Scans Claude network logs for outbound requests to unapproved domains | Disabled |
| [permissions](permissions.md) | Checks for tool permission restrictions in Claude settings | Enabled |
| [sandbox](sandbox.md) | Checks that sandbox/isolation settings are not explicitly disabled | Enabled |
| [settings](settings.md) | Validates that Claude settings files exist and contain valid JSON | Enabled |
| [tokens](tokens.md) | Detects anomalous token usage against a configured baseline | Disabled |

## Enabling and Disabling Modules

All module toggles live under `modules` in the [configuration file](../reference/config.md):

```json
{
  "modules": {
    "evals": true,
    "network": true,
    "tokens": false
  }
}
```

See the [Configuration Reference](../reference/config.md) for the full list of keys.
