# Tasks: MVP Release Binary & Distribution

**Input**: Design documents from `specs/002-mvp-release-binary/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/cli-contract.md ✅, quickstart.md ✅

**Tests**: Update unit test included for US3 (self-update) to verify atomic replace behavior.

**Organization**: Tasks grouped by user story — each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: User story label — US1/US2/US3/US4 (maps to spec.md priority P1–P4)
- All file paths are relative to repo root

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm build baseline and create CI directory structure

- [x] T001 Confirm `go build ./cmd/ai-check-guardrails` succeeds on current branch and create `.github/workflows/` directory at repo root

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Wire the version variable injection point and the no-update bypass — required by US2, US3, and US4 before any story-specific work begins

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T002 Add `var version = "dev"` package-level variable declaration to `cmd/ai-check-guardrails/main.go` (place before `func main()`; this is the `-ldflags "-X main.version=$VERSION"` injection point for all release builds)
- [x] T003 [P] Add `--no-update` boolean flag (default `false`) and `noAutoUpdate()` helper (checks `os.Getenv("NO_AUTO_UPDATE") == "1"`) to `cmd/ai-check-guardrails/main.go` flag parsing block — both must exist before US3 wires self-update into startup

**Checkpoint**: Foundation ready — `var version` exists, `--no-update` flag parses, build still passes

---

## Phase 3: User Story 1 — Install via curl (Priority: P1) 🎯 MVP

**Goal**: A developer copies one `curl | bash` command from the README, runs it, and has a working binary in `~/.local/bin` in under 30 seconds.

**Independent Test**: On a clean machine with only `curl` and `sh`, run the install command from `README.md` and confirm the binary is present and executable; verify SHA-256 checksum is validated during install.

### Implementation for User Story 1

- [x] T004 [P] [US1] Create `install.sh` at repo root: detect OS via `$(uname -s | tr '[:upper:]' '[:lower:]')` → `linux`/`darwin`; detect arch via `$(uname -m)` normalised to `amd64`/`arm64` (fail with message on unsupported arch); fetch `tag_name` from `https://api.github.com/repos/kaihendry/ai-check-guardrails/releases/latest` via `curl -fsSL ... | grep tag_name | sed ...`; download `ai-check-guardrails-{os}-{arch}` and `ai-check-guardrails-{os}-{arch}.sha256` sidecar; verify checksum (`sha256sum -c` on Linux, `shasum -a 256 -c` on macOS); install to `${INSTALL_DIR:-$HOME/.local/bin}` with `install -m 755`; print installed path and version on success; exit non-zero with message on any failure
- [x] T005 [P] [US1] Create `README.md` at repo root: one-liner install command block (`curl -fsSL https://raw.githubusercontent.com/kaihendry/ai-check-guardrails/main/install.sh | bash`); supported platforms table (Linux/macOS × amd64/arm64); usage example showing startup version line; `INSTALL_DIR` override example; `--no-update` and `NO_AUTO_UPDATE=1` disable examples; link to GitHub Releases page

**Checkpoint**: User Story 1 is independently testable — manually publish a GitHub Release with binaries present, then verify `install.sh` installs and verifies the binary

---

## Phase 4: User Story 2 — Version Logged on Startup (Priority: P2)

**Goal**: Every invocation of `ai-check-guardrails` prints the version string as the very first line of stderr output.

**Independent Test**: Run `go run ./cmd/ai-check-guardrails 2>&1 | head -1` and confirm output matches `ai-check-guardrails <version>`.

### Implementation for User Story 2

- [x] T006 [US2] Add `fmt.Fprintf(os.Stderr, "ai-check-guardrails %s\n", version)` as the absolute first statement of `func main()` in `cmd/ai-check-guardrails/main.go`, before flag parsing, before any other output — version line MUST precede all other stderr/stdout

**Checkpoint**: `go build -ldflags "-X main.version=v0.0.test" ./cmd/ai-check-guardrails && ./ai-check-guardrails 2>&1 | head -1` outputs `ai-check-guardrails v0.0.test`

---

