// # network
//
// Scans Claude network log files for outbound HTTP/HTTPS requests and compares observed
// domains against a configured allowlist. This helps detect unexpected data exfiltration
// or connections to unapproved services.
//
// # Findings
//
// | Type | Severity | Description | Remediation |
// |------|----------|-------------|-------------|
// | NETWORK_REQUEST | INFO | An outbound request was observed; no allowlist is configured so all domains are informational. | Verify this domain is on the approved list. Consider configuring allowlist.domains. |
// | NETWORK_REQUEST | WARN | An outbound request was made to a domain not on the configured allowlist. | Verify the domain is intentional and add it to allowlist.domains or investigate the request. |
//
// One finding is emitted per unique domain observed. Domains are extracted from
// ~/.claude/logs/network.log and ~/.claude/network.log.
//
// # Configuration
//
// | Key | Type | Default | Description |
// |-----|------|---------|-------------|
// | modules.network | bool | false | Enable or disable this module |
// | allowlist.domains | []string | *(empty)* | Approved outbound domains. When empty, all observed domains emit INFO. When set, unlisted domains emit WARN. |
//
// Example config:
//
//	{
//	  "modules": { "network": true },
//	  "allowlist": {
//	    "domains": ["api.anthropic.com", "github.com"]
//	  }
//	}
package network

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kaihendry/ai-check-guardrails/internal/config"
	"github.com/kaihendry/ai-check-guardrails/internal/modules"
)

func init() { modules.Register(&mod{}) }

type mod struct{}

func (m *mod) Name() string                   { return "network" }
func (m *mod) Enabled(cfg config.Config) bool { return cfg.Modules.Network }

var urlPattern = regexp.MustCompile(`https?://([a-zA-Z0-9._-]+)`)

func (m *mod) Run(cfg config.Config) ([]modules.Finding, error) {
	var findings []modules.Finding
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	logPaths := []string{
		filepath.Join(home, ".claude", "logs", "network.log"),
		filepath.Join(home, ".claude", "network.log"),
	}

	seen := make(map[string]bool)
	for _, p := range logPaths {
		f, err := os.Open(p)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			matches := urlPattern.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				if len(m) > 1 {
					domain := strings.ToLower(m[1])
					if !seen[domain] {
						seen[domain] = true
					}
				}
			}
		}
		f.Close()
	}

	allowed := make(map[string]bool, len(cfg.Allowlist.Domains))
	for _, d := range cfg.Allowlist.Domains {
		allowed[strings.ToLower(d)] = true
	}

	for domain := range seen {
		finding := modules.Finding{
			Type:        "NETWORK_REQUEST",
			Severity:    modules.SeverityInfo,
			Module:      m.Name(),
			Resource:    domain,
			Description: "Outbound network request observed to domain: " + domain,
			Remediation: "Verify this domain is on the approved list.",
		}
		if len(cfg.Allowlist.Domains) > 0 && !allowed[domain] {
			finding.Severity = modules.SeverityWarn
			finding.Description = "Outbound request to non-allowlisted domain: " + domain
		}
		findings = append(findings, finding)
	}
	return findings, nil
}
