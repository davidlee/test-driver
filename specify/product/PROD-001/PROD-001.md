---
id: PROD-001
slug: interstitial_mindfulness_logger
name: Interstitial Mindfulness Logger
created: '2026-03-04'
updated: '2026-03-04'
status: draft
kind: prod
aliases: [im]
relations: []
guiding_principles:
  - Friction-free capture — zero-config first run
  - Append-only daily files — never mutate past entries
  - Human-readable output — plain Markdown, no database
assumptions:
  - Single-user local tool, no network, no auth
  - POSIX-ish environment (Linux / macOS)
---

# PROD-001 – Interstitial Mindfulness Logger

```yaml supekku:spec.relationships@v1
schema: supekku.spec.relationships
version: 1
spec: PROD-001
requirements:
  primary:
    - PROD-001.FR-001
    - PROD-001.FR-002
    - PROD-001.FR-003
    - PROD-001.FR-004
    - PROD-001.FR-005
    - PROD-001.FR-006
    - PROD-001.FR-007
    - PROD-001.FR-008
    - PROD-001.FR-009
    - PROD-001.FR-010
    - PROD-001.FR-011
    - PROD-001.FR-012
    - PROD-001.NF-001
    - PROD-001.NF-002
    - PROD-001.NF-003
    - PROD-001.FR-013
    - PROD-001.FR-014
    - PROD-001.FR-015
    - PROD-001.FR-016
    - PROD-001.FR-017
  collaborators: []
interactions: []
```

```yaml supekku:spec.capabilities@v1
schema: supekku.spec.capabilities
version: 1
spec: PROD-001
capabilities:
  - id: daily-log-management
    name: Daily Log Management
    responsibilities:
      - Create daily Markdown file on first use each day
      - Append entries to existing daily file
    requirements: [FR-001, FR-002, FR-003]
    summary: >-
      Maintains one Markdown file per day under a configurable log directory.
      Each file has a date heading; entries are appended with coarse timestamps.
    success_criteria:
      - Running `im` on a new day creates a new file
      - Running `im` again on the same day appends to the existing file

  - id: timestamped-entries
    name: Timestamped Entry Capture
    responsibilities:
      - Emit timestamp subheadings using coarse/exact rules
      - Suppress subheadings within 10-min window unless on a round boundary
    requirements: [FR-004, FR-011]
    summary: >-
      Within 10 minutes of the last subheading, a new one is only emitted
      at exact 10-minute boundaries (minutes % 10 == 0), using the rounded
      time. After a 10+ minute gap, the exact unrounded time is used.
    success_criteria:
      - Entries at 10:03 and 10:07 share the "10:03" subheading (first was unrounded, gap >10m)
      - An entry at 10:10 gets a new "10:10" subheading (round boundary hit)
      - An entry at 10:47 after last subheading at 10:10 gets "10:47" (gap ≥10m, unrounded)

  - id: flexible-input
    name: Flexible Input Modes
    responsibilities:
      - Accept inline arguments
      - Accept input after -- separator
      - Read from stdin when piped
      - Open $EDITOR for interactive composition
    requirements: [FR-005, FR-006, FR-007, FR-008, FR-009, FR-010]
    summary: >-
      Multiple input modes so the user can capture thoughts however is
      most natural: quick one-liners from the shell, multi-line editor
      sessions, or piped output from other commands. The `-t` flag turns
      lines into Markdown task checkboxes.
    success_criteria:
      - Each input mode produces an appended entry in the daily log
      - `im -t buy milk` produces `- [ ] buy milk` in the log

  - id: configuration
    name: Minimal Configuration
    responsibilities:
      - Provide sensible defaults requiring no config file
      - Support optional TOML config at ~/.config/im/im.toml
    requirements: [NF-001]
    summary: >-
      Works out of the box with zero configuration. Optional TOML file
      lets users override the log directory and editor command.
    success_criteria:
      - First run with no config file succeeds
      - Config file overrides are respected

  - id: read-mode
    name: Read Mode
    responsibilities:
      - Render today's log to stdout (glow or cat)
      - Display today's log in a pager
      - Print today's file path for shell composition
    requirements: [FR-012, FR-013, FR-017]
    summary: >-
      The -r flag renders today's log to stdout via glow (if available)
      or cat — quick, no pager. The -p flag opens the log in a pager
      (glow --pager, rich --markdown --pager, $PAGER, or cat). The -f
      flag prints the absolute file path for shell composition.
    success_criteria:
      - `im -r` renders to stdout without a pager
      - `im -p` opens in a pager
      - `im -f` prints the absolute path to today's log file

  - id: file-template
    name: File Template Engine
    responsibilities:
      - Render daily file from a configurable Markdown template
      - Support dynamic variables in template context
      - Update frontmatter `updated_at` on entry append (when present)
    requirements: [FR-014, FR-015]
    summary: >-
      Daily files are created from a Go text/template. The template
      receives dynamic variables: created_at, updated_at, id (7-char
      nanoID), and title (from a configurable format string). When a
      file's frontmatter contains updated_at, it is patched on each
      append. If the field is absent, the append path is pure append.
    success_criteria:
      - Default template produces same format as current hardcoded output
      - Custom template with frontmatter renders correctly on file creation
      - updated_at is refreshed on append when present in frontmatter
      - updated_at is left alone (no mutation) when absent from file

  - id: time-format
    name: Time Format Configuration
    responsibilities:
      - Support 24-hour and 12-hour timestamp display
    requirements: [FR-016]
    summary: >-
      Timestamp subheadings support 24h (default) or 12h format with
      AM/PM, configurable via im.toml key time_format.
    success_criteria:
      - 24h mode produces `## 14:30`
      - 12h mode produces `## 2:30 PM`
