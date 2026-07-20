# ADR-007 — release-please + GoReleaser for versioning

**Status**: Accepted
**Date**: 2026-07-20

## Decision

Versioning and releases are automated with release-please (Conventional Commits → release PRs, tags, `CHANGELOG.md`) and GoReleaser (cross-compiles and publishes the backend binary on each `backend/vX.Y.Z` tag).

## Rationale

release-please is already used on other projects, and Conventional Commits are already required by `CONTRIBUTING.md`. `backend/` and `frontend/` are configured as independent packages since they version on different schedules. Go modules carry no in-repo version field (unlike `package.json`) — the version *is* the git tag — and since `backend/` is a module in a subdirectory rather than at the repo root, its tags are prefixed per Go's module-versioning spec: `backend/vX.Y.Z`, not bare `vX.Y.Z`. GoReleaser triggers off those tags to cross-compile the binary (linux/arm64 + linux/arm) and attach it to the GitHub Release, automating what Phase 6a of `docs/plans/2026-07-e2e-build-phases.md` would otherwise be a manual `go build` step.

## Alternatives considered

- Hand-rolled version bumping / manual tagging — ruled out in favor of automation already familiar from other projects.
- GoReleaser's own changelog generation — disabled, since release-please already owns the changelog.

## Consequences

- GoReleaser has no native config field for a subdirectory module's prefixed tags; the GitHub Actions workflow strips the `backend/` prefix into the `GORELEASER_CURRENT_TAG` environment variable before invoking it.
- Verified locally: `goreleaser check` validates the config, and a `--snapshot` build produced real, correctly cross-compiled ELF binaries for both linux/arm64 and linux/arm.
