# network

Scans Claude network log files for outbound HTTP/HTTPS requests and compares observed domains against a configured allowlist. This helps detect unexpected data exfiltration or connections to unapproved services.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `NETWORK_REQUEST` | INFO | An outbound request was observed; no allowlist is configured so all domains are informational. | Verify this domain is on the approved list. Consider configuring `allowlist.domains`. |
| `NETWORK_REQUEST` | WARN | An outbound request was made to a domain not on the configured allowlist. | Verify the domain is intentional and add it to `allowlist.domains` or investigate the request. |

One finding is emitted per unique domain observed. Domains are extracted from these log files:

- `~/.claude/logs/network.log`
- `~/.claude/network.log`

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.network` | bool | `false` | Enable or disable this module |
| `allowlist.domains` | `[]string` | *(empty)* | Approved outbound domains. When empty, all observed domains emit INFO findings. When set, unlisted domains emit WARN. |

**Example config**:

```json
{
  "modules": { "network": true },
  "allowlist": {
    "domains": ["api.anthropic.com", "github.com"]
  }
}
```
