# tokens

Detects anomalous token usage by comparing the estimated daily token count against a configured statistical baseline. Sudden spikes in token usage can indicate runaway agents, prompt injection attacks, or exfiltration attempts.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `MODULE_UNAVAILABLE` | INFO | No token baseline is configured; anomaly detection is skipped. | Add `token_baseline` to your config after establishing a normal usage baseline. |
| `TOKEN_ANOMALY` | WARN | Daily token usage exceeds the anomaly threshold (mean + multiplier × std_dev). Includes a `confidence` field (0–1). | Review recent Claude sessions for unusual activity. |

## Current Status

**Token reading is not yet implemented.** The internal function that would read Claude's usage data (`~/.claude/usage.json` or similar) is a placeholder that always returns `0`. As a result, the `TOKEN_ANOMALY` finding is never emitted in the current release — only `MODULE_UNAVAILABLE` is reachable.

Configuring `token_baseline` will not cause harm, but no anomaly detection will occur until the usage-reading function is implemented.

## Anomaly Threshold (future)

Once token reading is implemented, the threshold will be calculated as:

```
threshold = daily_mean + multiplier × std_dev
```

A `TOKEN_ANOMALY` finding will be emitted when the estimated daily usage exceeds this value. The `confidence` field reflects how many standard deviations above the mean the usage is, normalised to the range `[0, 1]`.

## Configuration

| Key | Type | Default | Required | Description |
|-----|------|---------|----------|-------------|
| `modules.tokens` | bool | `false` | — | Enable or disable this module |
| `token_baseline` | object | `null` | Required when module is enabled | Baseline statistics for anomaly detection |
| `token_baseline.daily_mean` | int | — | Yes | Expected average daily token count |
| `token_baseline.std_dev` | float | — | Yes | Standard deviation of daily token count |
| `token_baseline.multiplier` | float | `3.0` | No | Anomaly threshold multiplier (higher = less sensitive) |

**Note**: The tool returns a validation error if `modules.tokens` is `true` and `token_baseline` is absent.

**Example config**:

```json
{
  "modules": { "tokens": true },
  "token_baseline": {
    "daily_mean": 50000,
    "std_dev": 10000,
    "multiplier": 3.0
  }
}
```
