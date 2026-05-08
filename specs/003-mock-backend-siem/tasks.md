# Tasks: Mock Backend SIEM Endpoint

**Input**: Design documents from `/specs/003-mock-backend-siem/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/api.md ✅

**Tests**: Not explicitly requested in spec — no test tasks generated.

**Organization**: Tasks grouped by user story. Each phase is independently testable.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1/US2/US3)
- All paths relative to repository root unless noted

---

## Phase 1: Setup (Project Scaffolding)

**Purpose**: Create the `mock-backend/` directory and all project skeleton files. Nothing is functional yet — this is structure only.

- [x] T001 Create `mock-backend/go.mod` with `module github.com/kaihendry/ai-check-guardrails/mock-backend` and `go 1.26`
- [x] T002 Create `mock-backend/template.yml` — SAM CloudFormation template defining: `AWS::Serverless::Function` (ARM64, provided.al2, Makefile build, 10s timeout, 128MB), `AWS::Serverless::HttpApi` (catch-all `/{proxy+}` ANY route), `AWS::DynamoDB::Table` (PK=`pk` String, SK=`sk` String, PAY_PER_REQUEST), env var `DYNAMODB_TABLE` wired to `!Ref` the table, stack name `mock-siem-backend`, SAM deploy config using profile `AdministratorAccess-407461997746`
- [x] T003 [P] Create `mock-backend/Makefile` with targets: `build-MainFunction` (CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o ${ARTIFACTS_DIR}/bootstrap), `deploy` (sam build then sam deploy --profile AdministratorAccess-407461997746), `destroy` (sam delete), `local` (go run . for local HTTP server on port 8080)
- [x] T004 [P] Create `mock-backend/main.go` skeleton: package main, import stubs, empty `main()` that detects `AWS_LAMBDA_FUNCTION_NAME` env var and branches to Lambda or local mode (no handlers yet)

**Checkpoint**: `mock-backend/` directory exists with all four files; `go vet ./...` compiles with no errors (empty handlers are fine)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types, DynamoDB client, and HTTP routing. MUST be complete before any user story handler can be implemented.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Add dependencies to `mock-backend/go.mod` by running `go get github.com/apex/gateway/v2` and `go get github.com/aws/aws-sdk-go-v2/...` (config, service/dynamodb, feature/dynamodb/attributevalue) from within `mock-backend/`; commit resulting `go.sum`
- [x] T006 Define `AuditRun` and `Finding` structs in `mock-backend/main.go` matching the JSON fields emitted by `internal/siem.Emit()` (schema_version, run_id, timestamp, host, user, mode, version, findings, score, exit_code, duration_ms; Finding: type, severity, module, resource, description, remediation, confidence)
- [x] T007 Implement `newDynamoClient()` in `mock-backend/main.go`: load AWS config from environment, return `*dynamodb.Client`; read table name from `DYNAMODB_TABLE` env var into a package-level `tableName` variable; fall back to `"mock-siem-events"` if env var is absent
- [x] T008 [P] Complete dual-mode `main()` in `mock-backend/main.go`: register routes on `http.DefaultServeMux`, then if `AWS_LAMBDA_FUNCTION_NAME != ""` call `gateway.ListenAndServe("", nil)`, else call `http.ListenAndServe(":"+port, nil)` where port defaults to `"8080"`
- [x] T009 [P] Register three routes in `mock-backend/main.go`: `http.HandleFunc("POST /", handlePost)`, `http.HandleFunc("GET /", handleGet)`, `http.HandleFunc("GET /event/", handleDetail)` — stub each handler to return `http.StatusNotImplemented` for now

**Checkpoint**: `go build ./...` succeeds; `make local` starts a server; `curl -X POST http://localhost:8080/` returns 501; `curl http://localhost:8080/` returns 501

---

## Phase 3: User Story 1 — Submit Guardrails Findings (Priority: P1) 🎯 MVP

