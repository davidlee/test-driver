---
id: SPEC-001
slug: im_cli
name: im CLI
created: '2026-03-04'
updated: '2026-03-04'
status: draft
kind: spec
aliases: [im-tech]
relations:
- type: implements
  target: PROD-001
guiding_principles:
  - Small packages with clear boundaries — logfile, entry, reader
  - Testability via dependency injection (clock function), not interfaces
  - Append-only file operations — never mutate prior content
  - Pure functions where possible (entry package has zero I/O)
assumptions:
  - POSIX environment (Linux / macOS)
  - Single-user local tool — no concurrency protection beyond O_APPEND
  - Daily log files are small (< 10 KB)
category: assembly
c4_level: component
---

# SPEC-001 – im CLI

```yaml supekku:spec.relationships@v1
schema: supekku.spec.relationships
version: 1
spec: SPEC-001
requirements:
  primary:
    - SPEC-001.FR-001
    - SPEC-001.FR-002
    - SPEC-001.FR-003
    - SPEC-001.FR-004
    - SPEC-001.FR-005
    - SPEC-001.FR-006
    - SPEC-001.FR-007
    - SPEC-001.FR-008
    - SPEC-001.FR-009
    - SPEC-001.NF-001
    - SPEC-001.NF-002
    - SPEC-001.NF-003
    - SPEC-001.NF-004
  collaborators:
    - PROD-001
interactions: []
```

```yaml supekku:spec.capabilities@v1
schema: supekku.spec.capabilities
version: 1
spec: SPEC-001
capabilities:
  - id: entry-formatting
    name: Entry Text Formatting
    responsibilities:
      - Join and trim input text
      - Apply task checkbox prefix when requested
      - Return empty string for empty/whitespace-only input
    requirements: [FR-001]
    summary: >-
      Pure text formatting with no I/O dependencies. Receives a body string
      and task flag, returns formatted text ready for file append. Task mode
      prefixes each non-empty line with "- [ ] ".
    success_criteria:
      - Format("hello", false) returns "hello\n"
      - Format("buy milk", true) returns "- [ ] buy milk\n"
      - Format("", false) returns ""

  - id: daily-logfile
    name: Daily Log File Management
    responsibilities:
      - Create daily Markdown file with date heading on first write
      - Append entries without mutating prior content
      - Parse last timestamp subheading from file content
      - Determine whether to emit a new timestamp subheading
      - Support adaptive and round10 timestamp strategies
    requirements: [FR-002, FR-003, FR-004, FR-005]
    summary: >-
      Manages one Markdown file per day under a configurable directory.
      Each file has a level-1 date heading. Entries are appended under
      level-2 timestamp subheadings. The Appender uses an injected clock
      for testability and delegates text formatting to the entry package.
    success_criteria:
      - New file contains date heading + timestamp + entry
      - Subsequent appends preserve all prior content
      - Adaptive strategy suppresses headings within 10m unless on round boundary
      - Round10 strategy always rounds to nearest 10-minute boundary

  - id: timestamp-strategy
    name: Timestamp Strategy Engine
    responsibilities:
      - Implement adaptive strategy (FR-004 two-branch rule)
      - Implement round10 strategy (always round down to 10m boundary)
      - Compare current time against last heading to decide emission
    requirements: [FR-004, FR-005]
    summary: >-
      Pure decision function: given strategy, current time, last heading time,
      and whether a prior heading exists, returns whether to emit and what
      time to display. No I/O, no state — called by Appender.
    success_criteria:
      - Adaptive: gap < 10m + not on boundary → suppress
      - Adaptive: gap >= 10m or no prior → emit exact time
      - Adaptive: gap < 10m + on boundary → emit rounded time
      - Round10: same rounded slot as last → suppress
      - Round10: different slot → emit rounded time

  - id: viewer-dispatch
    name: Read Mode Viewer Dispatch
    responsibilities:
      - Find first available viewer on PATH
      - Replace process with viewer via syscall.Exec
      - Handle missing file gracefully
    requirements: [FR-006]
    summary: >-
      Dispatches today's log file to a terminal viewer. Priority order:
      glow --pager, rich --markdown --pager, $PAGER, cat. Uses exec to
      replace the im process. Returns ErrNoFile if the daily file doesn't
      exist.
    success_criteria:
      - First available viewer is selected
      - Missing file returns ErrNoFile
      - Process is replaced (not subprocess)

  - id: input-gathering
    name: Input Mode Gathering
    responsibilities:
      - Detect input mode (inline, pipe, editor)
      - Read piped stdin via io.ReadAll
      - Join positional args for inline mode
      - Guard against empty input
    requirements: [FR-007, FR-008, FR-009]
    summary: >-
      The CLI entrypoint gathers input from the detected mode, applies the
      empty input guard, then delegates to Appender or reader. Editor mode
      is deferred to DE-003.
    success_criteria:
      - Inline args joined with space
      - Piped stdin fully read
      - Empty/whitespace body produces no file mutation

  - id: configuration
    name: Configuration Loading
    responsibilities:
      - Load TOML config with sensible defaults
      - Expand ~ in paths
      - Validate enum fields
    requirements: [NF-001]
    summary: >-
      Zero-config on first run. Optional TOML at ~/.config/im/im.toml.
      Config struct populated from defaults when file missing, from file
      when present. Validation rejects unknown enum values.
    success_criteria:
      - No config file → defaults work
      - Config overrides respected
      - Invalid values rejected with clear error
```

