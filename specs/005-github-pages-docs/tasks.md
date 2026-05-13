---
description: "Tasks for GitHub Pages Module Documentation"
---

# Tasks: GitHub Pages Module Documentation

**Input**: Design documents from `specs/005-github-pages-docs/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/mkdocs-config.md

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create the mkdocs scaffolding — config, Makefile, CI workflow, .gitignore — so the site can build before any content is written.

- [x] T001 Create `mkdocs.yml` at repo root with site name, repo URL, Material theme, and full nav as defined in `specs/005-github-pages-docs/contracts/mkdocs-config.md`
- [x] T002 Create `Makefile` at repo root with `preview-docs` target (`uvx --with mkdocs-material mkdocs serve`) and `build-docs` target (`uvx --with mkdocs-material mkdocs build --strict`)
- [x] T003 [P] Create `.github/workflows/docs.yml` with build-and-deploy job: checkout → setup-python → pip install mkdocs-material → mkdocs build → upload-pages-artifact → deploy-pages; trigger on push to `main`; permissions: `pages: write`, `id-token: write`
- [x] T004 [P] Add `site/` entry to `.gitignore`

**Checkpoint**: `make build-docs` exits 0 (will fail until content files exist — proceed to Phase 2)

---

## Phase 2: Foundational (Placeholder Files)

**Purpose**: Create empty placeholder files for every page declared in `mkdocs.yml` so `make build-docs` exits 0 before content is authored.

**⚠️ CRITICAL**: Placeholders must exist before any user story content is authored; mkdocs fails fast on missing nav files.

- [x] T005 Create `docs/index.md` with placeholder heading `# ai-check-guardrails`
- [x] T006 [P] Create `docs/modules/index.md` with placeholder heading `# Modules`
- [x] T007 [P] Create `docs/modules/banner.md` with placeholder heading `# banner`
- [x] T008 [P] Create `docs/modules/bypass.md` with placeholder heading `# bypass`
- [x] T009 [P] Create `docs/modules/envkeys.md` with placeholder heading `# envkeys`
- [x] T010 [P] Create `docs/modules/evals.md` with placeholder heading `# evals`
- [x] T011 [P] Create `docs/modules/gamification.md` with placeholder heading `# gamification`
- [x] T012 [P] Create `docs/modules/hooks.md` with placeholder heading `# hooks`
- [x] T013 [P] Create `docs/modules/humanloop.md` with placeholder heading `# humanloop`
- [x] T014 [P] Create `docs/modules/mcp.md` with placeholder heading `# mcp`
- [x] T015 [P] Create `docs/modules/network.md` with placeholder heading `# network`
- [x] T016 [P] Create `docs/modules/permissions.md` with placeholder heading `# permissions`
- [x] T017 [P] Create `docs/modules/sandbox.md` with placeholder heading `# sandbox`
- [x] T018 [P] Create `docs/modules/settings.md` with placeholder heading `# settings`
- [x] T019 [P] Create `docs/modules/tokens.md` with placeholder heading `# tokens`
- [x] T020 [P] Create `docs/reference/module-interface.md` with placeholder heading `# Module Interface`
- [x] T021 [P] Create `docs/reference/config.md` with placeholder heading `# Configuration Reference`
- [x] T022 [P] Create `docs/reference/severity.md` with placeholder heading `# Severity Levels`

**Checkpoint**: `make build-docs` exits 0 — all 18 placeholder pages render without errors

---

## Phase 3: User Story 1 — Browse Module Reference (Priority: P1) 🎯 MVP

**Goal**: Each of the 13 modules has a complete documentation page covering purpose, findings table (type / severity / description / remediation), and configuration key. The module index lists all 13. The landing page links to the module index.

**Independent Test**: Run `make preview-docs`, open `http://127.0.0.1:8000/modules/`, click each module name, and confirm each page shows a Findings table with at least one row and a Configuration section listing the enable key.

### Implementation for User Story 1

Read `internal/modules/<name>/<name>.go` for each module to extract finding types, severities, remediation strings, and config keys before authoring its page.

