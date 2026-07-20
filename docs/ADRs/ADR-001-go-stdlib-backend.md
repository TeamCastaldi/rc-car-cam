# ADR-001 — Go + stdlib net/http for the backend

**Status**: Accepted
**Date**: 2026-07-20

## Decision

The car-side streaming server is written in Go 1.24+, using only the standard library's `net/http` — no web framework. chi, Gin, and Echo were all considered.

## Rationale

This is the first Go project for this team. Sticking to the standard library keeps the dependency surface at zero while learning the language, and Go 1.22+'s `ServeMux` (method + path-pattern matching) already covers the routing this project needs — a handful of routes for a streaming server, not a large JSON-CRUD API.

## Alternatives considered

- **chi** — a lightweight, idiomatic router on top of `net/http`. Ruled out for now, but it's stdlib-compatible, so it's a reasonable, low-risk upgrade later if middleware chaining becomes a real need.
- **Gin / Echo** — fuller-featured frameworks. Ruled out: heavier, and better suited to JSON-CRUD-heavy APIs than a streaming server; their abstractions would also work against the goal of learning Go itself.

## Consequences

- Adding routing/middleware conveniences later means either hand-rolling them or introducing chi — not a larger framework migration.
- Keeps the binary small and dependency-free, consistent with cross-compiling for a resource-constrained onboard computer (see ADR-006).
