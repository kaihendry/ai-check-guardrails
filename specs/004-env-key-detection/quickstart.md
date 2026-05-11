# Quickstart: Environment API Key Detection

## What this adds

A new `envkeys` detection module that checks whether `ANTHROPIC_API_KEY` (or
other configured variables) are present as plaintext environment variables and
emits `PLAINTEXT_CREDENTIAL` findings.

## Testing the module manually

```bash
# 1. Build
go build ./cmd/ai-check-guardrails/

# 2. Trigger a HIGH finding
ANTHROPIC_API_KEY=sk-ant-test123456789012345 ./ai-check-guardrails

# 3. Expect finding in output:
#    "type":"PLAINTEXT_CREDENTIAL","severity":"HIGH","resource":"ANTHROPIC_API_KEY"

# 4. Clean run — no finding
unset ANTHROPIC_API_KEY
./ai-check-guardrails
```

## Config options

```json
{
  "modules": {
    "env_keys": true
  },
  "env_key_watch_list": ["ANTHROPIC_API_KEY", "ANTHROPIC_ORG_TOKEN"]
}
```

Set `env_keys: false` to disable the module entirely.

## Score impact

Each `HIGH` finding reduces the posture score by 15 points (handled by existing
`internal/score` package — no change needed).
