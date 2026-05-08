# Feature Specification: MVP Release Binary & Distribution

**Feature Branch**: `002-mvp-release-binary`  
**Created**: 2026-05-08  
**Status**: Draft  
**Input**: User description: "I want an MVP ai-check-guardrails release binary, a README.md to show how to install (via curl), and an ability to self update on start to keep it low friction, on startup log the version, every new commit to main should trigger a pipeline for a new release"

## Clarifications

### Session 2026-05-08

- Q: Should the self-update mechanism verify binary checksum before replacing the current binary? → A: No — rely on HTTPS transport integrity only; checksum verification applies to `install.sh` only, not to self-update.
- Q: Should the version line be printed to stdout or stderr? → A: stderr — keeps stdout reserved for audit/JSON output and prevents version line from corrupting piped output.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Install Tool via Single Command (Priority: P1)

A developer discovers `ai-check-guardrails` and wants to try it immediately. They run a single `curl` command copied from the README and within seconds have the tool installed and ready to use — no package manager, no build tools, no configuration needed.

**Why this priority**: Installation is the critical first step — nothing else matters if users cannot easily get the tool. A frictionless install experience drives adoption and reduces abandonment.

**Independent Test**: Can be fully tested by running the curl install command on a clean machine and verifying the tool runs successfully.

**Acceptance Scenarios**:

1. **Given** a machine with `curl` and `sh` available, **When** the user runs the install command from the README, **Then** the tool is installed and executable within 30 seconds
2. **Given** the tool is already installed, **When** the user runs the install command again, **Then** the existing installation is updated to the latest version
3. **Given** an unsupported platform, **When** the user runs the install command, **Then** a clear error message explains what is not supported

---

### User Story 2 - Version Visible on Every Run (Priority: P2)

Every time a user runs `ai-check-guardrails`, the current version is displayed at startup. This allows users to know immediately which version they are running, making debugging and support requests straightforward.

**Why this priority**: Version visibility is a hygiene requirement for any distributed tool — it enables users to report issues accurately and know when they are out of date.

**Independent Test**: Run the tool and observe stderr — version string must appear on stderr before any other output; stdout must contain no version line.

**Acceptance Scenarios**:

1. **Given** the tool is installed, **When** the user runs the tool with any arguments, **Then** the version is printed to stderr as the first line; stdout contains only audit output
2. **Given** the tool is installed, **When** the user runs the tool with no arguments, **Then** the version is still printed to stderr

---

### User Story 3 - Self-Update on Startup (Priority: P3)

When a user runs the tool, it silently checks if a newer version is available and applies the update automatically before executing. The user always gets the latest improvements without manually re-running the install command.

**Why this priority**: Self-update keeps the user base current without friction, ensures security fixes propagate quickly, and reduces the support burden of users running outdated versions.

**Independent Test**: Install an older version, run the tool, and verify it updates to the latest version and continues execution normally.

**Acceptance Scenarios**:

1. **Given** an outdated version is installed and internet is available, **When** the user runs the tool, **Then** it downloads and applies the latest version and continues running
2. **Given** the latest version is already installed, **When** the user runs the tool, **Then** no update occurs and startup completes normally
3. **Given** internet is unavailable during startup, **When** the update check fails, **Then** the tool logs a non-fatal warning and continues running with the current version
4. **Given** an update is available, **When** the update is applied, **Then** the new version is logged before execution continues

---

### User Story 4 - Automated Release on Every Main Commit (Priority: P4)

Every time a developer merges or commits to the main branch, an automated pipeline builds and publishes a new versioned release. Contributors do not need to manually manage releases — pushing to main is sufficient.

**Why this priority**: Automated releases eliminate the maintenance burden of manual release management, ensure users always have access to the latest changes, and make the release process repeatable and auditable.

**Independent Test**: Commit to main, then verify a new release with the correct version is published within 15 minutes.

**Acceptance Scenarios**:

1. **Given** a commit is pushed to main, **When** the pipeline runs, **Then** binaries for all supported platforms are built and published as a versioned release
2. **Given** a pipeline build fails, **When** the failure occurs, **Then** no release is published and the failure is reported to the team
3. **Given** a new release is published, **When** a user with an older version runs the tool, **Then** the self-update mechanism detects and applies the new release

---

### Edge Cases

- What happens when the self-update download is interrupted mid-transfer?
- How does the tool behave when the user lacks write permissions to the install location?
- What occurs if the release server is unreachable during an update check?
- A corrupted download during self-update is mitigated by HTTPS transport integrity only; no additional checksum verification is performed during self-update (checksum verification is limited to the initial `install.sh` flow).
- What if two concurrent instances of the tool try to self-update simultaneously?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Users MUST be able to install the tool on a supported platform using a single `curl` command that requires no additional dependencies beyond a standard Unix shell
- **FR-002**: The tool MUST write its version string to stderr as the first line of output on every invocation; stdout is reserved for audit results
- **FR-003**: The tool MUST check for a newer release on every startup and, if found, download and apply the update before proceeding with normal execution
- **FR-004**: The self-update process MUST be non-blocking on failure — if the update check or download fails for any reason, the tool MUST continue running with the currently installed version
- **FR-005**: Every commit merged or pushed to the main branch MUST automatically trigger a release pipeline that builds distribution artifacts for all supported platforms
- **FR-006**: The release pipeline MUST publish versioned release artifacts to a publicly accessible location that the install script and self-update mechanism can reference
- **FR-007**: The README MUST contain a working install command, supported platforms, and basic usage instructions
- **FR-008**: The tool MUST log the version it is updating to when a self-update is applied

### Key Entities

- **Release**: A versioned set of platform-specific binaries published by the automated pipeline; identified by a version number derived from the commit or tag
- **Install Script**: A shell script fetched via `curl` that detects the user's platform, downloads the correct binary, and places it in the user's PATH
- **Supported Platform**: A combination of operating system and CPU architecture for which a binary is built and distributed (e.g., Linux x86-64, Linux ARM64, macOS x86-64, macOS ARM64)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new user can install and execute the tool on a supported platform in under 2 minutes using only the README install command
- **SC-002**: The tool version is displayed within 1 second of invocation on any supported platform
- **SC-003**: When an update is available and internet is accessible, the self-update completes and execution continues within 30 seconds of invoking the tool
- **SC-004**: A new release is publicly available within 15 minutes of a commit landing on the main branch
- **SC-005**: The self-update failure path (no internet, server down) does not interrupt normal tool operation — 100% of the time the tool continues running despite update failures
- **SC-006**: The install command works without modification on all supported platforms

## Assumptions

- Target platforms for the initial release are Linux (x86-64 and ARM64) and macOS (x86-64 and ARM64); Windows is out of scope for this MVP
- GitHub Releases is used as the distribution mechanism for release artifacts, and the install script fetches from there
- Version numbering follows a sequential or commit-based scheme automatically derived from the pipeline (no manual version bumping required)
- Users have outbound internet access to reach the release server for both initial install and self-update
- The tool is installed to a user-writable location (e.g., `~/.local/bin` or `/usr/local/bin`) so self-update does not require elevated privileges in the common case
- The automated release pipeline uses the project's existing version control hosting (GitHub) and its associated CI/CD capabilities
- Self-update is opt-out rather than opt-in — it is enabled by default to minimize friction, but users can disable it via an environment variable or flag
