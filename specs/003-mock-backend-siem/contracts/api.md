# API Contract: Mock Backend SIEM Endpoint

**Version**: 1.0 | **Date**: 2026-05-08

## Base URL

```
https://{api-gateway-id}.execute-api.{region}.amazonaws.com/
```

Resolved after `make deploy`. Also available locally at `http://localhost:8080/`.

---

## POST /

Receive an `AuditRun` submission from `ai-check-guardrails`.

### Request

```
POST /
Content-Type: application/json
```

**Body**: `AuditRun` JSON as emitted by `internal/siem.Emit()`:

```json
{
  "schema_version": "1.0",
  "run_id": "a1b2c3d4-e5f6-...",
  "timestamp": "2026-05-08T10:30:00Z",
  "host": "laptop.local",
  "user": "hendry",
  "mode": "audit",
  "version": "1.2.3",
  "findings": [
    {
      "type": "config_issue",
      "severity": "HIGH",
      "module": "mcp",
      "resource": "claude_desktop_config.json",
      "description": "MCP server with wildcard resource access",
      "remediation": "Restrict MCP server permissions",
      "confidence": 0.9
    }
  ],
  "score": 72,
  "exit_code": 1,
  "duration_ms": 145
}
```

### Responses

| Status | Condition | Body |
|--------|-----------|------|
| `201 Created` | Stored successfully | `{"run_id": "{run_id}", "sk": "{dynamodb_sk}"}` |
| `400 Bad Request` | Missing `run_id` or malformed JSON | `{"error": "invalid payload: {detail}"}` |
| `500 Internal Server Error` | DynamoDB write failure | `{"error": "storage error"}` |

---

## GET /

Render an HTML summary page of the most recent 50 submissions.

### Request

```
GET /
Accept: text/html
```

No parameters required.

### Response

```
200 OK
Content-Type: text/html; charset=utf-8
```

HTML page containing:
- Page title and submission count
- Table of recent submissions: columns `Received`, `Host`, `User`, `Score`, `Findings`, `Mode`
- Each row links to `GET /event/{sk}` for the detail view
- Empty-state message when no submissions exist

### Error Response

```
500 Internal Server Error
Content-Type: text/html; charset=utf-8
```

Plain error message if DynamoDB query fails.

---

## GET /event/{sk}

Render an HTML detail page for a single submission.

`{sk}` is the URL-safe base64-encoded DynamoDB sort key returned in the POST response body and embedded in summary page links.

### Request

```
GET /event/{sk}
Accept: text/html
```

### Response

```
200 OK
Content-Type: text/html; charset=utf-8
```

HTML page containing:
- Metadata block: RunID, Timestamp, Host, User, Mode, Version, Score, ExitCode, DurationMs
- Findings table: columns `Module`, `Severity`, `Type`, `Resource`, `Description`, `Remediation`, `Confidence`
- "No findings (clean run)" message when findings array is empty
- Back link to `GET /`

### Error Responses

| Status | Condition |
|--------|-----------|
| `400 Bad Request` | `{sk}` path parameter is missing or cannot be base64-decoded |
| `404 Not Found` | No item found in DynamoDB for the decoded SK |
| `500 Internal Server Error` | DynamoDB read failure |
