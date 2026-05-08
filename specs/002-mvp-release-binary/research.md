# Research: MVP Release Binary & Distribution

**Feature**: 002-mvp-release-binary
**Date**: 2026-05-08

## Decision 1: Release Tooling

**Decision**: Raw GitHub Actions cross-compile — no GoReleaser

**Rationale**: GoReleaser adds a `.goreleaser.yml` config file, a pinned external CI action, and an outbound dependency on goreleaser.com for zero functional gain. Go's built-in cross-compilation via `GOOS`/`GOARCH` covers all four targets natively from a single `ubuntu-latest` runner. The project constitution's Simplicity principle ("prefer stdlib and existing dependencies over new ones") makes this the correct choice.

**Alternatives considered**:
- GoReleaser: worthwhile only when adding Homebrew tap, Scoop manifest, Docker images, or Cosign signing — none of which apply here.
- Custom build scripts: unnecessary given GitHub Actions matrix syntax handles the same loop.

---

## Decision 2: Version Scheme

**Decision**: `v0.0.YYYYMMDDHHMMSS+<sha7>` timestamp-based version embedded via `-ldflags "-X main.version=$VERSION"`

**Rationale**: Per-commit releases on main need a monotonically increasing, human-readable, sortable version with no manual state to manage. A timestamp satisfies all three. The short git SHA as build metadata provides traceability.

**Alternatives considered**:
- Semantic versioning with tags: requires a human to tag; adds friction and violates the "every main commit = release" requirement.
- Sequential build number: requires a counter stored in the repo or an external service.
- Git commit SHA only: not sortable or human-readable.

---

## Decision 3: GitHub Release Tag Strategy

**Decision**: Rolling `latest` tag — assets overwritten on every release

**Rationale**: The install script and self-update mechanism need a stable, always-current download URL. A rolling `latest` tag provides this. Each push to main overwrites the release assets, keeping the tag pointer and download URL constant.

**Alternatives considered**:
- Per-commit SHA tags: creates tag clutter on main; all historical releases accumulate; install script URL changes each release. Suitable if audit trail is required — defer to a later iteration.

---

## Decision 4: Self-Update Implementation

**Decision**: stdlib-only HTTP + `os.Rename` atomic swap in a single `update.go` file

**Rationale**: Keeps the zero-dependency constraint. The atomic rename pattern (`os.CreateTemp` in the same directory as the executable → write → `os.Rename`) is a single `rename(2)` syscall on Linux/macOS and is safe without root when the user owns the install directory. `go-update` library would be the alternative but adds an external dependency.

**Key implementation patterns**:
- Version check: `GET https://api.github.com/repos/kaihendry/ai-check-guardrails/releases/latest` → parse `tag_name` with `encoding/json`
- Asset selection: binary named `ai-check-guardrails-{os}-{arch}` matched from `assets[].browser_download_url`
- Atomic replace: `os.CreateTemp(filepath.Dir(exe), ".update-*")` → write → `Chmod(0755)` → `os.Rename`
- Non-fatal: any error during update check/download is logged to stderr and execution continues

**Alternatives considered**:
- `inconshreveable/go-update`: handles more edge cases but adds a dependency.
- Subprocess to install script: fragile, requires trusting external shell execution.

---

## Decision 5: CI/CD Action for GitHub Releases

**Decision**: `softprops/action-gh-release@v2` with `tag_name: latest`

**Rationale**: Native create-or-update semantics without shell gymnastics. The `gh` CLI cannot atomically upsert a release (`gh release delete && gh release create` is racy). The action handles both creation and update in one step.

**Alternatives considered**:
- `gh release create` CLI: verbose, non-atomic upsert.
- `actions/upload-release-asset`: deprecated.

---

## Decision 6: Install Script Pattern

**Decision**: Single `install.sh` at repo root, fetched via `curl | bash`

**Pattern**:
1. Detect `$(uname -s)` → lowercase OS name (`linux` / `darwin`)
2. Detect `$(uname -m)` → normalize to `amd64` / `arm64`
3. Fetch latest tag from GitHub Releases API
4. Download `ai-check-guardrails-{os}-{arch}`
5. Verify SHA-256 checksum
6. Install to `${INSTALL_DIR:-$HOME/.local/bin}`

**Checksum note**: Linux uses `sha256sum`, macOS uses `shasum -a 256` — the script branches on `$OS`.

**Install location**: Default to `~/.local/bin` (user-writable, no sudo required) with a fallback suggestion for `/usr/local/bin` if the user prefers system-wide install.
