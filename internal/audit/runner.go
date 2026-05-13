package audit

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/lock"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
	"github.com/kaihendry/ai-check-guardrails/internal/modules/banner"
	"github.com/kaihendry/ai-check-guardrails/internal/score"
	"github.com/kaihendry/ai-check-guardrails/internal/siem"
)

var Version = "dev"

type AuditRun struct {
	SchemaVersion string           `json:"schema_version"`
	RunID         string           `json:"run_id"`
	Timestamp     time.Time        `json:"timestamp"`
	Host          string           `json:"host"`
	User          string           `json:"user"`
	Mode          config.RunMode   `json:"mode"`
	Version       string           `json:"version"`
	Findings      []modules.Finding `json:"findings"`
	Score         int              `json:"score"`
	ExitCode      int              `json:"exit_code"`
	DurationMs    int64            `json:"duration_ms"`
}

func Run(cfg config.Config) (AuditRun, int) {
	start := time.Now()

	l, err := lock.Acquire()
	if err != nil {
		if errors.Is(err, lock.ErrAlreadyRunning) {
			fmt.Fprintln(os.Stderr, "ALREADY_RUNNING: another instance is running")
			return AuditRun{}, 2
		}
		fmt.Fprintf(os.Stderr, "lock error: %v\n", err)
		return AuditRun{}, 2
	}
	defer l.Release()

	host, _ := os.Hostname()
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}

	var findings []modules.Finding
	for _, m := range modules.All() {
		if !m.Enabled(cfg) {
			continue
		}
		ff, err := m.Run(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "module %s error: %v\n", m.Name(), err)
		}
		findings = append(findings, ff...)
	}

	s := score.Calculate(findings)
	banner.Display(s, cfg)

	exitCode := 0
	for _, f := range findings {
		if f.Severity != modules.SeverityInfo {
			exitCode = 1
			break
		}
	}

	if cfg.Mode == config.ModeEnforce {
		for _, f := range findings {
			if f.Severity == modules.SeverityCritical {
				fmt.Fprintf(os.Stderr, "ENFORCEMENT: critical finding in module %s: %s\n", f.Module, f.Description)
			}
		}
	}

	run := AuditRun{
		SchemaVersion: "1.0",
		RunID:         newUUID(),
		Timestamp:     time.Now().UTC(),
		Host:          host,
		User:          user,
		Mode:          cfg.Mode,
		Version:       Version,
		Findings:      findings,
		Score:         s,
		ExitCode:      exitCode,
		DurationMs:    time.Since(start).Milliseconds(),
	}

	token := os.Getenv("AI_GUARDRAILS_SIEM_TOKEN")
	if err := siem.Emit(run, cfg.SIEMEndpoint, token); err != nil {
		fmt.Fprintf(os.Stderr, "siem emit error: %v\n", err)
		return run, 2
	}
	if cfg.SIEMEndpoint != "" {
		fmt.Fprintf(os.Stderr, "[siem] findings posted → %s\n", cfg.SIEMEndpoint)
	}

	return run, exitCode
}