- [x] T023 [P] [US1] Author `docs/modules/banner.md`: purpose (score-based banner display), findings (none from Run — module emits via Display), config key `modules.banner` (default: true), note that threshold comes from `score_threshold`
- [x] T024 [P] [US1] Author `docs/modules/bypass.md`: purpose (detects `--dangerously-skip-permissions` flag in Claude CLI history), findings table from `bypass.go` (`BYPASS_FLAG_DETECTED`, CRITICAL), config key `modules.bypass` (default: true)
- [x] T025 [P] [US1] Author `docs/modules/envkeys.md`: purpose (detects API keys in environment), findings table from `envkeys.go` (key present = INFO/WARN finding), config key `modules.env_keys` (default: true), extra key `env_key_watch_list` (default: `["ANTHROPIC_API_KEY"]`)
- [x] T026 [P] [US1] Author `docs/modules/evals.md`: purpose (checks for Anthropic-recommended hooks: PreToolUse, PostToolUse, Stop), findings table from `evals.go`, config key `modules.evals` (default: false)
- [x] T027 [P] [US1] Author `docs/modules/gamification.md`: purpose (emits posture score INFO finding to surface score in output), findings table from `gamification.go` (`SCORE_INFO`, INFO), config key `modules.gamification` (default: true)
- [x] T028 [P] [US1] Author `docs/modules/hooks.md`: purpose (checks that required pre-commit hooks are installed), findings table from `hooks.go`, config key `modules.hooks` (default: true), extra key `allowlist.precommit_hooks` (default: `["gitleaks"]`)
- [x] T029 [P] [US1] Author `docs/modules/humanloop.md`: purpose (checks Claude settings for human-in-the-loop approval configuration), findings table from `humanloop.go`, config key `modules.humanloop` (default: false)
- [x] T030 [P] [US1] Author `docs/modules/mcp.md`: purpose (audits MCP servers against approved allowlist), findings table from `mcp.go` (`UNCONFIGURED_ALLOWLIST` WARN, unapproved MCP HIGH), config key `modules.mcp` (default: true), extra key `allowlist.mcps`
- [x] T031 [P] [US1] Author `docs/modules/network.md`: purpose (scans Claude log files for outbound URLs and compares against approved domains), findings table from `network.go`, config key `modules.network` (default: false), extra key `allowlist.domains`
- [x] T032 [P] [US1] Author `docs/modules/permissions.go`: purpose (checks Claude settings.json for overly broad tool permissions), findings table from `permissions.go`, config key `modules.permissions` (default: true)
- [x] T033 [P] [US1] Author `docs/modules/sandbox.md`: purpose (checks whether Claude sandbox mode is enabled in settings), findings table from `sandbox.go`, config key `modules.sandbox` (default: true)
- [x] T034 [P] [US1] Author `docs/modules/settings.md`: purpose (validates Claude settings.json exists and is valid), findings table from `settings.go`, config key `modules.settings` (default: true)
- [x] T035 [P] [US1] Author `docs/modules/tokens.md`: purpose (detects anomalous token usage by comparing daily usage against a configured baseline), findings table from `tokens.go` (`MODULE_UNAVAILABLE` INFO when no baseline, anomaly HIGH when exceeded), config key `modules.tokens` (default: false), extra key `token_baseline` (object with `daily_mean`, `std_dev`, `multiplier`)
- [x] T036 [US1] Author `docs/modules/index.md`: module list table with columns: Module | Purpose (one line) | Default State; list all 13 modules in alphabetical order; link each name to its page
- [x] T037 [US1] Author `docs/index.md`: project overview (one paragraph), quick-start install command from README, link to Modules overview, link to GitHub Releases; no duplicate content from module pages

**Checkpoint**: All 13 module pages render in `make preview-docs`; each has a non-empty Findings table and a Configuration section; `docs/modules/index.md` links to all 13

---

## Phase 4: User Story 2 — Module Interface Reference (Priority: P2)

**Goal**: A contributor can understand the `Module` interface contract and `Finding` struct from the documentation site without opening source code.

**Independent Test**: Open `http://127.0.0.1:8000/reference/module-interface/`, confirm it documents `Name() string`, `Enabled(cfg Config) bool`, `Run(cfg Config) ([]Finding, error)`, and all 7 `Finding` fields with types.

### Implementation for User Story 2

Read `internal/modules/module.go` as the authoritative source for interface and struct definitions.

- [x] T038 [P] [US2] Author `docs/reference/module-interface.md`: document the `Module` interface (Name/Enabled/Run with parameter types and return types), document the `Finding` struct (all 7 fields: Type, Severity, Module, Resource, Description, Remediation, Confidence with their types and whether optional), include a minimal example of a compliant module (prose — no compilable code required), link to severity reference
- [x] T039 [P] [US2] Author `docs/reference/severity.md`: table of all 4 severity levels (INFO / WARN / HIGH / CRITICAL), meaning of each, when each is used, and how each contributes to the posture score; note that score threshold default is 70

**Checkpoint**: Both reference pages render; Finding struct table has 7 rows; Module interface section documents all 3 methods

---

## Phase 5: User Story 3 — Config Options per Module (Priority: P3)

