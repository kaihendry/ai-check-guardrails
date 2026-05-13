# Data Model: Transcript Token Usage Reading

## Existing Config Entities (modified)

### TokenBaseline (config.go)

```
TokenBaseline
├── DailyMean    int     (existing) — mean daily tokens for baseline
├── StdDev       float64 (existing) — standard deviation
├── Multiplier   float64 (existing) — anomaly multiplier (default 3.0 when 0)
└── LookbackHours int    (NEW, optional) — hours to look back; 0 → default 24
```

**Validation**: `LookbackHours` must be ≥ 0. No upper bound enforced.

---

## Internal Types (tokens package — unexported)

### transcriptRecord

Represents one parsed line from a `.jsonl` transcript file.

```
transcriptRecord
├── Timestamp string     — top-level RFC3339 timestamp (may be empty)
├── Type      string     — record type ("user", "assistant", etc.)
└── Message   struct
    └── Usage struct
        ├── InputTokens              int
        ├── OutputTokens             int
        ├── CacheCreationInputTokens int
        └── CacheReadInputTokens     int
```

Only records where `Message.Usage` is non-zero contribute to the token sum. Records without a `Timestamp` are included (conservative inclusion).

---

## State Transitions

```
tokens.Run(cfg)
  │
  ├─ cfg.TokenBaseline == nil
  │   └─► return MODULE_UNAVAILABLE (existing behaviour)
  │
  └─ cfg.TokenBaseline != nil
      │
      ├─ readDailyTokens(projectsDir, lookbackHours)
      │   ├─ projectsDir not found → return 0, nil
      │   ├─ per file unreadable  → skip, continue
      │   └─ per line malformed   → skip, continue
      │
      ├─ daily == 0 → return nil, nil
      │
      ├─ daily ≤ threshold → return nil, nil
      │
      └─ daily > threshold → return TOKEN_ANOMALY (WARN)
```

---

## Key Fields Mapping (JSONL → Go struct)

| JSONL field path | Go field | Notes |
|---|---|---|
| `timestamp` | `transcriptRecord.Timestamp` | RFC3339 string |
| `type` | `transcriptRecord.Type` | Not filtered on; usage guard is sufficient |
| `message.usage.input_tokens` | `Usage.InputTokens` | Summed |
| `message.usage.output_tokens` | `Usage.OutputTokens` | Summed |
| `message.usage.cache_creation_input_tokens` | `Usage.CacheCreationInputTokens` | Summed |
| `message.usage.cache_read_input_tokens` | `Usage.CacheReadInputTokens` | Summed |
