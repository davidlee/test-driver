---
id: IP-003.PHASE-01
slug: 003-editor_mode_and_timestamp_config-phase-01
name: Phase 1 - Editor package
created: '2026-03-04'
updated: '2026-03-04'
status: draft
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
    completed: false
exit_criteria:
  - item: "detectEditor returns first available from config → $EDITOR → defaults"
    completed: false
  - item: "Edit returns content on exit 0, empty on non-zero/abort"
    completed: false
  - item: "Content validation rejects non-UTF-8, oversized, binary"
    completed: false
  - item: "Table-driven tests with injected CommandChecker"
    completed: false
  - item: "just check exits 0"
    completed: false
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
- [ ] DR-003 reviewed

## 4. Exit Criteria / Done When
- [ ] `detectEditor` returns first available from config → `$EDITOR` → defaults
- [ ] `Edit` returns content on exit 0, `""` on non-zero/abort
- [ ] `validateContent` rejects non-UTF-8, > 10MB, binary content
- [ ] Table-driven tests with injected `CommandChecker`
- [ ] `just check` exits 0

## 5. Verification
- `go test ./internal/editor/...` — table-driven tests
- `just check` (lint + test)

## 6. Assumptions & STOP Conditions
- Assumptions: `go-editor` `LaunchTempFile` API matches DR-003 §6 sketch
  (`ed.LaunchTempFile(pattern, reader) ([]byte, string, error)`). If not,
  adapt — the interface is thin.
- STOP when: `go-editor` doesn't support blocking wait for GUI editors
  (vscode, zed) — would need an alternative lib or raw `os/exec`.

## 7. Tasks & Progress

| Status | ID | Description | Parallel? | Notes |
| --- | --- | --- | --- | --- |
| [ ] | 1.1 | Write tests for detectEditor | | TDD |
| [ ] | 1.2 | Implement detectEditor | | |
| [ ] | 1.3 | Write tests for shouldSave | | TDD |
| [ ] | 1.4 | Implement launchEditor + shouldSave | | |
| [ ] | 1.5 | Write tests for validateContent | | TDD |
| [ ] | 1.6 | Implement validateContent + containsBinary | | |
| [ ] | 1.7 | Implement Edit (orchestrator) | | Wires detect → launch → save → validate |
| [ ] | 1.8 | Lint and verify | | `just check` exits 0 |

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
| `go-editor` API differs from sketch | Adapt — interface is thin (create, set command, launch) | open |
| `go-editor` doesn't block for GUI editors | Document: user must set `--wait` flag in config | open |

## 9. Decisions & Outcomes
(populated during implementation)

## 10. Findings / Research Notes
(populated during implementation)

## 11. Wrap-up Checklist
- [ ] Exit criteria satisfied
- [ ] Verification evidence stored
- [ ] Spec/Delta/Plan updated with lessons
- [ ] Hand-off notes to Phase 2 (wiring + timestamp config)
