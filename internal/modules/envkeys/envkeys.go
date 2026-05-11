package envkeys

import (
	"os"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "envkeys" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.EnvKeys }

var defaultWatchList = []string{"ANTHROPIC_API_KEY"}

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	watchList := cfg.EnvKeyWatchList
	if len(watchList) == 0 {
		watchList = defaultWatchList
	}

	var findings []modules.Finding
	for _, name := range watchList {
		val, present := os.LookupEnv(name)
		if !present || val == "" {
			continue
		}
		severity := modules.SeverityWarn
		if len(val) >= 20 {
			severity = modules.SeverityHigh
		}
		findings = append(findings, modules.Finding{
			Type:        "PLAINTEXT_CREDENTIAL",
			Severity:    severity,
			Module:      m.Name(),
			Resource:    name,
			Description: "Credential variable " + name + " is set as a plaintext environment variable.",
			Remediation: "Remove " + name + " from the environment and load it from an approved secrets manager instead.",
		})
	}
	return findings, nil
}
