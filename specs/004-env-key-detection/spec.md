# Feature Specification: Environment API Key Detection

**Feature Branch**: `004-env-key-detection`
**Created**: 2026-05-11
**Status**: Draft
**Input**: User description: "Add a detection module that checks whether ANTHROPIC_API_KEY or other Anthropic credentials are present as plaintext environment variables, flagging this as a security risk since keys should be managed via a secrets manager rather than exposed in the process environment."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Security Team Detects Plaintext API Keys in Environment (Priority: P1)

A security engineer runs the audit tool on a developer workstation and receives a
finding when `ANTHROPIC_API_KEY` or similar Anthropic credential variables are
present as plaintext environment variables, indicating the key is not being
managed through an approved secrets manager.

**Why this priority**: Plaintext API keys in environment variables are directly
accessible to any process running under the same user, can be captured in shell
history, logged by process monitors, and are a primary credential exfiltration
vector. Early detection prevents key leakage before any enforcement is required.

**Independent Test**: Set `ANTHROPIC_API_KEY=sk-ant-test123` in the environment,
run the audit tool, and verify a `PLAINTEXT_CREDENTIAL` finding is emitted at
severity `HIGH`. Unset the variable and confirm no finding on the next run.

**Acceptance Scenarios**:

1. **Given** `ANTHROPIC_API_KEY` is set in the process environment, **When** the
   audit runs, **Then** a `PLAINTEXT_CREDENTIAL` finding is emitted with severity
   `HIGH` and the variable name (not the value) in the finding details.
2. **Given** no Anthropic credential variables are set, **When** the audit runs,
   **Then** no `PLAINTEXT_CREDENTIAL` findings are emitted.
3. **Given** a finding is emitted, **Then** the finding description includes a
   remediation hint pointing to the approved secrets manager workflow.

---

### User Story 2 - Detection Covers All Anthropic Credential Variables (Priority: P2)

The detection module checks a configurable list of credential variable names, not
just `ANTHROPIC_API_KEY`, so that future credential types (e.g., organisation
tokens, workspace keys) are also caught without a code change.

**Why this priority**: Anthropic may introduce additional credential types.
A hardcoded single-variable check would miss them; a configurable list future-proofs
the module at low cost.

**Independent Test**: Set `ANTHROPIC_ORG_TOKEN=test-token` in the environment.
Run the tool with a config that includes `ANTHROPIC_ORG_TOKEN` in the watched list.
Verify a `PLAINTEXT_CREDENTIAL` finding is emitted for `ANTHROPIC_ORG_TOKEN`.

**Acceptance Scenarios**:

1. **Given** the operator has added a custom variable name to the watched list,
   **When** that variable is present in the environment, **Then** a
   `PLAINTEXT_CREDENTIAL` finding is emitted for it.
2. **Given** no custom list is configured, **When** the audit runs, **Then** the
   module falls back to a sensible built-in default list (at minimum `ANTHROPIC_API_KEY`).
3. **Given** a variable on the watched list is not set, **When** the audit runs,
   **Then** no finding is emitted for that variable.

---

### User Story 3 - Score Impact and Terminal Feedback (Priority: P3)

Each `PLAINTEXT_CREDENTIAL` finding reduces the posture score and the terminal
banner informs the developer how to remediate by migrating the key to the approved
secrets manager.

**Why this priority**: Visibility without feedback leaves developers without a
clear remediation path. A score penalty and actionable banner complete the
feedback loop established in the core audit spec.

**Independent Test**: Run the tool with a clean environment (score 100). Set
`ANTHROPIC_API_KEY` and re-run. Verify the score decreases and the banner
includes remediation guidance.

**Acceptance Scenarios**:

1. **Given** a `PLAINTEXT_CREDENTIAL` finding exists, **When** the score is
   calculated, **Then** it is reduced by at least the weight assigned to `HIGH`
   severity findings.
2. **Given** a `PLAINTEXT_CREDENTIAL` finding exists, **When** the terminal
   banner is displayed, **Then** it includes a remediation hint describing how
   to move the key to the approved secrets manager.
3. **Given** the finding is remediated (variable unset), **When** the next audit
   runs, **Then** the score returns to its previous level and no banner hint is
   shown.

---

### Edge Cases

- What if the variable is set to an empty string? (Do not flag — an empty value
  is not a credential; emit at most an `INFO` finding if desired.)
- What if the variable name matches the watched list but the value does not look
  like a credential (e.g., a placeholder like `"changeme"`)? (Flag it anyway
  with severity `WARN`; the variable name alone indicates misconfiguration.)
- What if the module is disabled via feature toggle? (Skip all checks silently;
  no finding emitted, score unaffected.)
- What if the watched list in config is empty? (Use the built-in default list;
  log an `INFO` message noting the default is in use.)
- What if the audit runs inside a CI environment where `ANTHROPIC_API_KEY` is
  intentionally set? (Flag it; CI environments are a known leakage risk.
  Operators may suppress via `enforce`-mode allowlist, not by disabling detection.)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The module MUST detect any environment variable whose name matches
  the configured watched list and emit a `PLAINTEXT_CREDENTIAL` finding per
  matched variable.
- **FR-002**: The module MUST NOT include the variable's value in any finding,
  log, or SIEM event — only the variable name.
- **FR-003**: The default watched list MUST include at minimum `ANTHROPIC_API_KEY`.
- **FR-004**: The operator MUST be able to extend or replace the default watched
  list via the tool's configuration file without rebuilding.
- **FR-005**: Findings from this module MUST carry severity `HIGH` when the value
  is non-empty and appears to be a credential, and `WARN` when the value is
  present but does not resemble a credential.
- **FR-006**: The module MUST be toggleable via the existing feature-toggle
  mechanism (enabled/disabled without rebuilding).
- **FR-007**: Each `PLAINTEXT_CREDENTIAL` finding MUST include a remediation hint
  field describing how to migrate the credential to a secrets manager.
- **FR-008**: The module MUST integrate with the posture scoring system,
  contributing a score penalty proportional to the severity of each finding.
- **FR-009**: Variables set to empty strings MUST NOT produce `HIGH` or `WARN`
  findings.

### Key Entities

- **Watched Variable**: A named environment variable monitored by this module;
  defined in the default list or operator-supplied config.
- **Plaintext Credential Finding**: A `Finding` with `type: PLAINTEXT_CREDENTIAL`,
  severity `HIGH` or `WARN`, variable name, and a remediation hint.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The module detects a set `ANTHROPIC_API_KEY` within a single audit
  run and emits the finding within the 30-second tool completion SLA.
- **SC-002**: Zero false positives for variables that are absent or empty across
  100 consecutive audit runs on a clean workstation.
- **SC-003**: Operators can extend the watched variable list via config change
  alone, with no code change or rebuild required.
- **SC-004**: Developer posture score decreases measurably (by the `HIGH`-severity
  weight) whenever a plaintext credential is present, and returns to baseline
  when it is removed.
- **SC-005**: The credential value never appears in any log, SIEM event, or
  terminal output produced by the tool.

## Assumptions

- The tool reads environment variables from the process environment at the time
  of invocation; variables set after the tool starts are not detected.
- "Resembles a credential" is determined by a simple heuristic (e.g., minimum
  length threshold and non-trivial entropy); exact heuristic is an implementation
  detail.
- The approved secrets manager is documented separately by the security team;
  this module only detects the risk and links to that documentation.
- The module reuses the `Finding` and scoring types already defined in
  `001-ai-guardrails-audit`; no new data model types are required.
- CI/CD pipelines are not excluded from detection by default; suppression is an
  operator responsibility via existing `enforce`-mode configuration.
