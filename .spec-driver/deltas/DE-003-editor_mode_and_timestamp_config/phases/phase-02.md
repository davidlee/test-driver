---
id: IP-003.PHASE-02
slug: 003-editor_mode_and_timestamp_config-phase-02
name: Phase 2 - Wiring and timestamp config
created: '2026-03-04'
updated: '2026-03-04'
status: complete
kind: phase
---

```yaml supekku:phase.overview@v1
schema: supekku.phase.overview
version: 1
phase: IP-003.PHASE-02
plan: IP-003
delta: DE-003
objective: >-
  Wire editor.Edit into cmd/im/main.go ModeEditor case. Implement
  editor_timestamp start/end config per DR-003 §DEC-003.
entrance_criteria:
  - Phase 1 complete (editor package available)
exit_criteria:
  - im (no args, TTY) opens detected editor
  - Saved content appends to daily file under correct timestamp
  - Empty save exits cleanly, no file mutation
  - editor_timestamp=start uses pre-editor time (default)
  - editor_timestamp=end uses post-editor time
  - just check exits 0
verification:
  tests:
    - cmd/im integration test with mock editor script
  evidence: []
tasks:
  - id: 2.1
    summary: Write integration test with mock editor (TDD)
  - id: 2.2
    summary: Wire editor.Edit into runWrite ModeEditor case
  - id: 2.3
    summary: Implement editor_timestamp start/end selection
  - id: 2.4
    summary: Lint and verify
risks:
  - Mock editor script portability across CI environments
```

```yaml supekku:phase.tracking@v1
schema: supekku.phase.tracking
version: 1
phase: IP-003.PHASE-02
entrance_criteria:
  - item: "Phase 1 complete (editor package available)"
    completed: true
exit_criteria:
  - item: "im (no args, TTY) opens detected editor"
    completed: true
  - item: "Saved content appends to daily file under correct timestamp"
    completed: true
  - item: "Empty save exits cleanly, no file mutation"
    completed: true
  - item: "editor_timestamp=start uses pre-editor time (default)"
    completed: true
  - item: "editor_timestamp=end uses post-editor time"
    completed: true
  - item: "just check exits 0"
    completed: true
```

# Phase 2 – Wiring and timestamp config

## 1. Objective
Wire `editor.Edit` into `cmd/im/main.go` replacing the stub in `ModeEditor`.
Implement `editor_timestamp` start/end selection per DR-003 §DEC-003.

## 2. Links & References
- **Delta**: DE-003
- **Design Revision**: DR-003 §6 "Timestamp wiring in cmd/im/main.go", §DEC-003
- **Specs / PRODs**: PROD-001.FR-008, PROD-001.FR-010
- **Phase 1**: `internal/editor` package (commit `f82a44d`)

## 3. Entrance Criteria
- [x] Phase 1 complete — `internal/editor` package available

## 4. Exit Criteria / Done When
- [x] `im` (no args, TTY) opens detected editor
- [x] Saved content appends to daily file under correct timestamp
- [x] Empty save exits cleanly, no file mutation
- [x] `editor_timestamp=start` uses pre-editor time (default)
- [x] `editor_timestamp=end` uses post-editor time
- [x] `just check` exits 0

## 5. Verification
- Integration test: mock editor script writes known content, verify log file
- Integration test: mock editor that writes nothing (abort), verify no mutation
- Integration test: verify timestamp selection (start vs end)
- `just check` (lint + test)

## 6. Assumptions & STOP Conditions
- Assumptions: `editor.Edit` API is stable from Phase 1.
- STOP when: `Appender` clock injection doesn't support per-call override
  (it does — clock is set at construction time, so we create a new Appender
  with the chosen time).

## 7. Tasks & Progress

| Status | ID | Description | Notes |
| --- | --- | --- | --- |
| [x] | 2.1 | Write integration test with mock editor | 5 tests in main_test.go |
| [x] | 2.2 | Wire editor.Edit into runWrite ModeEditor | Replaced stub |
| [x] | 2.3 | Implement editor_timestamp start/end | Per DR-003 §DEC-003 |
| [x] | 2.4 | Lint and verify | 0 issues, 27 tests pass |

### Task Details

- **2.1 Integration test with mock editor**
  - Create mock editor shell script in testdata (writes known content to file arg)
  - Create empty mock (writes nothing — simulates abort)
  - Test: run `runEditor` helper → verify log file content and timestamp
  - Test: run with empty mock → verify no file created/mutated

- **2.2 Wire editor into runWrite**
  - Replace `fmt.Println("editor mode: not yet implemented")` with
    `editor.Edit(cfg, editor.DefaultCheck)` call
  - Empty content returns nil (no mutation)

- **2.3 Timestamp selection**
  - Capture `startTime` before `editor.Edit`, `endTime` after
  - Select based on `cfg.EditorTimestamp`
  - Pass closure to `logfile.NewAppender`

- **2.4 Lint and verify**
  - `just check` exits 0
