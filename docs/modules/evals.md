# evals

Checks whether the Anthropic-recommended evaluation hooks are configured in `~/.claude/settings.json`. These hooks enable adversarial review and private data protection by intercepting tool calls at key lifecycle events.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `EVAL_HOOK_ABSENT` | HIGH | An Anthropic-recommended eval hook is not configured. | Add the missing hook to the `hooks` section in `~/.claude/settings.json`. |

One finding is emitted for each missing hook. The three recommended hooks are:

| Hook | Purpose |
|------|---------|
| `PreToolUse` | Inspect or block tool calls before they execute |
| `PostToolUse` | Audit tool outputs after execution |
| `Stop` | Review model responses before they are finalized |

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.evals` | bool | `false` | Enable or disable this module |
