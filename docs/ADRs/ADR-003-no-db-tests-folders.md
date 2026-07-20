# ADR-003 — No top-level db/ or tests/ folders

**Status**: Accepted
**Date**: 2026-07-20

## Decision

The repo has no top-level `db/` folder and no top-level `tests/` folder.

## Rationale

- A database is not yet a settled requirement (see Open items) — scaffolding a `db/` folder ahead of that decision would be speculative.
- Tests are colocated with the source they cover in both stacks (`_test.go` next to Go source, Vitest convention in the frontend), following each language's own idiom rather than a shared top-level folder.

## Alternatives considered

- Scaffolding an empty `db/` now "for when it's needed" — ruled out as premature; not justified until a real persistence need exists.
- A shared top-level `tests/` folder — ruled out in favor of each stack's own colocation convention.

## Consequences

- Adding a database later means creating `db/` (or equivalent) at that time, plus updating `CLAUDE.md`'s Stack section.

## Open items

- Whether this project needs a persistent database at all, and for what (settings, clip metadata) — deferred until the streaming core works. See `docs/foundation.md`.