```yaml supekku:verification.coverage@v1
schema: supekku.verification.coverage
version: 1
subject: SPEC-001
entries:
  - artefact: VT-001
    kind: VT
    requirement: SPEC-001.FR-001
    status: verified
    notes: internal/entry/entry_test.go — 16 table-driven test cases
  - artefact: VT-002
    kind: VT
    requirement: SPEC-001.FR-002
    status: verified
    notes: internal/logfile/appender_test.go — TestAppender_NewFile
  - artefact: VT-003
    kind: VT
    requirement: SPEC-001.FR-003
    status: verified
    notes: internal/logfile/appender_test.go — TestAppender_PriorContentUnmodified
  - artefact: VT-004
    kind: VT
    requirement: SPEC-001.FR-004
    status: verified
    notes: internal/logfile/timestamp_test.go — 13 cases; appender_integration_test.go
  - artefact: VT-005
    kind: VT
    requirement: SPEC-001.FR-005
    status: verified
    notes: internal/logfile/timestamp_test.go — round10 cases; integration test
  - artefact: VT-006
    kind: VT
    requirement: SPEC-001.FR-006
    status: verified
    notes: internal/reader/reader_test.go — ResolveViewer + View tests
  - artefact: VT-007
    kind: VT
    requirement: SPEC-001.FR-007
    status: verified
    notes: Smoke tested — io.ReadAll wiring in cmd/im/main.go
  - artefact: VT-008
    kind: VT
    requirement: SPEC-001.FR-008
    status: verified
    notes: internal/logfile/parse_test.go — 9 regex test cases
  - artefact: VT-009
    kind: VT
    requirement: SPEC-001.FR-009
    status: verified
    notes: Smoke tested — empty pipe produces no file mutation
  - artefact: VT-010
    kind: VT
    requirement: SPEC-001.NF-001
    status: verified
    notes: internal/config/config_test.go — defaults, file loading, validation
  - artefact: VA-001
    kind: VA
    requirement: SPEC-001.NF-002
    status: planned
    notes: Informal — code path is trivial, formal benchmark deferred
```

## 1. Intent & Summary

- **Scope / Boundaries**: Technical implementation of the `im` CLI tool.
  Covers five internal packages (`config`, `cli`, `entry`, `logfile`, `reader`)
  and the `cmd/im` entrypoint. Does not cover editor mode (DE-003).

- **Value Signals**: All PROD-001 non-editor requirements implemented with
  test coverage. `just check` exits 0 (lint + tests).

