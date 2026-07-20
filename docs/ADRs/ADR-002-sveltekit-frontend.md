# ADR-002 — SvelteKit for the frontend

**Status**: Accepted
**Date**: 2026-07-20

## Decision

The browser viewer SPA is built with SvelteKit (Svelte 5) and TypeScript, scaffolded via the official `sv create` CLI (TypeScript, ESLint, Prettier, Vitest add-ons) rather than hand-written config.

## Rationale

Chosen over React+Vite and Vue+Vite for a smaller, more focused SPA with less boilerplate — this project's frontend is a single video viewer with a handful of controls, not a large application. Using the official scaffolding tool keeps the setup consistent with upstream defaults instead of a bespoke, hand-maintained config.

## Alternatives considered

- **React + Vite** — larger ecosystem, but more boilerplate than this project's scope calls for.
- **Vue + Vite** — a reasonable middle ground, but SvelteKit's smaller footprint won out.

## Consequences

- Future add-ons should go through `sv add` rather than hand-editing generated config, to stay consistent with upstream.
- Deployment depends on choosing a concrete SvelteKit adapter — see ADR-005.
