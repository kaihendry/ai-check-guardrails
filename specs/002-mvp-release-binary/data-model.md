# Data Model: MVP Release Binary & Distribution

**Feature**: 002-mvp-release-binary
**Date**: 2026-05-08

## Entities

### Version String

Represents the build identity of a released binary.

**Format**: `v0.0.YYYYMMDDHHMMSS+<sha7>`

| Field       | Type   | Description                                      | Example                      |
|-------------|--------|--------------------------------------------------|------------------------------|
| prefix      | string | Always `v0.0.` for MVP releases                  | `v0.0.`                      |
| timestamp   | string | UTC build time, `YYYYMMDDHHMMSS`                 | `20260508143022`             |
| sha         | string | 7-char git commit SHA (build metadata)           | `abc1234`                    |

**Embedded at build time** via `-ldflags "-X main.version=$VERSION"`. Default value in source is `"dev"` (used for local builds).

---

### Release Artifact

A single binary file published to GitHub Releases for a specific platform.

| Field    | Type   | Description                                            | Example                                  |
|----------|--------|--------------------------------------------------------|------------------------------------------|
| name     | string | Filename: `ai-check-guardrails-{os}-{arch}`            | `ai-check-guardrails-linux-amd64`        |
| os       | enum   | Target OS: `linux` or `darwin`                         | `linux`                                  |
| arch     | enum   | CPU architecture: `amd64` or `arm64`                   | `arm64`                                  |
| checksum | string | SHA-256 hex digest, stored in `{name}.sha256` sidecar  | `e3b0c44298fc1c14...`                    |
| version  | string | Version string embedded at build time                  | `v0.0.20260508143022+abc1234`            |

---

### GitHub Release

The published release record on GitHub, always pointed to by the `latest` tag.

| Field      | Type     | Description                                      |
|------------|----------|--------------------------------------------------|
| tag_name   | string   | Always `latest` (rolling tag)                    |
| name       | string   | `"Latest build (v0.0.YYYYMMDDHHMMSS+sha7)"`      |
| assets     | []Asset  | One artifact + one checksum file per platform    |
| created_at | datetime | Timestamp of most recent publish                 |

**Asset count per release**: 8 (4 binaries + 4 `.sha256` sidecar files)

---

### Update Check Result

In-memory state computed at startup by the self-update logic.

| Field          | Type    | Description                                        |
|----------------|---------|----------------------------------------------------|
| current_version | string | Version embedded at build time                    |
| latest_version  | string | `tag_name` from GitHub Releases API response      |
| update_needed   | bool   | `true` if `latest_version != current_version`     |
| download_url    | string | `browser_download_url` for the matching asset     |
| error           | error  | Non-nil if the check failed (non-fatal)           |

---

## State Transitions

### Self-Update Flow

```
startup
  │
  ├─► [check disabled?] ──► skip → execute normally
  │
  ├─► GET /releases/latest
  │     ├─ error ──────────────► log warning to stderr → execute normally
  │     └─ success
  │           ├─ current == latest ──► execute normally
  │           └─ current != latest
  │                 ├─ download error ──► log warning to stderr → execute normally
  │                 └─ download OK
  │                       ├─ atomic replace error ──► log warning to stderr → execute normally
  │                       └─ replaced OK ──► log "updated to <version>" → execute normally
  └─► (all paths converge to normal execution)
```

**Invariant**: The tool ALWAYS continues to normal execution regardless of self-update outcome.

---

## Configuration (env vars / flags)

No persistent config file is introduced by this feature. Self-update behavior is controlled by:

| Env Var / Flag           | Type | Default | Effect                                      |
|--------------------------|------|---------|---------------------------------------------|
| `NO_AUTO_UPDATE=1`       | env  | unset   | Skip self-update check entirely             |
| `--no-update` (flag)     | bool | false   | Same as `NO_AUTO_UPDATE=1` for one run      |

These are checked before the HTTP call, so CI environments can disable the check cleanly.
