# Contract: CLI Flags

**Binary**: `ai-check-guardrails`
**Flag package**: stdlib `flag`

## Usage

```
ai-check-guardrails [flags]
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--config` | string | `~/.config/ai-check-guardrails/config.json` | Path to JSON config file |
| `--mode` | string | value from config | Override run mode: `monitor` or `enforce` |
| `--output` | string | `stdout` | Output destination: `stdout` or `stderr` |
| `--install-launchd` | bool | false | Write launchd plist to `~/Library/LaunchAgents/` and load it |
| `--install-systemd` | bool | false | Write systemd user unit to `~/.config/systemd/user/` and enable it |
| `--uninstall` | bool | false | Remove installed launchd/systemd schedule and exit |
| `--version` | bool | false | Print version string and exit 0 |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Clean run — no findings |
| 1 | Run completed — one or more findings present |
| 2 | Tool error — misconfiguration, lock conflict, fatal module error |

## Constraints

- `--install-launchd` and `--install-systemd` are mutually exclusive.
- `--mode` overrides the value in the config file for that run only; it does not
  persist to the config file.
- All flags follow POSIX double-dash convention; single-dash short forms are not
  provided (Simplicity principle: no aliases).
