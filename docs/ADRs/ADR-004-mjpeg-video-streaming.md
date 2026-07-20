# ADR-004 — MJPEG over HTTP for the video stream

**Status**: Accepted
**Date**: 2026-07-20

## Decision

The camera → viewer video stream uses MJPEG over HTTP (`multipart/x-mixed-replace`), not WebRTC.

## Rationale

MJPEG is implementable entirely in the Go standard library (consistent with ADR-001) and renders natively in a browser `<img>` tag — no client-side decoding library needed on the frontend either.

## Alternatives considered

- **WebRTC** — the lowest-latency option of the two, but ruled out: it needs a third-party Go library (breaking the stdlib-only decision) plus signaling/ICE machinery that's disproportionate for a single local-network viewer.

## Consequences

- Higher latency and bandwidth use than WebRTC, but acceptable for a local-network feed — see the success metric in `docs/foundation.md`.
- If the real-world drive test (Phase 8 of `docs/plans/2026-07-e2e-build-phases.md`) shows MJPEG isn't good enough, the next step is a latency-tuning phase (resolution/framerate/bitrate), not necessarily a protocol change.

## Open items

None currently — revisit only if Phase 8 (real-world drive test) shows latency is a problem.
