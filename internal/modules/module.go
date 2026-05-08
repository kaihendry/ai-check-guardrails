package modules

import "github.com/kaihendry/ai-check-guardrails/internal/config"

type FindingType string
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarn     Severity = "WARN"
	SeverityHigh     Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type Finding struct {
	Type        FindingType `json:"type"`
	Severity    Severity    `json:"severity"`
	Module      string      `json:"module"`
	Resource    string      `json:"resource"`
	Description string      `json:"description"`
	Remediation string      `json:"remediation"`
	Confidence  *float64    `json:"confidence,omitempty"`
}

type Module interface {
	Name() string
	Enabled(cfg config.Config) bool
	Run(cfg config.Config) ([]Finding, error)
}

var registry []Module

func Register(m Module) {
	registry = append(registry, m)
}

func All() []Module {
	return registry
}
