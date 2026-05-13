# Contract: mkdocs.yml Configuration Schema

**Type**: Static site configuration  
**Location**: `mkdocs.yml` (repo root)

## Required Structure

```yaml
site_name: ai-check-guardrails
site_url: https://kaihendry.github.io/ai-check-guardrails/
repo_url: https://github.com/kaihendry/ai-check-guardrails
docs_dir: docs
theme:
  name: material

nav:
  - Home: index.md
  - Modules:
      - Overview: modules/index.md
      - banner: modules/banner.md
      - bypass: modules/bypass.md
      - envkeys: modules/envkeys.md
      - evals: modules/evals.md
      - gamification: modules/gamification.md
      - hooks: modules/hooks.md
      - humanloop: modules/humanloop.md
      - mcp: modules/mcp.md
      - network: modules/network.md
      - permissions: modules/permissions.md
      - sandbox: modules/sandbox.md
      - settings: modules/settings.md
      - tokens: modules/tokens.md
  - Reference:
      - Module Interface: reference/module-interface.md
      - Configuration: reference/config.md
      - Severity Levels: reference/severity.md
```

## Constraints

- `site_url` MUST match the GitHub Pages URL for correct canonical links
- `docs_dir` MUST be `docs` (standard mkdocs convention)
- All `nav` entries MUST correspond to actual files in `docs/`
- Theme MUST be `material` (provides search, responsive layout, navigation tabs)

## Makefile Targets

```makefile
preview-docs:
	uvx --with mkdocs-material mkdocs serve

build-docs:
	uvx --with mkdocs-material mkdocs build
```

## GitHub Actions Contract

**Trigger**: push to `main` branch  
**Workflow file**: `.github/workflows/docs.yml`  
**Output**: Built site deployed to `gh-pages` branch  
**GitHub Pages source**: `gh-pages` branch, root `/`

### Workflow steps (in order):
1. `actions/checkout@v4`
2. `actions/setup-python@v5` with Python 3.12
3. `pip install mkdocs-material`
4. `mkdocs build`
5. `actions/upload-pages-artifact@v3` from `./site`
6. `actions/deploy-pages@v4`

### Required GitHub repository settings:
- Pages source: GitHub Actions (not branch)
- `permissions.pages: write` and `permissions.id-token: write` on the deploy job
