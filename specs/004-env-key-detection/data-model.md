# Data Model: Environment API Key Detection

**Feature**: 004-env-key-detection | **Date**: 2026-05-11

## Reused Types (no changes)

| Type | Package | Notes |
|------|---------|-------|
| `Finding` | `internal/modules` | `Type="PLAINTEXT_CREDENTIAL"`, `Resource=<var name>`, value never stored |
| `Severity` | `internal/modules` | `HIGH` (len ≥ 20) or `WARN` (present, len < 20) |
| `Module` interface | `internal/modules` | `Name()`, `Enabled()`, `Run()` |

## Config Changes

### `ModuleToggles` — new field

```
EnvKeys bool `json:"env_keys"`
```

Default: `true` (enabled by default, consistent with `bypass`, `settings`, `hooks`).

### `Config` — new field

```
EnvKeyWatchList []string `json:"env_key_watch_list,omitempty"`
```

- When present and non-empty: replaces the built-in default list.
- When absent or empty: module uses built-in default `["ANTHROPIC_API_KEY"]`.

## New Module State

The `envkeys.mod` struct is stateless. No persistent storage.

## Finding Shape (example)

```json
{
  "type": "PLAINTEXT_CREDENTIAL",
  "severity": "HIGH",
  "module": "envkeys",
  "resource": "ANTHROPIC_API_KEY",
  "description": "Credential variable ANTHROPIC_API_KEY is set as a plaintext environment variable.",
  "remediation": "Remove ANTHROPIC_API_KEY from the environment and load it from an approved secrets manager instead."
}
```
