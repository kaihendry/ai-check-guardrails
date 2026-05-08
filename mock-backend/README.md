# mock-backend

AWS SAM serverless SIEM endpoint for [ai-check-guardrails](https://github.com/kaihendry/ai-check-guardrails).

POST findings to it, then open it in a browser to see a summary and pretty-printed findings.

See [specs/003-mock-backend-siem/quickstart.md](../specs/003-mock-backend-siem/quickstart.md) for the full guide.

## Quick start

```bash
# Deploy to AWS
make deploy

# Run locally (connects to real DynamoDB via AWS_PROFILE)
make local
```

Set the endpoint in your environment:

```bash
export AI_GUARDRAILS_SIEM_ENDPOINT=https://<api-id>.execute-api.<region>.amazonaws.com/
```

Then run `ai-check-guardrails` — each run posts findings automatically.

Open the endpoint URL in a browser to view submissions.

## Tear down

```bash
make destroy
```
