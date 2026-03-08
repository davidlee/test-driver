---
id: IP-002.PHASE-03
slug: 002-core_entry_writing_and_read_mode-phase-03
name: Phase 3 - Reader, wiring, and integration
created: '2026-03-04'
updated: '2026-03-04'
status: complete
kind: phase
---

```yaml supekku:phase.overview@v1
schema: supekku.phase.overview
version: 1
phase: IP-002.PHASE-03
plan: IP-002
delta: DE-002
objective: >-
  Viewer dispatch, stdin pipe read, empty guard, wire cmd/im/main.go, integration tests
entrance_criteria:
  - Phase 2 complete (logfile package available)
exit_criteria:
  - im -r dispatches to glow/rich/$PAGER/cat (first available)
  - im -r with no file prints "no entries for today"
  - Piped stdin captured via io.ReadAll
  - Empty input guard prevents file mutation
  - cmd/im/main.go wired — diagnostic stubs replaced
  - End-to-end integration tests pass
  - just check exits 0
verification:
  tests:
    - internal/reader — unit tests for viewer dispatch
    - cmd/im — integration tests
  evidence: []
tasks:
  - id: 3.1
    summary: Write tests for reader.View
  - id: 3.2
    summary: Implement reader.View
  - id: 3.3
    summary: Wire cmd/im/main.go
  - id: 3.4
    summary: Write integration tests
  - id: 3.5
    summary: Lint and verify
risks: []
```

```yaml supekku:phase.tracking@v1
schema: supekku.phase.tracking
version: 1
phase: IP-002.PHASE-03
entrance_criteria:
  - item: "Phase 2 complete (logfile package available)"
    completed: true
exit_criteria:
  - item: "im -r dispatches to glow/rich/$PAGER/cat (first available)"
    completed: true
  - item: "im -r with no file prints 'no entries for today'"
    completed: true
  - item: "Piped stdin captured via io.ReadAll"
    completed: true
  - item: "Empty input guard prevents file mutation"
    completed: true
  - item: "cmd/im/main.go wired — diagnostic stubs replaced"
    completed: true
  - item: "End-to-end integration tests pass"
    completed: true
  - item: "just check exits 0"
    completed: true
```

# Phase 3 – Reader, wiring, and integration

## 1. Objective
Implement `internal/reader` (viewer dispatch), wire all packages in
`cmd/im/main.go`, add stdin pipe reading and empty input guard, and write
integration tests.

## 2. Links & References
- **Delta**: DE-002
- **Design Revision**: DR-002 §DEC-005 (viewer dispatch), §DEC-006 (pipe read),
  §DEC-007 (empty guard)
- **Specs / PRODs**: PROD-001.FR-005–FR-007, FR-009, FR-012

## 3. Entrance Criteria
- [x] Phase 2 complete — `internal/logfile` available

## 4. Exit Criteria / Done When
- [x] `reader.View` dispatches to glow/rich/$PAGER/cat
- [x] "no entries for today" when file doesn't exist
- [x] Piped stdin captured via `io.ReadAll`
- [x] Empty input → silent exit 0, no file mutation
- [x] `cmd/im/main.go` wired — stubs replaced
- [x] Integration tests pass (2 full-sequence tests + unit tests)
- [x] `just check` exits 0

## 5. Verification
- `go test ./internal/reader/...`
- `go test ./cmd/im/...` (integration)
- `just check`

## 6. Assumptions & STOP Conditions
- Assumptions: `syscall.Exec` replaces process for viewers. Integration tests
  use temp dirs and mock executables on PATH where needed.
- STOP when: viewer dispatch needs cross-platform shims beyond POSIX.

## 7. Tasks & Progress

| Status | ID | Description | Parallel? | Notes |
| --- | --- | --- | --- | --- |
| [x] | 3.1 | Write tests for reader.View | | 3 tests |
| [x] | 3.2 | Implement reader.View | | exec chain + ErrNoFile |
| [x] | 3.3 | Wire cmd/im/main.go | | All stubs replaced |
| [x] | 3.4 | Write integration tests | | 2 full-sequence tests |
| [x] | 3.5 | Lint and verify | | 0 lint issues |

### Task Details

- **3.1–3.2 Reader**
  - **Design**: `func View(path string) error`. Check file exists. Exec chain:
    glow --pager → rich --markdown --pager → $PAGER → cat.
    Use `exec.LookPath` to find viewers. `syscall.Exec` to replace process.
  - **Files**: `internal/reader/reader.go`, `internal/reader/reader_test.go`

- **3.3 Wire main.go**
  - **Design**: Replace diagnostic stubs. Read mode → `reader.View`.
    Write mode → gather input (args join / stdin read), empty guard,
    `logfile.Appender.Append`.
  - **Files**: `cmd/im/main.go`

- **3.4 Integration tests**
  - **Design**: Build binary, run with temp dir, verify file output.
  - **Files**: `cmd/im/main_test.go` or `test/integration_test.go`

## 8. Risks & Mitigations
| Risk | Mitigation | Status |
| --- | --- | --- |
| syscall.Exec hard to unit test | Test lookup logic separately, integration test the full path | open |

## 9. Decisions & Outcomes
- `2026-03-04` — Separated ResolveViewer (testable lookup) from View (exec wrapper).
- `2026-03-04` — Integration tests placed in logfile package (exercises Appender end-to-end).
- `2026-03-04` — Editor mode stub prints message, deferred to DE-003.

## 10. Findings / Research Notes
- Smoke test confirmed: inline, --, pipe, -t, -r, and empty input all work.

## 11. Wrap-up Checklist
- [x] Exit criteria satisfied
- [x] Verification evidence — `just check` exits 0, smoke test passed
- [ ] Spec/Delta/Plan updated with lessons
- [x] Hand-off: DE-002 implementation complete, ready for audit
