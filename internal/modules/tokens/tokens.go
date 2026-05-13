package tokens

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct {
	projectsDirOverride string // non-empty value used in tests to inject a custom dir
}

func (m *mod) Name() string                   { return "tokens" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Tokens }

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	projectsDir := m.projectsDirOverride
	if projectsDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		projectsDir = filepath.Join(home, ".claude", "projects")
	}

	lookback := 24
	if cfg.TokenBaseline != nil && cfg.TokenBaseline.LookbackHours > 0 {
		lookback = cfg.TokenBaseline.LookbackHours
	}

	stats, err := collectStats(projectsDir, lookback)
	if err != nil {
		return nil, err
	}

	byModel := make(map[string]any, len(stats.ByModel))
	for model, ms := range stats.ByModel {
		byModel[model] = map[string]any{
			"input_tokens":       ms.InputTokens,
			"output_tokens":      ms.OutputTokens,
			"cache_write_tokens": ms.CacheWrite,
			"cache_read_tokens":  ms.CacheRead,
			"est_cost_usd":       ms.EstCostUSD,
		}
	}

	byProject := make(map[string]any, len(stats.ByProject))
	for project, tokens := range stats.ByProject {
		byProject[project] = tokens
	}

	metadata := map[string]any{
		"lookback_hours":      lookback,
		"total_tokens":        stats.TotalTokens,
		"input_tokens":        stats.InputTokens,
		"output_tokens":       stats.OutputTokens,
		"cache_write_tokens":  stats.CacheWrite,
		"cache_read_tokens":   stats.CacheRead,
		"est_cost_usd":        stats.EstCostUSD,
		"web_searches":        stats.WebSearches,
		"web_fetches":         stats.WebFetches,
		"by_model":            byModel,
		"by_project":          byProject,
	}

	findings := []modules.Finding{
		{
			Type:     "TOKEN_USAGE",
			Severity: modules.SeverityInfo,
			Module:   m.Name(),
			Resource: "token-usage",
			Description: fmt.Sprintf(
				"Token usage in the last %dh: %d tokens (~$%.4f USD).",
				lookback, stats.TotalTokens, stats.EstCostUSD,
			),
			Metadata: metadata,
		},
	}

	if cfg.TokenBaseline != nil && stats.TotalTokens > 0 {
		b := cfg.TokenBaseline
		mult := b.Multiplier
		if mult == 0 {
			mult = 3.0
		}
		threshold := float64(b.DailyMean) + mult*b.StdDev
		if float64(stats.TotalTokens) > threshold {
			raw := (float64(stats.TotalTokens) - float64(b.DailyMean)) / b.StdDev
			confidence := math.Min(raw/10.0, 1.0)
			findings = append(findings, modules.Finding{
				Type: "TOKEN_ANOMALY",
				Severity: modules.SeverityWarn,
				Module:   m.Name(),
				Resource: "token-usage",
				Description: fmt.Sprintf(
					"Token usage (%d) exceeds anomaly threshold (%.0f).",
					stats.TotalTokens, threshold,
				),
				Remediation: "Review recent Claude sessions for unusual activity.",
				Confidence:  &confidence,
			})
		}
	}

	return findings, nil
}

// modelPrice holds per-million-token rates for a model.
type modelPrice struct {
	input      float64
	output     float64
	cacheWrite float64
	cacheRead  float64
}

// pricing is keyed by model name prefix (longest match wins).
// Rates are USD per million tokens as of 2026-05.
var pricing = []struct {
	prefix string
	price  modelPrice
}{
	{"claude-opus-4", modelPrice{input: 15.0, output: 75.0, cacheWrite: 18.75, cacheRead: 1.50}},
	{"claude-sonnet-4", modelPrice{input: 3.0, output: 15.0, cacheWrite: 3.75, cacheRead: 0.30}},
	{"claude-haiku-4", modelPrice{input: 0.80, output: 4.0, cacheWrite: 1.00, cacheRead: 0.08}},
	{"claude-opus-3", modelPrice{input: 15.0, output: 75.0, cacheWrite: 18.75, cacheRead: 1.50}},
	{"claude-sonnet-3", modelPrice{input: 3.0, output: 15.0, cacheWrite: 3.75, cacheRead: 0.30}},
	{"claude-haiku-3", modelPrice{input: 0.25, output: 1.25, cacheWrite: 0.30, cacheRead: 0.03}},
}

