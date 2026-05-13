# Module Interface

Every guardrail check is a **module** — a self-contained unit that inspects one aspect of the Claude environment and returns zero or more **findings**.

## The Module Interface

Each module implements three methods:

| Method | Signature | Description |
|--------|-----------|-------------|
| `Name` | `Name() string` | Returns the module's stable identifier (e.g. `"bypass"`, `"envkeys"`). Used in finding output and config toggles. |
| `Enabled` | `Enabled(cfg Config) bool` | Returns `true` when the module should run. Reads the relevant `modules.*` key from the configuration. |
| `Run` | `Run(cfg Config) ([]Finding, error)` | Performs the audit and returns a slice of findings. An empty slice means no issues were detected. A non-nil error means the check could not be completed. |

Modules that encounter a non-fatal condition (e.g., a missing optional file) return `nil, nil` — no findings and no error.

## The Finding Struct

Every issue a module detects is represented as a `Finding`:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | Yes | A module-specific constant identifying the kind of issue (e.g. `"POLICY_BYPASS"`, `"MISSING_CONFIG"`). |
| `severity` | string | Yes | One of `INFO`, `WARN`, `HIGH`, or `CRITICAL`. See [Severity Levels](severity.md). |
| `module` | string | Yes | The name of the module that produced this finding. Matches the module's `Name()` return value. |
| `resource` | string | Yes | The specific file path, config key, or environment variable that was inspected (e.g. `"~/.zsh_history"`, `"allowlist.mcps"`). |
| `description` | string | Yes | A human-readable summary of what was found. |
| `remediation` | string | Yes | A concrete, actionable step to resolve the finding. |
| `confidence` | float (0–1) | No | Present only for heuristic findings (e.g. token anomaly detection). Indicates how certain the module is. |

## Writing a New Module

To add a module:

1. Create a package under `internal/modules/<name>/`.
2. Define a type implementing the three methods above.
3. Register it in an `init()` function: `modules.Register(&mod{})`.
4. Add a toggle to `ModuleToggles` in `internal/config/config.go` and set a default in `defaults()`.

The module will then be picked up automatically by the runner.
