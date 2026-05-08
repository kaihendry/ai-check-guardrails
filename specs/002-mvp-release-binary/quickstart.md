# Quickstart: MVP Release Binary & Distribution

**Feature**: 002-mvp-release-binary
**Date**: 2026-05-08

## Install (users)

```bash
curl -fsSL https://raw.githubusercontent.com/kaihendry/ai-check-guardrails/main/install.sh | bash
```

Installs to `~/.local/bin` by default. Override with:

```bash
INSTALL_DIR=/usr/local/bin curl -fsSL .../install.sh | bash
```

## Run

```bash
ai-check-guardrails [options]
```

On first run you'll see:
```
ai-check-guardrails v0.0.20260508143022+abc1234
[update] new version available: v0.0.20260509010101+def5678
[update] updated successfully, continuing...
```

## Disable self-update

```bash
NO_AUTO_UPDATE=1 ai-check-guardrails
# or
ai-check-guardrails --no-update
```

## Build locally (developers)

```bash
go build -ldflags "-X main.version=dev-local" -o ai-check-guardrails ./cmd/ai-check-guardrails
```

## Release a new version (automated)

Push or merge to `main`. The GitHub Actions `release.yml` workflow:
1. Builds binaries for all 4 platforms
2. Embeds a timestamp+SHA version string
3. Publishes to GitHub Releases under the `latest` tag

No manual tagging or version bumping required.

## Repository structure (this feature adds)

```
.github/
└── workflows/
    └── release.yml          # CI pipeline: build + publish on every main commit
cmd/
└── ai-check-guardrails/
    ├── main.go              # var version = "dev" + startup version log
    └── update.go            # self-update logic (stdlib only)
install.sh                   # curl-installable install script
README.md                    # install instructions + usage
```

## Key files to modify

| File | Change |
|------|--------|
| `cmd/ai-check-guardrails/main.go` | Add `var version = "dev"` and version log at startup; add `--no-update` flag |
| `cmd/ai-check-guardrails/update.go` | New file: self-update logic |
| `.github/workflows/release.yml` | New file: CI release pipeline |
| `install.sh` | New file: curl install script |
| `README.md` | New file: install instructions and usage |
