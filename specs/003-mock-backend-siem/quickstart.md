# Quickstart: Mock Backend SIEM Endpoint

## Prerequisites

- Go 1.22+
- AWS SAM CLI (`sam --version`)
- AWS CLI configured with profile `AdministratorAccess-407461997746`
- Docker (for `sam build`)

---

## Deploy to AWS

```bash
cd mock-backend
make deploy
```

This runs `sam build` then `sam deploy` and outputs the SIEM endpoint URL. Copy the URL and set it in your `ai-check-guardrails` config:

```bash
export AI_GUARDRAILS_SIEM_ENDPOINT=https://{api-id}.execute-api.{region}.amazonaws.com/
```

Or add it to your config file:
```yaml
siem_endpoint: "https://{api-id}.execute-api.{region}.amazonaws.com/"
```

---

## Run Locally

```bash
cd mock-backend
make local
```

Starts the server on `http://localhost:8080`. Point `ai-check-guardrails` at it:

```bash
export AI_GUARDRAILS_SIEM_ENDPOINT=http://localhost:8080/
```

> **Note**: Local mode uses an in-memory store (no DynamoDB). Data is lost on restart. For persistent local testing, deploy to AWS.

---

## Submit a Test Payload

```bash
ai-check-guardrails
```

With `AI_GUARDRAILS_SIEM_ENDPOINT` set, each run automatically posts to the endpoint.

Or manually with curl:

```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "schema_version": "1.0",
    "run_id": "test-001",
    "timestamp": "2026-05-08T10:00:00Z",
    "host": "test-host",
    "user": "tester",
    "mode": "audit",
    "version": "dev",
    "findings": [],
    "score": 100,
    "exit_code": 0,
    "duration_ms": 10
  }'
```

---

## View the Dashboard

Open in a browser:

```
http://localhost:8080/        # local
https://{api-id}.execute-api.{region}.amazonaws.com/   # AWS
```

You will see a summary table of recent submissions. Click any row to see the full pretty-printed findings for that run.

---

## Tear Down AWS Resources

```bash
cd mock-backend
make destroy
```

This deletes the CloudFormation stack, Lambda, DynamoDB table, and API Gateway.
