---
id: IP-004.PHASE-01
slug: 004-cli_read_mode_flags_f_path_p_pager_e_edit-phase-01
name: 'IP-004 Phase 01 — Implement all flags'
created: '2026-03-09'
updated: '2026-03-09'
status: draft
kind: phase
---

```yaml supekku:phase.overview@v1
schema: supekku.phase.overview
version: 1
phase: IP-004.PHASE-01
plan: IP-004
delta: DE-004
objective: >-
  Add -f, -p, -e flags to im CLI; split reader into Render (stdout) and
  View (pager); add editor.EditFile; fix -r to use Render.
entrance_criteria:
  - DR-004 reviewed
  - Current code surface read (main.go, reader.go, editor.go)
exit_criteria:
  - All four flags functional (-f, -p, -e, -r fix)
  - just check exits 0
  - Test coverage for each new flag and function
verification:
  tests:
    - reader.ResolveRenderer finds cat (stdout dispatch)
    - reader.Render returns ErrNoFile for missing file
    - runFile prints absolute path
    - runEdit prints message when no file
    - Flag dispatch precedence in run()
  evidence: []
tasks:
  - id: "1.1"
    description: "Split reader: add Render + ResolveRenderer"
  - id: "1.2"
    description: "Add editor.EditFile"
  - id: "1.3"
    description: "Wire -f, -p, -e flags + fix -r in main.go"
  - id: "1.4"
    description: "Tests for all new code"
risks:
  - "-r regression (VT-010 update needed)"
```

```yaml supekku:phase.tracking@v1
schema: supekku.phase.tracking
version: 1
phase: IP-004.PHASE-01
```

# Phase 01 — Implement all flags

## 1. Objective
Add `-f`, `-p`, `-e` flags and fix `-r` stdout/pager distinction.
Single phase — all changes are small, share the same code surface, no inter-dependencies.

## 2. Links & References
- **Delta**: DE-004
- **Design Revision**: DR-004 (full code sketches)
- **Requirements**: FR-013 (-f), FR-017 (-p), FR-019 (-e), FR-012 (-r fix)
- **Revision**: RE-006 (FR-019)

## 3. Entrance Criteria
- [x] DR-004 reviewed
- [x] Current code read (main.go, reader.go, editor.go, all tests)

## 4. Exit Criteria / Done When
- [x] `im -f` prints absolute path, exits 0
- [x] `im -r` renders to stdout (no pager)
- [x] `im -p` opens pager
- [x] `im -e` opens editor on existing file; message when no file
- [x] `just check` exits 0
- [x] Tests cover each new flag and function

## 5. Verification
- `go test ./...`
- `just check`
- Manual: `im -f` prints path; `im -r` renders stdout; `im -p` opens pager; `im -e` opens editor

## 6. Assumptions & STOP Conditions
- Assumptions: `glow` available for render testing (fallback to `cat` fine)
- STOP when: editor detection fails in unexpected ways (consult DE-003)

## 7. Tasks & Progress

| Status | ID | Description | Parallel? | Notes |
| --- | --- | --- | --- | --- |
| [x] | 1.1 | Split reader: Render + ResolveRenderer | [P] | Shared `resolve()` helper |
| [x] | 1.2 | Add editor.EditFile | [P] | `paramTypeCombine` lint fix |
| [x] | 1.3 | Wire flags in main.go | | `todayPath()` extracted |
| [x] | 1.4 | Tests for all new code | | 7 new tests |

### Task Details

- **1.1 Split reader: Render + ResolveRenderer**
  - Add `ResolveRenderer(filePath)` — dispatch chain: `glow` (no --pager) → `cat`
  - Add `Render(path)` — stat check + exec via ResolveRenderer
  - Existing `View` + `ResolveViewer` stay as-is (pager mode for `-p`)
  - Files: `internal/reader/reader.go`

- **1.2 Add editor.EditFile**
  - New exported func: `EditFile(path, cfgEditor, check)` — detects editor, runs on existing file
  - No temp file, no content validation — direct in-place edit
  - Files: `internal/editor/editor.go`

- **1.3 Wire flags in main.go**
  - Register `--file/-f`, `--pager/-p`, `--edit/-e` flags
  - Add `runFile`, `runPager`, `runEdit` functions
  - Update `run()` dispatch: file → read → pager → edit → write
  - Fix `-r`: call `reader.Render` instead of `reader.View`
  - Fix `-r` usage string: "Display today's log in a pager" → "Render today's log to stdout"
  - Files: `cmd/im/main.go`

- **1.4 Tests**
  - `reader_test.go`: ResolveRenderer finds cat, Render returns ErrNoFile
  - `main_test.go`: runFile prints path, runEdit message on no file
  - `editor_test.go`: EditFile test (mock exec)
  - Update VT-010 if needed

## 8. Risks & Mitigations
| Risk | Mitigation | Status |
| --- | --- | --- |
| -r regression | Update VT-010, verify existing tests | open |

## 9. Decisions & Outcomes
- Per DR-004 DEC-001: Split reader into Render/View (not parameterize)
- Per DR-004 DEC-002: EditFile uses os/exec directly (no go-editor)
- Per DR-004 DEC-003: Fix -r as part of this delta

## 10. Findings / Research Notes
- Current `reader.View` uses `glow --pager` — confirmed in code review
- `editor.detectEditor` is unexported but reusable from within the package
- Flag registration uses urfave/cli/v3 `BoolFlag` pattern

## 11. Wrap-up Checklist
- [ ] Exit criteria satisfied
- [ ] Verification evidence stored
- [ ] Delta/plan updated
