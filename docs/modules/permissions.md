# permissions

Checks whether `~/.claude/settings.json` contains explicit permission or tool-access configuration. Without defined restrictions, Claude may use any available tool without constraints.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `PERMISSION_MISCONFIGURED` | HIGH | No permission or tool-restriction configuration is found in settings.json. | Add permission restrictions to `~/.claude/settings.json` to limit tool access. |

The module looks for any of these keys in settings.json:

- `permissions`
- `allowedTools`
- `blockedTools`
- `allow`
- `deny`

If any one of these keys is present, no finding is emitted.

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.permissions` | bool | `true` | Enable or disable this module |
