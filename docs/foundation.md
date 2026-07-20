# RC Car Cam — Foundation
**Status**: Draft v0.1
**Date**: 2026-07-20

---

## The Problem

This is fundamentally a hands-on learning project rather than a response to an external pain point — the goal is real, practical experience across embedded video capture, live streaming, and networking, using an RC car as the vehicle for that learning. The problem it solves along the way is the ordinary limitation of line-of-sight RC driving: once the car is far away, behind an obstacle, or around a corner, you lose visibility and lose control confidence with it.

## The Solution

A camera mounted on the RC car captures video and feeds it to a small Go backend running on the car's onboard computer. That backend serves the live feed over HTTP to a SvelteKit single-page app, viewed in a phone or laptop browser on the same local network, so the car can be driven FPV (first-person view) beyond line of sight. The system is intentionally single-purpose — video in, video out — with no autonomy or control layer built on top of it.

## The User

Personal/family use only — the builder and their household. Not a commercial product, not built for other users or a public audience.

## What We Are Not Building

- No autonomous or computer-vision features — no object detection, no line-following, no self-driving. Manual FPV control only.
- No motor or steering control code — the RC car's drive electronics are entirely out of scope; this project is the camera + streaming layer only.
- No multi-user or cloud product — single car, single user/family, no accounts, no cloud backend.
- Video recording/storage was not explicitly ruled out during scoping — see Open Questions.

## Success Metric

You can drive the RC car using only the live video feed on a phone or laptop browser — with no direct line of sight to the car — and it feels usable: latency low enough that you're not crashing constantly because of lag.

## Open Questions

1. Does this need a persistent database, and for what (settings, clip metadata)? Deferred until the streaming core works.
2. Does this integrate with the homelab's Traefik/Authentik setup, or stay fully local-network/self-contained?
3. What onboard computer/SBC will the car run, and how does the backend get deployed to it (cross-compiled binary + systemd vs. Docker)?
4. Is video recording/storage ever in scope, or strictly live-only? Not explicitly excluded during scoping, but not committed to either.

---
*This document is the source of truth for product intent. Architecture and technology decisions live in `docs/ADRs/`; this file is about why, not how.*
