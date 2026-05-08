# Feature Specification: Mock Backend SIEM Endpoint

**Feature Branch**: `003-mock-backend-siem`
**Created**: 2026-05-08
**Status**: Draft
**Input**: User description: "Create in dir called mock-backend a AWS SAM project inspired by https://github.com/kaihendry/helloworld to serve as a SIEMEndpoint. Once posted viewing the backend, you should be able to see a summary of the last posts, the score and a pretty print of the findings"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Submit Guardrails Findings (Priority: P1)

A developer or CI pipeline runs the `ai-check-guardrails` tool against an AI system and submits the resulting findings and score to the mock SIEM endpoint via an HTTP POST. This is the core data ingestion path — without it, there is nothing to view.

**Why this priority**: All other user stories depend on data being submitted first. This is the entry point.

**Independent Test**: Can be fully tested by sending a POST request with a sample findings payload and verifying the endpoint returns a success acknowledgement.

**Acceptance Scenarios**:

1. **Given** the mock SIEM endpoint is running, **When** a client sends a POST with a valid findings payload including a score and detection results, **Then** the system accepts the submission and returns a success response.
2. **Given** the endpoint has received multiple submissions, **When** a new POST arrives, **Then** the new submission is stored alongside previous ones without overwriting them.
3. **Given** a POST is sent with a malformed or empty body, **When** the system receives it, **Then** it returns a clear error response indicating the payload was invalid.

---

### User Story 2 - View Summary of Recent Submissions (Priority: P2)

A security operator or developer opens the mock SIEM backend in a browser and sees a human-readable summary of the most recent submissions — showing how many have been received, when they arrived, and what scores were reported.

**Why this priority**: This is the primary observability use case. Operators need a quick overview before drilling into details.

**Independent Test**: Can be fully tested once at least one POST has been submitted — navigate to the root URL and confirm a summary table or list appears with timestamps and scores.

**Acceptance Scenarios**:

1. **Given** the backend has received at least one submission, **When** a user visits the root URL in a browser, **Then** they see a list of recent submissions with timestamp, score, and entry count.
2. **Given** many submissions have been received, **When** a user views the summary, **Then** only the most recent submissions (up to a reasonable limit) are displayed, ordered newest-first.
3. **Given** no submissions have been received yet, **When** a user visits the root URL, **Then** they see a friendly empty-state message rather than an error.

---

### User Story 3 - View Pretty-Printed Findings (Priority: P3)

A developer clicks into a specific submission from the summary page and sees a human-readable, formatted view of the full findings payload — all detection module results clearly labelled and easy to scan.

**Why this priority**: The summary gives a quick score; this story provides the detail needed to act on findings.

**Independent Test**: Can be fully tested by selecting any submission from the summary and verifying the full findings are displayed in a structured, readable format.

**Acceptance Scenarios**:

1. **Given** the user is viewing the summary page, **When** they select a specific submission, **Then** the full findings for that submission are displayed in a structured, readable layout.
2. **Given** the findings include multiple detection module results, **When** displayed, **Then** each module result is clearly labelled with its name, status, and any relevant detail.
3. **Given** the findings payload contains deeply nested data, **When** displayed, **Then** the layout remains readable without raw JSON clutter.

---

### Edge Cases

- What happens when the POST payload exceeds a reasonable size limit?
- How does the system handle concurrent submissions arriving at the same time?
- What if the backend is restarted — are previously submitted findings retained or lost?
- How does the summary page behave if a submission's findings are partially malformed?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST expose an HTTP endpoint that accepts POST requests containing AI guardrails findings and a score.
- **FR-002**: The system MUST store each submission with its received timestamp and full payload.
- **FR-003**: The system MUST expose a view accessible via a browser that lists recent submissions in reverse-chronological order.
- **FR-004**: Each submission entry in the summary view MUST display at minimum: the received timestamp and the score.
- **FR-005**: The system MUST provide a detailed view for each submission showing a human-readable, pretty-printed rendering of the full findings payload.
- **FR-006**: The system MUST return a meaningful success or error response for every POST request.
- **FR-007**: The system MUST handle at least the last 50 submissions without degradation in the summary view.
- **FR-008**: The project MUST be deployable as a self-contained unit from the `mock-backend` directory.

### Key Entities

- **Submission**: A single received POST, containing a score (numeric), a findings object (detection results from all modules), and a received-at timestamp.
- **Score**: A numeric value representing the overall guardrails assessment outcome for a submission.
- **Findings**: The structured collection of detection module results included in a submission payload.
- **Summary View**: The browser-facing page listing recent submissions with key metadata.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A submission POST completes and is acknowledged within 2 seconds under normal conditions.
- **SC-002**: The summary view loads and displays the most recent 50 submissions within 3 seconds in a browser.
- **SC-003**: 100% of submitted findings are retrievable in their original detail view without data loss.
- **SC-004**: A developer unfamiliar with the project can submit a test payload and view the result in under 5 minutes using only the project README.
- **SC-005**: The pretty-printed findings view requires no additional tools (no browser extensions, no JSON formatter) to read comfortably.

## Assumptions

- The mock backend is intended for local development and testing, not production use — therefore no authentication or authorisation is required.
- Findings payload structure follows the output format of the `ai-check-guardrails` binary defined in this repository.
- Storage is in-memory or file-based; persistence across restarts is desirable but not required for the initial version.
- The summary view displays a maximum of 50 recent submissions; older entries may be dropped without user notification.
- The project will live in a directory named `mock-backend` at the root of this repository.
- A minimal README will be provided explaining how to deploy and test the endpoint locally.
