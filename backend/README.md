# Backend

The car-side streaming server. Runs on whatever onboard computer ends up mounted on the car (SBC choice not yet made — see `docs/foundation.md`), captures the camera feed, and serves it to the frontend SPA over the local network.

## Stack

- Go 1.24+
- Standard library only (`net/http`) — no web framework. First Go project for this team; the stdlib's routing (method + path patterns in `http.ServeMux`) is enough for the handful of routes this needs, and it keeps the dependency surface at zero while learning the language. See `docs/ADRs/ADR-001-go-stdlib-backend.md` for the reasoning.

## Structure

```
backend/
  cmd/server/       # entrypoint (main.go) — HTTP server bootstrap
  internal/         # not created yet — added as real work starts:
    camera/           - camera capture, once the board/camera module is chosen
    stream/           - HTTP video streaming handlers (MJPEG/WebRTC/etc., TBD)
    api/              - control/status routes beyond /healthz
```

Only `cmd/server/main.go` exists today, with a `/healthz` endpoint — enough to build, test, and deploy a skeleton before the camera hardware decision is made. Don't add `internal/camera` or `internal/stream` speculatively; add them when there's a real board to target.

## Conventions

- Tests are colocated as `_test.go` files next to the code they cover — no separate top-level `tests/` folder.
- Format with `gofmt`; a change that isn't `gofmt`-clean shouldn't land.
- Lint with `golangci-lint run` (default rule set, no custom config yet).
- Config comes from environment variables (see `.env.example`), not flags or config files.

## Non-negotiable constraints

- Never expose the stream or control endpoints to the public internet without auth (local network or authenticated access only).
- Never add motor/steering control code here — this backend is camera + streaming only, driving the car is explicitly out of scope.

## Running locally

```
cp .env.example .env   # then edit as needed
go run ./cmd/server
```

Or via the root `Makefile`: `make dev-backend`.
