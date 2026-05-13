# Severity Levels

Every finding carries a severity level that indicates how urgently the issue should be addressed. The posture score is calculated from the severities of all findings in a run.

## Levels

| Level | Score Impact | Meaning |
|-------|-------------|---------|
| `INFO` | None | Informational observation. No action required, but worth being aware of. |
| `WARN` | Moderate | Advisory finding. An improvement is recommended but does not constitute a security gap on its own. |
| `HIGH` | Significant | A security gap that should be addressed before the next release or deployment. |
| `CRITICAL` | Maximum | A must-fix issue. Running in enforce mode with CRITICAL findings present will result in a non-zero exit code. |

## How the Score Works

The tool computes a posture score (0–100) from all findings after all modules run. Unresolved HIGH and CRITICAL findings lower the score. The score is compared against `score_threshold` (default: 70):

- **Score ≥ 100**: No banner displayed.
- **Score ≥ threshold**: A WARN banner is shown on stderr.
- **Score < threshold**: An ALERT banner is shown on stderr, optionally including a URL from `banner_url`.

The [gamification](../modules/gamification.md) module surfaces the target threshold as an INFO finding so it appears alongside other findings in the output.

## Enforce Mode

In `enforce` mode the tool exits with a non-zero exit code when any CRITICAL findings are present, making it suitable as a CI gate.