- **Guiding Principles**:
  - Minimal coupling: `entry` depends on nothing, `logfile` depends on `entry`
    and `config`, `reader` depends on nothing, `cmd/im` wires them all.
  - Testability via injected `func() time.Time` clock — no interfaces needed.
  - Append-only: `O_APPEND` flag, never read-modify-write.

- **Change History**: DE-001 (scaffold), DE-002 (core implementation).

## 2. Stakeholders & Journeys

- **Systems / Integrations**:
  - Filesystem: daily Markdown files under configurable `log_dir`.
  - External viewers: `glow`, `rich`, `$PAGER`, `cat` — discovered at runtime.
  - TOML config: `~/.config/im/im.toml` (optional).

- **Primary Journeys / Flows**:
  - **Inline write**: `im hello world` → detect inline mode → `entry.Format` →
    `logfile.Appender.Append` → file created/appended.
  - **Piped write**: `echo text | im` → detect pipe mode → `io.ReadAll` →
    same path as inline.
  - **Task write**: `im -t buy milk` → same path, `task=true` → checkbox prefix.
  - **Read**: `im -r` → resolve today's file path → `reader.View` →
    exec first available viewer.
  - **Empty input**: any mode producing empty body → silent exit 0, no mutation.

- **Edge Cases & Non-goals**:
  - Editor mode: deferred to DE-003 (stub in place).
  - Concurrent writers: not protected beyond `O_APPEND` atomicity.
  - Cross-day rollover: entry goes to current day's file (no lookback).

## 3. Responsibilities & Requirements

### Capability Overview

| Capability | Package | Dependencies |
|---|---|---|
| Entry text formatting | `internal/entry` | none |
| Daily log file management | `internal/logfile` | `entry`, `config` |
| Timestamp strategy engine | `internal/logfile` | `config` |
| Viewer dispatch | `internal/reader` | none |
| Input gathering + wiring | `cmd/im` | `cli`, `config`, `logfile`, `entry`, `reader` |
| Configuration loading | `internal/config` | none |

### Functional Requirements

- **FR-001**: `entry.Format` MUST return trimmed body with trailing newline
  for non-empty input, and apply `- [ ] ` prefix to each non-empty line when
  task flag is true. Empty/whitespace input MUST return empty string.

- **FR-002**: `logfile.Appender.Append` MUST create a new daily file with a
  level-1 date heading (`# Day, DD Month YYYY`) when the file does not exist.

- **FR-003**: `logfile.Appender.Append` MUST append to an existing daily file
  using `O_APPEND` without modifying prior content.

- **FR-004**: `logfile.ShouldEmitHeading` with adaptive strategy MUST:
  suppress heading when gap < 10m and not on a round-10 boundary; emit with
  rounded time when gap < 10m and on a round-10 boundary; emit with exact
  time when gap >= 10m or no prior heading exists.

- **FR-005**: `logfile.ShouldEmitHeading` with round10 strategy MUST always
  round down to the nearest 10-minute boundary and suppress when the rounded
  time matches the last heading's rounded time.

- **FR-006**: `reader.View` MUST dispatch to the first available viewer in
  order: `glow --pager`, `rich --markdown --pager`, `$PAGER`, `cat`. MUST
  return `ErrNoFile` when the daily file does not exist.

- **FR-007**: `cmd/im` MUST read piped stdin via `io.ReadAll` when input mode
  is `ModePipe`.

- **FR-008**: `logfile.ParseLastTimestamp` MUST extract the time from the last
  `## HH:MM` heading in file content using regex
  `^## (\d{2}:\d{2})\s*$`. MUST return false when no valid heading exists.

- **FR-009**: `cmd/im` MUST silently exit 0 with no file mutation when the
  gathered input body is empty or whitespace-only.

### Non-Functional Requirements

- **NF-001**: `config.Load` MUST return sensible defaults when config file is
  absent. MUST validate `editor_timestamp` and `timestamp_rounding` enum values.

