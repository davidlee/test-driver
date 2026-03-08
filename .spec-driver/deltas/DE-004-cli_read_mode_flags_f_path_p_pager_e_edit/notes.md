# Notes for DE-004

## 2026-03-09 — Phase 01 implementation

### Done
- **reader.go**: Split dispatch into `resolve()` helper with two candidate lists.
  `pagerCandidates()` (glow --pager, rich --markdown --pager, $PAGER, cat) and
  `rendererCandidates()` (glow, cat). `ResolveViewer`/`View` for pager,
  `ResolveRenderer`/`Render` for stdout. Shared `execViewer` handles stat+exec.
- **editor.go**: Added `EditFile(path, cfgEditor, check)` — detects editor via
  existing `detectEditor`, runs `os/exec.Command` with stdin/stdout/stderr
  connected. No temp file.
- **main.go**: Registered `--file/-f`, `--pager/-p`, `--edit/-e` flags.
  Extracted `todayPath()` helper. Switch dispatch: file → read → pager → edit → write.
  Fixed `-r` to call `reader.Render` (stdout) instead of `reader.View` (pager).
  Updated `-r` usage string.
- **Tests**: `ResolveRenderer` finds cat + no --pager assertion. `Render` returns
  `ErrNoFile`. `runFile` prints correct path (with and without existing file).
  `runEdit` prints "no entries" when no file. `EditFile` with `true` as editor.
  `EditFile` error when no editor available.

### Adaptations
- Lint caught `paramTypeCombine` on `EditFile` signature — combined `path, cfgEditor string`.
- DR-004 sketched separate resolve functions with duplicated loop logic; refactored
  into shared `resolve(filePath, candidates)` to stay DRY.

### Verification
- `just check` green (0 lint issues, all tests pass).

### Status
- Code + tests: complete, uncommitted.
- `.spec-driver` changes (DE-004 in-progress, IP-004, phase sheet): uncommitted.
- Ready for commit + audit.