```

```yaml supekku:verification.coverage@v1
schema: supekku.verification.coverage
version: 1
subject: PROD-001
entries:
  - artefact: VT-001
    kind: VT
    requirement: PROD-001.FR-001
    status: verified
    notes: internal/logfile/appender_test.go — TestAppender_NewFile creates YYYY-MM-DD.md with date heading
  - artefact: VT-002
    kind: VT
    requirement: PROD-001.FR-002
    status: verified
    notes: internal/logfile/appender_test.go — TestAppender_PriorContentUnmodified
  - artefact: VT-003
    kind: VT
    requirement: PROD-001.FR-003
    status: verified
    notes: internal/config/config_test.go — defaults to ~/log, overridable via im.toml
  - artefact: VT-004
    kind: VT
    requirement: PROD-001.FR-004
    status: verified
    notes: internal/logfile/timestamp_test.go — 13 cases covering adaptive strategy emission rules
  - artefact: VT-005
    kind: VT
    requirement: PROD-001.FR-005
    status: verified
    notes: Smoke tested — positional args joined as entry body
  - artefact: VT-006
    kind: VT
    requirement: PROD-001.FR-006
    status: verified
    notes: Smoke tested — -- separator passes remaining args as entry body
  - artefact: VT-007
    kind: VT
    requirement: PROD-001.FR-007
    status: verified
    notes: Smoke tested — piped stdin read via io.ReadAll in cmd/im/main.go
  - artefact: VT-008
    kind: VT
    requirement: PROD-001.FR-009
    status: verified
    notes: internal/entry/entry_test.go — task prefix "- [ ] " applied to each non-empty line
  - artefact: VT-009
    kind: VT
    requirement: PROD-001.FR-011
    status: verified
    notes: internal/logfile/timestamp_test.go — round10 strategy cases; config enum validated
  - artefact: VT-010
    kind: VT
    requirement: PROD-001.FR-012
    status: verified
    notes: internal/reader/reader_test.go — viewer dispatch chain, ErrNoFile handling
  - artefact: VT-011
    kind: VT
    requirement: PROD-001.NF-001
    status: verified
    notes: internal/config/config_test.go — zero-config defaults, optional TOML overrides
  - artefact: VT-012
    kind: VT
    requirement: PROD-001.NF-003
    status: verified
    notes: go build ./cmd/im produces single binary; just install places at ~/.local/bin/im
