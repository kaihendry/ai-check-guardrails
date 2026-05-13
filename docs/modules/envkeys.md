# envkeys

Detects API keys and credentials set as plaintext environment variables. Secrets loaded directly into the environment are visible to any process running in that session, increasing the risk of accidental exposure.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `PLAINTEXT_CREDENTIAL` | WARN | Credential variable is set but its value is short (fewer than 20 characters). | Remove the variable from the environment and load it from an approved secrets manager. |
| `PLAINTEXT_CREDENTIAL` | HIGH | Credential variable is set and its value is 20 or more characters (likely a real key). | Remove the variable from the environment and load it from an approved secrets manager. |

One finding is emitted per detected variable that is present and non-empty.

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.env_keys` | bool | `true` | Enable or disable this module |
| `env_key_watch_list` | `[]string` | `["ANTHROPIC_API_KEY"]` | List of environment variable names to watch. Overrides the default when non-empty. |

**Example config**:

```json
{
  "modules": { "env_keys": true },
  "env_key_watch_list": ["ANTHROPIC_API_KEY", "OPENAI_API_KEY", "MY_SECRET"]
}
```
