# ADR-005 — Frontend served by the backend

**Status**: Accepted
**Date**: 2026-07-20

## Decision

The built SvelteKit app (static output) is served by the Go backend on the onboard computer, rather than run as a separate dev/preview process on a laptop or phone.

## Rationale

Means there's one device to open on a phone browser, with nothing else that needs to be kept running elsewhere.

## Alternatives considered

- Running the frontend as its own process on a laptop/phone, pointed at the board's IP for the API only — ruled out in favor of a single self-contained device.

## Consequences

- Requires picking a concrete SvelteKit adapter (`adapter-static`) — `@sveltejs/adapter-auto` is a placeholder until Phase 6c of `docs/plans/2026-07-e2e-build-phases.md` implements this.
- The Go backend gains a responsibility (serving static files) beyond the video stream/control API.

## Open items

- Not yet implemented — tracked as Phase 6c of `docs/plans/2026-07-e2e-build-phases.md`.
