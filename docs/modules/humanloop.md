# humanloop

Checks whether a human-in-the-loop confirmation setting is present in `~/.claude/settings.json`. Without explicit confirmation requirements, Claude may take consequential actions without pausing for human review.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `HUMANLOOP_ABSENT` | WARN | No human-in-the-loop confirmation key is found in settings.json. | Add `requireConfirmation: true` to `~/.claude/settings.json`. |

The module looks for any of these keys in settings.json:

- `requireConfirmation`
- `humanInTheLoop`
- `confirmDangerousActions`

If any one of these keys is present, no finding is emitted.

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.humanloop` | bool | `false` | Enable or disable this module |