**Goal**: Accept a valid `AuditRun` JSON POST, store it in DynamoDB, return `{run_id, sk}`.

**Independent Test**: Run `make local`, POST the sample payload from `contracts/api.md` with `curl`, expect HTTP 201 and a JSON body containing `run_id` and `sk`. Confirm the item is visible in DynamoDB (via `aws dynamodb scan` or AWS console).

- [x] T010 [US1] Implement `handlePost()` in `mock-backend/main.go`: decode JSON request body into `AuditRun`; return `400` with `{"error": "invalid payload: ..."}` if body is malformed or `run_id` is empty; call `putEvent()` on success
- [x] T011 [US1] Implement `putEvent(ctx, run AuditRun)` in `mock-backend/main.go`: compute `sk = run.Timestamp.UTC().Format(time.RFC3339) + "#" + run.RunID`; marshal `run.Findings` to a JSON string; build DynamoDB attribute map with all fields from data-model.md; call `dynamoClient.PutItem()`; return error on failure
- [x] T012 [US1] Complete `handlePost()` response: on `putEvent()` error return `500` with `{"error":"storage error"}`; on success return `201` with `{"run_id": run.RunID, "sk": sk}`

**Checkpoint**: User Story 1 fully functional — POST stores data; 400 on bad input; 201 with correct body on success

---

## Phase 4: User Story 2 — View Summary of Recent Submissions (Priority: P2)

**Goal**: Browser GET on `/` shows a table of the 50 most recent submissions (newest first), each row linking to the detail page.

**Independent Test**: After submitting 2–3 test payloads via US1, open `http://localhost:8080/` in a browser; confirm a table appears with Received, Host, User, Score, Findings count, Mode columns; confirm each row is a clickable link; confirm empty-state message shows when no submissions exist.

- [x] T013 [US2] Define `SummaryRow` struct in `mock-backend/main.go` with fields: `SK`, `RunID`, `Timestamp`, `Host`, `User`, `Score`, `FindingCount`, `Mode`
- [x] T014 [US2] Implement `listEvents(ctx)` in `mock-backend/main.go`: call DynamoDB `Query` with `KeyConditionExpression: "pk = :all"`, `ScanIndexForward: false`, `Limit: 50`; unmarshal each item into `SummaryRow`; decode `FindingCount` from the length of the `findings` JSON string (unmarshal and count, or store count as attribute in T011)
- [x] T015 [US2] Define HTML summary template as a `const` string in `mock-backend/main.go`: `<html>` with a `<table>` rendering `SummaryRow` slice; each row links to `/event/` + url-safe base64 of `SK`; include a submission count header; XSS-safe via `html/template`
- [x] T016 [US2] Implement `handleGet()` in `mock-backend/main.go`: call `listEvents()`; render summary template with `html/template.Execute(w, rows)`; return 500 with plain error message if DynamoDB query fails; render empty-state paragraph when `len(rows) == 0`

**Checkpoint**: User Stories 1 and 2 both work — POST stores data; GET `/` renders a live table; empty state shows correctly on a fresh instance

---

## Phase 5: User Story 3 — View Pretty-Printed Findings (Priority: P3)

**Goal**: Clicking a row in the summary opens a detail page showing full metadata and a per-finding table with Module, Severity, Type, Resource, Description, Remediation, Confidence.

**Independent Test**: Submit a payload containing at least 2 findings; open the summary page; click the first row; confirm the detail page shows all metadata fields and a findings table with all columns populated; confirm "No findings" message on a clean-run payload; confirm the back link returns to `/`.