```

## 1. Intent & Summary

- **Problem / Purpose**: Developers and knowledge workers need a frictionless
  way to capture interstitial thoughts — the fleeting observations, decisions,
  and context that arise between tasks. Existing tools (notebooks, apps, wikis)
  impose too much ceremony. `im` is a CLI tool that appends timestamped entries
  to a daily Markdown file with zero setup.

- **Value Signals**: Adoption is measured by whether the user builds a habit.
  The tool succeeds if it takes < 5 seconds from thought to captured entry.

- **Guiding Principles**:
  - Friction-free capture — zero-config first run
  - Append-only daily files — never mutate past entries
  - Human-readable output — plain Markdown, no database
  - Do one thing well — capture, not search or analytics

- **Change History**: Initial specification. RE-003: added FR-013–FR-016
  (file path flag, template engine, dynamic variables, 12h time format).

## 2. Stakeholders & Journeys

- **Personas / Actors**:
  - *Developer* — wants to jot context between coding tasks without leaving
    the terminal. Values speed, keyboard-driven interaction, and files that
    integrate with existing note-taking (Obsidian, grep, git).

- **Primary Journeys / Flows**:
  - **Quick capture**: `im fixed the auth bug` → entry appended under today's
    timestamp heading.
  - **Multi-word with separator**: `im -- decided to use pgx instead of gorm,
    simpler for our case` → same behavior, `--` separates flags from text.
  - **Piped input**: `git log --oneline -5 | im` → piped content becomes the
    entry body.
  - **Editor mode**: `im` (no args, not piped) → opens `$EDITOR`, saved
    content becomes the entry.
  - **Read back**: `im -r` → display today's log in a pager / rich viewer.

- **Edge Cases & Non-goals**:
  - **Non-goal**: Search, tagging, or cross-day queries (use grep/rg).
  - **Non-goal**: Syncing, encryption, or multi-device support.
  - **Edge case**: Midnight rollover — an entry at 00:02 goes into the new
    day's file, not yesterday's.

## 3. Responsibilities & Requirements

### Functional Requirements

- **FR-001**: System MUST create `<log_dir>/YYYY-MM-DD.md` on first invocation
  each day, with a level-1 heading containing the formatted date.

- **FR-002**: System MUST append entries to an existing daily file without
  modifying prior content.

- **FR-003**: The log directory MUST default to `~/log` and be overridable
  via `im.toml`.

- **FR-004**: Timestamp subheading (level 2) emission rules:
  - If < 10 minutes since the last timestamp subheading: only emit a new
    subheading when `minutes % 10 == 0`, using the rounded time (e.g. 10:10,
    10:20).
  - If ≥ 10 minutes since the last timestamp subheading (or no prior
    subheading exists): emit a new subheading with the exact unrounded time
    (e.g. 10:03, 14:47).

- **FR-005**: System MUST accept positional arguments as the entry body
  (`im hello world` → entry "hello world").

- **FR-006**: System MUST accept entry text after a `--` separator
  (`im -- some text` → entry "some text").

- **FR-007**: System MUST read entry body from stdin when input is piped.

- **FR-008**: When invoked with no arguments and stdin is a TTY, system MUST
  open `$EDITOR` (or a configured editor) for entry composition. Editor mode
  MUST work correctly with neovim, emacs, helix, zed, and vscode.

- **FR-009**: When `-t` flag is provided, each line of the entry body MUST be
  prefixed with `- [ ] ` to produce Markdown task items.

- **FR-010**: In editor mode, the timestamp applied to the entry MUST be
  configurable via `im.toml` key `editor_timestamp` with values `"start"`
  (time editor was opened, default) or `"end"` (time editor was closed/saved).

- **FR-011**: Timestamp rounding strategy MUST be configurable via `im.toml`
  key `timestamp_rounding` with values:
  - `"adaptive"` (default) — the two-branch rule from FR-004 (coarse within
    10m, exact after gaps).
  - `"round10"` — always round down to the nearest 10-minute boundary.

- **FR-012**: When `-r` flag is provided, system MUST render today's log file
  to stdout using `glow` (if available) or `cat`. No pager. Output goes
  directly to the terminal.

- **FR-017**: When `-p` flag is provided, system MUST display today's log file
  in a pager using the first available viewer: `glow --pager`,
  `rich --markdown --pager`, `$PAGER`, or `cat`.

- **FR-013**: When `-f` flag is provided, system MUST print the absolute path
  to today's log file to stdout and exit. No entry is written.

- **FR-014**: Daily file format MUST be defined by an editable Markdown template
  using Go `text/template` syntax. Default template location:
  `~/.config/im/template.md`. A bundled default template MUST produce output
  equivalent to the current hardcoded format when no custom template exists.

- **FR-015**: The template engine MUST support the following dynamic variables:
  - `created_at` — file creation timestamp, format `YYYY-MM-DD HH:MM`
  - `updated_at` — last entry timestamp, format `YYYY-MM-DD HH:MM`
  - `id` — 7-character nanoID (`github.com/matoous/go-nanoid/v2`), generated
    once per file on creation
  - `title` — rendered from a configurable format string in `im.toml`
    (key: `title_format`, e.g. `title_format = "Log — {{ date }}"`)

  When a file's YAML frontmatter contains an `updated_at` field, it MUST be
  refreshed on each entry append. When `updated_at` is absent from the file,
  the append path MUST NOT mutate prior content.

- **FR-016**: Timestamp subheadings MUST support 12-hour format with AM/PM,
  configurable via `im.toml` key `time_format` with values:
  - `"24h"` (default) — `## 14:30`
  - `"12h"` — `## 2:30 PM`