## Phase 5: User Story 3 — Self-Update on Startup (Priority: P3)

**Goal**: When a newer release is available, the binary replaces itself atomically before continuing normal execution; failures are non-fatal and logged to stderr.

**Independent Test**: Install an older build (e.g. version `dev`), publish a real or mock newer release, run the binary, and observe it logs `[update] updated successfully` and continues to execute normally. Also verify that with `NO_AUTO_UPDATE=1` no HTTP call is made.

### Implementation for User Story 3

- [x] T007 [US3] Create `cmd/ai-check-guardrails/update.go` with the following: `type githubRelease struct { TagName string \`json:"tag_name"\`; Assets []struct{ Name string \`json:"name"\`; BrowserDownloadURL string \`json:"browser_download_url"\` } \`json:"assets"\` }` and `func checkAndUpdate(currentVersion string) error` — GET `https://api.github.com/repos/kaihendry/ai-check-guardrails/releases/latest` with 10s timeout, decode JSON, compare `rel.TagName` to `currentVersion` (skip if equal), find asset whose `Name` equals `ai-check-guardrails-{runtime.GOOS}-{runtime.GOARCH}`, download binary bytes via second HTTP GET, write to `os.CreateTemp(filepath.Dir(exe), ".update-*")` (same filesystem as exe for atomic rename), call `tmp.Chmod(0755)`, `tmp.Close()`, then `os.Rename(tmp.Name(), exe)`; log each step via `fmt.Fprintf(os.Stderr, "[update] ...\n")`; return error on any failure
- [x] T008 [US3] Wire `checkAndUpdate(version)` into `func main()` in `cmd/ai-check-guardrails/main.go` immediately after the version print (T006): call only when `!*noUpdate && !noAutoUpdate()`; if `checkAndUpdate` returns a non-nil error, log `fmt.Fprintf(os.Stderr, "[update] check failed: %v (continuing)\n", err)` and continue — NEVER exit or change exit code due to update failure
- [x] T009 [P] [US3] Create `cmd/ai-check-guardrails/update_test.go`: use `net/http/httptest.NewServer` to serve a fake GitHub Releases API response with `tag_name: "v0.0.99999999999999"` and a fake binary asset; call `checkAndUpdate("dev")` with the test server URL overridden; verify the temp file was written and renamed; add a second sub-test where `currentVersion` already matches `tag_name` and verify no download occurs; run with `go test ./cmd/ai-check-guardrails/...`

**Checkpoint**: `go test ./cmd/ai-check-guardrails/...` passes; manually run the binary and observe self-update notices on stderr followed by normal execution

---

## Phase 6: User Story 4 — Automated Release Pipeline (Priority: P4)

**Goal**: Every commit to main automatically produces a new GitHub Release with binaries for all four supported platforms, without any manual steps.

**Independent Test**: Push a commit to main and confirm a new GitHub Release appears within 15 minutes with 8 assets (4 binaries + 4 `.sha256` files) and the version embedded in the binary matches the release timestamp.

### Implementation for User Story 4

- [x] T010 [US4] Create `.github/workflows/release.yml`: trigger `on: push: branches: [main]`; `permissions: contents: write`; single job `release` on `ubuntu-latest`; steps — `actions/checkout@v4`, `actions/setup-go@v5` with `go-version-file: go.mod`, compute version (`echo "VERSION=v0.0.$(date -u '+%Y%m%d%H%M%S')+$(git rev-parse --short HEAD)" >> "$GITHUB_ENV"`), build loop over `linux/amd64 linux/arm64 darwin/amd64 darwin/arm64` (`GOOS={os} GOARCH={arch} go build -ldflags "-s -w -X main.version=$VERSION" -o "dist/ai-check-guardrails-{os}-{arch}" ./cmd/ai-check-guardrails`), generate SHA-256 sidecar for each binary (`sha256sum dist/ai-check-guardrails-{os}-{arch} > dist/ai-check-guardrails-{os}-{arch}.sha256`), publish via `softprops/action-gh-release@v2` with `tag_name: latest`, `name: "Latest (${{ env.VERSION }})"`, `make_latest: true`, `files: dist/*`
- [ ] T011 [P] [US4] Push a test commit to main (e.g., add a comment to `README.md`), observe the `release` workflow run in GitHub Actions, and confirm all 8 release assets appear in the `latest` GitHub Release and the embedded version string matches the expected `v0.0.YYYYMMDDHHMMSS+sha7` format

