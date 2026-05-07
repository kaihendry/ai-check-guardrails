# Implementation Plan: AI Guardrails Audit

**Branch**: `001-ai-guardrails-audit` | **Date**: 2026-05-07 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/001-ai-guardrails-audit/spec.md`

## Summary

A standalone Go binary that runs periodically (via `launchd` on macOS or a systemd
unit on Linux) to audit Claude AI guardrail compliance on developer workstations and
relay structured findings as JSON to a SIEM endpoint. The tool operates in two modes:
`monitor` (detect and log only) and `enforce` (detect, log, and alert/block). Each of
the 12 detection modules is individually toggle-able via config so the team can roll
out checks incrementally before enabling enforcement.

## Technical Context

**Language/Version**: Go 1.22+
**Primary Dependencies**: stdlib only (`flag`, `encoding/json`, `net/http`,
`log/slog`, `os/exec`, `embed`); no third-party runtime dependencies
**Storage**: Local JSON config file + local finding buffer file (for SIEM retry);
no database
**Testing**: `go test ./...` (stdlib); integration tests via `testscript` (stdlib-compatible)
**Target Platform**: macOS 13+ (primary, launchd scheduling), Linux (systemd
scheduling); single cross-compiled binary
**Project Type**: CLI binary (single self-contained executable)
**Performance Goals**: Full scan completes in < 30 seconds on a standard
developer workstation (aligned with SC-001)
**Constraints**: Single binary with no external runtime dependencies; < 20MB binary
size; safe to run as the developer's own user account (no root required for
monitor mode)
**Scale/Scope**: Single workstation per invocation; ~100 enrolled workstations
per enterprise deployment

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Simplicity

| Gate | Status | Notes |
|------|--------|-------|
| Single binary, no framework overhead | ✅ PASS | stdlib `flag` only; no cobra/viper |
| Each module does one thing | ✅ PASS | One module per detection concern |
| No abstraction before demonstrated need | ✅ PASS | Module interface justified by 12 modules; no plugin system |
| YAGNI: feature toggles scoped to listed modules only | ✅ PASS | One toggle per module in config; no open extension mechanism |
| Dependency count minimal | ✅ PASS | Zero runtime third-party deps; `go:embed` for schedule templates |

### Integrity

| Gate | Status | Notes |
|------|--------|-------|
| All errors surface to stderr + non-zero exit | ✅ PASS | FR-017: exit 0/1/2; all errors → stderr |
| No silent failures | ✅ PASS | Each module returns `error`; runner logs and counts errors |
| Output deterministic for same input | ✅ PASS | No randomness; timestamp-only non-determinism in log event |
| Security findings never suppressed | ✅ PASS | All findings in JSON output regardless of mode |
| Confidence level communicated on partial detection | ✅ PASS | Token anomaly module emits `confidence` field |

**Post-Phase-1 re-check**: No violations introduced by design artifacts. No
Complexity Tracking entries required.

## Project Structure

### Documentation (this feature)

```text
specs/001-ai-guardrails-audit/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   ├── cli-flags.md     # CLI flag schema
│   ├── siem-event.md    # SIEM JSON event schema
│   └── config.md        # Tool config file schema
└── tasks.md             # Phase 2 output (/speckit-tasks — not created here)
```

### Source Code (repository root)

```text
cmd/
└── ai-check-guardrails/
    └── main.go              # Entry point: parse flags, load config, run audit

internal/
├── audit/
│   └── runner.go            # Orchestrates module execution, scoring, SIEM emit
├── config/
│   └── config.go            # Load/validate JSON config; feature-toggle resolution
├── lock/
│   └── lock.go              # Run-lock via syscall.Flock; exit 2 if already running
├── siem/
│   └── transport.go         # Emit JSON to stdout; optional HTTP POST with retry
├── score/
│   └── score.go             # Posture score calculation (weighted findings)
└── modules/
    ├── module.go            # Module interface + registry
    ├── settings/            # FR-006: Verify Claude settings.json
    ├── mcp/                 # FR-007: MCP/skill allowlist check
    ├── permissions/         # FR-008: Permission validation tests
    ├── tokens/              # FR-013: Token usage anomaly detection
    ├── network/             # FR-014: Outbound network request logging
    ├── evals/               # FR-015: Anthropic eval hook integration
    ├── sandbox/             # FR-007 (scope): Sandboxing checks
    ├── humanloop/           # FR-008 (scope): Human-in-the-loop detection
    ├── bypass/              # FR-009: Policy-bypass flag detection
    ├── banner/              # FR-012: Terminal banner display
    ├── hooks/               # FR-010: Pre-commit hook presence check
    └── gamification/        # FR-011: Score display + training link

embed/
├── launchd.plist.tmpl       # macOS launchd plist template (go:embed)
└── systemd.service.tmpl     # Linux systemd unit template (go:embed)

go.mod
go.sum                       # empty on first run (no third-party deps)
```

**Structure Decision**: Single-project Go layout. All packages are `internal/` to
prevent accidental import by external tools. `cmd/ai-check-guardrails/main.go` is
the sole binary entry point. No `pkg/` layer — complexity not yet justified.
