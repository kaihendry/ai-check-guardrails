# bypass

Scans shell history files for use of permission-bypass flags in Claude CLI invocations. Running Claude with `--dangerously-skip-permissions` disables safety guardrails and is a critical policy violation.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `POLICY_BYPASS` | CRITICAL | A permission-bypass flag was found in shell history. | Remove use of permission-bypass flags; request a security exception if genuinely needed. |

Scanned history files (in order):
- `~/.zsh_history`
- `~/.bash_history`
- `~/.history`

Detected flags:
- `--dangerously-skip-permissions`
- `--dangerously_skip_permissions`

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.bypass` | bool | `true` | Enable or disable this module |
