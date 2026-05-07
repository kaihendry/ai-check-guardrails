---
description: "Task list for AI Guardrails Audit implementation"
---

# Tasks: AI Guardrails Audit

**Input**: Design documents from `specs/001-ai-guardrails-audit/`
**Prerequisites**: plan.md ✅ spec.md ✅ research.md ✅ data-model.md ✅ contracts/ ✅

**Tests**: Not explicitly requested in spec — no TDD tasks generated.

**Organization**: Tasks grouped by user story to enable independent delivery.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1–US4)

---

## Phase 1: Setup

**Purpose**: Initialize Go project and embed assets.

- [x] T001 Initialize Go module: run `go mod init github.com/YOUR_ORG/ai-check-guardrails` and create `go.mod`
- [x] T002 Create directory structure per plan.md: `cmd/ai-check-guardrails/`, `internal/audit/`, `internal/config/`, `internal/lock/`, `internal/modules/`, `internal/score/`, `internal/siem/`, `embed/`
- [x] T003 [P] Create macOS launchd plist template at `embed/launchd.plist.tmpl` — interval 1800s, run as current user, stdout/stderr to `/tmp/ai-check-guardrails.log`
- [x] T004 [P] Create Linux systemd unit + timer templates at `embed/systemd.service.tmpl` and `embed/systemd.timer.tmpl` — OnCalendar=*:0/30 (every 30 min)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that ALL user stories depend on. No story work begins until this phase is complete.

**⚠️ CRITICAL**: Complete all T005–T012 before any Phase 3+ work.

- [x] T005 Implement `internal/config/config.go`: Config struct (matches contracts/config.md schema), JSON loading from `~/.config/ai-check-guardrails/config.json`, defaults application, validation (HTTPS-only siem_endpoint, absolute scan_root, tokens baseline required when tokens module enabled)
- [x] T006 [P] Implement `internal/lock/lock.go`: acquire exclusive lockfile at `$TMPDIR/ai-check-guardrails.lock` via `syscall.Flock(LOCK_EX|LOCK_NB)`; return `ErrAlreadyRunning` if lock unavailable; release on process exit
- [x] T007 [P] Implement `internal/modules/module.go`: `Module` interface (`Name() string; Enabled(Config) bool; Run(Config) ([]Finding, error)`); package-level `registry []Module`; exported `Register(Module)` called from each module's `init()`; exported `All() []Module`
- [x] T008 [P] Implement `internal/score/score.go`: `Calculate(findings []Finding) int` — start 100, deduct CRITICAL −25 / HIGH −15 / WARN −5 / INFO 0, floor at 0; return int in [0,100]
- [x] T009 [P] Implement `internal/siem/transport.go`: `Emit(run AuditRun, endpoint string, token string) error` — marshal AuditRun to single-line JSON on stdout; if endpoint non-empty, HTTP POST with `Authorization: Bearer` header, 10s timeout, log failure to stderr, do not retry
- [x] T010 Implement `internal/audit/runner.go`: `Run(cfg Config) AuditRun` — acquire lock (exit 2 on `ErrAlreadyRunning`); iterate `modules.All()`, call `Enabled()` then `Run()`; collect `[]Finding`; calculate score; populate AuditRun fields (RunID uuid, Timestamp UTC, Host, User, Mode, Version, DurationMs); call `siem.Emit`; return AuditRun
- [x] T011 Implement `cmd/ai-check-guardrails/main.go`: parse flags (`--config`, `--mode`, `--output`, `--install-launchd`, `--install-systemd`, `--uninstall`, `--version`) per contracts/cli-flags.md; load config; call `audit.Run`; exit with code from AuditRun.ExitCode

**Checkpoint**: `go build ./...` succeeds with no modules yet (empty registry). Binary exits 0 on clean config with 0 findings.

---

## Phase 3: User Story 1 — Security Team Gains Visibility (Priority: P1) 🎯 MVP

