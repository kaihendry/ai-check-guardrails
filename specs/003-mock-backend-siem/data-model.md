# Data Model: Mock Backend SIEM Endpoint

**Phase**: 1 | **Date**: 2026-05-08

## DynamoDB Table: `mock-siem-events`

Managed by the SAM stack. Deleted when the stack is destroyed.

### Key Schema

| Attribute | Type | Role | Example |
|-----------|------|------|---------|
| `pk` | String | Partition Key | `"all"` (constant for all items) |
| `sk` | String | Sort Key | `"2026-05-08T10:30:00Z#a1b2c3d4-..."` |

Sort key format: `{RFC3339_timestamp}#{run_id}` — RFC3339 sorts lexicographically; appending `#run_id` ensures uniqueness within the same second.

### Item Attributes

| Attribute | DynamoDB Type | Source | Notes |
|-----------|---------------|--------|-------|
| `pk` | S | constant | Always `"all"` |
| `sk` | S | `{run.Timestamp.UTC().Format(time.RFC3339)}#{run.RunID}` | Sort + unique key |
| `run_id` | S | `AuditRun.RunID` | UUID for linking |
| `schema_version` | S | `AuditRun.SchemaVersion` | |
| `timestamp` | S | `AuditRun.Timestamp` | ISO8601 UTC |
| `host` | S | `AuditRun.Host` | |
| `user` | S | `AuditRun.User` | |
| `mode` | S | `AuditRun.Mode` | |
| `version` | S | `AuditRun.Version` | Binary version |
| `score` | N | `AuditRun.Score` | Integer 0-100 |
| `exit_code` | N | `AuditRun.ExitCode` | |
| `duration_ms` | N | `AuditRun.DurationMs` | |
| `findings` | S | `json.Marshal(AuditRun.Findings)` | Full findings JSON |

`findings` is stored as a JSON string (not a DynamoDB Map) to avoid attribute name conflicts and simplify unmarshalling on read.

### Access Patterns

| Pattern | DynamoDB Operation | Parameters |
|---------|--------------------|------------|
| List 50 most recent | `Query` | `pk="all"`, `ScanIndexForward=false`, `Limit=50` |
| Get single event | `GetItem` | `pk="all"`, `sk={sk_value}` |
| Store new event | `PutItem` | All attributes above |

## Entity: AuditRun (received payload)

Mirrors `internal/audit/AuditRun` in the parent repository. Decoded from POST body.

```
AuditRun {
    schema_version  string
    run_id          string       // UUID
    timestamp       time.Time    // UTC
    host            string
    user            string
    mode            string       // "audit" | "enforce"
    version         string
    findings        []Finding
    score           int          // 0–100
    exit_code       int          // 0 = clean, 1 = findings, 2 = error
    duration_ms     int64
}

Finding {
    type         string
    severity     string          // INFO | WARN | HIGH | CRITICAL
    module       string
    resource     string
    description  string
    remediation  string
    confidence   *float64        // optional, 0.0–1.0
}
```

### Validation Rules

- `run_id` MUST be non-empty; reject with 400 if missing
- `score` MUST be in range 0–100; values outside are stored as-is with a logged warning
- `findings` MAY be empty (clean run)
- Unrecognised JSON fields are silently ignored (forward-compatible)

## Summary View Entity (computed on GET /)

Derived from DynamoDB Query result. No separate storage.

```
SummaryRow {
    SK             string    // DynamoDB SK — used as URL token
    RunID          string
    Timestamp      string    // formatted for display
    Host           string
    User           string
    Score          int
    FindingCount   int       // len(findings)
    Mode           string
}
```
