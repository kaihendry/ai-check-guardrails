# Data Model: GitHub Pages Module Documentation

**Branch**: `005-github-pages-docs` | **Date**: 2026-05-13

## Entities

### Module (documentation entity)

Each module maps 1-to-1 to a markdown page under `docs/modules/`.

| Field            | Type     | Source                                  | Notes |
|------------------|----------|-----------------------------------------|-------|
| name             | string   | `module.Name()`                         | Slug used for page filename |
| enabled_default  | bool     | `config.defaults()` ModuleToggles      | Whether on by default |
| config_key       | string   | `ModuleToggles` JSON tag                | e.g. `modules.env_keys` |
| purpose          | prose    | Authored                                | One paragraph |
| finding_types    | []FindingType | Module source                      | All types emitted by `Run()` |
| extra_config     | []ConfigKey   | Module source                      | Keys beyond the enable toggle |

**Modules inventory**:

| Module        | Config Key           | Enabled by Default |
|---------------|----------------------|--------------------|
| banner        | `modules.banner`     | true               |
| bypass        | `modules.bypass`     | true               |
| envkeys       | `modules.env_keys`   | true               |
| evals         | `modules.evals`      | false              |
| gamification  | `modules.gamification` | true             |
| hooks         | `modules.hooks`      | true               |
| humanloop     | `modules.humanloop`  | false              |
| mcp           | `modules.mcp`        | true               |
| network       | `modules.network`    | false              |
| permissions   | `modules.permissions`| true               |
| sandbox       | `modules.sandbox`    | true               |
| settings      | `modules.settings`   | true               |
| tokens        | `modules.tokens`     | false              |

---

### Finding (reference entity)

Described on the Module Interface reference page.

| Field        | Type      | Required | Notes |
|--------------|-----------|----------|-------|
| type         | string    | yes      | Module-specific string constant |
| severity     | Severity  | yes      | One of: INFO, WARN, HIGH, CRITICAL |
| module       | string    | yes      | Module name (matches `Name()`) |
| resource     | string    | yes      | The config key or file path inspected |
| description  | string    | yes      | Human-readable finding summary |
| remediation  | string    | yes      | Actionable fix guidance |
| confidence   | float64   | no       | 0–1, present only for heuristic findings |

---

### Severity (reference enum)

Described on the Severity reference page.

| Value    | Meaning                          | Score Impact |
|----------|----------------------------------|--------------|
| INFO     | Informational; no action needed  | 0            |
| WARN     | Advisory; improvement recommended | Moderate    |
| HIGH     | Security gap requiring attention  | Significant |
| CRITICAL | Must-fix before enforcement mode | Maximum     |

---

### Config (reference entity)

Described in `docs/reference/config.md`.

| Key                  | Type          | Default          | Description |
|----------------------|---------------|------------------|-------------|
| mode                 | string        | `monitor`        | `monitor` or `enforce` |
| siem_endpoint        | string        | (empty)          | HTTPS endpoint for findings export |
| scan_root            | string        | `$HOME`          | Root path for file scans |
| banner_url           | string        | (empty)          | Custom banner URL |
| score_threshold      | int           | `70`             | Minimum score to pass |
| modules.*            | bool          | see table above  | Per-module enable toggle |
| allowlist.mcps       | []string      | (empty)          | Approved MCP identifiers |
| allowlist.skills     | []string      | (empty)          | Approved skill names |
| allowlist.domains    | []string      | (empty)          | Approved network domains |
| allowlist.precommit_hooks | []string | `["gitleaks"]`   | Required pre-commit hooks |
| token_baseline       | object        | null             | Required if tokens module enabled |
| token_baseline.daily_mean | int    | —                | Expected daily token mean |
| token_baseline.std_dev    | float64 | —               | Standard deviation |
| token_baseline.multiplier | float64 | —               | Anomaly threshold multiplier |
| env_key_watch_list   | []string      | `["ANTHROPIC_API_KEY"]` | Keys to detect in env |

---

## Documentation Page Relationships

```
docs/
├── index.md                    # Landing: project intro + quick link to modules
├── modules/
│   ├── index.md                # Module list table (all 13 modules, one-line purpose, default state)
│   ├── banner.md               # Per-module pages (×13)
│   ├── bypass.md
│   ├── envkeys.md
│   ├── evals.md
│   ├── gamification.md
│   ├── hooks.md
│   ├── humanloop.md
│   ├── mcp.md
│   ├── network.md
│   ├── permissions.md
│   ├── sandbox.md
│   ├── settings.md
│   └── tokens.md
└── reference/
    ├── module-interface.md     # Module interface contract + Finding struct
    ├── config.md               # Full config key reference
    └── severity.md             # Severity levels explained
```
