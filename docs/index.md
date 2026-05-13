# ai-check-guardrails

Audit Claude AI guardrail compliance on developer workstations. The tool runs a set of focused [modules](modules/index.md), each checking a specific security concern, and produces a scored findings report.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/kaihendry/ai-check-guardrails/main/install.sh | bash
```

## Quick start

```bash
ai-check-guardrails
```

The tool reads its configuration from `~/.config/ai-check-guardrails/config.json` and outputs findings to stdout. A security score banner is written to stderr.

## Modules

The tool ships with 13 modules covering the most common Claude security concerns:

- [Browse all modules →](modules/index.md)

## Reference

- [Module Interface](reference/module-interface.md) — how modules and findings are structured
- [Configuration](reference/config.md) — all configuration keys and their defaults
- [Severity Levels](reference/severity.md) — what INFO, WARN, HIGH, and CRITICAL mean

## Releases

Pre-built binaries for Linux and macOS (amd64 and arm64) are published automatically on every commit to `main`. See [GitHub Releases](https://github.com/kaihendry/ai-check-guardrails/releases/tag/latest) for the latest build.