**Checkpoint**: GitHub Release at `https://github.com/kaihendry/ai-check-guardrails/releases/tag/latest` shows 8 assets, download + run binary prints correct version

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Validate correctness and verify end-to-end integration across all user stories

- [x] T012 [P] Run `go vet ./...` and `go test ./...` from repo root and fix any issues introduced by changes to `cmd/ai-check-guardrails/main.go`, new `update.go`, or `update_test.go`
- [ ] T013 [P] Execute end-to-end validation per `quickstart.md`: run the `curl | bash` install command against the live `latest` release, verify version printed on startup, verify self-update applies when a locally-built older binary is installed
- [x] T014 [P] Update `specs/002-mvp-release-binary/checklists/requirements.md` — mark all Feature Readiness items complete and add post-implementation verification notes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 — **BLOCKS all user stories**
- **US1 (Phase 3)**: Depends on Foundational — no dependency on US2/US3/US4
- **US2 (Phase 4)**: Depends on Foundational (needs `var version`) — no dependency on US1/US3/US4
- **US3 (Phase 5)**: Depends on Foundational AND US2 (version print must precede update wire-up) — no dependency on US1/US4
- **US4 (Phase 6)**: Depends on Foundational AND US2 (version variable must exist for `-ldflags`) — no dependency on US1/US3
- **Polish (Phase 7)**: Depends on all desired stories being complete

### User Story Dependencies

- **US1 (P1)**: Independent after Foundational — requires a published release to fully validate (provided by US4)
- **US2 (P2)**: Independent after Foundational — testable with `go run` alone
- **US3 (P3)**: Depends on US2 (version print ordering in main.go) — testable with `httptest` mock
- **US4 (P4)**: Depends on US2 (version variable embedding) — testable by pushing to main

### Within Each User Story

- US3: T007 (create update.go) → T008 (wire into main.go) → T009 (test) in that order
- US4: T010 (create release.yml) → T011 (validate by pushing) in that order
- All other stories: tasks within each story are independently orderable

---

## Parallel Opportunities

```bash
# Phase 3 (US1): T004 and T005 can run in parallel
Task T004: "Create install.sh at repo root"
Task T005: "Create README.md at repo root"

# Phase 2: T002 and T003 can run in parallel
Task T002: "Add var version to main.go"
Task T003: "Add --no-update flag to main.go"

# Polish phase: T012, T013, T014 can all run in parallel
```

---

## Implementation Strategy

### MVP First (User Story 1 + US2 Only)

1. Complete Phase 1: Setup (T001)
2. Complete Phase 2: Foundational (T002, T003)
3. Complete Phase 4: US2 — Version logging (T006)
4. Complete Phase 6: US4 — CI pipeline (T010, T011)
5. Complete Phase 3: US1 — Install script (T004, T005)
6. **STOP and VALIDATE**: Curl-install from published release, verify version prints
7. Demo if ready

### Incremental Delivery

1. Setup + Foundational → baseline builds
2. US2 + US4 → binary builds in CI, version prints, published release
3. US1 → users can install via curl (MVP milestone)
4. US3 → users get automatic updates on every run
5. Polish → end-to-end validation complete

---

## Notes

- `[P]` tasks touch different files and have no incomplete dependencies — safe to parallelize
- US3 (self-update) is the most complex story; its test (T009) uses `httptest` to avoid real network calls
- The `update.go` test must override the GitHub API URL — define a `var githubAPIBase = "https://api.github.com"` in `update.go` that tests can swap to point at `httptest.NewServer`
- CI pipeline uses a rolling `latest` tag — each push overwrites assets; no tag clutter
- Self-update never changes exit code; constitution gate "no silent failures" is satisfied by stderr logging
