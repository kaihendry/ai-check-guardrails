# gamification

Emits an informational finding that anchors the security posture score in the output. This makes the score visible alongside other findings and reminds users of the target threshold, encouraging improvement over time.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `SCORE_INFO` | INFO | Security posture score will be calculated from findings in this run. Shows the target threshold. | Resolve HIGH and CRITICAL findings to improve your score. |

This module always emits exactly one finding per run.

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.gamification` | bool | `true` | Enable or disable this module |
| `score_threshold` | int | `70` | The target score shown in the finding's description |
