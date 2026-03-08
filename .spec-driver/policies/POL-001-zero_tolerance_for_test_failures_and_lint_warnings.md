---
id: POL-001
title: 'POL-001: Zero tolerance for test failures and lint warnings'
status: required
created: '2026-03-04'
updated: '2026-03-04'
reviewed: '2026-03-04'
owners: []
supersedes: []
superseded_by: []
standards: []
specs: []
requirements: []
deltas: []
related_policies: []
related_standards: []
tags: [quality, ci]
summary: All test failures and lint warnings must be resolved before proceeding with any work.
---

# POL-001: Zero tolerance for test failures and lint warnings

## Statement

All test failures and all lint warnings MUST be resolved immediately before
proceeding with further work. This applies regardless of whether the failures
or warnings are pre-existing or newly introduced.

There are no exceptions. "It was already broken" is not a valid reason to
continue.

## Rationale

Broken windows breed broken windows. Allowing pre-existing failures to
persist means every agent must distinguish "expected" failures from real
regressions — a cognitive tax that leads to missed bugs. Fixing forward
keeps the codebase honest and the feedback loop tight.

## Scope

- **Applies to**: All Go source code in this repository.
- **Tools**: `go test ./...` (tests), `golangci-lint run ./...` (lint).
- **Agents**: Every agent (human or AI) working in this repo.
- **Timing**: Before committing, before moving to the next task, and before
  declaring any work complete.

## Verification

- `just check` runs both `lint` and `test` — must exit 0.
- Agents MUST run `just check` after any code change and before announcing
  completion.
- If `just check` fails, the agent MUST fix the issues before proceeding,
  even if the failures are unrelated to the current task.
