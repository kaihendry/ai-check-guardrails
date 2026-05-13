# banner

Displays a security posture summary banner on stderr after all findings are collected. The banner shows the current score and guides users toward remediation or alerts them when the score falls below the configured threshold.

## Findings

This module emits no findings itself — it operates on the final score after all other modules run.

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| *(none)* | — | Banner output goes to stderr; findings come from other modules. | Resolve HIGH and CRITICAL findings from other modules to raise your score. |

## Behavior

| Score | Outcome |
|-------|---------|
| 100 | No banner displayed |
| ≥ threshold | WARN banner on stderr showing current score |
| < threshold | ALERT banner on stderr with score, threshold, and action URL |

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.banner` | bool | `true` | Enable or disable this module |
| `score_threshold` | int | `70` | Minimum passing score; scores below this trigger the ALERT banner |
| `banner_url` | string | *(empty)* | URL shown in the ALERT banner (e.g., your security training portal) |
