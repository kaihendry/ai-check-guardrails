package bypass

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "bypass" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Bypass }

var bypassFlags = []string{
	"--dangerously-skip-permissions",
	"--dangerously_skip_permissions",
}

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	var findings []modules.Finding
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	historyFiles := []string{
		filepath.Join(home, ".zsh_history"),
		filepath.Join(home, ".bash_history"),
		filepath.Join(home, ".history"),
	}

	for _, hf := range historyFiles {
		f, err := os.Open(hf)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			for _, flag := range bypassFlags {
				if strings.Contains(line, flag) {
					findings = append(findings, modules.Finding{
						Type:        "POLICY_BYPASS",
						Severity:    modules.SeverityCritical,
						Module:      m.Name(),
						Resource:    hf,
						Description: "Policy bypass flag '" + flag + "' found in shell history.",
						Remediation: "Remove use of permission-bypass flags; request a security exception if needed.",
					})
					break
				}
			}
		}
		f.Close()
	}
	return findings, nil
}