**Goal**: An operator can find every config key and its default value from a single reference page, and each module page cross-references its own keys.

**Independent Test**: Open `http://127.0.0.1:8000/reference/config/`, confirm table lists all config keys including `token_baseline`, `env_key_watch_list`, `allowlist.*`, and `modules.*` keys with their types and defaults.

### Implementation for User Story 3

Read `internal/config/config.go` as the authoritative source for all config keys, types, and defaults.

- [x] T040 [US3] Author `docs/reference/config.md`: full config key reference table with columns: Key | Type | Default | Description; cover all top-level keys (`mode`, `siem_endpoint`, `scan_root`, `banner_url`, `score_threshold`, `modules.*`, `allowlist.*`, `token_baseline`, `env_key_watch_list`); note env override `AI_GUARDRAILS_SIEM_ENDPOINT`; note config file location (`~/.config/ai-check-guardrails/config.json`)

**Checkpoint**: Config reference page renders; table has at least 20 rows; `token_baseline` sub-keys documented

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Validate completeness and wire up CI.

- [x] T041 Run `make build-docs` and confirm exit code 0 with `--strict` flag (no broken links, no missing files)
- [x] T042 [P] Cross-check each module page's Findings table against its source file in `internal/modules/<name>/<name>.go` — confirm no finding types are omitted; update any that are missing
- [x] T043 [P] Verify `docs/modules/index.md` lists exactly 13 modules and all links resolve
- [x] T044 [P] Verify GitHub Actions workflow `.github/workflows/docs.yml` uses `actions/deploy-pages@v4` and has correct permissions (`pages: write`, `id-token: write`)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on T001 (mkdocs.yml must exist before `make build-docs` can be tested)
- **User Stories (Phase 3, 4, 5)**: All depend on Phase 2 completion (placeholder files must exist)
  - US1 (Phase 3), US2 (Phase 4), US3 (Phase 5) are independent of each other after Phase 2
- **Polish (Phase 6)**: Depends on all user story phases completing

### User Story Dependencies

- **US1 (P1)**: Can start after Phase 2 — no dependency on US2 or US3
- **US2 (P2)**: Can start after Phase 2 — no dependency on US1 or US3
- **US3 (P3)**: Can start after Phase 2 — no dependency on US1 or US2; T040 does not touch module pages

### Within Each User Story

- All 13 module page tasks (T023–T035) are fully parallel — different files, no shared state
- T036 (module index) depends on T023–T035 being drafted so content can be summarized
- T037 (landing page) depends on T036 (links to module index)

### Parallel Opportunities

- T003, T004 (Phase 1) parallel with T001 and T002
- T006–T022 (Phase 2 placeholders) all parallel after T001
- T023–T035 (US1 module pages) all parallel after Phase 2
- T038, T039 (US2 reference pages) parallel with each other and with US1 tasks
- T042, T043, T044 (Phase 6) parallel with each other after prior phases

---

## Parallel Example: User Story 1

```text
# All 13 module pages can be authored simultaneously:
Task T023: docs/modules/banner.md
Task T024: docs/modules/bypass.md
Task T025: docs/modules/envkeys.md
Task T026: docs/modules/evals.md
Task T027: docs/modules/gamification.md
Task T028: docs/modules/hooks.md
Task T029: docs/modules/humanloop.md
Task T030: docs/modules/mcp.md
Task T031: docs/modules/network.md
Task T032: docs/modules/permissions.md
Task T033: docs/modules/sandbox.md
Task T034: docs/modules/settings.md
Task T035: docs/modules/tokens.md

# Then sequentially:
Task T036: docs/modules/index.md  (after module pages drafted)
Task T037: docs/index.md          (after module index drafted)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational placeholders
3. Complete Phase 3: US1 module pages (13 pages + index + landing)
4. **STOP and VALIDATE**: `make preview-docs` — all module pages browsable, Findings tables complete
5. Merge to `main` — GitHub Pages publishes automatically

### Incremental Delivery

1. Phase 1 + 2 → infrastructure ready
2. Phase 3 → Module reference browsable (US1 done) → **MVP deployed**
3. Phase 4 → Module interface documented (US2 done)
4. Phase 5 → Config reference published (US3 done)
5. Phase 6 → Validation passes, CI confirmed green

---

## Notes

- [P] tasks operate on different files — no conflicts
- Module page content is derived from source files in `internal/modules/<name>/<name>.go` — read each file before authoring its doc page
- `make build-docs` uses `--strict` to catch broken links at build time
- Do not reference Go implementation details (package names, imports) in user-facing doc pages — describe behaviour, not code
