# Feature Specification: Fix Self-Update Version Check

**Feature Branch**: `007-fix-update-check`  
**Created**: 2026-05-13  
**Status**: Draft  
**Input**: User description: "only download if there is a new version... 002-mvp-release-binary has a bug"

## Background

The `ai-check-guardrails` tool performs a self-update check on every startup. The bug: the tool always downloads and replaces itself, even when the currently installed version is already the latest. The observed output demonstrates this:

```
ai-check-guardrails v0.0.20260513135826+5118f3a
[update] new version available: latest
[update] downloading...
[update] updated to latest, continuing...
v0.0.20260513135826+5118f3a
```

The version before and after the update is identical, confirming the tool re-downloads itself unnecessarily on every run. This is caused by the version embedded in the binary (`v0.0.20260513135826+5118f3a`) never matching the release tag returned by the release provider (`latest`), so the comparison always treats the current version as outdated.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - No Unnecessary Downloads When Already Up to Date (Priority: P1)

A user who already has the latest version runs the tool. The tool starts without downloading anything — no network traffic for a binary update, no delay, no update messages on stderr.

**Why this priority**: This is the core bug. Every user is affected on every run. Unnecessary downloads waste bandwidth, slow startup, and erode trust. Fixing this is the minimum viable outcome for the feature.

**Independent Test**: Install the latest version, run the tool, and confirm no `[update]` lines appear on stderr and no network request is made to download a binary.

**Acceptance Scenarios**:

1. **Given** the latest version is installed, **When** the user runs the tool, **Then** no update messages appear and the binary is not replaced
2. **Given** the latest version is installed, **When** the user runs the tool multiple times in succession, **Then** no update occurs on any run
3. **Given** the tool starts up, **When** the version check is performed, **Then** the comparison correctly identifies "already up to date" and skips the download

---

### User Story 2 - Update Still Occurs When a Newer Version Exists (Priority: P2)

A user with an older version runs the tool. The tool detects that a newer release is available, downloads it, and continues — exactly as the original feature intended.

**Why this priority**: The fix must not break the valid update path. The self-update feature from spec 002 is a core user-facing capability and must remain functional.

**Independent Test**: Install a known older version, run the tool, and verify it updates to the latest version and reports the new version on stderr.

**Acceptance Scenarios**:

1. **Given** an outdated version is installed, **When** the user runs the tool, **Then** the update is downloaded, applied, and the new version is logged before execution continues
2. **Given** the update completes, **When** the user runs the tool again immediately, **Then** no further update occurs

---

### User Story 3 - Version Identity Is Stable and Comparable (Priority: P3)

The version string embedded in the binary and the version identifier used by the release system refer to the same release, enabling accurate comparisons.

**Why this priority**: The root cause of the bug is a mismatch between what the binary reports as its version and what the release system returns as the "latest" tag. Fixing the comparison requires both sides to use the same identifier.

**Independent Test**: Run `ai-check-guardrails --version`, note the version string, then query the release system for the latest release identifier — both values must be equal for the same release.

**Acceptance Scenarios**:

1. **Given** a binary is built from a specific release, **When** the version is checked, **Then** the reported version exactly matches the release identifier used by the update check
2. **Given** a new release is published, **When** a user installs it and runs the tool, **Then** the version string matches the release tag and no spurious update is triggered

---

### Edge Cases

- What happens when the release system returns an unexpected or malformed version identifier?
- What happens when the network is unavailable and the version check cannot complete?
- What happens when the binary is run from a location where it cannot be replaced (e.g., read-only filesystem)?
- What happens if the version embedded in the binary is empty or a development build value?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The tool MUST skip the binary download when the currently installed version matches the latest available release
- **FR-002**: The tool MUST still download and apply an update when the currently installed version does not match the latest available release
- **FR-003**: The version identifier embedded in a released binary MUST use the same format and value as the release tag used by the update check
- **FR-004**: The tool MUST NOT print `[update]` messages when no update is required
- **FR-005**: When no update is needed, the tool MUST proceed to normal execution without any additional delay from update logic
- **FR-006**: A development or locally built binary (without an official release tag) MUST either skip the update check or handle the version mismatch gracefully without downloading

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of runs with the latest version installed produce zero `[update]` lines on stderr
- **SC-002**: 100% of runs with an outdated version installed trigger a successful update
- **SC-003**: The version reported by `--version` matches the release tag for any officially released binary
- **SC-004**: Startup time for users on the latest version is not increased by the update check (check completes without download in under 3 seconds on a standard connection)

## Assumptions

- The release system provides a stable, unique identifier per release that can be used for exact string comparison
- The build pipeline that produces release binaries has the ability to embed the release tag into the binary at build time
- Development builds (built locally without going through the release pipeline) are expected to have a different version format; the update check for such builds may skip or behave differently without this being considered a bug
- The install script (`install.sh`) is not in scope for this fix — it installs the correct binary and this fix is only about the runtime self-update comparison
