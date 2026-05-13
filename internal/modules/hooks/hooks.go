// # hooks
//
// Scans Git repositories under the configured scan root for required pre-commit hooks.
// Missing hooks leave repositories exposed to commits that bypass security checks such
// as secret scanning.
//
// # Findings
//
// | Type | Severity | Description | Remediation |
// |------|----------|-------------|-------------|
// | MISSING_PRECOMMIT_HOOK | HIGH | A required pre-commit hook is not installed in a repository. | Install the required hook in the repository's .git/hooks/ directory. |
//
// One finding is emitted per missing hook per repository. The module scans up to 3
// directory levels below the scan root for Git repositories.
//
// # Configuration
//
// | Key | Type | Default | Description |
// |-----|------|---------|-------------|
// | modules.hooks | bool | true | Enable or disable this module |
// | allowlist.precommit_hooks | []string | ["gitleaks"] | Names of hooks that must be installed in every repository |
// | scan_root | string | $HOME | Root directory to search for Git repositories |
//
// Example config:
//
//	{
//	  "modules": { "hooks": true },
//	  "allowlist": {
//	    "precommit_hooks": ["gitleaks", "detect-secrets"]
//	  }
//	}
package hooks

import (
	"os"
	"path/filepath"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "hooks" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Hooks }

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	var findings []modules.Finding
	scanRoot := cfg.ScanRoot
	if scanRoot == "" {
		var err error
		scanRoot, err = os.UserHomeDir()
		if err != nil {
			return nil, err
		}
	}

	required := cfg.Allowlist.PreCommitHooks
	if len(required) == 0 {
		required = []string{"gitleaks"}
	}

	// Walk up to 3 directory levels to find .git dirs.
	repos := findGitRepos(scanRoot, 3)
	for _, repo := range repos {
		hooksDir := filepath.Join(repo, ".git", "hooks")
		for _, hookName := range required {
			hookPath := filepath.Join(hooksDir, hookName)
			if _, err := os.Stat(hookPath); os.IsNotExist(err) {
				findings = append(findings, modules.Finding{
					Type:        "MISSING_PRECOMMIT_HOOK",
					Severity:    modules.SeverityHigh,
					Module:      m.Name(),
					Resource:    repo,
					Description: "Required pre-commit hook '" + hookName + "' not found in repository.",
					Remediation: "Install " + hookName + " in " + hooksDir,
				})
			}
		}
	}
	return findings, nil
}

func findGitRepos(root string, maxDepth int) []string {
	var repos []string
	var walk func(path string, depth int)
	walk = func(path string, depth int) {
		if depth < 0 {
			return
		}
		gitDir := filepath.Join(path, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			repos = append(repos, path)
			return
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return
		}
		for _, e := range entries {
			if e.IsDir() && e.Name() != ".git" {
				walk(filepath.Join(path, e.Name()), depth-1)
			}
		}
	}
	walk(root, maxDepth)
	return repos
}
