# Contract: SIEM Event Schema

Every audit run emits exactly one SIEM event as a single-line JSON object to
stdout and optionally via HTTP POST to `siem_endpoint`. The schema below is the
normative definition; implementations MUST NOT add or remove top-level fields
without incrementing the `schema_version`.

## JSON Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12",
  "title": "AuditRun",
  "type": "object",
  "required": ["schema_version","run_id","timestamp","host","user","mode",
               "version","findings","score","exit_code","duration_ms"],
  "properties": {
    "schema_version": { "type": "string", "const": "1.0" },
    "run_id":         { "type": "string", "format": "uuid" },
    "timestamp":      { "type": "string", "format": "date-time" },
    "host":           { "type": "string" },
    "user":           { "type": "string" },
    "mode":           { "type": "string", "enum": ["monitor","enforce"] },
    "version":        { "type": "string" },
    "findings": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["type","severity","module","resource","description","remediation"],
        "properties": {
          "type":        { "type": "string" },
          "severity":    { "type": "string", "enum": ["INFO","WARN","HIGH","CRITICAL"] },
          "module":      { "type": "string" },
          "resource":    { "type": "string" },
          "description": { "type": "string" },
          "remediation": { "type": "string" },
          "confidence":  { "type": "number", "minimum": 0, "maximum": 1 }
        }
      }
    },
    "score":       { "type": "integer", "minimum": 0, "maximum": 100 },
    "exit_code":   { "type": "integer", "enum": [0, 1, 2] },
    "duration_ms": { "type": "integer", "minimum": 0 }
  }
}
```

## Example (clean run)

```json
{
  "schema_version": "1.0",
  "run_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2026-05-07T09:00:00Z",
  "host": "macbook-pro.local",
  "user": "alice",
  "mode": "monitor",
  "version": "0.1.0",
  "findings": [],
  "score": 100,
  "exit_code": 0,
  "duration_ms": 4820
}
```

## Example (findings present)

```json
{
  "schema_version": "1.0",
  "run_id": "660e8400-e29b-41d4-a716-446655440001",
  "timestamp": "2026-05-07T09:05:00Z",
  "host": "macbook-pro.local",
  "user": "bob",
  "mode": "monitor",
  "version": "0.1.0",
  "findings": [
    {
      "type": "UNAPPROVED_MCP",
      "severity": "HIGH",
      "module": "mcp",
      "resource": "mcp://custom-data-exporter",
      "description": "MCP 'custom-data-exporter' is not on the approved allowlist.",
      "remediation": "Contact the security team to request allowlist approval or remove the MCP."
    },
    {
      "type": "MISSING_PRECOMMIT_HOOK",
      "severity": "HIGH",
      "module": "hooks",
      "resource": "/home/bob/projects/myrepo/.git/hooks/pre-commit",
      "description": "Required pre-commit hook 'gitleaks' not found in repository.",
      "remediation": "Run 'gitleaks install' in the repository root."
    }
  ],
  "score": 70,
  "exit_code": 1,
  "duration_ms": 6120
}
```

## HTTP POST Transport

When `siem_endpoint` is configured:

- Method: `POST`
- Content-Type: `application/json`
- Body: single JSON object (same schema as above, compact/single-line)
- Auth: Bearer token from env var `AI_GUARDRAILS_SIEM_TOKEN`
- Timeout: 10 seconds
- On failure: log error to stderr; do not retry in v1; exit code unchanged
