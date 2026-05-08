# CLI Contract: ai-check-guardrails

**Feature**: 002-mvp-release-binary
**Date**: 2026-05-08

This document defines the observable interface contract for the `ai-check-guardrails` binary as extended by the release binary feature. Downstream tools (scripts, CI pipelines, system integrations) must be able to rely on these contracts remaining stable.

---

## Startup Output Contract

On every invocation, the first line written to stderr MUST be the version string in this format:

```
ai-check-guardrails v0.0.20260508143022+abc1234
```

**Guarantees**:
- Version line is always the first output (before any audit output or update notices)
- Written to stderr so it does not pollute JSON stdout output
- Format: `ai-check-guardrails <version>` — single space, no brackets

---

## Self-Update Notices

When a self-update occurs, additional lines are written to stderr AFTER the version line:

```
ai-check-guardrails v0.0.20260508143022+abc1234
[update] new version available: v0.0.20260509010101+def5678
[update] downloading...
[update] updated successfully, continuing...
```

When an update check fails (non-fatal):

```
ai-check-guardrails v0.0.20260508143022+abc1234
[update] check failed: <reason> (continuing with current version)
```

**Guarantees**:
- Update notices are always on stderr, never stdout
- Update failures do NOT affect exit code
- After any update path, normal execution always follows

---

## Exit Codes

| Code | Meaning                                    |
|------|--------------------------------------------|
| 0    | Success — audit completed                  |
| 1    | Failure — audit detected issues or runtime error |
| 2    | Misuse — bad arguments or invalid config   |

Self-update failures do NOT change the exit code.

---

## Flags Added by This Feature

```
--no-update    Disable self-update check for this invocation (boolean, default false)
```

Environment variable: `NO_AUTO_UPDATE=1` disables self-update check.

---

## Install Script Contract

The `install.sh` script at the repository root is the canonical installation method.

**Invocation**:
```bash
curl -fsSL https://raw.githubusercontent.com/kaihendry/ai-check-guardrails/main/install.sh | bash
```

**Guarantees**:
- Installs to `$HOME/.local/bin` by default (no sudo required)
- Install location overridable via `INSTALL_DIR` environment variable
- Prints the installed version and path on success
- Exits non-zero with a message on any failure
- Verifies binary SHA-256 checksum before installing
- Supported platforms: Linux amd64/arm64, macOS amd64/arm64

**Binary naming convention** (GitHub Releases asset names):
```
ai-check-guardrails-linux-amd64
ai-check-guardrails-linux-arm64
ai-check-guardrails-darwin-amd64
ai-check-guardrails-darwin-arm64
```

Each binary has a sidecar checksum file: `ai-check-guardrails-{os}-{arch}.sha256`

---

## GitHub Releases API Contract

Self-update checks via:

```
GET https://api.github.com/repos/kaihendry/ai-check-guardrails/releases/latest
Accept: application/vnd.github+json
```

The tool parses `tag_name` and `assets[].browser_download_url` from the response.

**Assumed stable**: This GitHub API endpoint and response shape are treated as a stable external contract.
