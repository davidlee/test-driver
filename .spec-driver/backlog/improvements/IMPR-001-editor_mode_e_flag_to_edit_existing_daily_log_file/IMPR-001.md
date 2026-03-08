---
id: IMPR-001
name: 'Editor mode: -e flag to edit existing daily log file'
created: '2026-03-05'
updated: '2026-03-05'
status: idea
kind: improvement
---

# Editor mode: -e flag to edit existing daily log file

## Context
User feedback during DE-003 Phase 2. Currently `im` (no args, TTY) opens an
editor to compose a **new** entry. The `-e`/`--edit` flag would instead open
today's existing log file directly in the editor for free-form editing.

## Desired Behaviour
- `im -e` opens `<log_dir>/<today>.md` in the detected editor
- If no file exists for today, either create an empty one or error gracefully
- On save, the file is updated in place (no append logic)

## Origin
User feedback on DE-003 (editor mode and timestamp config).
