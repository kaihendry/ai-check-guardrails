# Research: Mock Backend SIEM Endpoint

**Phase**: 0 | **Date**: 2026-05-08

## DynamoDB Table Design for Time-Ordered Queries

**Decision**: Single table, constant partition key `"all"`, sort key `{RFC3339_timestamp}#{run_id}`

**Rationale**: Lambda is ephemeral — in-memory storage loses data on cold start. DynamoDB with a constant PK allows `Query(pk="all", ScanIndexForward=false, Limit=50)` to return the 50 most recent submissions without a GSI or Scan. RFC3339 timestamps sort lexicographically, so newest-first is achieved by reversing the scan direction. Appending `#run_id` ensures uniqueness when two runs complete within the same second.

**Alternatives considered**:
- In-memory (map/slice): Lost on every cold start; unusable for a persistent mock
- S3 per-object: Requires 50 separate GetObject calls to build summary; ListObjectsV2 alone can't return score
- DynamoDB with GSI on timestamp: Adds a GSI (extra cost/complexity) when a compound SK achieves the same result
- DynamoDB Scan: Works for small tables but degrades at scale; Query is better practice

---

## apex/gateway Pattern for Local + Lambda

**Decision**: Use `github.com/apex/gateway/v2` to serve standard Go `http.Handler` in both Lambda and local modes, identical to helloworld

**Rationale**: Single codebase runs locally with `go run .` and in Lambda with zero code divergence. The only branch is `os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""` to select the listen method.

**Alternatives considered**:
- AWS Lambda Go SDK directly: Requires Lambda-specific event types; breaks local testing
- Custom shim: Reinventing apex/gateway; violates Simplicity principle

---

## HTML Rendering in Go

**Decision**: Use stdlib `html/template` with embedded templates; no external template engine

**Rationale**: `html/template` provides XSS-safe HTML escaping automatically. No additional dependency needed. Template is embedded via `//go:embed` or defined as a string constant in main.go to keep the project single-file.

**Alternatives considered**:
- External templating library (templ, jet): Adds dependency; unnecessary for two simple pages
- JSON API + separate frontend: Increases complexity; contradicts "view in a browser" requirement
- `text/template`: Does not auto-escape HTML; security risk per constitution

---

## SAM Template Structure

**Decision**: Follow helloworld `template.yml` pattern: `AWS::Serverless::Function` with `Events.Api` using `HttpApi`, ARM64, `provided.al2`, Makefile build. Add `AWS::DynamoDB::Table` resource with `BillingMode: PAY_PER_REQUEST`.

**Rationale**: Proven pattern from helloworld. On-demand DynamoDB removes capacity planning. ARM64 reduces Lambda cost ~20%.

**Alternatives considered**:
- Provisioned DynamoDB: Unnecessary for a mock backend with unpredictable/low traffic
- x86_64 Lambda: Higher cost with no benefit for this workload

---

## Environment Variable Wiring

**Decision**: Pass DynamoDB table name to Lambda via environment variable `DYNAMODB_TABLE` set from a `!Ref` to the table resource in the SAM template.

**Rationale**: Standard SAM pattern; keeps table name out of compiled code.

---

## Stack Name and Region

**Decision**: Stack `mock-siem-backend`, default region `ap-southeast-1` (Singapore — consistent with account locality), AWS profile `AdministratorAccess-407461997746`.

**Rationale**: User confirmed this profile via `aws sts get-caller-identity`.

**Alternatives considered**: `us-east-1` is common default but ap-southeast-1 is closer to the account owner's likely location.

---

## Findings Display Format

**Decision**: Summary page shows a table: Received, Host, User, Score, Finding Count. Each row links to a detail page `/event/{sk}` (URL-safe base64 of the DynamoDB SK). Detail page shows: metadata block + findings table (Module, Severity, Type, Description, Remediation).

**Rationale**: Matches the spec requirement for "summary of last posts, the score, and a pretty print of the findings" in a browser-accessible HTML view. No JavaScript required.

**Alternatives considered**:
- Raw JSON pretty-print with `<pre>`: Meets "pretty print" literally but poor UX for severities/modules
- Single page with accordion: Adds complexity; two-page approach is simpler
