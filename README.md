# RC Car Cam

A camera mounted on an RC car streams live FPV video to a phone/browser SPA so you can drive it beyond line of sight. It's a hands-on learning project — the goal is the build itself (embedded video, streaming, networking), not solving an external pain point. Built for personal/family use. Driving stays on the RC car's own RF transmitter — this project is camera + streaming only, never motor/steering control (see `CLAUDE.md` Constraints); that's a deliberate, permanent scope boundary, not a "not built yet."

---

## Where things stand (updated 2026-07-24)

**All the software that can be built without hardware is done.** Phases 0–3 of `docs/plans/2026-07-e2e-build-phases.md` are complete and merged: mock camera → MJPEG streaming → shared-secret auth → frontend viewer with connection states, all verified working end-to-end in a browser.

**The project is paused here, waiting on a hardware purchase (Phase 4a, `docs/hardware.md`) — not a technical blocker, a budget one.** There's nothing more to build in software until that happens — resist the urge to invent more Phase-3-adjacent work; the honest next step is "wait," and that's fine.

**When it's time to resume, the prompt is simply:**
> "I bought the hardware for RC Car Cam — let's start Phase 4a (bring-up)."

That's enough context for a fresh Claude session to pick this back up: read this README section, then `CLAUDE.md`'s Current State, then `docs/plans/2026-07-e2e-build-phases.md` starting at Phase 4, and go from there (OS install/SSH/Wi-Fi SOP, verifying the camera via the vendor's own tool, then swapping the real camera into `internal/camera` behind the same `Source` interface the mock already uses — no other layer changes).

---

## Stack

- **Backend** (`backend/`): Go 1.24+, standard library only (`net/http`, no framework) — a deliberate choice for a first Go project, see `docs/ADRs/ADR-001-go-stdlib-backend.md`
- **Frontend** (`frontend/`): SvelteKit (Svelte 5), TypeScript, scaffolded with `sv create`
- **Database**: none yet — undecided, see `docs/foundation.md` open questions
- **Test runner**: `go test` (backend), Vitest (frontend)
- **Lint/format**: `go vet` + golangci-lint + gofmt (backend); ESLint + Prettier (frontend)
- **Deployment target**: undecided — onboard computer for the car hasn't been chosen yet

## Quick Start

Prerequisites: Go 1.24+, Node 22+, [golangci-lint](https://golangci-lint.run/welcome/install/) (for `make lint`).

```sh
# Backend
cp backend/.env.example backend/.env
make dev-backend      # or: cd backend && go run ./cmd/server

# Frontend (separate terminal)
cp frontend/.env.example frontend/.env
cd frontend && npm install
make dev-frontend      # or: cd frontend && npm run dev -- --open
```

`STREAM_AUTH_TOKEN` (backend `.env`) and `PUBLIC_STREAM_AUTH_TOKEN` (frontend `.env`) must be set to the **same** value — the backend refuses to start without it, and the frontend won't be able to load the stream if they don't match. The `.env.example` defaults (`changeme`) already match each other out of the box for local dev; generate a real random value (e.g. `openssl rand -hex 32`) for anything beyond quick local testing.

Run `make` targets from the repo root, not from inside `backend/`/`frontend/` — `make dev-frontend` from inside `frontend/` will fail with "No rule to make target."

Other common commands: `make build`, `make test`, `make lint` (run both stacks).

## Testing

From the repo root, `make test` runs both stacks (`go test ./...` in `backend/`, `npm test` in `frontend/`). `make lint` similarly runs `go vet ./... && golangci-lint run` plus `npm run lint`.

To run backend tests directly:

```sh
cd backend
go test ./...                        # run all tests
go test ./... -v                     # verbose, see each test name pass/fail
go test ./internal/camera/... -v     # scope to one package
go vet ./...                         # static checks
golangci-lint run                    # linter (see install instructions below)
gofmt -l .                           # formatting check; no output means clean
```

To run frontend tests directly: `cd frontend && npm test`.

`golangci-lint` isn't packaged for a current version via `apt`; install it with `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest` (needs `$(go env GOPATH)/bin` on your `PATH`) or `sudo snap install golangci-lint --classic` (it's published as a classic-confinement snap, so `--classic` is required).

## Project Structure

```
backend/        Go streaming server (runs on the car's onboard computer)
frontend/       SvelteKit viewer SPA (live video only — no driving controls, see above)
docs/           All project documentation
scripts/        Dev-time utilities (not shipped)
```

For architecture decisions, see `docs/ADRs/`.
For technical specs, see `docs/specs/`.
For SOPs and runbooks, see `docs/SOPs/`.
For the founding brief, see `docs/foundation.md`.
For the hardware parts list, see `docs/hardware.md`.