**Goal**: Tool runs, scans Claude settings, and emits a valid JSON SIEM event per contracts/siem-event.md.

**Independent Test**: Run `./ai-check-guardrails --config testdata/minimal-config.json`; confirm single-line JSON on stdout with `schema_version:"1.0"`, non-empty `run_id`, `score` in [0,100], and `exit_code` 0 or 1.

### Implementation for User Story 1

- [x] T012 [P] [US1] Implement `internal/modules/settings/settings.go`: locate Claude `settings.json` at `~/.claude/settings.json` (and `~/.claude/settings.local.json`); check required keys present; detect any keys that differ from defaults (emit `SETTINGS_OVERRIDE` WARN per key); emit `MISSING_CONFIG` CRITICAL if file absent; register via `init()`
- [x] T013 [P] [US1] Implement `internal/modules/banner/banner.go`: if score < `cfg.ScoreThreshold` and `cfg.BannerURL` non-empty, print formatted banner to stderr with score, threshold, and training URL; register via `init()`; module emits no Findings (banner is side-effect only)
- [x] T014 [P] [US1] Implement `internal/modules/gamification/gamification.go`: emit `INFO` Finding with description containing score and trend hint ("Your security score is X/100"); register via `init()`
- [x] T015 [US1] Import all US1 modules in `cmd/ai-check-guardrails/main.go` via blank imports (`_ "github.com/.../settings"` etc.) to trigger `init()` registration
- [x] T016 [US1] Create `testdata/minimal-config.json` per contracts/config.md minimal example; run `./ai-check-guardrails --config testdata/minimal-config.json` and confirm output validates against SIEM event JSON schema in contracts/siem-event.md

**Checkpoint**: US1 independently functional — binary produces valid SIEM event with settings findings and gamification score.

---

## Phase 4: User Story 2 — Detection of Risky MCP / Skill Usage (Priority: P2)

**Goal**: Tool detects unapproved MCPs and access-control violations; UNAPPROVED_MCP findings appear in log.

**Independent Test**: Add a non-allowlisted MCP identifier to `testdata/mcp-test-config.json`; run binary; confirm `UNAPPROVED_MCP` HIGH finding present in JSON output with correct `resource` field.

### Implementation for User Story 2

- [x] T017 [P] [US2] Implement `internal/modules/mcp/mcp.go`: read active MCPs from Claude's MCP configuration (default: `~/.claude/claude_desktop_config.json` or `settings.json` `mcpServers` key); compare each against `cfg.Allowlist.MCPs`; emit `UNAPPROVED_MCP` HIGH per unknown MCP; if allowlist empty emit `UNCONFIGURED_ALLOWLIST` WARN; register via `init()`
- [x] T018 [P] [US2] Implement `internal/modules/permissions/permissions.go`: for each permission restriction in `settings.json`, attempt a controlled read/stat of a path that should be blocked; emit `PERMISSION_MISCONFIGURED` HIGH if access unexpectedly succeeds; register via `init()`
- [x] T019 [P] [US2] Implement `internal/modules/sandbox/sandbox.go`: verify Claude's sandboxing config keys are present and set to expected restrictive values; emit `SANDBOX_VIOLATION` HIGH if any sandbox config key is absent or permissive; register via `init()`
- [x] T020 [P] [US2] Implement `internal/modules/network/network.go`: read Claude's network log or last-run metadata if available; log each unique destination domain as `NETWORK_REQUEST` INFO Finding with `resource` set to domain; emit nothing if no network log is accessible; register via `init()`
- [x] T021 [US2] Import all US2 modules in `cmd/ai-check-guardrails/main.go`; create `testdata/mcp-test-config.json` with one non-allowlisted MCP; run binary and confirm `UNAPPROVED_MCP` HIGH in output

**Checkpoint**: US1 + US2 independently functional. MCP findings present; approved MCPs produce no findings.

---

