# Research: GitHub Pages Module Documentation

**Branch**: `005-github-pages-docs` | **Date**: 2026-05-13

## Decision 1: Static Site Generator

**Decision**: mkdocs with Material theme

**Rationale**:
- Pure markdown content with a single YAML config file — no templating language to learn
- `mkdocs serve` gives hot-reload local preview with one command
- Material theme provides search, navigation, and responsive layout out of the box
- Deployable to GitHub Pages via `mkdocs gh-deploy` or the standard build+push workflow
- Installable without polluting the system via `uvx` (uv tool runner already in global config)
- Widely used for Go CLI documentation (e.g., goreleaser, golangci-lint)

**Alternatives considered**:
- Hugo: Go-native binary, faster builds, but requires learning Hugo templates and frontmatter; more setup overhead for a simple reference site
- Jekyll: GitHub Pages default, but requires Ruby and is slower to set up
- Docusaurus: Excellent UX but pulls in Node.js / npm — disproportionate for this scope
- godoc/pkgsite: Documents exported symbols only; `internal/` packages are excluded from public godoc; also produces API reference not user-facing prose

## Decision 2: Local Preview Command

**Decision**: `make preview-docs` runs `uvx --with mkdocs-material mkdocs serve`

**Rationale**:
- `uvx` creates an isolated ephemeral environment without a global `pip install`
- The global CLAUDE.md mandates using `uv` for Python scripts; `uvx` is the zero-install variant
- No `requirements.txt` or `pyproject.toml` needed — dependencies declared inline in the Makefile target
- A developer running `make preview-docs` for the first time gets a working server with no extra setup

## Decision 3: GitHub Pages Deployment

**Decision**: New GitHub Actions workflow (`.github/workflows/docs.yml`) that builds with mkdocs and deploys to the `gh-pages` branch on every push to `main`

**Rationale**:
- The existing `release.yml` handles binary releases; docs warrant a separate, independently triggerable workflow
- `actions/deploy-pages` (official GitHub action) is now stable and handles branch/permissions automatically
- Alternatively: `peaceiris/actions-gh-pages` is battle-tested but adds a third-party dep; `actions/deploy-pages` is first-party and preferred per the Simplicity principle

## Decision 4: Content Authoring

**Decision**: Hand-written markdown files in `docs/` directory derived from source code at authoring time

**Rationale**:
- Documentation goal is operator-facing prose (purpose, configuration, remediation guidance), not raw API reference
- `internal/` packages are invisible to godoc by design; a generated doc approach cannot cover them without custom tooling
- Modules are small and stable; hand-authoring is low-effort and produces higher-quality prose than auto-generation
- A future iteration could layer in auto-generation if modules proliferate

## Decision 5: Documentation Root Layout

**Decision**: `docs/` directory at repo root, `mkdocs.yml` at repo root, `Makefile` at repo root

**Rationale**:
- Standard mkdocs convention (`mkdocs.yml` at root, content under `docs/`)
- `Makefile` at root follows Go project convention and allows `make preview-docs` from checkout root
- No extra nesting that would complicate GitHub Actions paths

## Resolved Unknowns

| Unknown | Resolution |
|---------|-----------|
| Auto-generate vs. hand-write | Hand-write prose from source inspection |
| Tool to serve locally | `uvx mkdocs-material` via Makefile |
| CI deploy target | `gh-pages` branch via `actions/deploy-pages` |
| Module count | 13 modules (banner, bypass, envkeys, evals, gamification, hooks, humanloop, mcp, network, permissions, sandbox, settings, tokens) |
| Hugo vs mkdocs | mkdocs — lower friction, simpler config |
