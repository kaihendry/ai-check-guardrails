---

description: "Tasks for fixing the self-update version check"
---

# Tasks: Fix Self-Update Version Check

**Input**: Design documents from `/specs/007-fix-update-check/`
**Prerequisites**: plan.md, spec.md, research.md

**Organization**: Tasks are grouped by user story. This is a surgical one-file fix — no setup or foundational phases are needed since the existing project infrastructure is unchanged.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to

---

## Phase 1: User Story 1 - No Unnecessary Downloads When Already Up to Date (Priority: P1) 🎯 MVP

**Goal**: Fix the version comparison so the tool skips downloading when already on the latest release.

**Independent Test**: Install the latest binary, run it, and confirm zero `[update]` lines appear on stderr and the binary is not replaced.

### Implementation for User Story 1

- [x] T001 [US1] Change `tag_name: latest` to `tag_name: ${{ env.VERSION }}` in `.github/workflows/release.yml` (line 49)
- [x] T002 [US1] Run existing tests to confirm no-op path still passes: `go test ./cmd/ai-check-guardrails/...`

**Checkpoint**: After T001 merges to main, the next CI release will tag the release with the actual version string. A freshly installed binary will no longer trigger a spurious re-download.

---

## Phase 2: User Story 2 - Update Still Occurs When a Newer Version Exists (Priority: P2)

**Goal**: Confirm the valid update path is not broken by the fix.

**Independent Test**: `TestCheckAndUpdate_UpdatesWhenNewVersionAvailable` passes — an outdated binary triggers a download and is replaced with the new binary content.

### Implementation for User Story 2

- [x] T003 [US2] Verify `TestCheckAndUpdate_UpdatesWhenNewVersionAvailable` in `cmd/ai-check-guardrails/update_test.go` passes after T001 (no code change needed — test is already correct)

**Checkpoint**: Both existing tests pass. The update path is confirmed to work for both the no-op and the actual-update cases.

---

## Phase 3: User Story 3 - Version Identity Is Stable and Comparable (Priority: P3)

**Goal**: Confirm the version string embedded in a released binary matches the release tag returned by the update check.

**Independent Test**: Run `ai-check-guardrails --version` and query the GitHub Releases API — both return the same version string.

### Implementation for User Story 3

- [ ] T004 [US3] After merging, install the newly built binary via `install.sh` and run `ai-check-guardrails --version` — confirm the output matches the GitHub release `tag_name` returned by the API

**Checkpoint**: Version string in binary matches release tag; no spurious update occurs.

---

## Phase 4: Polish & Cross-Cutting Concerns

- [ ] T005 [P] Optionally simplify release name in `.github/workflows/release.yml` from `"Latest (${{ env.VERSION }})"` to `${{ env.VERSION }}` for consistency (cosmetic only, no functional impact)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (US1)**: No dependencies — start immediately
- **Phase 2 (US2)**: Depends on T001 (needs the fix in place to confirm tests pass correctly)
- **Phase 3 (US3)**: Depends on T001 merged and a new CI release published
- **Phase 4 (Polish)**: Optional, can happen any time after T001

### User Story Dependencies

- **US1 (P1)**: No dependencies — the one-line fix is self-contained
- **US2 (P2)**: Verification only — depends on US1 fix being present
- **US3 (P3)**: Post-deploy verification — depends on CI building a new release after T001 merges

### Parallel Opportunities

No meaningful parallelism — there is only one code change (T001). T002, T003, T004 are sequential verification steps.

---

## Parallel Example: User Story 1

```bash
# Single task — no parallelism within US1
Task: "Change tag_name in .github/workflows/release.yml"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Make the one-line change in `.github/workflows/release.yml` (T001)
2. Run `go test ./cmd/ai-check-guardrails/...` locally (T002)
3. Open PR, merge to main
4. **STOP and VALIDATE**: Observe CI publish a release tagged with the version string; install and run — confirm no `[update]` messages

### Incremental Delivery

1. T001 → fixes the bug (merge immediately)
2. T003 → confirms no regression (CI run confirms)
3. T004 → post-deploy smoke test (manual verification)
4. T005 → optional cosmetic cleanup (low priority)

---

## Notes

- Total tasks: 5 (4 substantive + 1 optional cosmetic)
- The entire fix is T001 — one line changed in one file
- All existing tests in `update_test.go` remain valid and must pass unchanged
- No new test tasks: existing tests already cover both the no-op and update paths
