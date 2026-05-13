# sandbox

Checks whether sandbox or isolation settings are explicitly disabled in `~/.claude/settings.json`. Disabling sandboxing removes restrictions on Claude's data access and network behaviour.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `SANDBOX_VIOLATION` | HIGH | A sandbox-related config key is explicitly set to `false`. | Enable sandboxing in `~/.claude/settings.json` to restrict Claude's data access. |

One finding is emitted per disabled key. The module checks these settings.json keys:

- `sandbox`
- `isolatedEnv`
- `restrictedMode`
- `networkAccess`

A finding is only emitted when the key is present **and** explicitly set to `false`. Absent keys do not trigger findings.

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.sandbox` | bool | `true` | Enable or disable this module |
