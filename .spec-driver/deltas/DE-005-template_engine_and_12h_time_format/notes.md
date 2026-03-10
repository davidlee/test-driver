# Notes for DE-005

## 2026-03-09 — Phase 01 implementation

### Done
- **config.go**: Added `TimeFormat` type (`24h`/`12h`), `TitleFormat string`,
  `DefaultTemplatePath()`, validation for `time_format`.
- **parse.go**: Dual regex `(\d{1,2}):(\d{2})(?:\s*(AM|PM))?`. Extracted
  `convert12hTo24h` to satisfy cyclomatic complexity lint.
- **appender.go**: Refactored `Append` → `buildContent` + `renderNewFileHeader`
  + `appendFile` to reduce complexity. Uses `tmpl.Render` for new-file path.
  `FormatTime` for 12h/24h subheadings. Appender struct gains `timeFormat`,
  `templatePath`, `titleFormat`.
- **internal/tmpl** (new): Package with embedded `default.md` (`# {{ .Date }}\n`).
  `Render(path, Context)` with fallback to default on missing file or parse error.
  `missingkey=zero` option.
- **Tests**: Template (default, custom, fallback, missing key), 12h parse
  (AM, PM, noon, midnight, mixed), FormatTime (8 cases), config validation
  (invalid/valid time_format, default fields).

### Adaptations
- Renamed `internal/template` → `internal/tmpl` (lint: conflicts with stdlib).
- Extracted `buildContent`, `renderNewFileHeader`, `appendFile` from `Append` to
  satisfy cyclop (max 10).
- Extracted `convert12hTo24h` from `ParseLastTimestamp` for same reason.

### Verification
- `just check` green (0 lint issues, all tests pass including backward compat).

### Status
- Code + tests: complete, uncommitted.
- Ready for commit. P02 next (nanoID, title, frontmatter patching).
