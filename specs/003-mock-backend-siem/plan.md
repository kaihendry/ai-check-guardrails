# Implementation Plan: Mock Backend SIEM Endpoint

**Branch**: `003-mock-backend-siem` | **Date**: 2026-05-08 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/003-mock-backend-siem/spec.md`

## Summary

Deploy a serverless Go HTTP service in `mock-backend/` that receives `AuditRun` JSON payloads from `ai-check-guardrails` via POST and renders a browser-accessible summary of the most recent 50 submissions with scores and pretty-printed findings. Follows the [kaihendry/helloworld](https://github.com/kaihendry/helloworld) AWS SAM pattern: Go binary on Lambda (ARM64, provided.al2), apex/gateway for HTTP bridging, Makefile build, DynamoDB for persistence.

## Technical Context

**Language/Version**: Go 1.26  
**Primary Dependencies**: `github.com/apex/gateway/v2` (Lambda HTTP adapter), `github.com/aws/aws-sdk-go-v2` (DynamoDB), `html/template` (stdlib)  
**Storage**: DynamoDB single table — PK `"all"`, SK `{RFC3339}#{run_id}` for time-ordered queries without a GSI  
**Testing**: `go test ./...` with `net/http/httptest`  
**Target Platform**: AWS Lambda `provided.al2` ARM64 + local HTTP server (port 8080)  
**Project Type**: serverless web-service  
**Performance Goals**: POST acknowledged within 500 ms; summary of 50 records rendered within 2 s  
**Constraints**: Lambda timeout 10 s; DynamoDB on-demand billing; no authentication (mock/dev use only)  
**Scale/Scope**: Development and testing use; up to 50 recent submissions visible; AWS account 407461997746, profile `AdministratorAccess-407461997746`

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Status | Notes |
|-----------|------|--------|-------|
| I. Simplicity | One package, minimal dependencies | ✅ PASS | Single `main` package; apex/gateway + AWS SDK v2 only; html/template from stdlib |
| I. Simplicity | No speculative abstractions | ✅ PASS | Two handlers (`handlePost`, `handleGet`); no interfaces beyond what Lambda requires |
| I. Simplicity | New dependency justified | ✅ PASS | DynamoDB justified: Lambda is ephemeral — in-memory storage loses data on cold start; S3 would need 50 separate GetObject calls for summary |
| II. Integrity | Errors surface explicitly | ✅ PASS | POST returns 400 on bad payload, 201 on success; GET returns 500 with message if DynamoDB fails |
| II. Integrity | Output is deterministic | ✅ PASS | Summary sorted by stored timestamp; no hidden side effects |
| Security | Input validated before storage | ✅ PASS | JSON payload decoded into typed struct; unrecognised fields ignored |
| Security | No sensitive data in output | ✅ PASS | HTML output escapes via html/template; no tokens or credentials stored |

**Complexity Tracking**: no violations.

## Project Structure

### Documentation (this feature)

```text
specs/003-mock-backend-siem/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── api.md
└── tasks.md             # Phase 2 output (/speckit-tasks — NOT created here)
```

### Source Code (repository root)

```text
mock-backend/
├── main.go              # Lambda + local server entry; HTTP handlers
├── go.mod               # module github.com/kaihendry/ai-check-guardrails/mock-backend
├── go.sum
├── template.yml         # SAM CloudFormation template
├── Makefile             # build-MainFunction, deploy, destroy, local
└── README.md
```

**Structure Decision**: Single-directory, single-package project following the helloworld pattern exactly. No `cmd/`, no `internal/` — the entire service is one Go file analogous to `helloworld/main.go`.
