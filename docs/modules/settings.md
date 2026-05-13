# settings

Validates that `~/.claude/settings.json` (and `settings.local.json` if present) exist and contain valid configuration. It also surfaces non-default overrides as advisory findings so operators can review unexpected changes.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `MISSING_CONFIG` | CRITICAL | `~/.claude/settings.json` does not exist. | Create `~/.claude/settings.json` with required security keys. |
| `SETTINGS_INVALID` | HIGH | A settings file exists but is not valid JSON. | Fix the JSON syntax in the affected file. |
| `SETTINGS_OVERRIDE` | WARN | A key in settings.json is set to a non-default (non-null, non-false, non-empty) value. | Review whether this override is approved by your security policy. |

`SETTINGS_OVERRIDE` emits one finding per non-default key. Both `settings.json` and `settings.local.json` are checked when the local file exists.

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.settings` | bool | `true` | Enable or disable this module |
