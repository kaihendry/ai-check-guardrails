# Feature Specification: AI Guardrails Audit

**Feature Branch**: `001-ai-guardrails-audit`
**Created**: 2026-05-07
**Status**: Draft
**Input**: User description: "A SecEng tool that would run periodically via launchd or a startup hook..."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Security Team Gains Visibility into Claude Usage (Priority: P1)

A security engineer runs the audit tool on a developer workstation (or it runs
automatically at login/on a schedule) and receives a structured report — sent to
the SIEM and optionally printed to the terminal — showing the current posture of
Claude AI configuration, active MCPs, permission state, and detected anomalies.

**Why this priority**: Without baseline visibility, no enforcement is possible.
This story delivers the foundational detection capability that all other stories
depend on.

**Independent Test**: Run the tool with all feature-toggles in monitor mode.
Confirm a structured JSON log event is emitted covering settings, MCPs,
permissions, and token-usage status. No enforcement actions are taken.

**Acceptance Scenarios**:

1. **Given** the tool is invoked on a workstation, **When** it completes a scan,
   **Then** it emits a structured log event containing findings from every enabled
   detection module, with no false-blocking of normal developer workflows.
2. **Given** a misconfigured settings file (e.g., missing required keys),
   **When** the tool scans the settings, **Then** the finding is recorded in the
   log with severity `WARN` or `CRITICAL` and the developer is notified via
   the terminal banner.
3. **Given** the tool is scheduled via `launchd` or a startup hook,
   **When** the schedule fires, **Then** the tool runs unattended, logs results to
   the SIEM endpoint, and exits cleanly without user interaction.

---

### User Story 2 - Detection of Risky MCP / Skill Usage (Priority: P2)

The security team maintains an approved list of MCPs and skills. The tool detects
any MCP or skill not on that list and flags it as a potential exfiltration risk in
the audit log.

**Why this priority**: Unapproved MCPs are a primary exfiltration vector. Early
detection — before enforcement — allows the team to baseline and update the
allowlist without blocking productivity.

**Independent Test**: Add a non-allowlisted MCP to the test configuration. Run
the tool. Verify a `POLICY_VIOLATION` finding is emitted for the unapproved MCP
and no approved MCPs are flagged.

**Acceptance Scenarios**:

1. **Given** a developer has an MCP not on the approved list, **When** the audit
   runs, **Then** the tool emits a finding with `type: UNAPPROVED_MCP` and the
   MCP identifier.
2. **Given** all active MCPs are on the approved list, **When** the audit runs,
   **Then** no MCP-related findings are emitted.
3. **Given** the approved list is empty (not configured), **When** the audit runs,
   **Then** all MCPs are flagged with severity `WARN` and a configuration-needed
   notice is included.

---

### User Story 3 - Policy-Bypass Detection (Priority: P3)

The tool detects indicators of policy bypass attempts — such as the use of
`--dangerously-skip-permissions` flags, absent pre-commit hooks (e.g., gitleaks),
or missing human-in-the-loop confirmations — and logs them as high-severity
findings.

**Why this priority**: Bypass detection is a high-value, low-false-positive check
that catches deliberate evasion even before enforcement is enabled.

**Independent Test**: Run the tool on a workstation where `gitleaks` pre-commit
hook is absent. Verify a `MISSING_PRECOMMIT_HOOK` finding is emitted at severity
`HIGH`. Restore the hook and confirm no finding on the next run.

**Acceptance Scenarios**:

1. **Given** a project repo lacks the required pre-commit hooks, **When** the
   audit runs, **Then** a `MISSING_PRECOMMIT_HOOK` finding is logged.
2. **Given** the tool detects recent use of `--dangerously-skip-permissions`,
   **Then** a `POLICY_BYPASS` finding is logged at severity `CRITICAL`.
3. **Given** no bypass indicators are present, **When** the audit runs,
   **Then** no bypass findings are emitted.

---

### User Story 4 - Gamification: Per-User Security Score (Priority: P4)

Each audit run calculates a security posture score (0–100) for the developer based
on weighted findings. The score and trend are included in the structured log and
optionally displayed as a terminal banner, encouraging self-service improvement.

**Why this priority**: A score creates a feedback loop that motivates compliance
without requiring enforcement, easing adoption during the phased rollout.

**Independent Test**: Run the tool on a known-clean configuration and verify a
score of 100 is returned. Introduce two policy violations and verify the score
decreases proportionally and is present in the log event.

**Acceptance Scenarios**:

1. **Given** a clean workstation configuration, **When** the audit runs, **Then**
   a score of 100 is returned and logged.
2. **Given** findings of varying severity exist, **When** the score is calculated,
   **Then** higher-severity findings cause larger score decreases than lower ones.
3. **Given** a score below a configurable threshold (default: 70), **When** the
   audit completes, **Then** the terminal banner includes a link to the security
   training policy page.

---

### Edge Cases

