# Contract: Config File Schema

**Default path**: `~/.config/ai-check-guardrails/config.json`
**Override**: `--config <path>` flag

The config file is JSON. All fields are optional; omitted fields use the defaults
shown below.

## JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12",
  "title": "Config",
  "type": "object",
  "properties": {
    "mode": {
      "type": "string",
      "enum": ["monitor", "enforce"],
      "default": "monitor"
    },
    "siem_endpoint": {
      "type": "string",
      "format": "uri",
      "description": "HTTPS URL for HTTP POST delivery. Omit to use stdout only."
    },
    "scan_root": {
      "type": "string",
      "default": "$HOME",
      "description": "Root directory for pre-commit hook and repo scanning."
    },
    "banner_url": {
      "type": "string",
      "format": "uri",
      "description": "URL shown in terminal banner when score < score_threshold."
    },
    "score_threshold": {
      "type": "integer",
      "minimum": 0,
      "maximum": 100,
      "default": 70
    },
    "modules": {
      "type": "object",
      "properties": {
        "settings":     { "type": "boolean", "default": true },
        "mcp":          { "type": "boolean", "default": true },
        "permissions":  { "type": "boolean", "default": true },
        "tokens":       { "type": "boolean", "default": false },
        "network":      { "type": "boolean", "default": false },
        "evals":        { "type": "boolean", "default": false },
        "sandbox":      { "type": "boolean", "default": true },
        "humanloop":    { "type": "boolean", "default": false },
        "bypass":       { "type": "boolean", "default": true },
        "banner":       { "type": "boolean", "default": true },
        "hooks":        { "type": "boolean", "default": true },
        "gamification": { "type": "boolean", "default": true }
      }
    },
    "allowlist": {
      "type": "object",
      "properties": {
        "mcps":            { "type": "array", "items": { "type": "string" } },
        "skills":          { "type": "array", "items": { "type": "string" } },
        "domains":         { "type": "array", "items": { "type": "string" } },
        "precommit_hooks": { "type": "array", "items": { "type": "string" },
                             "default": ["gitleaks"] }
      }
    },
    "token_baseline": {
      "type": "object",
      "description": "Required when tokens module is enabled.",
      "properties": {
        "daily_mean": { "type": "integer" },
        "std_dev":    { "type": "number" },
        "multiplier": { "type": "number", "default": 3.0 }
      },
      "required": ["daily_mean", "std_dev"]
    }
  }
}
```

## Minimal Example (monitor mode, core modules only)

```json
{
  "mode": "monitor",
  "allowlist": {
    "mcps": ["mcp://filesystem", "mcp://brave-search"],
    "skills": ["speckit-plan", "speckit-specify"],
    "precommit_hooks": ["gitleaks"]
  },
  "banner_url": "https://intranet.example.com/security-training"
}
```

## Full Example

```json
{
  "mode": "monitor",
  "siem_endpoint": "https://siem.example.com/api/v1/events",
  "scan_root": "/home/alice",
  "banner_url": "https://intranet.example.com/security-training",
  "score_threshold": 80,
  "modules": {
    "settings": true,
    "mcp": true,
    "permissions": true,
    "tokens": true,
    "network": true,
    "evals": false,
    "sandbox": true,
    "humanloop": false,
    "bypass": true,
    "banner": true,
    "hooks": true,
    "gamification": true
  },
  "allowlist": {
    "mcps": ["mcp://filesystem", "mcp://brave-search"],
    "skills": ["speckit-plan", "speckit-specify"],
    "domains": ["api.anthropic.com", "github.com"],
    "precommit_hooks": ["gitleaks", "detect-secrets"]
  },
  "token_baseline": {
    "daily_mean": 50000,
    "std_dev": 12000,
    "multiplier": 3.0
  }
}
```

## Validation Rules

- If `tokens` module is enabled and `token_baseline` is absent, the tool MUST
  emit a `MODULE_UNAVAILABLE` finding for the tokens module and skip it.
- If `mcps` allowlist is empty and `mcp` module is enabled, all MCPs are flagged
  as `UNAPPROVED_MCP` with severity WARN (not HIGH) and include an
  `UNCONFIGURED_ALLOWLIST` finding.
- `siem_endpoint` MUST use HTTPS; an HTTP URL MUST be rejected with exit code 2.
- `scan_root` MUST be an absolute path; relative paths MUST be rejected with
  exit code 2.
