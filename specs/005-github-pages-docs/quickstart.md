# Quickstart: GitHub Pages Module Documentation

## Prerequisites

- [uv](https://docs.astral.sh/uv/) installed (provides `uvx`)
- Git repository cloned locally

## Local Preview

```bash
make preview-docs
```

Opens a local server at `http://127.0.0.1:8000`. Changes to files in `docs/` are reflected live.

## Build (without serving)

```bash
make build-docs
```

Output is written to `./site/`. This directory is gitignored.

## Adding or Updating Module Documentation

1. Edit `docs/modules/<module-name>.md`
2. Run `make preview-docs` to verify layout
3. Commit and push to `main` — the docs workflow publishes automatically

## CI Publishing

The `.github/workflows/docs.yml` workflow runs on every push to `main`:
1. Builds the site with mkdocs-material
2. Deploys the built output to GitHub Pages

The live site is available at:  
`https://kaihendry.github.io/ai-check-guardrails/`
