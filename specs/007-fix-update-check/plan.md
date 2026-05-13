# Implementation Plan: Fix Self-Update Version Check

**Branch**: `007-fix-update-check` | **Date**: 2026-05-13 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/007-fix-update-check/spec.md`

## Summary

The self-update mechanism always re-downloads the binary because the GitHub release tag (`latest`) never matches the version embedded in the binary (`v0.0.TIMESTAMP+HASH`). The fix is a one-line change in the release pipeline: use the computed `$VERSION` string as the release tag instead of the static string `"latest"`.

## Technical Context

**Language/Version**: Go 1.24 (see go.mod)
**Primary Dependencies**: stdlib only (`net/http`, `os`, `encoding/json`) — no new dependencies
**Storage**: N/A
**Testing**: `go test ./cmd/ai-check-guardrails/...` (existing test suite)
**Target Platform**: Linux/macOS, amd64/arm64 (GitHub Actions, Ubuntu runner)
**Project Type**: CLI tool
**Performance Goals**: Update check must complete without download in under 3 seconds
**Constraints**: No new dependencies; fix must be a surgical change
**Scale/Scope**: Single file change in `.github/workflows/release.yml`

## Constitution Check

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Simplicity — one change, no new abstractions | ✅ Pass | Single line changed in release.yml; no new code |
| II. Integrity — no silent failures | ✅ Pass | Existing error handling in update.go is unchanged |
| III. Documentation — no module doc changes | ✅ Pass | No module added or modified |
| Security — no new external input surfaces | ✅ Pass | No change to input handling |
| CLI Standards — no interface changes | ✅ Pass | No flags, output, or exit codes affected |

No violations. Complexity Tracking table omitted.

## Project Structure

### Documentation (this feature)

```text
specs/007-fix-update-check/
├── plan.md              # This file
├── research.md          # Phase 0 output
└── tasks.md             # Phase 2 output (/speckit-tasks command)
```

No data-model.md or contracts/ — this fix touches no data entities or external interfaces.

### Source Code (affected files only)

```text
.github/workflows/
└── release.yml          # tag_name: latest → tag_name: ${{ env.VERSION }}
```

No changes to `cmd/`, `internal/`, or `install.sh`.

## Implementation Steps

### Step 1 — Fix release pipeline (the only change required)

In `.github/workflows/release.yml`, change:

```yaml
      - name: Publish release
        uses: softprops/action-gh-release@v2
        with:
          tag_name: latest
```

to:

```yaml
      - name: Publish release
        uses: softprops/action-gh-release@v2
        with:
          tag_name: ${{ env.VERSION }}
```

This ensures:
- Each push to main creates a release tagged with the exact version string embedded in the binary
- GitHub's `/releases/latest` API returns the most recent release with a `tag_name` that matches `currentVersion`
- The comparison in `update.go` (`rel.TagName == currentVersion`) evaluates correctly

### Step 2 — Verify existing tests still pass

Run:
```
go test ./cmd/ai-check-guardrails/...
```

Both existing tests must pass:
- `TestCheckAndUpdate_UpdatesWhenNewVersionAvailable` — confirms update still runs for an outdated binary
- `TestCheckAndUpdate_NoopWhenAlreadyCurrent` — confirms no download when already up to date

No new tests needed: the existing test for the no-op path already validates the fixed behaviour end-to-end.

### Step 3 — Verify post-merge behaviour

After the fix is merged to main:
1. A new release is published with `tag_name = v0.0.TIMESTAMP+HASH`
2. Install the new binary via `install.sh`
3. Run the binary — confirm no `[update]` lines appear on stderr
4. Run it again — confirm still no `[update]` lines

## Out of Scope

- The release name field (`name: "Latest (${{ env.VERSION }})"`) — cosmetic, no functional impact
- Cleanup of old `latest`-tagged releases on GitHub — stale but harmless
- Changes to `install.sh` — already works by fetching the tag dynamically
- Changes to `update.go` — the comparison logic is correct; only the data it compares is wrong
