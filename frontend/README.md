# Frontend

The viewer SPA: shows the car's live FPV video feed in the browser (phone or laptop). No driving controls — never will be; the RC car's own RF transmitter is the only way to drive it, a permanent scope boundary (see root `CLAUDE.md` Constraints), not a future feature. Talks to `backend/` over HTTP on the local network — nothing else.

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

`src/routes/+page.svelte` is the real viewer (Phases 2a/2b/3b) — an `<img>` streaming the live MJPEG feed with loading/connected/error states, auto-recovery via a `/healthz` liveness poll, and the shared-secret auth token sent as a query param. No component library yet — `src/lib/` still only has the scaffold's `vitest-examples/` placeholder.

## Conventions

- Tests are colocated next to the code they cover (Vitest convention) — no separate top-level `tests/` folder.
- `npm run lint` before committing; `npm run check` for type errors.
- Deployment adapter is decided but not yet wired up: `adapter-static`, served by the Go backend on the onboard computer (see `docs/ADRs/ADR-005-frontend-served-by-backend.md`; `docs/plans/2026-07-e2e-build-phases.md` Phase 6c). `@sveltejs/adapter-auto` is a placeholder until that phase lands.

## Environment

Reads `PUBLIC_API_BASE_URL` and `PUBLIC_STREAM_AUTH_TOKEN` (SvelteKit's `$env/static/public`) to know where the backend is and authenticate stream requests — the token must match the backend's `STREAM_AUTH_TOKEN`. See `.env.example`.

## Running locally

```
cp .env.example .env   # then edit as needed
npm install
npm run dev -- --open
```

Or via the root `Makefile`: `make dev-frontend`.
