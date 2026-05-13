# Research: Transcript Token Usage Reading

## Transcript File Format

**Decision**: Parse `~/.claude/projects/**/*.jsonl` line-by-line using `bufio.Scanner`.

**Rationale**: Each file is a newline-delimited JSON stream. Reading line-by-line avoids loading the entire file into memory. The Claude Code session-report plugin (reference implementation) confirmed this format: each line is a JSON object; assistant messages contain a top-level `timestamp` (RFC3339) and a `message.usage` object with `input_tokens`, `output_tokens`, `cache_creation_input_tokens`, `cache_read_input_tokens`.

**Alternatives considered**:
- `json.Decoder` stream: equivalent but slightly more complex for line-skip-on-error; rejected in favour of `bufio.Scanner` + per-line `json.Unmarshal` for simpler error isolation.
- Read entire file then split: uses more memory for large files; rejected.

---

## Token Fields to Sum

**Decision**: Sum all four fields — `input_tokens`, `output_tokens`, `cache_creation_input_tokens`, `cache_read_input_tokens` — from `message.usage` in assistant-type records.

**Rationale**: The baseline `daily_mean` in the session-report reference represents total billable token consumption. Cache creation tokens are billable (at a higher rate) and cache read tokens represent consumption too. Including all four gives the most complete picture for anomaly detection.

**Alternatives considered**:
- Sum only `input_tokens + output_tokens`: excludes cache tokens, underestimates real usage; rejected.
- Use only `output_tokens`: too narrow; rejected.

---

## Lookback Window Configuration

**Decision**: Add optional `lookback_hours` int field to `TokenBaseline` config struct (default 24 when zero).

**Rationale**: Placing it inside `TokenBaseline` keeps token-related config cohesive. Using 0-value-as-default means no breaking config change for existing users.

**Alternatives considered**:
- New top-level config field: more visible but scatters related config; rejected.
- Hard-code 24 h: can't fulfil User Story 3; rejected.

---

## Identifying Records Within Lookback Window

**Decision**: Use the top-level `timestamp` field (RFC3339 string) on each JSONL record for time filtering. Records without a timestamp field are included (conservative — don't silently drop).

**Rationale**: The `timestamp` is present on all message-type records in the observed format. Using it is simpler than parsing message IDs or inferring time from file modification time.

**Alternatives considered**:
- File `mtime`: would include files touched outside the window; rejected.
- Session start time from first record: harder to parse; rejected.

---

## Directory Walk Strategy

**Decision**: Use `os.ReadDir` for the top-level projects directory, then `filepath.Glob` for `*.jsonl` within each project subdirectory.

**Rationale**: The structure is exactly two levels deep (`~/.claude/projects/<project>/*.jsonl`). A two-level read is simpler and faster than `filepath.WalkDir` for this fixed depth.

**Alternatives considered**:
- `filepath.WalkDir`: handles arbitrary depth but adds complexity for a known-depth structure; rejected per Simplicity principle.

---

## Error Handling

**Decision**: Unreadable files → skip silently (no error returned). Malformed JSONL lines → skip that line, continue. Missing top-level projects directory → return 0, nil.

**Rationale**: The tool runs on machines where Claude Code may not be installed (CI, shared boxes). Constitution Integrity requires no silent failures that suppress real issues, but skipping absent/unreadable inputs is correct behaviour here — it's not an error state.
