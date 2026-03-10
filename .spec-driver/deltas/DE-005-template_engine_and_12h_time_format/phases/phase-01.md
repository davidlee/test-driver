---
id: IP-005.PHASE-01
slug: 005-template_engine_and_12h_time_format-phase-01
name: 'IP-005 Phase 01 — Config, 12h time, template package'
created: '2026-03-09'
updated: '2026-03-09'
status: draft
kind: phase
---

```yaml supekku:phase.overview@v1
schema: supekku.phase.overview
version: 1
phase: IP-005.PHASE-01
plan: IP-005
delta: DE-005
objective: >-
  Add TimeFormat/TitleFormat config fields, implement 12h time formatting
  with dual-format parsing, create template package with embedded default,
  and wire template into appender new-file path.
entrance_criteria:
  - DR-005 reviewed
  - Current code read (appender.go, parse.go, timestamp.go, config.go)
exit_criteria:
  - FR-014 functional (template-driven new file creation)
  - FR-016 functional (12h subheadings)
  - Default template backward compatible with current output
  - ParseLastTimestamp handles both 24h and 12h
  - just check exits 0
verification:
  tests:
    - Template Render with default produces dateHeading-equivalent output
    - Template Render with custom template
    - Template Render falls back to default on missing file
    - formatTime 24h and 12h output
    - ParseLastTimestamp 12h cases (AM, PM, 12 AM, 12 PM)
    - Config validation for time_format
    - Existing appender tests pass unchanged
  evidence: []
tasks:
  - id: "1.1"
    description: "Config: add TimeFormat, TitleFormat fields"
  - id: "1.2"
    description: "12h time: formatTime + dual-format ParseLastTimestamp"
  - id: "1.3"
    description: "Template package: load, parse, render with embedded default"
  - id: "1.4"
    description: "Wire template into appender new-file path"
  - id: "1.5"
    description: "Tests for all new code"
risks:
  - Regex change to ParseLastTimestamp must not regress existing 24h parsing
```

```yaml supekku:phase.tracking@v1
schema: supekku.phase.tracking
version: 1
phase: IP-005.PHASE-01
```

# Phase 01 — Config, 12h time, template package

## 1. Objective
Foundation for FR-014 and FR-016. New template package, config extensions,
12h formatting, dual-format parsing. After this phase, daily files are created
from templates and 12h time is functional.

## 2. Links & References
- **Delta**: DE-005
- **Design Revision**: DR-005 §6.1 (template), §6.6 (12h), §6.7 (parse), §6.8 (config)
- **Requirements**: FR-014 (template), FR-016 (12h time)

## 3. Entrance Criteria
- [x] DR-005 reviewed
- [x] Current code read

## 4. Exit Criteria / Done When
- [x] `internal/tmpl` package exists with embedded default
- [x] Default template produces output identical to current `dateHeading()`
- [x] Custom template path loads and renders
- [x] `TimeFormat` config field with validation
- [x] 12h subheading format (`## 2:30 PM`)
- [x] `ParseLastTimestamp` handles both 24h and 12h
- [x] Appender uses template for new-file creation
- [x] All existing tests pass unchanged
- [x] `just check` exits 0

## 5. Verification
- `go test ./...`
- `just check`

## 6. Assumptions & STOP Conditions
- Assumptions: Template path fixed at `~/.config/im/template.md`
- STOP when: Appender config surface grows unwieldy (consider refactoring)

## 7. Tasks & Progress

| Status | ID | Description | Parallel? | Notes |
| --- | --- | --- | --- | --- |
| [x] | 1.1 | Config: TimeFormat, TitleFormat | [P] | + DefaultTemplatePath |
| [x] | 1.2 | 12h time: formatTime + dual parse | [P] | + convert12hTo24h extract |
| [x] | 1.3 | Template package | [P] | Renamed to `internal/tmpl` |
| [x] | 1.4 | Wire template into appender | | Refactored Append for cyclop |
| [x] | 1.5 | Tests | | 4 template + 5 parse + 8 format + 2 config |

### Task Details

- **1.1 Config: TimeFormat, TitleFormat**
  - Add `TimeFormat` type + constants (`24h`, `12h`)
  - Add `TitleFormat string` field
  - Default: `TimeFormat24h`, `TitleFormat: ""`
  - Validation: reject invalid `time_format`
  - Files: `internal/config/config.go`, `internal/config/config_test.go`

- **1.2 12h time: formatTime + dual-format parse**
  - Add `formatTime(t, format)` in appender or logfile package
  - Update subheading write to use `formatTime`
  - Change `timestampRe` to `(\d{1,2}):(\d{2})(?:\s*(AM|PM))?`
  - AM/PM→24h conversion in `ParseLastTimestamp`
  - Files: `internal/logfile/appender.go`, `internal/logfile/parse.go`, `internal/logfile/parse_test.go`, `internal/logfile/timestamp_test.go`

- **1.3 Template package**
  - New `internal/template/template.go`
  - Embed `internal/template/default.md` (`# {{ .Date }}\n`)
  - `Context` struct: `Date`, `CreatedAt`, `UpdatedAt`, `ID`, `Title`
  - `Render(templatePath string, ctx Context) (string, error)`
  - Fallback to default on missing file or parse error
  - Files: `internal/template/template.go`, `internal/template/default.md`, `internal/template/template_test.go`

- **1.4 Wire template into appender**
  - `Appender` gains `templatePath` and `timeFormat` from config
  - `NewAppender` signature: still `(clock, cfg)` — derive template path internally
  - Replace `dateHeading(now)` with `template.Render(path, ctx)` on new-file path
  - Template context: populate `Date` and `CreatedAt`/`UpdatedAt` with `now`; `ID` and `Title` as empty for now (P02)
  - Files: `internal/logfile/appender.go`

- **1.5 Tests**
  - Template: default render, custom render, fallback
  - Config: TimeFormat validation
  - Parse: 12h cases (1:00 AM, 12:00 PM, 2:30 PM, 12:00 AM)
  - Appender: existing tests pass, 12h subheading output
  - Integration: new file with default template ≡ old output

## 8. Risks & Mitigations
| Risk | Mitigation | Status |
| --- | --- | --- |
| ParseLastTimestamp regex regression | Existing test cases unchanged, add 12h cases | open |
| Template path resolution | Derive from config dir, not hardcode $HOME | open |

## 9. Decisions & Outcomes
- Per DR-005: all design decisions accepted

## 10. Findings / Research Notes
- (to be filled during execution)

## 11. Wrap-up Checklist
- [ ] Exit criteria satisfied
- [ ] Verification evidence stored
- [ ] Hand-off notes to P02 (nanoID, title, frontmatter patching)
