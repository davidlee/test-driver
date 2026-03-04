---
id: mem.pattern.project.agent-pitfalls
name: Agent pitfalls when using spec-driver
kind: memory
status: active
memory_type: pattern
created: '2026-03-04'
updated: '2026-03-04'
tags: []
summary: Common wrong turns agents take with spec-driver tooling
---

# Agent pitfalls when using spec-driver

## Summary

Wrong turns observed during agent sessions, with corrections.

## Pitfalls

- **`doctrine` is a skill, not a file.** Tried to `Read .spec-driver/doctrine.md` — doesn't exist. Use `/doctrine` (the skill) instead.
- **Requirements must be synced before validation.** After writing FR/NF requirements in a PROD/SPEC, run `spec-driver sync` to register them in the requirements registry. Without sync, `validate` rejects `applies_to.requirements` references as "not found".
- **Don't remove valid references to appease the validator.** When validation fails on requirement refs, the fix is `sync`, not deleting the references.
- **`resolve` is for memory links only.** `spec-driver resolve links` resolves `[[...]]` links in memory bodies — it does not look up requirements or specs.
- **`find` and `schema` have subcommands, not positional args.** Use `--help` to check syntax before guessing.
- **`sync` registers requirements automatically.** It parses `FR-NNN` / `NF-NNN` patterns from spec Markdown and populates the registry. No manual step needed beyond writing the spec correctly.
- **Don't use `spec-driver` domain for project memories.** Memories under `mem.*.spec-driver.*` are managed by the spec-driver installer and get removed on reinstall. Use `mem.*.project.*` for project-specific memories.
- **`create phase` appends a duplicate phase entry to the plan YAML.** Running `spec-driver create phase "..." --plan IP-XXX` appends a stray `- id: IP-XXX.PHASE-XX` line at the end of the `phases:` list in the plan file. You must manually remove this duplicate after each `create phase` call.
