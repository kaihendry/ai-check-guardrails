# hooks

Scans Git repositories under the configured scan root for required pre-commit hooks. Missing hooks leave repositories exposed to commits that bypass security checks such as secret scanning.

## Findings

| Type | Severity | Description | Remediation |
|------|----------|-------------|-------------|
| `MISSING_PRECOMMIT_HOOK` | HIGH | A required pre-commit hook is not installed in a repository. | Install the required hook in the repository's `.git/hooks/` directory. |

One finding is emitted per missing hook per repository. The module scans up to 3 directory levels below the scan root for Git repositories.

## Configuration

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `modules.hooks` | bool | `true` | Enable or disable this module |
| `allowlist.precommit_hooks` | `[]string` | `["gitleaks"]` | Names of hooks that must be installed in every repository |
| `scan_root` | string | `$HOME` | Root directory to search for Git repositories |

**Example config**:

```json
{
  "modules": { "hooks": true },
  "allowlist": {
    "precommit_hooks": ["gitleaks", "detect-secrets"]
  }
}
```
