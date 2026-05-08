# Quickstart: AI Guardrails Audit

## Prerequisites

- Go 1.22+ installed (build only; the produced binary has no runtime dependencies)
- macOS 13+ or Linux with systemd

## 1. Build

```bash
git clone <repo>
cd ai-check-guardrails
go build -o ai-check-guardrails ./cmd/ai-check-guardrails/
```

## 2. Create minimal config

```bash
mkdir -p ~/.config/ai-check-guardrails
cat > ~/.config/ai-check-guardrails/config.json <<'EOF'
{
  "mode": "monitor",
  "allowlist": {
    "mcps": [],
    "precommit_hooks": ["gitleaks"]
  },
  "banner_url": "https://your-intranet/security-training"
}
EOF
```

> **Note**: An empty `mcps` list means all MCPs will be flagged as WARN until you
> populate the approved list. This is intentional for the initial baseline phase.

## 3. Run manually

```bash
./ai-check-guardrails
```

Expected output (clean workstation, monitor mode, empty allowlist):

```json
{"schema_version":"1.0","run_id":"...","timestamp":"...","host":"...","user":"...","mode":"monitor","version":"0.1.0","findings":[{"type":"UNCONFIGURED_ALLOWLIST","severity":"WARN",...}],"score":95,"exit_code":1,"duration_ms":3200}
```

Exit code 1 indicates findings are present (the unconfigured allowlist warning).

## 4. Schedule (macOS launchd)

```bash
./ai-check-guardrails --install-launchd
```

This writes `~/Library/LaunchAgents/com.example.ai-check-guardrails.plist` and
loads it. The tool will run every 30 minutes by default.

## 5. Schedule (Linux systemd)

```bash
./ai-check-guardrails --install-systemd
systemctl --user daemon-reload
systemctl --user enable --now ai-check-guardrails.timer
```

## 6. Verify SIEM delivery (optional)

Add the endpoint to config:

```json
{
  "siem_endpoint": "https://siem.example.com/api/v1/events"
}
```

Set the auth token:

```bash
export AI_GUARDRAILS_SIEM_TOKEN=your-bearer-token
./ai-check-guardrails
```

The tool will POST the JSON event to the endpoint in addition to stdout.

## 7. Populate the MCP allowlist

After running in monitor mode for a few days, review the `UNAPPROVED_MCP` findings
in your SIEM to identify which MCPs are in active use. Add approved ones to config:

```json
{
  "allowlist": {
    "mcps": ["mcp://filesystem", "mcp://brave-search"],
    "skills": ["speckit-plan", "speckit-specify"]
  }
}
```

## 8. Enable enforcement (when ready)

Change mode in config:

```json
{ "mode": "enforce" }
```

In enforce mode, `CRITICAL` findings cause the tool to exit with code 1 and print
a blocking banner. Non-CRITICAL findings are still logged but do not block.

## Phased Rollout Guide

| Phase | Config changes | Expected outcome |
|-------|---------------|-----------------|
| 1 — Baseline | `mode: monitor`, all default modules | Visibility; no disruption |
| 2 — MCP control | Populate `allowlist.mcps` | Unapproved MCPs flagged HIGH |
| 3 — Hooks | `allowlist.precommit_hooks` populated | Missing hooks flagged HIGH |
| 4 — Bypass detection | `bypass: true` (default) | `--dangerously-skip-permissions` flagged CRITICAL |
| 5 — Enforcement | `mode: enforce` | CRITICAL findings produce blocking exit |
| 6 — Token anomalies | Enable `tokens`, configure baseline | Anomalous usage flagged |

## Uninstall

```bash
./ai-check-guardrails --uninstall
rm ~/.config/ai-check-guardrails/config.json
rm ./ai-check-guardrails
```
