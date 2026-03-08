# Notes for DE-001

## Phase 1 — Go module and build skeleton

**Status**: complete

**Done**:
- `go mod init im` with Go 1.24 target
- `github.com/urfave/cli/v3` v3.7.0 added
- `cmd/im/main.go` — stub entrypoint with `-t`, `-r` flags, version injection point
- `internal/config/`, `internal/cli/` dirs created (empty, ready for Phase 2/3)
- `.golangci.yml` v2 config in place (from vice, zk exclusion removed)
- `justfile` recipes: `build`, `install`, `lint`, `test`, `check`

**Surprises / adaptations**:
- Initial golangci-lint was v1 (nix); config is v2 format. Resolved by nix update to golangci-lint v2.10.1.
- `go.mod` initially had `go 1.25.7` (auto-detected); pinned to `go 1.24` for broader compat.

**Verification**: `just check` exits 0 (lint: 0 issues, test: no test files yet).

**Uncommitted work**: all changes are uncommitted.

**Follow-up**: Phase 2 (config loading) is next.

## Phase 2 — Config loading

**Status**: complete

**Done**:
- `internal/config/config.go` — `Config` struct, `Load()`, `DefaultConfig()`,
  `ResolvedLogDir()`, `ResolvedEditor()`, validation, `~` expansion
- `BurntSushi/toml` v1.6.0 for TOML parsing
- 10 tests covering: defaults, full config, partial config, invalid values,
  invalid TOML, path expansion, editor fallback chain

**Surprises**: `t.Setenv` cannot be used with `t.Parallel` — removed parallel
from env-dependent tests.

**Verification**: `just check` exits 0.

## Phase 3 — CLI flags and input mode detection

**Status**: complete

**Done**:
- `internal/cli/cli.go` — `InputMode` type, `DetectInputMode()`, `IsTerminal` var
- `golang.org/x/term` for TTY detection
- 6 tests: inline (single/multi arg), pipe, editor, args-win-over-pipe, String()
- `cmd/im/main.go` wired up: loads config, detects mode, prints diagnostic output
- All modes verified: `im hello`, `im -t buy milk`, `im -- text`, `echo | im`

**Design note**: `IsTerminal` exported as a var for test stubbing. Acceptable
for an internal package; would use an interface if this were public API.

**Verification**: `just check` exits 0.

**All phases complete. Uncommitted work.**
