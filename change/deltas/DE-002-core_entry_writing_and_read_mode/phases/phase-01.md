---
id: IP-002.PHASE-01
slug: 002-core_entry_writing_and_read_mode-phase-01
name: Phase 1 - Entry text formatting
created: '2026-03-04'
updated: '2026-03-04'
status: complete
kind: phase
---

```yaml supekku:phase.overview@v1
schema: supekku.phase.overview
version: 1
phase: IP-002.PHASE-01
plan: IP-002
delta: DE-002
objective: >-
  Pure text composition — body join, trim, task prefix
entrance_criteria:
  - DR-002 reviewed
exit_criteria:
  - entry.Format returns correctly joined, trimmed text
  - entry.Format with task=true prefixes each line with "- [ ] "
  - Empty input returns empty string
  - Table-driven tests cover edge cases
  - just check exits 0
verification:
  tests:
    - internal/entry — table-driven unit tests
  evidence: []
tasks:
  - id: 1.1
    summary: Create internal/entry package with Format function
  - id: 1.2
    summary: Write table-driven tests for Format
  - id: 1.3
    summary: Lint and verify
risks: []
```

```yaml supekku:phase.tracking@v1
schema: supekku.phase.tracking
version: 1
phase: IP-002.PHASE-01
entrance_criteria:
  - item: "DR-002 reviewed"
    completed: true
exit_criteria:
  - item: "entry.Format returns correctly joined, trimmed text"
    completed: true
  - item: "entry.Format with task=true prefixes each line with '- [ ] '"
    completed: true
  - item: "Empty input returns empty string"
    completed: true
  - item: "Table-driven tests cover edge cases"
    completed: true
  - item: "just check exits 0"
    completed: true
```

# Phase 1 – Entry text formatting

## 1. Objective
Implement `internal/entry` — a pure text formatting package with no I/O
dependencies. Receives raw input lines and a task flag, returns a formatted
string ready for `logfile` to append.

## 2. Links & References
- **Delta**: DE-002
- **Design Revision**: DR-002 §DEC-004 (entry vs logfile split)
- **Specs / PRODs**: PROD-001.FR-005, FR-006, FR-007, FR-009
- **Support Docs**: PROD-001 §4 (output format), §5 (empty input)

## 3. Entrance Criteria
- [x] DR-002 reviewed

## 4. Exit Criteria / Done When
- [x] `entry.Format(text, false)` returns trimmed body with trailing newline
- [x] `entry.Format(text, true)` prefixes each non-empty line with `- [ ] `
- [x] Empty/whitespace-only input returns `""`
- [x] Table-driven tests cover: single line, multi-line, task prefix,
  empty input, whitespace-only, leading/trailing whitespace
- [x] `just check` exits 0

## 5. Verification
- `go test ./internal/entry/...` — table-driven tests
- `just check` (lint + test)

## 6. Assumptions & STOP Conditions
- Assumptions: `entry.Format` receives already-gathered text (a single string),
  not raw args or stdin. The caller joins args before passing them in.
- STOP when: ambiguity in how multi-line task prefix should handle blank lines
  (current assumption: skip blank lines, don't prefix them).

## 7. Tasks & Progress

| Status | ID | Description | Parallel? | Notes |
| --- | --- | --- | --- | --- |
| [x] | 1.1 | Write tests for entry.Format (TDD) | | 16 test cases |
| [x] | 1.2 | Implement entry.Format | | All 16 pass |
| [x] | 1.3 | Lint and verify | | just check exits 0, 0 lint issues |

### Task Details

- **1.1 Write tests for entry.Format**
  - **Design / Approach**: Table-driven tests. Cases: single word, multi-word,
    multi-line, task=true single line, task=true multi-line, task=true with
    blank lines (skip), empty string, whitespace-only string.
  - **Files / Components**: `internal/entry/entry_test.go`

- **1.2 Implement entry.Format**
  - **Design / Approach**: `func Format(body string, task bool) string`.
    Trim whitespace. If empty after trim, return `""`. If task, split on
    newlines, prefix non-empty lines with `- [ ] `, rejoin. Ensure trailing
    newline on non-empty result.
  - **Files / Components**: `internal/entry/entry.go`

- **1.3 Lint and verify**
  - **Design / Approach**: Run `just check`. Fix any lint issues.
  - **Files / Components**: n/a

## 8. Risks & Mitigations
| Risk | Mitigation | Status |
| --- | --- | --- |
| Blank line handling in task mode unclear | Skip blank lines (don't prefix) — revisit if spec clarifies | open |

## 9. Decisions & Outcomes
- `2026-03-04` — `entry.Format` signature: `func Format(body string, task bool) string`. Caller joins args. Entry package is pure.

## 10. Findings / Research Notes
- n/a

## 11. Wrap-up Checklist
- [x] Exit criteria satisfied
- [x] Verification evidence stored — `just check` exits 0, 16/16 tests pass
- [ ] Spec/Delta/Plan updated with lessons
- [x] Hand-off notes: `internal/entry` ready for `logfile` to consume in Phase 2
