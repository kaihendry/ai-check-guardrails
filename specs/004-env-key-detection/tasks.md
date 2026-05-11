# Tasks: Environment API Key Detection

**Input**: Design documents from `/specs/004-env-key-detection/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅

**Organization**: Tasks grouped by user story for independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: User story this task belongs to

---

## Phase 1: Setup

**Purpose**: No new project structure or dependencies required — existing module tree is reused.

- [x] T001 Confirm `go build ./cmd/ai-check-guardrails/` passes clean before any changes

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Config additions that MUST exist before the module file can compile.

⚠️ **CRITICAL**: US1, US2, and US3 all depend on this phase completing first.

- [x] T002 Add `EnvKeys bool \`json:"env_keys"\`` to `ModuleToggles` struct in `internal/config/config.go`
- [x] T003 Add `EnvKeyWatchList []string \`json:"env_key_watch_list,omitempty"\`` field to `Config` struct in `internal/config/config.go`
- [x] T004 Set `EnvKeys: true` in the `defaults()` function in `internal/config/config.go`
- [x] T005 Confirm `go build ./...` still passes after config changes

**Checkpoint**: Config compiles — module implementation can now begin.

---

## Phase 3: User Story 1 — Detect Plaintext API Keys (Priority: P1) 🎯 MVP

**Goal**: Emit a `PLAINTEXT_CREDENTIAL` finding when `ANTHROPIC_API_KEY` is set in the environment.

**Independent Test**: `ANTHROPIC_API_KEY=sk-ant-test123456789012345 ./ai-check-guardrails` emits a JSON finding with `"type":"PLAINTEXT_CREDENTIAL"` and `"severity":"HIGH"`. Unset the variable and confirm no finding.

- [x] T006 [US1] Create `internal/modules/envkeys/envkeys.go` with `init()` registration, `Name() = "envkeys"`, `Enabled()` checking `cfg.Modules.EnvKeys`, and `Run()` that iterates a built-in default list `["ANTHROPIC_API_KEY"]`, calls `os.LookupEnv`, skips empty values, emits `PLAINTEXT_CREDENTIAL` at `HIGH` (len ≥ 20) or `WARN` (len < 20), with `Resource` = variable name and `Remediation` = "Remove <VAR> from the environment and load it from an approved secrets manager instead."
- [x] T007 [US1] Add blank import `_ "github.com/kaihendry/ai-check-guardrails/internal/modules/envkeys"` to `cmd/ai-check-guardrails/main.go`
- [x] T008 [US1] Write `internal/modules/envkeys/envkeys_test.go` covering: HIGH finding for value len ≥ 20, WARN for short value, no finding when variable absent, no finding when value is empty string
- [x] T009 [US1] Run `go test ./internal/modules/envkeys/...` and confirm all tests pass

**Checkpoint**: US1 fully functional — binary emits correct finding for set `ANTHROPIC_API_KEY`.

---

## Phase 4: User Story 2 — Configurable Watched Variable List (Priority: P2)

**Goal**: Module reads `cfg.EnvKeyWatchList` when non-empty; falls back to built-in default otherwise.

**Independent Test**: Config with `"env_key_watch_list": ["ANTHROPIC_ORG_TOKEN"]` and `ANTHROPIC_ORG_TOKEN=sk-ant-test123456789012345` set produces a `PLAINTEXT_CREDENTIAL` finding for `ANTHROPIC_ORG_TOKEN` and none for `ANTHROPIC_API_KEY` (unset).

- [x] T010 [P] [US2] Update `Run()` in `internal/modules/envkeys/envkeys.go` to use `cfg.EnvKeyWatchList` when `len(cfg.EnvKeyWatchList) > 0`, else use the built-in default list
- [x] T011 [P] [US2] Add test cases to `internal/modules/envkeys/envkeys_test.go` covering: custom list replaces default, empty config list falls back to default
- [x] T012 [US2] Run `go test ./internal/modules/envkeys/...` and confirm all tests pass

**Checkpoint**: US2 complete — operators can supply a custom watch list via config without a rebuild.

---

## Phase 5: User Story 3 — Score Impact and Remediation Hint (Priority: P3)

**Goal**: Each `PLAINTEXT_CREDENTIAL` finding reduces the posture score, and the finding's remediation field provides actionable guidance.

**Independent Test**: Clean run scores 100. Set `ANTHROPIC_API_KEY=sk-ant-test123456789012345`; score drops to 85 (one HIGH = −15). Remediation field in the finding contains text directing to a secrets manager.

- [x] T013 [US3] Add test in `internal/modules/envkeys/envkeys_test.go` asserting `Finding.Remediation` is non-empty and mentions "secrets manager" for every emitted finding
- [x] T014 [US3] Add integration test in `internal/score/score_test.go` (or a new test function) asserting one `HIGH` `PLAINTEXT_CREDENTIAL` finding reduces score from 100 to 85 using the existing `score.Calculate()` — no code change needed if score logic already handles HIGH severity correctly; this task is verification only
- [x] T015 [US3] Run `go test ./internal/score/... ./internal/modules/envkeys/...` and confirm all pass

**Checkpoint**: US3 complete — score penalty confirmed, remediation text present in all findings.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [x] T016 [P] Run `go test ./...` (full suite) and confirm no regressions
- [x] T017 [P] Run `go build ./cmd/ai-check-guardrails/` and confirm binary builds cleanly
- [x] T018 Run the quickstart.md manual test steps and confirm expected output

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — BLOCKS all user stories
- **US1 (Phase 3)**: Depends on Phase 2
- **US2 (Phase 4)**: Depends on Phase 3 (T006 must exist before T010 modifies it)
- **US3 (Phase 5)**: Depends on Phase 3 (finding shape must be final)
- **Polish (Phase 6)**: Depends on Phases 3–5

### Within Each Phase

- T010 and T011 are marked [P] — they touch different parts of the same file but can be drafted in parallel

---

## Parallel Example: Phase 2 Config Changes

```
# T002, T003, T004 touch the same file — run sequentially
T002 → T003 → T004 → T005 (build verify)
```

## Parallel Example: Phase 4 (US2)

```
# T010 (implementation) and T011 (tests) can be drafted in parallel:
Task: "Update Run() to use cfg.EnvKeyWatchList..."
Task: "Add test cases for custom list..."
# Then T012 (run tests) sequentially after both complete
```

---

## Implementation Strategy

### MVP (US1 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T002–T005)
3. Complete Phase 3: US1 (T006–T009)
4. **STOP and VALIDATE**: manual test per quickstart.md
5. Ship or demo

### Incremental Delivery

1. Setup + Foundational → config compiles
2. US1 → detection works for `ANTHROPIC_API_KEY`
3. US2 → operator-configurable list
4. US3 → score and remediation confirmed
5. Polish → full test suite green

---

## Notes

- [P] tasks operate on different files or independent sections — safe to parallelize
- No new dependencies: `os.LookupEnv` is stdlib
- Credential value MUST NOT appear in any output — only variable name in `Resource`
- `init()` registration pattern: module file calls `modules.Register(...)` in its `init()` — main.go blank-imports the package to activate it (matches all 12 existing modules)
- Score deduction for HIGH severity (−15) is already implemented in `internal/score/score.go` — no change required
