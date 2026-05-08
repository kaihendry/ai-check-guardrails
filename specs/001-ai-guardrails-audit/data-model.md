# Data Model: AI Guardrails Audit

**Branch**: `001-ai-guardrails-audit` | **Date**: 2026-05-07

All types are Go structs serialised to JSON. No external database; state lives
in the config file and in the emitted SIEM events.

---

## Finding

The atomic unit of audit output. One finding per detected issue.

```go
type Finding struct {
    Type        FindingType `json:"type"`
    Severity    Severity    `json:"severity"`
    Module      string      `json:"module"`       // module Name()
    Resource    string      `json:"resource"`     // affected file, MCP ID, etc.
    Description string      `json:"description"`
    Remediation string      `json:"remediation"`
    Confidence  *float64    `json:"confidence,omitempty"` // 0.0–1.0; nil = certain
}
```

### FindingType (string enum)

| Value | Source module | Severity default |
|-------|--------------|-----------------|
| `MISSING_CONFIG` | settings | CRITICAL |
| `SETTINGS_OVERRIDE` | settings | WARN |
| `UNAPPROVED_MCP` | mcp | HIGH |
| `UNCONFIGURED_ALLOWLIST` | mcp | WARN |
| `PERMISSION_MISCONFIGURED` | permissions | HIGH |
| `TOKEN_ANOMALY` | tokens | WARN |
| `NETWORK_REQUEST` | network | INFO |
| `EVAL_HOOK_ABSENT` | evals | HIGH |
| `SANDBOX_VIOLATION` | sandbox | HIGH |
| `HUMANLOOP_ABSENT` | humanloop | WARN |
| `POLICY_BYPASS` | bypass | CRITICAL |
| `MISSING_PRECOMMIT_HOOK` | hooks | HIGH |
| `MODULE_UNAVAILABLE` | runner | INFO |
| `ALREADY_RUNNING` | lock | INFO |

### Severity (string enum)

`INFO` | `WARN` | `HIGH` | `CRITICAL`

---

## AuditRun

Top-level record for a single execution. This is the root of every SIEM event.

```go
type AuditRun struct {
    RunID       string    `json:"run_id"`       // UUID v4
    Timestamp   time.Time `json:"timestamp"`    // UTC RFC3339
    Host        string    `json:"host"`         // os.Hostname()
    User        string    `json:"user"`         // os.Getenv("USER")
    Mode        RunMode   `json:"mode"`         // "monitor" | "enforce"
    Version     string    `json:"version"`      // binary semver
    Findings    []Finding `json:"findings"`
    Score       int       `json:"score"`        // 0–100
    ExitCode    int       `json:"exit_code"`    // 0 | 1 | 2
    DurationMs  int64     `json:"duration_ms"`
}
```

### RunMode (string enum)

`monitor` | `enforce`

---

## PostureScore

Derived from `AuditRun.Findings`. Not a stored entity — calculated at run end.

**Scoring algorithm**:
- Start at 100.
- Deduct per finding by severity: CRITICAL −25, HIGH −15, WARN −5, INFO 0.
- Floor at 0 (cannot go negative).
- Store result in `AuditRun.Score`.

| Threshold | Banner behaviour (monitor mode) |
|-----------|--------------------------------|
| 100 | No banner |
| 70–99 | Advisory banner with improvement tips |
| < 70 | Warning banner with link to security training |

---

## Config

Tool configuration read from `~/.config/ai-check-guardrails/config.json`.

```go
type Config struct {
    Mode           RunMode        `json:"mode"`             // default: "monitor"
    SIEMEndpoint   string         `json:"siem_endpoint"`    // optional HTTPS URL
    ScanRoot       string         `json:"scan_root"`        // default: $HOME
    Modules        ModuleToggles  `json:"modules"`
    Allowlist      Allowlist      `json:"allowlist"`
    TokenBaseline  *TokenBaseline `json:"token_baseline,omitempty"`
    BannerURL      string         `json:"banner_url"`       // training link
    ScoreThreshold int            `json:"score_threshold"`  // default: 70
}

type ModuleToggles struct {
    Settings    bool `json:"settings"`     // default: true
    MCP         bool `json:"mcp"`          // default: true
    Permissions bool `json:"permissions"`  // default: true
    Tokens      bool `json:"tokens"`       // default: false (needs baseline)
    Network     bool `json:"network"`      // default: false (passive only)
    Evals       bool `json:"evals"`        // default: false
    Sandbox     bool `json:"sandbox"`      // default: true
    HumanLoop   bool `json:"humanloop"`    // default: false
    Bypass      bool `json:"bypass"`       // default: true
    Banner      bool `json:"banner"`       // default: true
    Hooks       bool `json:"hooks"`        // default: true
    Gamification bool `json:"gamification"` // default: true
}

type Allowlist struct {
    MCPs        []string `json:"mcps"`         // approved MCP identifiers
    Skills      []string `json:"skills"`       // approved skill names
    Domains     []string `json:"domains"`      // approved outbound domains
    PreCommitHooks []string `json:"precommit_hooks"` // required hook names
}

type TokenBaseline struct {
    DailyMean   int     `json:"daily_mean"`
    StdDev      float64 `json:"std_dev"`
    Multiplier  float64 `json:"multiplier"`  // anomaly = mean + multiplier*stddev
}
```

---

## SIEMEvent

The JSON record forwarded to the SIEM endpoint (HTTP POST body) and printed to
stdout. This is `AuditRun` serialised to JSON — no additional wrapper.

Validation rules:
- `run_id` MUST be a valid UUID v4.
- `timestamp` MUST be UTC RFC3339.
- `findings` MAY be empty (clean run).
- `score` MUST be in range [0, 100].
- `exit_code` MUST be 0, 1, or 2.

---

## State Transitions

### RunMode progression (phased rollout)

```
monitor ──(operator changes config)──► enforce
```

There is no automatic promotion. The operator edits `config.json` to change mode.

### Scheduled execution lifecycle

```
launchd/systemd fires
    │
    ▼
lock acquired? ──No──► exit 2 (ALREADY_RUNNING to stderr)
    │ Yes
    ▼
load config → validate
    │
    ▼
run enabled modules → collect []Finding
    │
    ▼
calculate score
    │
    ▼
emit SIEM event (stdout + optional HTTP POST)
    │
    ▼
display banner if score < threshold
    │
    ▼
enforce mode: block/alert on CRITICAL findings?
    │
    ▼
release lock → exit (0 clean, 1 findings, 2 error)
```