## Phase 5: User Story 3 — Policy-Bypass Detection (Priority: P3)

**Goal**: Tool detects `--dangerously-skip-permissions` usage history and missing pre-commit hooks.

**Independent Test**: Remove `gitleaks` hook from a test repo under `testdata/nohooks-repo/.git/`; run binary with `testdata/hooks-test-config.json`; confirm `MISSING_PRECOMMIT_HOOK` HIGH in output. Restore hook; confirm no finding.

### Implementation for User Story 3

- [x] T022 [P] [US3] Implement `internal/modules/bypass/bypass.go`: inspect shell history files (`~/.zsh_history`, `~/.bash_history`) and Claude session logs for occurrences of `--dangerously-skip-permissions`; emit `POLICY_BYPASS` CRITICAL per occurrence found; register via `init()`
- [x] T023 [P] [US3] Implement `internal/modules/hooks/hooks.go`: walk Git repos under `cfg.ScanRoot` (depth limited to 3 levels); for each repo check `.git/hooks/` for each hook name in `cfg.Allowlist.PreCommitHooks`; emit `MISSING_PRECOMMIT_HOOK` HIGH per missing hook with `resource` set to repo path; register via `init()`
- [x] T024 [P] [US3] Implement `internal/modules/humanloop/humanloop.go`: check Claude `settings.json` for human-in-the-loop confirmation keys (e.g., `requireConfirmation`); emit `HUMANLOOP_ABSENT` WARN if key absent or false; register via `init()`
- [x] T025 [P] [US3] Implement `internal/modules/evals/evals.go`: check Claude `settings.json` for Anthropic-recommended eval hook keys; emit `EVAL_HOOK_ABSENT` HIGH for each missing recommended hook; register via `init()`
- [x] T026 [P] [US3] Implement `internal/modules/tokens/tokens.go`: if `cfg.TokenBaseline` nil → emit `MODULE_UNAVAILABLE` INFO and return; otherwise read Claude usage log for today's token count; if count > `mean + multiplier*stddev` emit `TOKEN_ANOMALY` WARN with `confidence` field set to (count−mean)/stddev normalised to [0,1]; register via `init()`
- [x] T027 [US3] Import all US3 modules in `cmd/ai-check-guardrails/main.go`; create `testdata/nohooks-repo/.git/` without gitleaks hook; run binary; confirm `MISSING_PRECOMMIT_HOOK` HIGH in output

**Checkpoint**: US1 + US2 + US3 independently functional. Bypass and missing-hook findings detected correctly.

---

## Phase 6: User Story 4 — Gamification: Per-User Security Score (Priority: P4)

**Goal**: Score reflects weighted findings; banner displays training link when score drops below threshold.

**Independent Test**: Run binary against clean config → confirm score 100, no banner. Introduce two CRITICAL findings by pointing config at broken settings → confirm score ≤ 50, banner appears on stderr with training URL.

### Implementation for User Story 4

- [x] T028 [US4] Validate scoring weights in `internal/score/score.go` against data-model.md thresholds (CRITICAL −25, HIGH −15, WARN −5, INFO 0, floor 0); add edge-case logic: if total deductions > 100 return 0 not negative
- [x] T029 [US4] Update `internal/modules/banner/banner.go` to conditionally display: score 100 → no banner; 70–99 → advisory message; < 70 → warning with `cfg.BannerURL` training link; verify banner outputs to stderr not stdout (stdout is reserved for SIEM JSON)
- [x] T030 [US4] Add `testdata/low-score-config.json` pointing `siem_endpoint` to `http://` (invalid, rejected at config load → exit 2) and a `testdata/two-critical-config.json` that triggers two CRITICAL bypass findings; run against `two-critical-config.json`; confirm score ≤ 50 and banner present

**Checkpoint**: All four user stories independently functional and testable.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Scheduling, enforce mode, and final validation.

