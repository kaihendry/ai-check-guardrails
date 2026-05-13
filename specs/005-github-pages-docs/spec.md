# Feature Specification: GitHub Pages Module Documentation

**Feature Branch**: `005-github-pages-docs`  
**Created**: 2026-05-13  
**Status**: Draft  
**Input**: User description: "documentation of the internal/modules of the guard rails published to github pages"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Browse Module Reference (Priority: P1)

A developer or security engineer wants to understand what each guardrail module does, what findings it emits, and how to enable or configure it. They visit the GitHub Pages site and can navigate to any module page and read its purpose, finding types, severities, and remediation guidance.

**Why this priority**: The documentation site is the primary entry point for new users evaluating or adopting the tool. Without it, users must read source code to understand module behaviour.

**Independent Test**: Navigate to the published GitHub Pages URL, click any module link (e.g., "bypass"), and confirm the page describes the module's purpose, its finding types, severities, and remediation text.

**Acceptance Scenarios**:

1. **Given** the GitHub Pages site is live, **When** a user visits the module index page, **Then** all 13 modules (banner, bypass, envkeys, evals, gamification, hooks, humanloop, mcp, network, permissions, sandbox, settings, tokens) appear as navigable entries.
2. **Given** the user is on a module page, **When** they read the findings table, **Then** each finding type lists its severity level (INFO / WARN / HIGH / CRITICAL) and a plain-language remediation step.

---

### User Story 2 - Understand the Module Interface (Priority: P2)

A contributor wants to understand the `Module` interface (Name, Enabled, Run) and the `Finding` structure (Type, Severity, Module, Resource, Description, Remediation, Confidence) in order to write a new module. The documentation explains these contracts without requiring source code access.

**Why this priority**: Contributor experience is important for the project's growth, but it is secondary to end-user documentation.

**Independent Test**: From the documentation site, locate the "Module Interface" section and verify it defines the `Name()`, `Enabled()`, and `Run()` methods and the `Finding` struct fields.

**Acceptance Scenarios**:

1. **Given** the documentation site is live, **When** the user reads the Module Interface page, **Then** they see the contract for each required method and the full `Finding` struct with all fields described.
2. **Given** a contributor is reading the "How to add a module" guide, **When** they follow the documented steps, **Then** the guide references real field names and severity constants used in the codebase.

---

### User Story 3 - Find Config Options per Module (Priority: P3)

An operator wants to know which config keys toggle each module on or off and what additional configuration a module accepts (e.g., `env_key_watch_list`, `token_baseline`, `allowlist.mcps`). The documentation links config keys to the modules that consume them.

**Why this priority**: Config reference is useful but users can survive with the YAML config file comments in the short term.

**Independent Test**: On the envkeys module page, confirm the documentation lists the `env_key_watch_list` config key and explains its default value (`ANTHROPIC_API_KEY`).

**Acceptance Scenarios**:

1. **Given** the user is on a module's documentation page, **When** they look at the Configuration section, **Then** they see which config key enables the module and any additional config keys it reads.
2. **Given** the module has a default value for a config option, **When** the user reads the config table, **Then** the default value is shown alongside the key name.

---

### Edge Cases

- What happens when the GitHub Actions workflow runs on a pull request — does it publish a preview or only publish on merge to main?
- How does the documentation stay in sync if a new module is added without updating the docs source?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The documentation site MUST have a landing page listing all modules with brief descriptions.
- **FR-002**: Each module MUST have a dedicated page covering: purpose, finding types, severities, remediation guidance, and configuration keys.
- **FR-003**: The site MUST document the `Module` interface contract (Name, Enabled, Run) and the `Finding` struct fields.
- **FR-004**: The site MUST be published automatically to GitHub Pages on every merge to the `main` branch via a CI workflow.
- **FR-005**: The documentation MUST include a severity reference explaining INFO, WARN, HIGH, and CRITICAL levels and their meaning.
- **FR-006**: Each module page MUST cross-reference the config key that enables or disables it.
- **FR-007**: The site MUST be publicly accessible without authentication at a stable URL.

### Key Entities

- **Module**: A guardrail check; has a name, an enabled flag driven by config, and produces zero or more Findings.
- **Finding**: An individual audit result; has type, severity, module name, resource, description, remediation text, and optional confidence score.
- **Config**: The YAML configuration consumed by the tool; controls which modules run and supplies per-module parameters.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All 13 modules have their own documentation page accessible from the site index within one navigation click.
- **SC-002**: A new contributor can identify the correct config key and severity constant for any module without opening source code, confirmed by a task-completion test.
- **SC-003**: The site is automatically rebuilt and published within 5 minutes of a merge to `main`, with no manual steps required.
- **SC-004**: Every Finding type emitted by each module is listed in its documentation page, with zero omissions detectable by comparing documentation against the codebase.

## Assumptions

- Documentation content will be authored from the existing source code (module names, finding types, severity constants, remediation strings) rather than from external requirements.
- The project already has a GitHub repository (`github.com/kaihendry/ai-check-guardrails`) with Actions enabled, so GitHub Pages publishing is available.
- Documentation will cover the 13 modules present at the time of writing; a process for keeping docs current with future modules is out of scope for this feature.
- Mobile-friendly layout is desirable but not a hard requirement for this initial publish.
- No authentication or access control is needed; the documentation is entirely public.