### Non-Functional Requirements

- **NF-001**: System MUST work with zero configuration on first run. Optional
  config file at `~/.config/im/im.toml`.

- **NF-002**: Time from invocation to entry written MUST be < 100 ms for
  inline argument mode (excluding editor/pager).

- **NF-003**: Binary MUST be a single statically-linked Go executable named
  `im`.

### Success Metrics / Signals

- **Adoption**: User invokes `im` ≥ 3 times per working day after first week.
- **Quality**: Zero data loss — entries are never dropped or corrupted.

## 4. Solution Outline

- **Binary**: `im`, built with Go, single-file distribution.
- **Config**: `~/.config/im/im.toml` — optional, keys: `log_dir`, `editor`,
  `editor_timestamp` (`"start"` | `"end"`),
  `timestamp_rounding` (`"adaptive"` | `"round10"`),
  `time_format` (`"24h"` | `"12h"`), `title_format` (Go template string).
- **Template**: `~/.config/im/template.md` — optional Go `text/template`
  defining daily file structure. Variables: `created_at`, `updated_at`, `id`,
  `title`, `date`.
- **Output format**:
  ```markdown
  # Wednesday, 23 January 2026

  ## 10:04

  fixed the auth bug

  ## 10:22

  decided to use pgx instead of gorm, simpler for our case

  ## 14:30

  reviewed PR #42, left comments on error handling
  ```

## 5. Behaviour & Scenarios

- **File creation**: If `<log_dir>/YYYY-MM-DD.md` does not exist, create it
  with the day heading, then append the entry under a timestamp subheading.
- **File append**: If the file exists, read the last `## HH:MM` heading. If
  the current rounded time matches, append the entry below it. Otherwise,
  append a new `## HH:MM` heading followed by the entry.
- **Read mode** (`-r`): Display today's file using `$PAGER`, `glow`, or `cat`
  (first available).
- **Empty input**: If the user provides no text (empty editor save, empty
  pipe), do nothing — exit cleanly without appending.

## 6. Quality & Verification

- **Testing Strategy**: Unit tests for timestamp rounding, file creation logic,
  and input mode detection. Integration tests for end-to-end append scenarios.
- **Verification Coverage**: To be populated as tech specs and implementation
  proceed.

## 7. Backlog Hooks & Dependencies

- **Open Decisions / Questions**:
  - Exact `im.toml` schema (keep minimal for v1).
  - Whether `-r` should support viewing arbitrary dates (`-r 2026-01-20`).
