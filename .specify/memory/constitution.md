<!--
SYNC IMPACT REPORT
==================
Version change: 1.0.0 → 1.1.0

Modified principles: N/A

Added sections:
  - Core Principles: I. Simplicity, II. Integrity, III. Documentation (new)
  - Security Engineering Standards
  - CLI Interface Standards
  - Governance

Removed sections: N/A

Templates requiring updates:
  - .specify/templates/plan-template.md  ✅ Constitution Check table picks up new principle automatically
  - .specify/templates/spec-template.md  ✅ No changes needed
  - .specify/templates/tasks-template.md ✅ No changes needed

Deferred TODOs: none
-->

# ai-check-guardrails Constitution

## Core Principles

### I. Simplicity

Every component MUST do one thing and do it well. Complexity MUST be justified with
a documented, present-day need — not anticipated future use.

- Implementations MUST prefer stdlib and existing dependencies over new ones.
- Abstractions MUST emerge from duplication, never from speculation (YAGNI).
- A simpler solution that covers 90% of cases is preferred over a perfect solution
  that covers 100%.
- Complexity violations MUST be recorded in the plan's Complexity Tracking table.

**Rationale**: CLI tools accumulate debt fast. Keeping scope minimal ensures the tool
remains auditable, maintainable, and trustworthy to its security-focused users.

### II. Integrity

The tool's output MUST be accurate, auditable, and honest about uncertainty.

- Silent failures are PROHIBITED. All errors MUST surface to stderr with a non-zero
  exit code.
- Output MUST be deterministic for the same input; no hidden side effects.
- Where certainty is partial (e.g., heuristic detection), the tool MUST communicate
  confidence level explicitly.
- Security-relevant findings MUST NOT be suppressed for UX convenience.

**Rationale**: A security tool that hides failures or hedges silently is worse than
no tool. Users rely on this output to make trust decisions.

### III. Documentation

The package-level Go doc comment in each module is the **single source of truth** for
that module's documentation. `docs/modules/<name>.md` is a generated file — produced by
`make gen-docs` (`cmd/gendocs`) — and MUST NOT be edited by hand.

- Every module in `internal/modules/` MUST have a package-level doc comment written as
  markdown (Go 1.19+ doc comment format). It MUST cover: what the module reads/detects,
  all findings (type, severity, description), configuration keys, and an example config.
- `docs/modules/<name>.md` is generated from the Go doc comment. Regenerate with
  `make gen-docs` and commit the result alongside any module change.
- Any spec that adds or modifies a module MUST include a task to update the package doc
  comment and regenerate the docs. The plan's Constitution Check MUST confirm this.
- Stale status notes (e.g., "not yet implemented") MUST be removed when the feature
  ships.

**Rationale**: Keeping docs in the source file means they are reviewed in the same PR as
the code change. Generated markdown files stay in sync mechanically, not by convention.

## Security Engineering Standards

- All input from external sources (files, stdin, env vars, APIs) MUST be validated
  and sanitized before use.
- The tool MUST NOT execute arbitrary shell commands derived from user-controlled
  input without explicit sandboxing or a documented threat model.
- Dependency supply-chain risk MUST be considered; prefer audited, minimal deps.
- Sensitive data (tokens, keys, PII) MUST NOT be written to stdout, log files, or
  temp files without explicit opt-in from the operator.
- Security-relevant decisions MUST be logged at an auditable verbosity level.

## CLI Interface Standards

- Protocol: stdin / positional args → stdout; errors → stderr.
- Exit codes MUST follow POSIX convention: 0 = success, 1 = failure, 2 = misuse.
- Both human-readable and machine-readable (JSON) output formats MUST be supported
  where output is non-trivial.
- Flags and subcommands MUST follow the principle of least surprise (align with
  common UNIX tooling conventions).

## Governance

This constitution supersedes all other project practices. Any deviation requires
an amendment — no silent exceptions.

**Amendment procedure**:
1. Open a PR describing the proposed change and its rationale.
2. Increment `CONSTITUTION_VERSION` per semantic versioning rules.
3. Update `LAST_AMENDED_DATE` to the merge date.
4. Run the Sync Impact Report checklist above before merging.

**Versioning policy**:
- MAJOR: principle removed, renamed, or fundamentally redefined.
- MINOR: new principle or section added / materially expanded.
- PATCH: wording clarification, typo fix, non-semantic refinement.

**Compliance review**: All PRs MUST reference applicable constitution gates in the
plan's Constitution Check section. Complexity violations require explicit sign-off.

**Version**: 1.1.0 | **Ratified**: 2026-05-07 | **Last Amended**: 2026-05-13