- What happens when the SIEM endpoint is unreachable? (Log locally, retry on
  next scheduled run; do not block the developer's workflow.)
- What if the settings file does not exist at all? (Emit a `MISSING_CONFIG`
  finding at severity `CRITICAL`; do not crash.)
- What if a feature toggle references a module that is not yet implemented?
  (Skip silently and log a `MODULE_UNAVAILABLE` debug message.)
- What if the tool is run by a non-developer service account with no MCP
  configuration? (All MCP checks return no findings; score unaffected.)
- What if two simultaneous scheduled runs overlap? (Second run detects a lock
  and exits cleanly with an `ALREADY_RUNNING` status code.)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The tool MUST run as a standalone binary invocable from the command
  line with no required arguments.
- **FR-002**: The tool MUST support scheduling via `launchd` (macOS) and via a
  startup hook (Linux/Windows) for unattended periodic execution.
- **FR-003**: The tool MUST emit all findings as structured (JSON) log events to
  a configurable SIEM endpoint and to local stdout/stderr.
- **FR-004**: The tool MUST include feature toggles for each detection module,
  allowing individual modules to be enabled or disabled without rebuilding.
- **FR-005**: The tool MUST operate in two modes: `monitor` (detect and log only)
  and `enforce` (detect, log, and block or alert on violations), switchable via
  configuration or CLI flag.
- **FR-006**: The tool MUST verify the Claude `settings.json` file for presence
  of required keys and detect any unauthorized overrides.
- **FR-007**: The tool MUST compare active MCPs and skills against a configurable
  approved allowlist and flag any not on the list.
- **FR-008**: The tool MUST validate that configured access restrictions are
  actively enforced by attempting controlled read/access tests and verifying
  denial.
- **FR-009**: The tool MUST detect usage of `--dangerously-skip-permissions` and
  similar policy-bypass flags in recent Claude invocation history.
- **FR-010**: The tool MUST check that required pre-commit hooks (e.g., gitleaks)
  are present and active in all local Git repositories under a configurable scan
  root.
- **FR-011**: The tool MUST calculate a numeric security posture score (0–100)
  per run and include it in the structured log event.
- **FR-012**: The tool MUST display a configurable terminal banner that can link
  to training resources and updated security policies.
- **FR-013**: The tool MUST detect irregular or anomalous token usage patterns
  relative to a configurable baseline and flag deviations as findings.
- **FR-014**: The tool MUST log outbound network requests made by Claude agents
  (webfetch/curl equivalents) and record destination domains and categories.
- **FR-015**: The tool MUST support Anthropic-recommended evaluation hooks for
  adversarial prompt review and private data protection checks.
- **FR-016**: The tool MUST prevent concurrent executions via a run-lock
  mechanism.
- **FR-017**: The tool MUST exit with distinct exit codes: `0` (clean), `1`
  (findings present), `2` (tool error/misconfiguration).

### Key Entities

- **Audit Run**: A single execution of the tool; contains a timestamp, host
  identity, list of findings, posture score, and mode.
- **Finding**: An individual detected issue; has a type, severity
  (`INFO/WARN/HIGH/CRITICAL`), affected resource, description, and remediation
  hint.
- **Detection Module**: A toggleable unit of audit logic (e.g., Settings Check,
  MCP Monitor, Permission Validator). Enabled/disabled via feature toggles.
- **Allowlist**: Operator-maintained list of approved MCPs, skills, and network
  destinations.
- **Posture Score**: A derived 0–100 integer calculated from weighted findings in
  an audit run.
- **SIEM Event**: The structured log record forwarded to the security information
  and event management system.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of enabled detection modules produce findings within 30 seconds
  of tool invocation on a standard developer workstation.
- **SC-002**: The tool runs unattended (scheduled) on 100% of enrolled workstations
  without requiring manual intervention for at least 30 consecutive days.
- **SC-003**: Security team achieves full visibility into active MCPs across all
  enrolled workstations within one week of rollout, with zero missed detections of
  unapproved MCPs.
- **SC-004**: Policy-bypass findings (e.g., `--dangerously-skip-permissions`) are
  detected and logged within one scheduled run interval of the bypass occurring.
- **SC-005**: Developer posture scores improve by an average of 15 points across
  the enrolled population within 60 days of gamification rollout.
- **SC-006**: Zero developer workflows are blocked during the `monitor` phase;
  all tool-related interruptions are attributable to `enforce` mode being active.
- **SC-007**: SIEM event delivery success rate of ≥99% when the SIEM endpoint is
  reachable; failed deliveries are retried and locally buffered without data loss.

## Assumptions

- The tool is deployed on developer workstations by the security team; end users
  do not install it themselves.
- `launchd` is available on macOS targets; equivalent startup mechanisms exist on
  other supported platforms.
- The SIEM endpoint accepts structured JSON over HTTPS; authentication credentials
  are supplied via environment variable or a secrets manager, not hardcoded config.
- Claude's `settings.json` location follows the standard path conventions for
  Claude Code; the tool will look in the default locations and allow override via
  an environment variable.
- The approved MCP/skill allowlist is maintained by the security team and
  distributed as a versioned config file; the tool reads it at each run.
- Token usage baselines will be established from historical data during the initial
  monitor-only phase before anomaly thresholds are configured.
- The gamification score is informational only in `monitor` mode; it does not gate
  any developer action.
- Network request monitoring is passive (log-and-report); active blocking of
  network requests is an `enforce`-mode capability for a future iteration.
- Pre-commit hook checks scan only local Git repositories under the user's home
  directory by default; the scan root is configurable.