- [x] T017 [US3] Define HTML detail template as a `const` string in `mock-backend/main.go`: metadata `<dl>` block (RunID, Timestamp, Host, User, Mode, Version, Score, ExitCode, DurationMs) followed by a findings `<table>` (Module, Severity, Type, Resource, Description, Remediation, Confidence); "No findings (clean run)" paragraph when slice is empty; back link `<a href="/">← Back</a>`
- [x] T018 [US3] Implement `getEvent(ctx, sk string)` in `mock-backend/main.go`: call DynamoDB `GetItem` with `pk="all"` and the provided `sk`; return `nil, nil` if item not found; unmarshal all fields including JSON-decode of the `findings` attribute back into `[]Finding`
- [x] T019 [US3] Implement `handleDetail()` in `mock-backend/main.go`: extract path suffix after `/event/`; URL-safe base64-decode to recover `sk`; return 400 if decode fails; call `getEvent()`; return 404 if result is nil; render detail template via `html/template`

**Checkpoint**: All three user stories fully functional end-to-end — POST → GET / (summary) → GET /event/{sk} (detail)

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: README, validation, and smoke test before deploy

- [x] T020 [P] Write `mock-backend/README.md` covering: prerequisites, `make deploy` flow, `make local` usage, curl example POST, browser view URL, `make destroy`; reference `specs/003-mock-backend-siem/quickstart.md` for full detail
- [x] T021 [P] Run `sam validate --template mock-backend/template.yml --profile AdministratorAccess-407461997746` and resolve any template errors
- [x] T022 [P] Run `go vet ./...` and `go build ./...` from `mock-backend/`; fix any compilation warnings
- [ ] T023 End-to-end smoke test per `quickstart.md`: `make local`, POST sample payload via curl (expect 201), open browser at `http://localhost:8080/` (expect table), click a row (expect detail page with findings), POST a clean-run payload (expect "No findings" on detail page)
- [ ] T024 Deploy to AWS: `make deploy` from `mock-backend/`; capture the API Gateway URL from stack outputs; run the smoke test against the live URL; confirm DynamoDB table exists in AWS console

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 completion — **BLOCKS all user stories**
- **US1 (Phase 3)**: Depends on Phase 2 — no dependency on US2 or US3
- **US2 (Phase 4)**: Depends on Phase 2 + US1 data being available for testing (US1 must be done first)
- **US3 (Phase 5)**: Depends on Phase 2 + US2 (detail links come from summary page rows)
- **Polish (Phase 6)**: Depends on all user story phases complete

### Within Each Phase

- Models/structs before client usage
- Client initialisation before handler implementation
- Handler stub before full implementation

### Parallel Opportunities

- T003 and T004 (Phase 1): Different files — run in parallel
- T008 and T009 (Phase 2): Different concerns — run in parallel after T006/T007
- T021, T022, T023 (Phase 6): Independent — run in parallel

---

## Parallel Example: Phase 1

```bash
# Launch in parallel (different files):
Task T003: "Create mock-backend/Makefile"
Task T004: "Create mock-backend/main.go skeleton"
```

## Parallel Example: Phase 2

```bash
# After T006/T007 complete, launch in parallel:
Task T008: "Complete dual-mode main()"
Task T009: "Register three route stubs"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — blocks all stories)
3. Complete Phase 3: User Story 1 (POST endpoint)
4. **STOP and VALIDATE**: POST sample payload, confirm 201, confirm DynamoDB item exists
5. Deploy MVP to AWS if needed

### Incremental Delivery

1. Setup + Foundational → Server starts, routes stubbed
2. US1 → POST works, data stored
3. US2 → Summary view live in browser
4. US3 → Detail view with pretty findings
5. Polish → README, deploy, smoke test

---

## Notes

- All handlers live in a single `mock-backend/main.go` file — no sub-packages per constitution Simplicity principle
- HTML templates defined as `const` strings in main.go — no separate template files needed
- `html/template` auto-escapes all user data — no manual XSS protection needed
- Local mode uses the same DynamoDB table as AWS (via `AWS_PROFILE` env) unless `DYNAMODB_TABLE` is unset, in which case it errors; run with a real AWS profile or point to DynamoDB Local
- `[P]` tasks = different files or independent concerns with no shared state dependency
