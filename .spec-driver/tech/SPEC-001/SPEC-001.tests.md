---
id: SPEC-001.TESTS
slug: im_cli-tests
name: im CLI Testing Guide
created: '2026-03-04'
updated: '2026-03-04'
status: draft
kind: guidance
aliases: []
relations:
- type: tests
  target: SPEC-001
guiding_principles: []
assumptions: []
---

# SPEC-001 Testing Guide

## 1. Overview
- **Tech Spec**: SPEC-001 – im CLI
- **Purpose**: Document testing strategy, suite inventory, and conventions
  for the `im` CLI implementation.
- **Test Owners**: Contributors to `im`.

## 2. Guidance & Conventions
- **Frameworks / Libraries**: Go standard `testing` package. No external
  test libraries (no testify).
- **Structure**: Tests colocated with source — `_test.go` files in each
  package directory.
- **Factories & Helpers**:
  - `fixedClock(h, m int) func() time.Time` — returns a clock function
    for deterministic timestamp tests.
  - `t.TempDir()` — all file tests use Go's temporary directory cleanup.
  - `t.Setenv()` — for `$PAGER` and environment variable tests.
- **Mocking Strategy**: No mocks. Dependency injection via function values
  (clock). File tests use real filesystem in temp dirs. Viewer tests use
  `exec.LookPath` against real PATH.

## 3. Strategy Matrix

| Scenario / Capability | Level | Rationale | Notes |
| --- | --- | --- | --- |
| Entry formatting (plain + task) | Unit | Pure function, no I/O | Table-driven, 16 cases |
| Last-timestamp parsing | Unit | Pure function on string input | Regex edge cases, 9 cases |
| Timestamp strategies (adaptive, round10) | Unit | Pure decision logic | PROD-001 §3 examples as test cases |
| File creation (new day) | Unit | Real filesystem in temp dir | Verifies date heading + structure |
| File append (same day) | Unit | Real filesystem in temp dir | Verifies append-only invariant |
| Heading suppression | Unit | Appender with fixed clock | Adaptive: within 10m, not on boundary |
| Heading emission on boundary | Unit | Appender with fixed clock | Adaptive: within 10m, on boundary |
| Full write sequence (adaptive) | Integration | Multi-entry scenario in temp dir | 5 entries + empty guard |
| Full write sequence (round10) | Integration | Multi-entry scenario in temp dir | 3 entries with rounding |
| Viewer resolution | Unit | LookPath on real PATH | Falls through to cat |
| Missing file for read mode | Unit | os.Stat check | Returns ErrNoFile |
| Config loading | Unit | Existing from DE-001 | Defaults, file, validation |
| Input mode detection | Unit | Existing from DE-001 | Inline, pipe, editor |

## 4. Test Suite Inventory

- **Suite**: `internal/entry/entry_test.go`
  - **Purpose**: Text formatting and task prefix
  - **Key Cases**:
    1. Single word → "hello\n"
    2. Multi-line → preserved with trailing newline
    3. Task single line → "- [ ] buy milk\n"
    4. Task multi-line → each line prefixed, blank lines skipped
    5. Empty/whitespace → ""
    6. Leading/trailing whitespace → trimmed
  - **Dependencies**: None

- **Suite**: `internal/logfile/parse_test.go`
  - **Purpose**: Regex-based last-timestamp extraction
  - **Key Cases**:
    1. No headings → (zero, false)
    2. Single heading → correct H:M
    3. Multiple headings → returns last
    4. Trailing space on heading → still matches
    5. Malformed headings (10:3, 25:00, not-a-time) → ignored
    6. Heading not at line start → ignored
    7. Midnight (00:00) and end-of-day (23:59)
  - **Dependencies**: None

- **Suite**: `internal/logfile/timestamp_test.go`
  - **Purpose**: Adaptive and round10 heading emission decisions
  - **Key Cases**:
    1. Adaptive: no prior heading → emit exact
    2. Adaptive: gap >= 10m → emit exact
    3. Adaptive: gap < 10m, not on boundary → suppress
    4. Adaptive: gap < 10m, on :10/:20/:00 boundary → emit rounded
    5. Adaptive: exactly 10m gap → emit exact
    6. Round10: rounds to 10m boundary
    7. Round10: same slot as last → suppress
    8. Round10: different slot → emit
  - **Dependencies**: `config.TimestampRounding` type

- **Suite**: `internal/logfile/appender_test.go`
  - **Purpose**: File creation, append, task prefix, round10, directory creation
  - **Key Cases**:
    1. New file → date heading + timestamp + entry
    2. Append to existing → new heading after gap
    3. Suppress heading within 10m window
    4. Task prefix in file output
    5. Round10 strategy rounding
    6. Nested directory creation
    7. Prior content not modified (append-only invariant)
  - **Dependencies**: `t.TempDir()`, `fixedClock`

- **Suite**: `internal/logfile/appender_integration_test.go`
  - **Purpose**: Full multi-entry write sequences
  - **Key Cases**:
    1. Adaptive: 5 entries spanning gaps, suppressions, boundaries, and task mode
    2. Round10: 3 entries with rounding and suppression
  - **Dependencies**: `t.TempDir()`, stepping clocks

- **Suite**: `internal/reader/reader_test.go`
  - **Purpose**: Viewer resolution and missing file handling
  - **Key Cases**:
    1. ResolveViewer finds cat (always available on POSIX)
    2. PAGER env var picked up
    3. View returns ErrNoFile for nonexistent path
  - **Dependencies**: `t.Setenv`, real PATH

## 5. Regression & Edge Cases

- **Timestamp date reconstruction**: `ParseLastTimestamp` returns time on
  zero date; Appender must reconstruct on today's date for correct gap
  calculation. Regression covered by `TestAppender_SuppressHeading`.
- **Blank lines in task mode**: `entry.Format` skips blank lines when
  applying task prefix. Covered by "task skips blank lines" test case.
- **File permissions**: gosec requires 0750 (dir) / 0600 (file). Lint
  catches regressions automatically.

## 6. Infrastructure & Amenities

- **Run all tests**: `just test` or `go test ./...`
- **Run with lint**: `just check` (= `just lint` + `just test`)
- **Run specific package**: `go test ./internal/logfile/...`
- **Verbose**: `go test ./... -v`
- **No external services required** — all tests use temp dirs and real
  filesystem.
- **Known flakiness**: None. All tests are deterministic (injected clocks,
  temp dirs).

## 7. Coverage Expectations

- **Target**: All exported functions have at least one test case. All
  branching logic (timestamp strategies) has exhaustive case coverage.
- **Critical behaviours**: Append-only invariant, timestamp emission rules,
  empty input guard.
- **Gaps**: `cmd/im/main.go` has no test file (integration tested via
  logfile package). Reader `syscall.Exec` path not unit-testable — covered
  by smoke test.

## 8. Backlog Hooks

- Editor mode tests (DE-003) — will need `$EDITOR` mocking or temp file
  verification.
- Formal NF-002 benchmark (< 100ms) — `testing.B` benchmark deferred.
- Reader integration test with mock viewers on PATH — would strengthen
  dispatch chain coverage.
