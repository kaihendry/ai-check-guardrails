# Research: Environment API Key Detection

**Feature**: 004-env-key-detection | **Date**: 2026-05-11

## Decisions

### 1. Environment variable inspection

- **Decision**: Use `os.LookupEnv(name)` per variable name in the watched list.
- **Rationale**: Returns `(value, present bool)`, cleanly distinguishing "not set" from "set to empty". No shell subprocess needed.
- **Alternatives considered**: `os.Environ()` (full scan) — rejected because it scans all vars unnecessarily when the watched list is small.

### 2. Credential heuristic (HIGH vs WARN severity)

- **Decision**: A value is classified `HIGH` if `len(value) >= 20` (consistent with Anthropic key format `sk-ant-…`). Values present but shorter are classified `WARN`.
- **Rationale**: Avoids false `HIGH` findings for placeholder values like `"changeme"` or `"todo"` while still flagging them at `WARN`.
- **Alternatives considered**: Entropy calculation — rejected; adds complexity with marginal benefit for this use case.

### 3. Watched variable list defaults

- **Decision**: Built-in default list is `["ANTHROPIC_API_KEY"]`. Operator can set `env_key_watch_list` in `config.json` to replace the default entirely.
- **Rationale**: A replace-not-merge approach avoids surprises where an operator adds a custom list but the default keys keep firing. Explicit is safer.
- **Alternatives considered**: Merge (additive) — rejected per Simplicity principle; operator intent is clearer with replace semantics.

### 4. Empty-value handling

- **Decision**: Variables set to empty string produce no finding.
- **Rationale**: An empty value carries no credential; emitting `INFO` findings for empty vars would be noise.

### 5. Module registration

- **Decision**: `init()` function in `envkeys.go` registers the module, matching the pattern used by all 12 existing modules. `main.go` receives one blank import `_ "github.com/kaihendry/ai-check-guardrails/internal/modules/envkeys"`.
- **Rationale**: Consistent with existing codebase; no changes to runner or registry types.
