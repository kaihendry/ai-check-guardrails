package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RunMode string

const (
	ModeMonitor RunMode = "monitor"
	ModeEnforce RunMode = "enforce"
)

type ModuleToggles struct {
	Settings     bool `json:"settings"`
	MCP          bool `json:"mcp"`
	Permissions  bool `json:"permissions"`
	Tokens       bool `json:"tokens"`
	Network      bool `json:"network"`
	Evals        bool `json:"evals"`
	Sandbox      bool `json:"sandbox"`
	HumanLoop    bool `json:"humanloop"`
	Bypass       bool `json:"bypass"`
	Banner       bool `json:"banner"`
	Hooks        bool `json:"hooks"`
	Gamification bool `json:"gamification"`
}

type Allowlist struct {
	MCPs           []string `json:"mcps"`
	Skills         []string `json:"skills"`
	Domains        []string `json:"domains"`
	PreCommitHooks []string `json:"precommit_hooks"`
}

type TokenBaseline struct {
	DailyMean  int     `json:"daily_mean"`
	StdDev     float64 `json:"std_dev"`
	Multiplier float64 `json:"multiplier"`
}

type Config struct {
	Mode           RunMode        `json:"mode"`
	SIEMEndpoint   string         `json:"siem_endpoint"`
	ScanRoot       string         `json:"scan_root"`
	BannerURL      string         `json:"banner_url"`
	ScoreThreshold int            `json:"score_threshold"`
	Modules        ModuleToggles  `json:"modules"`
	Allowlist      Allowlist      `json:"allowlist"`
	TokenBaseline  *TokenBaseline `json:"token_baseline,omitempty"`
}

func defaults() Config {
	home, _ := os.UserHomeDir()
	return Config{
		Mode:           ModeMonitor,
		ScanRoot:       home,
		ScoreThreshold: 70,
		Modules: ModuleToggles{
			Settings:     true,
			MCP:          true,
			Permissions:  true,
			Tokens:       false,
			Network:      false,
			Evals:        false,
			Sandbox:      true,
			HumanLoop:    false,
			Bypass:       true,
			Banner:       true,
			Hooks:        true,
			Gamification: true,
		},
		Allowlist: Allowlist{
			PreCommitHooks: []string{"gitleaks"},
		},
	}
}

func Load(path string) (Config, error) {
	cfg := defaults()
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return cfg, err
		}
		path = filepath.Join(home, ".config", "ai-check-guardrails", "config.json")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("reading config %s: %w", path, err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing config %s: %w", path, err)
	}
	return cfg, validate(cfg)
}

func validate(cfg Config) error {
	if cfg.Mode != ModeMonitor && cfg.Mode != ModeEnforce {
		return fmt.Errorf("invalid mode %q: must be monitor or enforce", cfg.Mode)
	}
	if cfg.SIEMEndpoint != "" && !strings.HasPrefix(cfg.SIEMEndpoint, "https://") {
		return fmt.Errorf("siem_endpoint must use HTTPS")
	}
	if cfg.ScanRoot != "" && !filepath.IsAbs(cfg.ScanRoot) {
		return fmt.Errorf("scan_root must be an absolute path")
	}
	if cfg.Modules.Tokens && cfg.TokenBaseline == nil {
		return fmt.Errorf("token_baseline required when tokens module is enabled")
	}
	return nil
}
