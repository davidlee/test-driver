---
id: mem.fact.project.stdlib-naming
name: Go stdlib package name conflicts
kind: memory
status: active
memory_type: fact
created: '2026-03-10'
updated: '2026-03-10'
tags:
- project
- lint
- sharp-edge
summary: revive var-naming rejects internal package names matching Go stdlib (e.g.
  template). Renamed internal/template to internal/tmpl.
---

# Go stdlib package name conflicts

The `revive` linter (via `golangci-lint`) enforces `var-naming`: package
names must not match Go standard library package names.

- `internal/template` → renamed `internal/tmpl` (conflicts with `text/template`)
- Check before naming new packages: `go doc <name>` — if it resolves, pick a different name.

Scope: any new `internal/` package in this project.

Provenance: DE-005 P01 lint failure, commit a58f540.
