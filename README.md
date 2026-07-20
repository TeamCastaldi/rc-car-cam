# RC Car Cam

A camera mounted on an RC car streams live FPV video to a phone/browser SPA so you can drive it beyond line of sight. It's a hands-on learning project — the goal is the build itself (embedded video, streaming, networking), not solving an external pain point. Built for personal/family use.

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
frontend/       SvelteKit viewer SPA (video + controls)
docs/           All project documentation
scripts/        Dev-time utilities (not shipped)
```

For architecture decisions, see `docs/ADRs/`.
For technical specs, see `docs/specs/`.
For SOPs and runbooks, see `docs/SOPs/`.
For the founding brief, see `docs/foundation.md`.
For the hardware parts list, see `docs/hardware.md`.
