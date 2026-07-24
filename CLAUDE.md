# CLAUDE.md

This file is the primary context document for Claude (and other LLM assistants) working in this repository. Fill it in as the project takes shape. An incomplete CLAUDE.md is better than none — add sections as you know them.

---

## How to fill this in

This file should be filled in during or just after initial project setup, then kept current as the project evolves. The best time to update it is at the end of a working session, before you close the repo.

Each section below has a comment explaining what to put there. Remove the comments as you fill in real content.

A well-maintained CLAUDE.md means every new LLM session starts with full context instead of having to rediscover the project from scratch.

---

## Project identity

RC Car Cam — a camera mounted on an RC car streams live FPV video to a phone/browser so you can drive it beyond line of sight. It's a hands-on learning project: the point is the build itself (embedded video, streaming, networking), not solving an external pain point. Built for personal/family use, not a commercial product.

This is also Nathan's first Go project. When implementing or modifying backend Go code, explain new language features, stdlib idioms, and *why* an approach was chosen as part of the response — not just a diff. Keep the explaining in conversation, not inline code comments (this repo already prefers minimal comments — see Code style). Frontend/TypeScript work doesn't need the same treatment unless asked.

## Stack

- **Backend** (`backend/`): Go 1.24+, standard library only (`net/http`) — no framework
- **Frontend** (`frontend/`): SvelteKit (Svelte 5), TypeScript
- **Database**: none yet — see Open questions
- **Test runner**: `go test` (backend), Vitest (frontend)
- **Linter/formatter**: `go vet` + golangci-lint + `gofmt` (backend); ESLint + Prettier (frontend)
- **Versioning/releases**: release-please (Conventional Commits → release PRs, tags, `CHANGELOG.md`, versioned independently per package) + GoReleaser (cross-compiles and publishes the backend binary on each `backend/vX.Y.Z` tag)
- **Deployment target**: undecided — the car's onboard computer hasn't been chosen yet

## Architecture

