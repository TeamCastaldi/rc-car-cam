# RC Car Cam — End-to-End Build Phases

Drafted 2026-07-20. Sequenced software-first: Phases 1–3 need no hardware at all, so parts can be ordered in parallel with that work (shipping lead time isn't on the critical path). Each sub-phase is scoped to be one branch/PR — small enough that "does this belong in this branch?" has an obvious answer.

Decisions baked into this plan (see `CLAUDE.md` decision log for the ADR entries):
- **Streaming protocol: MJPEG over HTTP** (`multipart/x-mixed-replace`) — stdlib-only, no third-party dependency, plenty for a local-network feed.
- **Frontend deploy target: served by the onboard computer** — SvelteKit built static, served by the Go backend. One device to open on your phone; nothing to run on a laptop.
- **Deployment mechanism: cross-compiled binary + systemd, not Docker** — this is a call made while finalizing this plan, not one of the three you were asked about. Flagging it explicitly: a single Go binary on a single resource-constrained board doesn't need a container runtime's overhead, and it keeps with the stdlib/zero-dependency ethos already set in ADR-001. Easy to revisit later if you'd rather containerize — say so and this plan (and CLAUDE.md's open question #3) updates.
- **Branch granularity: one branch/PR per sub-phase.**

---

## Phase 0 — Foundations *(done)*
Repo scaffold, CI, `/healthz`, default frontend page, docs. No action needed.

## Phase 1 — Mock video pipeline *(no hardware)*
| Sub-phase | Branch | Work | Exit criteria |
|---|---|---|---|
| 1a | `feature/mock_camera_source` | Define a small `Source` interface (the seam between "frame producer" and "HTTP handler") + a mock implementation (looping test image). | Interface + mock compile, have unit tests, nothing else wired up yet. |
| 1b | `feature/mjpeg_stream_handler` | Implement the `internal/stream` HTTP handler: consumes a `Source`, writes an MJPEG (`multipart/x-mixed-replace`) response. | Handler unit-tested against the mock `Source`. |
| 1c | `feature/wire_mock_stream_route` | Wire a `/stream` route into `main.go` using the mock source. | `curl`/browser hitting `/stream` shows the looping mock feed. |

**Not in scope for Phase 1:** real camera code, auth, frontend changes, anything hardware.

## Phase 2 — Frontend viewer *(no hardware)*
| Sub-phase | Branch | Work | Exit criteria |
|---|---|---|---|
| 2a | `feature/frontend_viewer_page` | Replace the default SvelteKit page with a minimal viewer (`<img>` pointed at `PUBLIC_API_BASE_URL` + `/stream` — MJPEG streams natively via a plain `<img src>`, no `<video>`/JS decoding needed). | Opening the frontend shows the live mock feed end-to-end. |
| 2b | `feature/viewer_connection_states` | Basic loading/error/connected states. | Killing the backend shows an error state, not a broken image icon. |

**Not in scope:** driving controls, visual polish, auth UI.

## Phase 3 — Auth *(no hardware)*
| Sub-phase | Branch | Work | Exit criteria |
|---|---|---|---|
| 3a | `feature/stream_auth_middleware` | Shared-secret token middleware on `/stream` (and any control routes), reading `STREAM_AUTH_TOKEN` (already stubbed in `backend/.env.example` — uncomment it here). | Request without token fails (401/403); with it, succeeds. Both paths tested. |
| 3b | `feature/frontend_auth_token` | Frontend sends the token on stream requests. | Viewer works end-to-end with auth enabled. |

**Not in scope:** accounts, multi-user, Authentik/Traefik integration (separate deferred decision — open question #2).

---

*Parts can be ordered any time during Phases 1–3.*

## Phase 4 — Hardware bring-up *(the purchase + first boot, no custom code)*
| Sub-phase | Branch | Work | Exit criteria |
|---|---|---|---|
| 4a | — (purchase, not a branch) | Buy per `docs/hardware.md`. | Parts in hand. |
| 4b | `docs/onboard_os_bringup_sop` | OS install, SSH, Wi-Fi join, hostname — written up as a SOP. | New `docs/SOPs/SOP-onboard-computer-setup.md`; can SSH into the board. |
| 4c | (fold into 4b's SOP) | Verify the camera works via the vendor's own CLI tool (e.g. `rpicam-hello`), before any of our code touches it. | Vendor tool shows camera output — proves the hardware itself is good. |

**Not in scope:** our Go code, mounting to the car.

## Phase 5 — Real camera integration
| Sub-phase | Branch | Work | Exit criteria |
|---|---|---|---|
| 5a | `feature/real_camera_source` | Implement the real `Source` in `internal/camera` (shape depends on what 4c found — likely shelling out to a vendor capture tool piping MJPEG frames). | Unit-testable in isolation (fake the subprocess/reader). |
| 5b | `feature/swap_real_camera_source` | Swap `main.go` to use the real source (decide then: env-flag toggle to keep the mock for laptop dev, or board-only). | Running on the board serves real live video over `/stream`. |

**Not in scope:** deployment automation, mounting.

## Phase 6 — Deploy to the board
| Sub-phase | Branch | Work | Exit criteria |
|---|---|---|---|
| 6a | `feature/cross_compile_build` | Cross-compile the backend for the board's arch. | A real binary for that arch, built and inspected locally. |
| 6b | `feature/systemd_deploy` | systemd unit + a deploy SOP. | Reboot the board, backend comes up unattended. |
| 6c | `feature/serve_frontend_static` | SvelteKit `adapter-static` build, served by the Go backend. | Opening the board's IP on a phone shows the full app, nothing running on a laptop. |

**Not in scope:** CI/CD automation of this deploy — a single car doesn't need a pipeline, the SOP is the right amount of tooling.

## Phase 7 — Physical integration *(assembly only, no code)*
| Sub-phase | Work | Exit criteria |
|---|---|---|
| 7a | Mount per `docs/hardware.md` §5. | Rig physically secured to the car. |
| 7b | Power/heat validation. | Session-length battery life confirmed, no thermal throttling in the enclosed body. |

**Resist the urge to "fix" software here** — problems in this phase are almost always assembly/power, not code, and vice versa in earlier phases.

## Phase 8 — Real-world drive test
Drive it with no direct line of sight, judge it against the success metric in `docs/foundation.md`. Not a build phase — if it doesn't feel usable, the likely next step is a latency/protocol-tuning phase (resolution, framerate, bitrate), scoped separately, not features bolted on here.

---

## Deliberately not phased yet
- Video recording/storage — open question #4, only gets a phase once decided.
- Homelab Traefik/Authentik integration — open question #2, same treatment.
- Latency tuning — only becomes real work if Phase 8 shows it's needed.