func priceFor(model string) (modelPrice, bool) {
	for _, p := range pricing {
		if strings.HasPrefix(model, p.prefix) {
			return p.price, true
		}
	}
	return modelPrice{}, false
}

func costUSD(p modelPrice, input, output, cacheWrite, cacheRead int64) float64 {
	const M = 1_000_000.0
	return float64(input)*p.input/M +
		float64(output)*p.output/M +
		float64(cacheWrite)*p.cacheWrite/M +
		float64(cacheRead)*p.cacheRead/M
}

type modelStats struct {
	InputTokens int64
	OutputTokens int64
	CacheWrite  int64
	CacheRead   int64
	EstCostUSD  float64
}

type usageStats struct {
	TotalTokens int64
	InputTokens int64
	OutputTokens int64
	CacheWrite  int64
	CacheRead   int64
	EstCostUSD  float64
	WebSearches int
	WebFetches  int
	ByModel     map[string]*modelStats
	ByProject   map[string]int64
}

type transcriptRecord struct {
	Timestamp string `json:"timestamp"`
	CWD       string `json:"cwd"`
	Message   struct {
		Model string `json:"model"`
		Usage struct {
			InputTokens              int `json:"input_tokens"`
			OutputTokens             int `json:"output_tokens"`
			CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
			CacheReadInputTokens     int `json:"cache_read_input_tokens"`
			ServerToolUse            struct {
				WebSearchRequests int `json:"web_search_requests"`
				WebFetchRequests  int `json:"web_fetch_requests"`
			} `json:"server_tool_use"`
		} `json:"usage"`
	} `json:"message"`
}

func collectStats(projectsDir string, lookbackHours int) (usageStats, error) {
	stats := usageStats{
		ByModel:   make(map[string]*modelStats),
		ByProject: make(map[string]int64),
	}

	projects, err := os.ReadDir(projectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return stats, nil
		}
		return stats, err
	}

	cutoff := time.Now().UTC().Add(-time.Duration(lookbackHours) * time.Hour)

	for _, project := range projects {
		if !project.IsDir() {
			continue
		}
		dir := filepath.Join(projectsDir, project.Name())
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
				continue
			}
			accumulateFile(filepath.Join(dir, entry.Name()), cutoff, &stats)
		}
	}
	return stats, nil
}

func accumulateFile(path string, cutoff time.Time, stats *usageStats) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		var rec transcriptRecord
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			continue
		}
		u := rec.Message.Usage
		if u.InputTokens == 0 && u.OutputTokens == 0 {
			continue
		}
		if rec.Timestamp != "" {
			t, err := time.Parse(time.RFC3339, rec.Timestamp)
			if err == nil && t.Before(cutoff) {
				continue
			}
		}

		in := int64(u.InputTokens)
		out := int64(u.OutputTokens)
		cw := int64(u.CacheCreationInputTokens)
		cr := int64(u.CacheReadInputTokens)
		total := in + out + cw + cr

		stats.TotalTokens += total
		stats.InputTokens += in
		stats.OutputTokens += out
		stats.CacheWrite += cw
		stats.CacheRead += cr
		stats.WebSearches += u.ServerToolUse.WebSearchRequests
		stats.WebFetches += u.ServerToolUse.WebFetchRequests

		model := rec.Message.Model
		if model == "" {
			model = "unknown"
		}
		ms, ok := stats.ByModel[model]
		if !ok {
			ms = &modelStats{}
			stats.ByModel[model] = ms
		}
		ms.InputTokens += in
		ms.OutputTokens += out
		ms.CacheWrite += cw
		ms.CacheRead += cr
		if p, ok := priceFor(model); ok {
			c := costUSD(p, in, out, cw, cr)
			ms.EstCostUSD += c
			stats.EstCostUSD += c
		}

		project := rec.CWD
		if project == "" {
			project = "unknown"
		}
		stats.ByProject[project] += total
	}
}

// sumTokensFromFile is kept for direct use in tests.
func sumTokensFromFile(path string, cutoff time.Time) int64 {
	var stats usageStats
	stats.ByModel = make(map[string]*modelStats)
	stats.ByProject = make(map[string]int64)
	accumulateFile(path, cutoff, &stats)
	return stats.TotalTokens
}
