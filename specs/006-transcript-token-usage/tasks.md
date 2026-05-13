# Tasks: Transcript Token Usage Reading

**Input**: Design documents from `/specs/006-transcript-token-usage/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create test fixture files that all user story tests depend on.

- [x] T001 Create testdata/tokens/ directory with recent-high.jsonl — JSONL file containing assistant messages timestamped within the last hour with token counts that exceed any reasonable baseline (e.g., 200,000 total tokens)
- [x] T002 [P] Create testdata/tokens/recent-low.jsonl — JSONL file containing assistant messages timestamped within the last hour with low token counts (e.g., 500 total tokens) that fall within a normal baseline
- [x] T002 [P] Create testdata/tokens/old-session.jsonl — JSONL file containing assistant messages timestamped 48 hours ago with high token counts (to verify lookback window exclusion)

**Checkpoint**: Fixture files ready — all tests can reference them via os.DirFS or relative path

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Config change and internal struct that all user story implementations share.

**⚠️ CRITICAL**: Must complete before user story implementation begins

- [x] T002 Add `LookbackHours int` field to `TokenBaseline` struct in `internal/config/config.go` (zero value defaults to 24 h; no validation change needed — zero is valid)
- [x] T002 [P] Define unexported `transcriptRecord` struct in `internal/modules/tokens/tokens.go` with fields mapping to JSONL format per `specs/006-transcript-token-usage/data-model.md` (Timestamp string, Message.Usage with InputTokens, OutputTokens, CacheCreationInputTokens, CacheReadInputTokens int)

**Checkpoint**: Foundation ready — user story tasks can begin

---

## Phase 3: User Story 1 - Detect Token Anomalies from Real Usage Data (Priority: P1) 🎯 MVP

**Goal**: Replace the always-zero `estimateDailyTokens()` placeholder with real transcript reading so TOKEN_ANOMALY findings are produced when actual daily token usage exceeds the configured threshold.

**Independent Test**: Set a baseline (DailyMean=1000, StdDev=100, Multiplier=3.0), point the module at `testdata/tokens/` containing only `recent-high.jsonl`, run `mod.Run(cfg)`, and verify a TOKEN_ANOMALY finding with WARN severity is returned.

### Implementation for User Story 1

- [x] T002 [US1] Implement `readDailyTokens(projectsDir string, lookbackHours int) (int64, error)` in `internal/modules/tokens/tokens.go` — walk `<projectsDir>/*/**.jsonl` two levels deep using `os.ReadDir`, open each `.jsonl` file, scan line-by-line with `bufio.Scanner`, unmarshal each line into `transcriptRecord`, filter by timestamp within lookback window, sum all four token fields; skip unreadable files and malformed lines; return 0, nil when dir absent
- [x] T002 [US1] Replace the `estimateDailyTokens()` call in `(m *mod) Run()` with `readDailyTokens(projectsDir, lookbackHours)` in `internal/modules/tokens/tokens.go` — derive `projectsDir` from `os.UserHomeDir()` joined with `.claude/projects`; derive `lookbackHours` from `cfg.TokenBaseline.LookbackHours` (default 24 when 0)
- [x] T002 [US1] Add test `TestTokens_AboveThreshold` in `internal/modules/tokens/tokens_test.go` — inject testdata dir path, set baseline below the fixture total, assert TOKEN_ANOMALY WARN finding is returned
- [x] T002 [P] [US1] Add test `TestTokens_BelowThreshold` in `internal/modules/tokens/tokens_test.go` — inject testdata dir path with only `recent-low.jsonl`, set baseline well above fixture total, assert zero findings returned

**Checkpoint**: User Story 1 fully functional — `go test ./internal/modules/tokens/` passes with new tests

---

## Phase 4: User Story 2 - Handle Missing or Inaccessible Transcripts Gracefully (Priority: P2)

**Goal**: Ensure the tokens module never errors when transcripts are absent or unreadable, preserving correct behaviour in CI and non-Claude environments.

**Independent Test**: Call `readDailyTokens("/nonexistent/path", 24)` and assert it returns `(0, nil)`.

### Implementation for User Story 2

- [x] T010 [US2] Add test `TestTokens_NoTranscriptsDir` in `internal/modules/tokens/tokens_test.go` — call `mod.Run()` with a config pointing `projectsDir` to a temp dir that does not exist, assert zero findings and nil error (requires passing projectsDir through config or a helper override; see T011)
- [x] T011 [US2] Extend `readDailyTokens` in `internal/modules/tokens/tokens.go` to accept a configurable `projectsDir` string (already in T006 signature) and verify the `os.IsNotExist` path returns `(0, nil)` — add test `TestTokens_MalformedLines` that creates a temp JSONL file with some valid and some invalid lines and asserts only valid tokens are counted

**Checkpoint**: User Stories 1 AND 2 independently tested and passing

---

## Phase 5: User Story 3 - Configurable Lookback Window (Priority: P3)

**Goal**: Allow operators to tune the lookback window via `token_baseline.lookback_hours` in config, so that only sessions within that window contribute to the anomaly check.

**Independent Test**: Set `LookbackHours=1`, point the module at a dir containing only `old-session.jsonl` (timestamped 48 h ago), and assert zero findings even though the token count would otherwise exceed the baseline.

### Implementation for User Story 3

- [x] T012 [US3] Add test `TestTokens_LookbackExcludesOldSessions` in `internal/modules/tokens/tokens_test.go` — use `testdata/tokens/old-session.jsonl`, set `LookbackHours=1`, set a low baseline, assert zero findings (old session excluded)
- [x] T013 [P] [US3] Add test `TestTokens_LookbackIncludesRecentSessions` in `internal/modules/tokens/tokens_test.go` — use `testdata/tokens/recent-high.jsonl`, set `LookbackHours=24`, set a low baseline, assert TOKEN_ANOMALY is returned (recent session included)

**Checkpoint**: All three user stories independently functional and tested

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Verify integration with existing tests and clean up the placeholder.

- [x] T014 Remove `estimateDailyTokens()` function stub from `internal/modules/tokens/tokens.go` (replaced by `readDailyTokens`)
- [x] T015 [P] Run `go test ./...` from repo root and confirm all existing tests still pass (TestTokens_NoBaseline and TestTokens_WithBaselineNoUsageLog must continue to pass unchanged)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately; T002 and T003 can run in parallel
- **Foundational (Phase 2)**: Depends on Phase 1 completion; T004 and T005 can run in parallel
- **User Story 1 (Phase 3)**: Depends on Phase 2 completion; T008 and T009 can run in parallel
- **User Story 2 (Phase 4)**: Depends on T006/T007 completion (readDailyTokens signature established)
- **User Story 3 (Phase 5)**: Depends on T007 completion (LookbackHours wired in Run()); T012 and T013 can run in parallel
- **Polish (Phase 6)**: Depends on all user stories complete; T015 depends on T014

### User Story Dependencies

- **US1 (P1)**: Requires Foundational complete — core implementation, no dependency on US2/US3
- **US2 (P2)**: Requires T006 signature to be stable (readDailyTokens accepts projectsDir)
- **US3 (P3)**: Requires T007 (LookbackHours wired into Run())

### Parallel Opportunities

- T002 and T003 (fixture creation) run in parallel
- T004 and T005 (config + struct) run in parallel
- T008 and T009 (US1 tests) run in parallel
- T012 and T013 (US3 tests) run in parallel

---

## Parallel Example: User Story 1

```bash
# After T006+T007, launch both US1 tests together:
Task: "TestTokens_AboveThreshold in internal/modules/tokens/tokens_test.go"  # T008
Task: "TestTokens_BelowThreshold in internal/modules/tokens/tokens_test.go"  # T009
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Create fixture JSONL files
2. Complete Phase 2: Config field + transcriptRecord struct
3. Complete Phase 3: Implement readDailyTokens + wire into Run() + tests
4. **STOP and VALIDATE**: `go test ./internal/modules/tokens/` — all tests pass
5. Tokens module now produces real TOKEN_ANOMALY findings

### Incremental Delivery

1. Phase 1 + 2 → Fixtures and types ready
2. Phase 3 → US1 complete: real anomaly detection works → MVP
3. Phase 4 → US2 complete: graceful degradation verified
4. Phase 5 → US3 complete: lookback window configurable
5. Phase 6 → All existing tests still pass

---

## Notes

- [P] tasks touch different files and have no incomplete dependencies
- Each user story is independently runnable with `go test ./internal/modules/tokens/ -run TestTokens_<Name>`
- Fixture JSONL files must use RFC3339 timestamps; use `time.Now().UTC().Format(time.RFC3339)` pattern for recent files and `time.Now().Add(-48*time.Hour)` for old files — hardcode plausible values in fixtures
- The `readDailyTokens` function must accept `projectsDir` as a parameter (not hardcode `~/.claude/projects`) to allow tests to inject testdata paths
