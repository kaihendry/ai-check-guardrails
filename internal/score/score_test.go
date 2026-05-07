package score

import (
	"testing"

	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		findings []modules.Finding
		want     int
	}{
		{"no findings", nil, 100},
		{"one INFO", []modules.Finding{{Severity: modules.SeverityInfo}}, 100},
		{"one WARN", []modules.Finding{{Severity: modules.SeverityWarn}}, 95},
		{"one HIGH", []modules.Finding{{Severity: modules.SeverityHigh}}, 85},
		{"one CRITICAL", []modules.Finding{{Severity: modules.SeverityCritical}}, 75},
		{"two CRITICAL", []modules.Finding{
			{Severity: modules.SeverityCritical},
			{Severity: modules.SeverityCritical},
		}, 50},
		{"floor at zero", []modules.Finding{
			{Severity: modules.SeverityCritical},
			{Severity: modules.SeverityCritical},
			{Severity: modules.SeverityCritical},
			{Severity: modules.SeverityCritical},
			{Severity: modules.SeverityCritical},
		}, 0},
		{"mixed severities", []modules.Finding{
			{Severity: modules.SeverityCritical},
			{Severity: modules.SeverityHigh},
			{Severity: modules.SeverityWarn},
			{Severity: modules.SeverityInfo},
		}, 55},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Calculate(tt.findings); got != tt.want {
				t.Errorf("Calculate() = %d, want %d", got, tt.want)
			}
		})
	}
}
