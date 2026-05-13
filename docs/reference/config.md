# Configuration Reference

The tool reads its configuration from `~/.config/ai-check-guardrails/config.json`. All keys are optional; the defaults shown below apply when a key is absent or when no config file exists.

## Top-Level Keys

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `mode` | string | `"monitor"` | Run mode: `monitor` (findings reported, non-zero exit only on error) or `enforce` (non-zero exit when CRITICAL findings are present). |
| `siem_endpoint` | string | *(empty)* | HTTPS endpoint to forward findings to a SIEM. Must start with `https://`. Can be overridden with the `AI_GUARDRAILS_SIEM_ENDPOINT` environment variable. |
| `scan_root` | string | `$HOME` | Root directory for file-based scans (e.g., searching for Git repositories). Must be an absolute path. |
| `banner_url` | string | *(empty)* | URL shown in the ALERT banner when the score falls below the threshold (e.g., a security training portal). |
| `score_threshold` | int | `70` | Minimum passing score (0–100). Scores below this trigger an ALERT banner. |
| `token_baseline` | object | `null` | Required when the `tokens` module is enabled. See [Token Baseline](#token-baseline) below. |
| `env_key_watch_list` | `[]string` | `["ANTHROPIC_API_KEY"]` | Environment variable names the `envkeys` module watches. Overrides the default when non-empty. |

## Module Toggles (`modules`)

| Key | Type | Default | Module |
|-----|------|---------|--------|
| `modules.banner` | bool | `true` | [banner](../modules/banner.md) |
| `modules.bypass` | bool | `true` | [bypass](../modules/bypass.md) |
| `modules.env_keys` | bool | `true` | [envkeys](../modules/envkeys.md) |
| `modules.evals` | bool | `false` | [evals](../modules/evals.md) |
| `modules.gamification` | bool | `true` | [gamification](../modules/gamification.md) |
| `modules.hooks` | bool | `true` | [hooks](../modules/hooks.md) |
| `modules.humanloop` | bool | `false` | [humanloop](../modules/humanloop.md) |
| `modules.mcp` | bool | `true` | [mcp](../modules/mcp.md) |
| `modules.network` | bool | `false` | [network](../modules/network.md) |
| `modules.permissions` | bool | `true` | [permissions](../modules/permissions.md) |
| `modules.sandbox` | bool | `true` | [sandbox](../modules/sandbox.md) |
| `modules.settings` | bool | `true` | [settings](../modules/settings.md) |
| `modules.tokens` | bool | `false` | [tokens](../modules/tokens.md) |

## Allowlist (`allowlist`)

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `allowlist.mcps` | `[]string` | *(empty)* | Approved MCP server identifiers. Used by the [mcp](../modules/mcp.md) module. |
| `allowlist.skills` | `[]string` | *(empty)* | Approved skill names (reserved for future use). |
| `allowlist.domains` | `[]string` | *(empty)* | Approved outbound domains. Used by the [network](../modules/network.md) module. |
| `allowlist.precommit_hooks` | `[]string` | `["gitleaks"]` | Required pre-commit hook names. Used by the [hooks](../modules/hooks.md) module. |

## Token Baseline

Required when `modules.tokens` is `true`. Omitting it while the module is enabled causes a validation error at startup.

| Key | Type | Required | Default | Description |
|-----|------|----------|---------|-------------|
| `token_baseline.daily_mean` | int | Yes | — | Expected average daily token count based on normal usage. |
| `token_baseline.std_dev` | float | Yes | — | Standard deviation of daily token count. |
| `token_baseline.multiplier` | float | No | `3.0` | Anomaly threshold multiplier. Threshold = `daily_mean + multiplier × std_dev`. |

## Environment Variables

| Variable | Overrides | Description |
|----------|-----------|-------------|
| `AI_GUARDRAILS_SIEM_ENDPOINT` | `siem_endpoint` | Overrides the SIEM endpoint from the config file. Useful in CI environments where secrets are injected as environment variables. |

## Example Configuration

```json
{
  "mode": "monitor",
  "score_threshold": 80,
  "modules": {
    "evals": true,
    "humanloop": true,
    "network": true,
    "tokens": true
  },
  "allowlist": {
    "mcps": ["mcp://filesystem", "mcp://github"],
    "domains": ["api.anthropic.com"],
    "precommit_hooks": ["gitleaks", "detect-secrets"]
  },
  "token_baseline": {
    "daily_mean": 50000,
    "std_dev": 10000,
    "multiplier": 3.0
  },
  "env_key_watch_list": ["ANTHROPIC_API_KEY", "OPENAI_API_KEY"]
}
```