- **NF-002**: Inline write path (no editor, no pager) SHOULD complete in
  < 100ms.

- **NF-003**: `logfile.Appender` MUST use an injected `func() time.Time`
  clock. No direct `time.Now` calls in the `logfile` package.

- **NF-004**: `logfile.Appender` MUST create the log directory
  (`os.MkdirAll`) if it does not exist, with permissions 0750.

### Operational Targets

- **Performance**: < 100ms for inline mode (informal, not benchmarked).
- **Reliability**: Zero data loss — append-only, no read-modify-write.
- **Maintainability**: Zero lint warnings (`golangci-lint`), table-driven tests.

## 4. Solution Outline

### Architecture / Components

```
cmd/im/main.go          — CLI entrypoint, input gathering, wiring
  ├── internal/config    — TOML config loading, defaults, validation
  ├── internal/cli       — Input mode detection (inline/pipe/editor)
  ├── internal/entry     — Pure text formatting, task prefix
  ├── internal/logfile   — File create/append, timestamp strategies
  └── internal/reader    — Viewer dispatch (exec)
```

### Data & Contracts

**Daily log file format** (`YYYY-MM-DD.md`):
```markdown
# Wednesday, 23 January 2026

## 10:03

fixed the auth bug

## 10:47

decided to use pgx instead of gorm
```

**Key types**:
- `config.Config` — LogDir, Editor, EditorTimestamp, TimestampRounding
- `logfile.Appender` — clock + strategy, method `Append(dir, body, task)`
- `entry.Format(body, task) string` — pure function
- `reader.View(path) error` — exec-based dispatch

## 5. Behaviour & Scenarios

### Primary Flows

**Write flow**:
1. Load config → resolve log dir.
2. Detect input mode from args/stdin.
3. Gather body (join args / ReadAll stdin).
4. Empty guard: if `strings.TrimSpace(body) == ""`, exit 0.
5. `entry.Format(body, task)` → formatted text.
6. `Appender.Append(dir, formatted, task)`:
   a. `MkdirAll(dir)`.
   b. `ReadFile` existing content (empty if new).
   c. `ParseLastTimestamp` → last heading time.
   d. `ShouldEmitHeading` → emit decision + display time.
   e. Build buffer: date heading (if new) + timestamp heading (if emit) + body.
   f. `OpenFile(O_APPEND|O_CREATE|O_WRONLY)` → write → close.

**Read flow**:
1. Load config → resolve log dir.
2. Build today's file path.
3. `reader.View(path)`:
   a. `os.Stat` — if not exist, return `ErrNoFile`.
   b. `ResolveViewer` — find first available on PATH.
   c. `syscall.Exec` — replace process.

### Error Handling / Guards

- Missing config file → defaults (not an error).
- Invalid config values → error with clear message.
- Missing log dir → created automatically.
- Empty input → silent exit 0.
- No viewer found → error "no viewer found".
- No daily file for `-r` → "no entries for today", exit 0.

## 6. Quality & Verification

- **Testing Strategy**: See `SPEC-001.tests.md` companion.
  - Unit: table-driven tests per package.
  - Integration: full write sequences in temp dirs.
  - Smoke: manual binary runs for all input modes.

- **Observability & Analysis**: N/A — local CLI tool, no telemetry.

- **Security & Compliance**: File permissions 0750 (dir) / 0600 (file).
  No network, no auth, no secrets.

- **Verification Coverage**: See `supekku:verification.coverage@v1` block.

- **Acceptance Gates**: `just check` exits 0.

## 7. Backlog Hooks & Dependencies

- **Related Specs / PROD**: PROD-001 (product spec, owns user-facing
  requirements). SPEC-001 implements the technical architecture.
- **Risks & Mitigations**: None beyond what's in DE-002.
- **Known Gaps / Debt**:
  - Editor mode (FR-008, FR-010) — deferred to DE-003.
  - Formal NF-002 benchmark — deferred.
  - Subheading level parameterization (OQ-001) — constant for now.
- **Open Decisions / Questions**: None active.
