---
id: IP-003.PHASE-01
slug: 003-editor_mode_and_timestamp_config-phase-01
name: Phase 1 - Editor package
created: '2026-03-04'
updated: '2026-03-04'
status: complete
kind: phase
---

```yaml supekku:phase.overview@v1
schema: supekku.phase.overview
version: 1
phase: IP-003.PHASE-01
plan: IP-003
delta: DE-003
objective: >-
  Editor detection, launch, save logic, content validation
entrance_criteria:
  - DR-003 reviewed
exit_criteria:
  - detectEditor returns first available from config → $EDITOR → defaults
  - Edit returns content on exit 0, empty on non-zero/abort
  - Content validation rejects non-UTF-8, oversized, binary
  - Table-driven tests with injected CommandChecker
  - just check exits 0
verification:
  tests:
    - internal/editor — table-driven unit tests
  evidence: []
tasks:
  - id: 1.1
    summary: Write tests for editor detection
  - id: 1.2
    summary: Implement detectEditor
  - id: 1.3
    summary: Write tests for launchEditor and shouldSave
  - id: 1.4
    summary: Implement launchEditor and shouldSave
  - id: 1.5
    summary: Write tests for content validation
  - id: 1.6
    summary: Implement validateContent and containsBinary
  - id: 1.7
    summary: Implement Edit (top-level orchestrator)
  - id: 1.8
    summary: Lint and verify
risks:
  - go-editor lib API may differ from sketched usage
```

```yaml supekku:phase.tracking@v1
schema: supekku.phase.tracking
version: 1
phase: IP-003.PHASE-01
entrance_criteria:
  - item: "DR-003 reviewed"
    completed: true
exit_criteria:
  - item: "detectEditor returns first available from config → $EDITOR → defaults"
    completed: true
  - item: "Edit returns content on exit 0, empty on non-zero/abort"
    completed: true
  - item: "Content validation rejects non-UTF-8, oversized, binary"
    completed: true
  - item: "Table-driven tests with injected CommandChecker"
    completed: true
  - item: "just check exits 0"
    completed: true
```

# Phase 1 – Editor package

## 1. Objective
Implement `internal/editor` — editor detection, launch via temp file, save
logic, and content validation. Pure package with no coupling to `cmd/im`.
All behaviour covered by table-driven tests with injected `CommandChecker`.

## 2. Links & References
- **Delta**: DE-003
- **Design Revision**: DR-003 §6 (code sketches — full implementations for
  all functions), §DEC-001 (simplification rationale), §DEC-002 (go-editor
  lib), §DEC-004 (fallback chain)
- **Specs / PRODs**: PROD-001.FR-008
- **Dep**: `github.com/confluentinc/go-editor`

## 3. Entrance Criteria
- [x] DR-003 reviewed

## 4. Exit Criteria / Done When
- [x] `detectEditor` returns first available from config → `$EDITOR` → defaults
- [x] `Edit` returns content on exit 0, `""` on non-zero/abort
- [x] `validateContent` rejects non-UTF-8, > 10MB, binary content
- [x] Table-driven tests with injected `CommandChecker`
- [x] `just check` exits 0

## 5. Verification
- `go test ./internal/editor/...` — table-driven tests
- `just check` (lint + test)

## 6. Assumptions & STOP Conditions
- Assumptions: `go-editor` `LaunchTempFile` API matches DR-003 §6 sketch
  (`ed.LaunchTempFile(pattern, reader) ([]byte, string, error)`). Confirmed.
- STOP when: `go-editor` doesn't support blocking wait for GUI editors
  (vscode, zed) — would need an alternative lib or raw `os/exec`.

## 7. Tasks & Progress

| Status | ID | Description | Parallel? | Notes |
| --- | --- | --- | --- | --- |
| [x] | 1.1 | Write tests for detectEditor | | 9 cases |
| [x] | 1.2 | Implement detectEditor | | |
| [x] | 1.3 | Write tests for shouldSave | | 6 cases |
| [x] | 1.4 | Implement launchEditor + shouldSave | | |
| [x] | 1.5 | Write tests for validateContent | | 6 cases + oversized |
| [x] | 1.6 | Implement validateContent + containsBinary | | |
| [x] | 1.7 | Implement Edit (orchestrator) | | |
| [x] | 1.8 | Lint and verify | | 0 lint issues |

### Task Details

- **1.1–1.2 Editor detection**
  - **Design**: DR-003 §6 `detectEditor` sketch. Table-driven: config set,
    `$EDITOR` set, neither set, none available → error. Inject
    `CommandChecker` that returns true for specific names.
  - **Files**: `internal/editor/editor.go`, `internal/editor/editor_test.go`

- **1.3–1.4 Launch + save logic**
  - **Design**: DR-003 §6 `launchEditor` and `shouldSave` sketches.
    `launchEditor` uses `go-editor` `LaunchTempFile`. `shouldSave`: exit 0
    + non-empty → true, else false. Test save logic as pure function.
    Launch tested via integration in Phase 2 (mock editor script).
  - **Files**: `internal/editor/editor.go`, `internal/editor/editor_test.go`

- **1.5–1.6 Content validation**
  - **Design**: DR-003 §6 `validateContent` + `containsBinary` sketches.
    Cases: valid UTF-8, invalid UTF-8, > 10MB, null bytes, > 5%
    non-printable, normal text with tabs/newlines (should pass).
  - **Files**: `internal/editor/editor.go`, `internal/editor/editor_test.go`

- **1.7 Edit orchestrator**
  - **Design**: `Edit(cfg, check)` → `detectEditor` → `launchEditor` →
    `shouldSave` → `validateContent` → return content or `""`.
  - **Files**: `internal/editor/editor.go`

- **1.8 Lint and verify**
  - **Design**: `just check`. Fix any lint issues.

## 8. Risks & Mitigations
| Risk | Mitigation | Status |
| --- | --- | --- |
| `go-editor` API differs from sketch | Adapted: returns `[]byte` not `string`; cast in `launchEditor` | resolved |
| `go-editor` doesn't block for GUI editors | Document: user must set `--wait` flag in config | open |

## 9. Decisions & Outcomes
- `2026-03-04` — `go-editor` `LaunchTempFile` returns `([]byte, string, error)` as sketched. Content cast to `string` in `launchEditor`.
- `2026-03-04` — Lint required named return values on `launchEditor` (gocritic).

## 10. Findings / Research Notes
- `go-editor` v0.11.0 API confirmed: `NewEditor()`, `.Command`, `.LaunchTempFile(pattern, reader)`.

## 11. Wrap-up Checklist
- [x] Exit criteria satisfied
- [x] Verification evidence — `just check` exits 0, 22 tests across editor package
- [ ] Spec/Delta/Plan updated with lessons
- [x] Hand-off: `internal/editor` ready for wiring in Phase 2
