# Research: AI Guardrails Audit

**Branch**: `001-ai-guardrails-audit` | **Date**: 2026-05-07

## Decision 1: CLI Framework

**Decision**: Use Go stdlib `flag` package only; no third-party CLI framework.

**Rationale**: The tool has a flat flag set (`--mode`, `--config`, `--output`,
`--install-launchd`, `--install-systemd`, `--version`) with no subcommands.
`cobra` adds binary bloat (~100 KB+), supply-chain risk (critical for a security
tool), and complexity that YAGNI prohibits. stdlib `flag` handles the full
required surface area.

**Alternatives considered**:
- `cobra` — feature-rich but 5× complexity overhead for a flat flag set; violates
  Simplicity principle.
- `urfave/cli` — lighter than cobra but still a third-party dependency with no
  gain over stdlib for this use case.

## Decision 2: SIEM Transport

**Decision**: Structured JSON to stdout as primary transport; optional HTTP POST
to a configurable endpoint using stdlib `net/http`. No persistent retry queue in
v1 — failed HTTP delivery is logged to stderr and the run exits with code 1.

**Rationale**: stdout JSON works out of the box with any log aggregator (`tee`,
`jq`, syslog forwarding) and requires zero configuration. HTTP POST is opt-in;
adding a persistent queue in v1 violates Simplicity. If the SIEM is unreachable,
the operator's log pipeline captures stdout; the finding is not lost at the
workstation level.

**Alternatives considered**:
- Always-on HTTP with persistent retry queue — adds queue storage, state
  management, and complexity disproportionate to v1 scope.
- stdout-only forever — limits operational SIEM integration at scale.

## Decision 3: Config File Format

**Decision**: JSON, read from `~/.config/ai-check-guardrails/config.json` by
default (overridable via `--config` flag).

**Rationale**: stdlib `encoding/json` requires no external dependencies. JSON
is already the format of Claude's `settings.json`, establishing organizational
precedent and keeping the operator's mental model consistent. Smaller binary,
zero supply-chain risk.

**Alternatives considered**:
- TOML — requires `BurntSushi/toml` or `pelletier/go-toml`; violates zero-dep goal.
- YAML — requires `gopkg.in/yaml.v3`; complex parser, error-prone user configs.

## Decision 4: Scheduled Execution

**Decision**: Embed `launchd` plist and systemd unit templates via `go:embed`.
Provide `--install-launchd` and `--install-systemd` flags that write the rendered
template to the correct system location and activate the schedule.

**Rationale**: A single binary with embedded templates eliminates missing-file
failures and path breakage during deployment. `go:embed` (GA since Go 1.16, mature
in 1.22) adds < 2 KB to binary size. One command (`ai-check-guardrails --install-launchd`)
is the complete onboarding step.

**Alternatives considered**:
- Ship templates as separate files — adds deployment complexity and breaks if
  files are moved.
- Generate templates from code strings — less readable than actual XML/INI files
  maintained in the `embed/` directory.

## Decision 5: Module Architecture

**Decision**: `Module` interface (`Name() string; Enabled(Config) bool;
Run(Config) ([]Finding, error)`) with `init()`-time registration into a package-level
registry slice. Each module package's `init()` calls `modules.Register(...)`.

**Rationale**: 12 modules is exactly the count where an interface + registry pays
for itself (vs. hardcoded calls) without needing a plugin system. Compile-time
composition: enable/disable by import. Testable in isolation. No `.so` loading,
no factory boilerplate, no filesystem scans for plugins.

**Alternatives considered**:
- Dynamic plugin system (`.so` loading) — version-lock problems, OS-specific build
  complexity, violates Simplicity principle.
- Hardcoded sequential calls in runner — no toggle-ability without code changes;
  scales poorly beyond ~5 modules.

## Decision 6: Run-Lock Mechanism

**Decision**: Lockfile at `$TMPDIR/ai-check-guardrails.lock` using
`syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB)`. If the lock cannot be
acquired, exit immediately with code 2 and emit `ALREADY_RUNNING` to stderr.

**Rationale**: `syscall.Flock` is available on macOS and Linux (primary targets),
requires zero external dependencies, and is automatically released by the kernel
on process exit or crash — no stale lock cleanup needed. The non-blocking
`LOCK_NB` flag gives an immediate, clean exit rather than a hang.

**Alternatives considered**:
- PID file — stale after crashes; requires manual cleanup logic.
- `gofrs/flock` wrapper — adds a third-party dependency for a thin stdlib wrapper.
- Named mutex — Windows-only primitive; not portable.

## Decision 7: Go Version Target

**Decision**: Go 1.22+.

**Rationale**: Go 1.22 introduces improved loop variable semantics and is the
current LTS-equivalent stable release. `go:embed` (1.16+), `log/slog` (1.21+),
and `testing/testscript`-compatible test tooling are all available. macOS 13+
ships compatible toolchains.

**Alternatives considered**:
- Go 1.21 — lacks minor stdlib improvements; no material difference for this tool.
- Go 1.20 or earlier — missing `log/slog` (requires third-party structured logger).
