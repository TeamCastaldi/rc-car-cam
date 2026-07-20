# ADR-006 — Cross-compiled binary + systemd, not Docker

**Status**: Accepted
**Date**: 2026-07-20 (reconfirmed same day after explicit review)

## Decision

The backend deploys to the onboard computer as a cross-compiled Go binary, run via systemd — not as a Docker container.

## Rationale

A single Go binary on one resource-constrained board doesn't need a container runtime's overhead. Keeps with the stdlib/zero-dependency ethos from ADR-001.

## Alternatives considered

- **Docker** — reasonable if this ever grows into a multi-service or multi-board setup, but not justified for a single binary on a single car today. Explicitly reconsidered on 2026-07-20 against homelab consistency (Traefik/Authentik/Plex etc. all run as Docker Compose there) — kept as binary + systemd, since the overhead cost on a board already doing camera/video work outweighs the tooling-consistency benefit.

## Consequences

- Deployment tooling (cross-compilation, service management) differs from the rest of the homelab's Docker Compose pattern — a deliberate, reconsidered exception, not an oversight.
- GoReleaser (ADR-007) automates the cross-compilation step this decision implies.

## Open items

- Which onboard computer/SBC this actually targets is still open — see `docs/hardware.md` and `docs/foundation.md`.
