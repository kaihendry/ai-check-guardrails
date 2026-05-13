# Implementation Plan: Transcript Token Usage Reading

**Branch**: `006-transcript-token-usage` | **Date**: 2026-05-13 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/006-transcript-token-usage/spec.md`

## Summary

Replace the placeholder `estimateDailyTokens()` stub in the tokens module with real transcript reading from `~/.claude/projects/**/*.jsonl`. Sum all token fields from assistant messages within a configurable lookback window (default 24 h) and compare against the configured anomaly threshold.

## Technical Context

**Language/Version**: Go 1.26.2  
**Primary Dependencies**: stdlib only (`encoding/json`, `os`, `path/filepath`, `time`, `bufio`)  
**Storage**: Local filesystem — `~/.claude/projects/<project>/*.jsonl` (one JSON object per line)  
**Testing**: `go test ./internal/modules/tokens/...`  
**Target Platform**: Linux / macOS (same as existing modules)  
**Project Type**: CLI tool — internal module  
**Performance Goals**: < 2 s to process up to 100 JSONL files totalling 50 MB  
**Constraints**: No new dependencies; must use stdlib only per constitution Simplicity principle  
**Scale/Scope**: Single user's local machine; files fit in memory line-by-line (streaming read)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Check | Status |
|-----------|-------|--------|
| **Simplicity** | Only `tokens.go` and `tokens_test.go` change; stdlib-only; no new abstraction layer | ✅ Pass |
| **Integrity** | Unreadable files are skipped (not silently zeroed), caller sees partial count; MODULE_UNAVAILABLE preserved when no baseline | ✅ Pass |
| **Security: file input validation** | JSONL lines are parsed with `encoding/json`; malformed lines are skipped; no shell execution | ✅ Pass |
| **Security: sensitive data** | Token counts are integers; no key/secret data is read or emitted | ✅ Pass |
| **Integrity: determinism** | Same input files → same output; no hidden side effects | ✅ Pass |

No complexity violations to track.

## Project Structure

### Documentation (this feature)

```text
specs/006-transcript-token-usage/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit-tasks)
```

### Source Code (repository root)

```text
internal/modules/tokens/
├── tokens.go             # Replace estimateDailyTokens(); add readDailyTokens()
└── tokens_test.go        # Add 3+ new test cases using temp testdata

internal/config/
└── config.go             # Add optional LookbackHours to TokenBaseline struct

testdata/tokens/           # New: fixture JSONL files for testing
├── recent-high.jsonl      # Usage above threshold within 24 h
├── recent-low.jsonl       # Usage below threshold within 24 h
└── old-session.jsonl      # Session outside lookback window
```

**Structure Decision**: Single-project layout. Only two existing files change; fixture testdata is added to a new `testdata/tokens/` directory following Go conventions.

## Complexity Tracking

No constitution violations.