- [x] T031 [P] Implement `--install-launchd` in `cmd/ai-check-guardrails/main.go`: render `embed/launchd.plist.tmpl` with binary path and username; write to `~/Library/LaunchAgents/com.example.ai-check-guardrails.plist`; run `launchctl load`; exit 0
- [x] T032 [P] Implement `--install-systemd` in `cmd/ai-check-guardrails/main.go`: write rendered `embed/systemd.service.tmpl` and `embed/systemd.timer.tmpl` to `~/.config/systemd/user/`; print `systemctl --user daemon-reload && systemctl --user enable --now ai-check-guardrails.timer`; exit 0
- [x] T033 [P] Implement `--uninstall` in `cmd/ai-check-guardrails/main.go`: detect platform; run `launchctl unload` + remove plist, or `systemctl --user disable` + remove unit files; exit 0
- [x] T034 Implement enforce mode in `internal/audit/runner.go`: when `cfg.Mode == "enforce"`, after collecting findings check if any CRITICAL finding exists; if so set `AuditRun.ExitCode = 1` (findings already cause exit 1; enforce adds a stderr alert: "ENFORCEMENT: [N] critical finding(s) detected")
- [x] T035 [P] Run `go vet ./...` and fix all reported issues
- [x] T036 Validate quickstart.md scenarios end-to-end: build binary, create minimal config, run, confirm JSON output is valid per contracts/siem-event.md schema, confirm exit codes match contracts/cli-flags.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 completion — BLOCKS all user stories
- **US1 (Phase 3)**: Depends on Phase 2
- **US2 (Phase 4)**: Depends on Phase 2; can run in parallel with US1 if staffed
- **US3 (Phase 5)**: Depends on Phase 2; can run in parallel with US1/US2 if staffed
- **US4 (Phase 6)**: Depends on Phase 3 (score module must exist); sequential
- **Polish (Phase 7)**: Depends on all user story phases complete

### User Story Dependencies

- **US1 (P1)**: Only Phase 2 prerequisite — independent
- **US2 (P2)**: Only Phase 2 prerequisite — independent of US1
- **US3 (P3)**: Only Phase 2 prerequisite — independent of US1/US2
- **US4 (P4)**: Depends on score module from Phase 2 and banner from US1 (T013) — partial dependency

### Parallel Opportunities

- T003, T004 — parallel (different embed files)
- T006, T007, T008, T009 — parallel within Phase 2 (different packages)
- T012, T013, T014 — parallel within US1 (different module files)
- T017, T018, T019, T020 — parallel within US2 (different module files)
- T022, T023, T024, T025, T026 — parallel within US3 (different module files)
- T031, T032, T033, T035 — parallel in Polish phase

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: US1 (settings check + banner + gamification score)
4. **STOP and VALIDATE**: `./ai-check-guardrails` emits valid SIEM JSON event
5. Deploy/demo to security team for feedback

### Incremental Delivery

1. Setup + Foundational → binary builds, exits cleanly
2. Add US1 → first real SIEM event; settings findings visible
3. Add US2 → MCP detection live; allowlist feedback loop begins
4. Add US3 → bypass and hook enforcement detection live
5. Add US4 → score and banner active; training loop begins
6. Polish → scheduling, enforce mode, production-ready

### Parallel Team Strategy

With two engineers after Foundational phase:

- Engineer A: US1 + US4 (visibility + score refinement)
- Engineer B: US2 + US3 (MCP detection + bypass detection)

---

## Notes

- [P] tasks operate on different files with no shared state — safe to run in parallel
- [Story] label maps each task to the user story it delivers
- Each user story phase is independently completable and testable before the next
- Module `init()` registration pattern: each module file calls `modules.Register(...)` in its `init()` — main.go blank-imports it to activate
- Token anomaly module (T026) requires `token_baseline` in config; it self-disables with an INFO finding if missing — safe to enable by default in config once baseline is established
- `go:embed` for schedule templates avoids shipping separate files; binary is fully self-contained
