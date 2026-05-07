package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"text/template"

	"github.com/kaihendry/ai-check-guardrails/internal/audit"
	"github.com/kaihendry/ai-check-guardrails/internal/config"

	// Register all detection modules.
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/banner"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/bypass"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/evals"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/gamification"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/hooks"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/humanloop"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/mcp"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/network"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/permissions"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/sandbox"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/settings"
	_ "github.com/kaihendry/ai-check-guardrails/internal/modules/tokens"
)

//go:embed templates/launchd.plist.tmpl
var launchdTmpl string

//go:embed templates/systemd.service.tmpl
var systemdServiceTmpl string

//go:embed templates/systemd.timer.tmpl
var systemdTimerTmpl string

func main() {
	var (
		cfgPath        = flag.String("config", "", "path to config JSON")
		modeOverride   = flag.String("mode", "", "override run mode: monitor or enforce")
		installLaunchd = flag.Bool("install-launchd", false, "install macOS launchd schedule")
		installSystemd = flag.Bool("install-systemd", false, "install Linux systemd schedule")
		uninstall      = flag.Bool("uninstall", false, "remove installed schedule")
		showVersion    = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(audit.Version)
		os.Exit(0)
	}

	if *installLaunchd && *installSystemd {
		fmt.Fprintln(os.Stderr, "error: --install-launchd and --install-systemd are mutually exclusive")
		os.Exit(2)
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(2)
	}
	if *modeOverride != "" {
		cfg.Mode = config.RunMode(*modeOverride)
	}

	if *installLaunchd {
		os.Exit(installLaunchdSchedule())
	}
	if *installSystemd {
		os.Exit(installSystemdSchedule())
	}
	if *uninstall {
		os.Exit(uninstallSchedule())
	}

	_, exitCode := audit.Run(cfg)
	os.Exit(exitCode)
}

type scheduleData struct {
	BinaryPath string
	ConfigPath string
	Username   string
}

func resolveScheduleData() scheduleData {
	binary, _ := os.Executable()
	home, _ := os.UserHomeDir()
	cfgPath := filepath.Join(home, ".config", "ai-check-guardrails", "config.json")
	u, _ := user.Current()
	return scheduleData{BinaryPath: binary, ConfigPath: cfgPath, Username: u.Username}
}

func installLaunchdSchedule() int {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "install-launchd: %v\n", err)
		return 2
	}
	data := resolveScheduleData()
	dest := filepath.Join(home, "Library", "LaunchAgents", "com.example.ai-check-guardrails.plist")
	if err := renderTemplate(launchdTmpl, dest, data); err != nil {
		fmt.Fprintf(os.Stderr, "install-launchd: %v\n", err)
		return 2
	}
	fmt.Printf("launchd plist written to %s\nRun: launchctl load %s\n", dest, dest)
	return 0
}

func installSystemdSchedule() int {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "install-systemd: %v\n", err)
		return 2
	}
	data := resolveScheduleData()
	unitDir := filepath.Join(home, ".config", "systemd", "user")
	if err := os.MkdirAll(unitDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "install-systemd: %v\n", err)
		return 2
	}
	pairs := []struct{ tmpl, name string }{
		{systemdServiceTmpl, "ai-check-guardrails.service"},
		{systemdTimerTmpl, "ai-check-guardrails.timer"},
	}
	for _, p := range pairs {
		if err := renderTemplate(p.tmpl, filepath.Join(unitDir, p.name), data); err != nil {
			fmt.Fprintf(os.Stderr, "install-systemd %s: %v\n", p.name, err)
			return 2
		}
	}
	fmt.Printf("systemd units written to %s\nRun:\n  systemctl --user daemon-reload\n  systemctl --user enable --now ai-check-guardrails.timer\n", unitDir)
	return 0
}

func uninstallSchedule() int {
	home, _ := os.UserHomeDir()
	unitDir := filepath.Join(home, ".config", "systemd", "user")
	paths := []string{
		filepath.Join(home, "Library", "LaunchAgents", "com.example.ai-check-guardrails.plist"),
		filepath.Join(unitDir, "ai-check-guardrails.service"),
		filepath.Join(unitDir, "ai-check-guardrails.timer"),
	}
	removed := false
	for _, p := range paths {
		if err := os.Remove(p); err == nil {
			fmt.Printf("removed %s\n", p)
			removed = true
		}
	}
	if !removed {
		fmt.Println("no schedule files found")
	}
	return 0
}

func renderTemplate(tmplStr, dest string, data any) error {
	t, err := template.New("").Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating %s: %w", dest, err)
	}
	defer f.Close()
	return t.Execute(f, data)
}
