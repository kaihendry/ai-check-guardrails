# Implementation Plan: Environment API Key Detection

**Branch**: `004-env-key-detection` | **Date**: 2026-05-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/004-env-key-detection/spec.md`

## Summary

Add a `envkeys` detection module to `ai-check-guardrails` that inspects the
process environment for plaintext Anthropic credential variables (defaulting to
`ANTHROPIC_API_KEY`) and emits `PLAINTEXT_CREDENTIAL` findings. The module
follows the existing `init()`-registration pattern, adds one new toggle to
`ModuleToggles` and one optional field to `Config`, and requires no new
dependencies.

## Technical Context

**Language/Version**: Go 1.26
**Primary Dependencies**: stdlib only (`os`)
**Storage**: N/A
**Testing**: `go test ./...`
**Target Platform**: Linux/macOS developer workstation (same as existing tool)
**Project Type**: CLI tool — new detection module inside existing binary
**Performance Goals**: Module completes in < 1 ms (pure env-var lookup)
**Constraints**: No new dependencies; credential value must never appear in output
**Scale/Scope**: Single new file + config additions; no structural changes

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Status | Notes |
|-----------|------|--------|-------|
| I. Simplicity | One package, minimal deps | ✅ PASS | Single new file `internal/modules/envkeys/envkeys.go`; stdlib only (`os`) |
| I. Simplicity | No speculative abstractions | ✅ PASS | No new interfaces; reuses `Module`, `Finding`, `Config` types |
| I. Simplicity | New dependency justified | ✅ N/A | No new dependency |
| II. Integrity | Errors surface explicitly | ✅ PASS | `os.LookupEnv` never errors; module returns `nil` error on clean run |
| II. Integrity | Output is deterministic | ✅ PASS | Env-var snapshot at invocation time; no hidden state |
| II. Integrity | Confidence communicated | ✅ PASS | `WARN` vs `HIGH` severity distinguishes confirmed credential from placeholder |
| Security | Sensitive data not written to output | ✅ PASS | Only variable *name* included in `Finding.Resource`; value never logged |
| Security | Input validated | ✅ PASS | Env values treated as opaque strings; length/entropy heuristic applied locally |
| CLI | Exit codes follow spec | ✅ PASS | Handled by existing runner; no change needed |

**Complexity Tracking**: no violations.

## Project Structure

### Documentation (this feature)

```text
specs/004-env-key-detection/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit-tasks)
```

### Source Code Changes

```text
internal/
├── config/
│   └── config.go            # +EnvKeys bool toggle; +EnvKeyWatchList []string field
└── modules/
    └── envkeys/
        └── envkeys.go       # new — PLAINTEXT_CREDENTIAL detection module

cmd/ai-check-guardrails/
└── main.go                  # +blank import for envkeys module registration
```

No new top-level directories. No new files outside the existing module tree.
