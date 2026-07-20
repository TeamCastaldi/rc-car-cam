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
- No deployment adapter has been chosen yet (`@sveltejs/adapter-auto` is the placeholder) — swap it once the deployment target is decided (see `docs/foundation.md` open questions).

## Environment

Reads `PUBLIC_API_BASE_URL` (SvelteKit's `$env/static/public`) to know where the backend is. See `.env.example`.

## Running locally

```
cp .env.example .env   # then edit as needed
npm install
npm run dev -- --open
```

Or via the root `Makefile`: `make dev-frontend`.