- `backend/` — the streaming server, runs on the car's onboard computer. Entry point is `cmd/server/main.go`. `internal/camera` (mock `Source` + `MockSource`), `internal/stream` (MJPEG handler), and `internal/auth` (shared-secret token middleware) all exist and are wired in. `internal/api` is still sketched-only, not created — add it once real control/status routes are needed, not speculatively.
- `frontend/` — the browser viewer SPA. `src/routes/+page.svelte` is the real viewer now (not the SvelteKit default template) — an `<img>` pointed at `PUBLIC_API_BASE_URL` + `/stream` (with the auth token as a query param), loading/connected/error states with auto-recovery via a `/healthz` liveness poll. Talks to the backend over HTTP on the local network only. Once built for real deployment, it's served as static files by the backend itself (see ADR-005) — not run separately. **No driving controls exist or are planned in the frontend** — the RC car's physical RF transmitter remains the only way to drive; this was explicitly confirmed as a permanent scope boundary, not a "not started yet" item (see Constraints).
- Data flow (working end-to-end with the mock camera): `camera.MockSource` → `stream.Handler` (MJPEG, `multipart/x-mixed-replace`, see ADR-004) → `auth.RequireToken` middleware on `/stream` → frontend `<img>` element. No autonomous/CV processing sits in this path — see Constraints. Real camera hardware swaps in for the mock at Phase 5 with no other layer changing.
- Auth exists and is enforced: `STREAM_AUTH_TOKEN` (backend) / `PUBLIC_STREAM_AUTH_TOKEN` (frontend, must match) are required — the backend fails to start without a token set. No database exists yet (see Open questions).
- Both stacks are driven through the root `Makefile` (`dev-backend`, `dev-frontend`, `build`, `test`, `lint`) — **run `make` targets from the repo root, not from inside `backend/`/`frontend/`** (a recurring gotcha — the Makefile's targets assume the root as cwd). The `.github/prompts/*.prompt.md` files and CI reference `make test`/`make lint` rather than stack-specific commands directly, so adding a stack later means updating the Makefile in one place.
- Full build sequence (mock pipeline → hardware bring-up → real integration → deploy → physical assembly → drive test) is phased in `docs/plans/2026-07-e2e-build-phases.md`. **Phases 0–3 (the entire software-only mock pipeline) are done as of 2026-07-24.** Phase 4 onward requires physical hardware that hasn't been purchased yet — see "Current state" and the root `README.md`'s "Where things stand" section before starting new feature work.

## Constraints (non-negotiable)

- Never expose the stream or control endpoints to the public internet without auth — local network or authenticated access only.
- Never add motor/steering control code — this project is camera + streaming only; driving the car itself is out of scope.
- No autonomous/CV features (object detection, line-following, self-driving) — manual FPV control only.
- No multi-user/cloud product surface — single car, single user/family, no accounts, no cloud backend.

## Code style

- Go: `gofmt`-clean, default `golangci-lint` rule set, no custom config.
- Svelte/TypeScript: ESLint + Prettier defaults as generated by the `sv create` scaffold, no overrides.
- Tests are colocated next to the source they cover (`_test.go` in Go, Vitest convention in the frontend) — there is no top-level `tests/` folder.
- Config comes from environment variables (`.env`, gitignored — see each service's `.env.example`), not flags or config files.

## Current state

**Phases 0–3 of `docs/plans/2026-07-e2e-build-phases.md` are complete — the entire software-only mock pipeline works end-to-end, verified live in a browser.** Everything remaining is blocked on Phase 4a: buying the physical hardware (`docs/hardware.md`). See the root `README.md`'s "Where things stand" section for the resume plan — **read that section, and this one, before starting any new feature work**, since the honest next step is "wait for hardware," not "find more software to build."

### Done

- Repo scaffolded from template; `docs/foundation.md` and this file written.
- `backend/`: Go module, `/healthz`, `camera.Source`/`MockSource` (mock frame producer), `stream.Handler` (MJPEG `multipart/x-mixed-replace` streaming), `auth.RequireToken` (shared-secret token middleware, fails closed if `STREAM_AUTH_TOKEN` unset), a stdlib-only `.env` loader (`cmd/server/env.go` — before this, nothing in the *backend* loaded `.env` files at all; the frontend has always loaded `.env` on its own via SvelteKit/Vite). `/stream` wired in, auth-gated. Build, vet, lint, and full test suite (24 tests) verified clean.
- `frontend/`: SvelteKit app scaffolded and the real viewer built — `+page.svelte` shows the live MJPEG feed via `<img>`, with loading/connected/error states, a `/healthz` liveness poll with auto-recovery, and the shared-secret token sent as a `?token=` query param (`PUBLIC_STREAM_AUTH_TOKEN`, must match the backend's `STREAM_AUTH_TOKEN`). Lint, type-check, test, and build verified clean.
- Root tooling in place: CI workflow (`.github/workflows/ci.yml`), `Makefile`, `dependabot.yml` entries, per-service `.env.example` files, prompt `Config` blocks updated, root README's "Testing" section documents the exact commands.
- release-please cutting real releases per-component: `backend` and `frontend` both have live tags/CHANGELOGs now (first frontend release, `v0.1.0`, shipped alongside Phase 2/3 work).

### In progress

(nothing — paused pending hardware purchase, see above)

### Not started (blocked on hardware, Phase 4+)

- Hardware purchase (`docs/hardware.md`) — Phase 4a. **This is the actual blocker**, not a technical unknown.
- OS bring-up SOP, real camera capture swap, systemd deploy, static frontend serving, physical mounting, real-world drive test — Phases 4 through 8, all sequenced after purchase in `docs/plans/2026-07-e2e-build-phases.md`.
- A deployment adapter for SvelteKit (`@sveltejs/adapter-auto` is a placeholder until the deployment target — Phase 6c — is reached).
- `backend/internal/api` (control/status routes beyond `/healthz`) — not created, no concrete need for it yet.

## Open questions

1. Does this need a persistent database, and for what (settings, clip metadata)? Deferred until the streaming core works.
2. Does this integrate with the homelab's Traefik/Authentik setup, or stay fully local-network/self-contained?
3. What onboard computer/SBC will the car run? A candidate parts list (Raspberry Pi 5, Pi Camera Module 3 Wide) is in `docs/hardware.md`, but nothing is purchased/confirmed yet. The deployment mechanism itself is decided — cross-compiled binary + systemd, not Docker (ADR-006) — only the specific board remains open.
4. Is video recording/storage ever in scope, or strictly live-only? Not explicitly excluded during scoping, but not committed to either.

## Decision log

See `docs/ADRs/` for architecture decision records.

---

*Last updated: 2026-07-24 | Session: phases-0-3-complete-docs-sync — Phases 1a through 3b all shipped and merged this session (PRs #13–#25); the mock pipeline + frontend viewer + auth all work end-to-end, verified live. Confirmed no driving-controls scope creep (RF transmitter stays the only control path — permanent decision, not deferred). This update marks Phases 0–3 done across CLAUDE.md, the build-phases doc, and the root README, and adds a "Where things stand" resume note since the project is now paused on Phase 4a (hardware purchase) for budget reasons.*
