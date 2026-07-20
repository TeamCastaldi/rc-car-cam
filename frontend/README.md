# Frontend

The viewer SPA: shows the car's live FPV video feed and driving controls in the browser (phone or laptop). Talks to `backend/` over HTTP on the local network — nothing else.

## Stack

- SvelteKit (Svelte 5), TypeScript
- Vitest for unit tests, ESLint + Prettier for lint/format
- Scaffolded with `npx sv create` (see the recreate command this generates for the exact flags) rather than hand-written — keep using `sv add` for future add-ons instead of hand-editing config where possible

## Structure

```
frontend/
  src/
    routes/     # SvelteKit pages — the video viewer page lives here
    lib/        # shared components/utilities
  static/       # static assets
```

The default template page (`src/routes/+page.svelte`) hasn't been replaced with the real viewer yet — that's the first real frontend work (see `CLAUDE.md` Current State).

## Conventions

- Tests are colocated next to the code they cover (Vitest convention) — no separate top-level `tests/` folder.
- `npm run lint` before committing; `npm run check` for type errors.
- Deployment adapter is decided but not yet wired up: `adapter-static`, served by the Go backend on the onboard computer (see `docs/ADRs/ADR-005-frontend-served-by-backend.md`; `docs/plans/2026-07-e2e-build-phases.md` Phase 6c). `@sveltejs/adapter-auto` is a placeholder until that phase lands.

## Environment

Reads `PUBLIC_API_BASE_URL` (SvelteKit's `$env/static/public`) to know where the backend is. See `.env.example`.

## Running locally

```
cp .env.example .env   # then edit as needed
npm install
npm run dev -- --open
```

Or via the root `Makefile`: `make dev-frontend`.
