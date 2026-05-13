# Feature Specification: Transcript Token Usage Reading

**Feature Branch**: `006-transcript-token-usage`  
**Created**: 2026-05-13  
**Status**: Draft  
**Input**: User description: "implement token usage reading from ~/.claude/projects transcripts for the tokens module"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Detect Token Anomalies from Real Usage Data (Priority: P1)

A developer runs `ai-check-guardrails` and the tokens module reads actual daily token consumption from the local Claude Code session transcripts. If consumption exceeds the configured baseline threshold, a warning finding is raised — replacing the current placeholder that always returns zero.

**Why this priority**: Without real data, the tokens module cannot detect anomalous Claude usage. This is the core value of the module and its current gap.

**Independent Test**: Configure a token baseline and ensure at least one transcript file exists in `~/.claude/projects`. Run the tool; verify a TOKEN_ANOMALY finding appears when simulated usage exceeds the threshold and no finding appears when usage is normal.

**Acceptance Scenarios**:

1. **Given** transcripts exist in `~/.claude/projects` from the past 24 hours with total token counts above the configured anomaly threshold, **When** the tokens module runs, **Then** a TOKEN_ANOMALY finding with WARN severity is returned.
2. **Given** transcripts exist but total token usage is within the normal range, **When** the tokens module runs, **Then** no anomaly findings are returned.
3. **Given** no transcripts exist or the transcripts directory is absent, **When** the tokens module runs, **Then** the module returns no anomaly findings and does not error.

---

### User Story 2 - Handle Missing or Inaccessible Transcripts Gracefully (Priority: P2)

When the transcript directory does not exist (e.g., Claude Code is not installed, or this is a CI environment), the tokens module degrades gracefully rather than failing.

**Why this priority**: The tool runs in diverse environments; it must never crash due to an absent transcript directory.

**Independent Test**: Run the module with a config pointing to a non-existent transcript path. Verify zero findings are returned and no error is propagated.

**Acceptance Scenarios**:

1. **Given** `~/.claude/projects` does not exist, **When** the tokens module runs, **Then** it returns zero findings with no error.
2. **Given** the transcript files exist but are unreadable (permissions), **When** the tokens module runs, **Then** it skips unreadable files and continues without error.

---

### User Story 3 - Configurable Lookback Window (Priority: P3)

Operators can configure how far back (in hours) the module looks when summing token usage, defaulting to the last 24 hours. This allows tuning whether a daily or rolling window is used for anomaly detection.

**Why this priority**: Different teams audit usage on different cadences; a fixed 24-hour window may miss patterns or double-count sessions.

**Independent Test**: Set a lookback of 1 hour and populate only old transcripts (>1 hour ago). Verify zero tokens are counted and no anomaly is raised regardless of baseline.

**Acceptance Scenarios**:

1. **Given** a lookback window of 24 hours (default), **When** transcripts from the past 23 hours are present, **Then** those sessions are included in the token count.
2. **Given** a lookback window of 24 hours, **When** a transcript is timestamped 25 hours ago, **Then** it is excluded from the count.

---

### Edge Cases

- What happens when a transcript file is corrupt or contains invalid records? (Skip invalid entries; count valid ones.)
- What if the same session appears in multiple projects? (Sum all projects; no deduplication by session is needed since each project directory is independent.)
- What if `token_baseline` is configured but the transcript directory is empty? (Return zero usage, no anomaly.)
- What if token counts overflow a 32-bit integer? (Use 64-bit integers for accumulation.)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The tokens module MUST read session transcript files from the user's Claude Code transcript directory (`~/.claude/projects/**/*.jsonl`) to determine actual daily token consumption.
- **FR-002**: The module MUST sum `input_tokens`, `output_tokens`, `cache_creation_input_tokens`, and `cache_read_input_tokens` from all assistant messages within the configured lookback window.
- **FR-003**: The module MUST default to a 24-hour lookback window when none is explicitly configured.
- **FR-004**: The module MUST compare the summed token total against the anomaly threshold derived from the configured `token_baseline` and emit a TOKEN_ANOMALY finding when the total exceeds the threshold.
- **FR-005**: The module MUST return zero findings (not an error) when the transcript directory is absent, empty, or all files are unreadable.
- **FR-006**: The module MUST skip individual malformed or unreadable transcript entries without failing the entire run.
- **FR-007**: The module MUST emit a MODULE_UNAVAILABLE finding when no `token_baseline` is configured (existing behaviour, preserved).

### Key Entities

- **Transcript File**: A structured log file under `~/.claude/projects/<project>/` containing one record per line, each representing a message event in a Claude Code session.
- **Usage Record**: An assistant message entry in a transcript file that contains token count fields (`input_tokens`, `output_tokens`, `cache_creation_input_tokens`, `cache_read_input_tokens`) and a timestamp.
- **Token Baseline**: Configuration specifying `daily_mean`, `std_dev`, and optional `multiplier` used to compute the anomaly threshold.
- **Lookback Window**: The time range (default 24 hours) used to filter which transcript records contribute to the daily token total.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: The tokens module produces a TOKEN_ANOMALY finding when real transcript data shows daily usage exceeding `daily_mean + multiplier × std_dev`, replacing the current always-zero placeholder.
- **SC-002**: The module processes all transcript files across all project directories within the lookback window in under 2 seconds on a machine with up to 100 session files totalling 50 MB.
- **SC-003**: The module produces zero false positives (no TOKEN_ANOMALY) when token usage is within the normal range defined by the baseline.
- **SC-004**: The module never returns a non-nil error due to missing, empty, or corrupted transcript data.
- **SC-005**: Existing tests continue to pass without modification; new tests cover the real-data reading path with at least 3 distinct scenarios (above threshold, below threshold, no transcripts).

## Assumptions

- Claude Code stores session transcripts as structured log files under `~/.claude/projects/<project-slug>/` with one JSON object per line; this structure is stable and does not require authentication to read.
- Assistant messages that contain token usage are identified by having a `message.usage` object with numeric `input_tokens` and `output_tokens` fields.
- Each transcript record contains a top-level `timestamp` field in RFC3339 format used for lookback window filtering.
- The transcript directory path defaults to `~/.claude/projects` and does not need to be user-configurable in this feature scope.
- The tokens module is only concerned with total token volume across all projects, not per-project or per-model breakdowns.
- Running `ai-check-guardrails` on a machine where no Claude Code transcripts exist is a valid scenario; silent no-op is the correct behaviour.
- A future spec may add per-project breakdown reporting; this feature only implements the total daily sum needed for anomaly detection.
