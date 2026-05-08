# Implementation Plan: MVP Release Binary & Distribution

**Branch**: `002-mvp-release-binary` | **Date**: 2026-05-08 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `specs/002-mvp-release-binary/spec.md`

## Summary

Add a release pipeline, install script, self-update mechanism, and README so that `ai-check-guardrails` can be distributed as a pre-built binary. Every commit to main automatically builds and publishes a new release to GitHub Releases. Users install with a single `curl` command; the binary self-updates on startup using stdlib-only HTTP + atomic rename.

## Technical Context

**Language/Version**: Go 1.26.2 (existing project)
**Primary Dependencies**: stdlib only (`net/http`, `encoding/json`, `os`, `os/exec`); no new runtime dependencies added
**Storage**: None — no persistent state introduced by this feature
**Testing**: `go test ./...` (existing); new `update.go` covered by unit tests using an HTTP test server
**Target Platform**: Linux amd64/arm64, macOS amd64/arm64; Windows out of scope for MVP
**Project Type**: CLI binary (existing — this feature adds distribution infrastructure)
**Performance Goals**: Self-update check + download completes in < 30 seconds on standard internet connection (aligned with SC-003)
**Constraints**: Zero new runtime dependencies; binary self-replaces without root when installed to `~/.local/bin`; CI pipeline completes in < 15 minutes (aligned with SC-004)
**Scale/Scope**: Single binary per platform; 4 release artifacts per pipeline run; curl install script for individual developers

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Simplicity

| Gate | Status | Notes |
|------|--------|-------|
| No new dependencies introduced | ✅ PASS | stdlib `net/http` + `encoding/json` only; no `go-update` library |
| Self-update logic in one file | ✅ PASS | `cmd/ai-check-guardrails/update.go`, ~60 lines |
| CI pipeline minimal | ✅ PASS | Single workflow file, raw `go build` matrix; no GoReleaser |
| YAGNI: no Homebrew tap, Docker, Cosign | ✅ PASS | Deferred — not in spec |
| Install script uses standard Unix primitives | ✅ PASS | `curl`, `uname`, `sha256sum`/`shasum` only |

### Integrity

| Gate | Status | Notes |
|------|--------|-------|
| Self-update failures are non-fatal and surfaced to stderr | ✅ PASS | FR-004: all error paths log to stderr and continue |
| Version always printed to stderr on startup | ✅ PASS | FR-002: first line of stderr output |
| Binary checksum verified during install | ✅ PASS | install.sh verifies SHA-256 before placing binary |
| Atomic replace — no partial binary state | ✅ PASS | `os.CreateTemp` same-dir → `os.Rename` (single rename syscall) |
| Exit code unaffected by update failures | ✅ PASS | Update errors do not change exit code |

**Post-Phase-1 re-check**: No violations introduced by design artifacts. No Complexity Tracking entries required.

## Project Structure

### Documentation (this feature)

```text
specs/002-mvp-release-binary/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── cli-contract.md  # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit-tasks command)
```

### Source Code (changes to repository root)

```text
.github/
└── workflows/
    └── release.yml          # New: CI pipeline triggered on every main commit

cmd/
└── ai-check-guardrails/
    ├── main.go              # Modified: add var version, startup log, --no-update flag
    └── update.go            # New: self-update logic (~60 lines, stdlib only)

install.sh                   # New: curl-installable install script
README.md                    # New: install instructions, usage, badges
```

**Structure Decision**: Single project layout (existing). This feature adds one new file per concern (`update.go`, `release.yml`, `install.sh`, `README.md`) without restructuring the existing tree.
