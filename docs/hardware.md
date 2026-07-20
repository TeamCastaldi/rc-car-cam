# RC Car Cam — Hardware Parts List

**Status**: Draft — recommendation, nothing purchased/confirmed yet
**Date**: 2026-07-20

---

This is a physical build, not just software — this doc is the bill of materials for it. It complements `docs/foundation.md` (the why) and `CLAUDE.md` (the software architecture); this is the what-to-buy.

Prices below are rough, indicative bands only (`$` / `$$` / `$$$`), not quotes — verify current price and stock before buying anything. Model names are starting points for what to search for, not the only acceptable option.

## Assumptions this list is built on

- No RC car yet — a chassis is included below.
- Onboard computer: "more headroom" tier (Raspberry Pi 4/5 class, not a Pi Zero 2 W).
- Camera: CSI ribbon camera module, not a USB webcam.
- No pan/tilt mount — camera is fixed-mount for now. (A pan/tilt servo would be a camera-aiming mechanism, not vehicle control, so it wouldn't violate the "no motor/steering control" constraint — it's just left out here since it wasn't asked for. Easy to add later.)
- The car's own radio/ESC continues to drive the motor and steering exactly as normal RC — nothing in this project's software touches that, per the constraint in `CLAUDE.md`.

## 1. RC car (chassis + radio)

| Item | Example | Why |
|---|---|---|
| RTR basher-style truck, 1/10 scale, 2WD | Traxxas Stampede 2WD RTR | Long-running, extremely durable product line — good for a first hardware project where crashes while learning FPV driving are expected. Flat, roomy chassis top deck for mounting electronics. RTR means it comes pre-built with radio/ESC/motor already installed. |
| Alternative | ARRMA Granite RTR | Similarly durable basher-class truck if the Traxxas isn't available. |

Cost band: **$$** (RTR basher trucks). 2WD is fine and cheaper than 4WD — this build doesn't need extra traction, just a stable camera platform.

## 2. Onboard computer

| Item | Example | Why |
|---|---|---|
| SBC | Raspberry Pi 5 (4GB) | Matches the "more headroom" choice — Cortex-A76 CPU has plenty of margin for encoding/serving video plus running the Go backend. |
| Storage | microSD, 32–64GB, A2-rated | A2 rating matters for random I/O performance under sustained writes (logs, any future recording). |
| Cooling | Official Raspberry Pi 5 Active Cooler (or equivalent) | The Pi 5 runs hot under sustained load, more so inside a mostly-enclosed RC car body — passive-only cooling risks thermal throttling on long drives. |
| Power (see §4) | — | Pi 5 wants 5V/5A via USB-C PD — more current than a Pi 4 needs. |

**Worth knowing before buying:** the Raspberry Pi 5 dropped the dedicated hardware H.264 encoder block that the Pi 4 had — Pi 5 leans on its faster CPU instead. If you'd rather offload video encoding to dedicated silicon and don't need the extra CPU headroom, a **Raspberry Pi 4B (4GB)** is the alternative — cheaper too. Either works for MJPEG-style streaming (the simpler, lower-latency option for a local-network feed like this); H.264 only matters if bandwidth becomes a real constraint later.

Cost band: **$$** for the Pi 5 + cooler + microSD.

## 3. Camera

| Item | Example | Why |
|---|---|---|
| Camera module | Raspberry Pi Camera Module 3 — **Wide** variant | Wider field of view suits FPV driving (more situational awareness than the standard lens's narrower FOV). |
| Cable | CSI cable matching your Pi's connector | Pi 5 and Pi 4 use *different-size* CSI connectors (Pi 5 uses the smaller connector shared with Pi Zero 2 W; Pi 4 uses the older full-size one) — buy the cable length/variant for the board you actually pick, not just whatever ships in the camera box. |

Cost band: **$** (camera module + cable).

## 4. Power for the electronics

The RC car's own battery (typically 2S/3S LiPo, 7.4–11.1V) feeds the motor via the ESC — it's not clean regulated 5V for a Pi. Two ways to power the electronics, in order of how much I'd recommend starting with:

1. **Separate USB-C PD power bank** (10,000mAh+, rated for 5V/3A+ output) dedicated to the Pi. Simplest, electrically isolates the Pi from motor noise/spikes, easy to recharge independently of the car's battery. Recommended starting point.
2. **UBEC (switching 5V regulator) tapped off the main battery/ESC.** One battery for everything, but needs wiring/soldering, and a cheap UBEC can pass through enough motor electrical noise to brown out the Pi mid-drive. A reasonable upgrade once the rest of the system is working and you want to simplify charging.

Cost band: **$** (power bank) or **$** (UBEC + wiring).

## 5. Mounting and protection

- Mounting plate/tray (aluminum, acrylic, or 3D-printed) sized to the truck's chassis deck, to hold the Pi + camera + power bank as one removable unit.
- Nylon standoffs/screws and hook-and-loop (Velcro) strips for securing the tray and cables.
- A small enclosure or 3D-printed housing for the Pi board — the car will crash occasionally while you're learning FPV driving; bare boards don't love that.

Cost band: **$**.

## 6. Networking

No extra part needed here — the Pi 4/5's built-in Wi-Fi covers the local-network streaming this project is built around. One real limitation worth flagging, not solving with a part: Wi-Fi range is what actually caps how far "beyond line of sight" can go before the feed degrades. If that becomes the binding constraint later, that's a networking/architecture question (external antenna, travel router, etc.), not a parts-list item — revisit it in `docs/foundation.md`'s open questions if it comes up.

## Rough total

Roughly **$350–500** all-in at these tiers (RTR truck being the largest line item), before shipping/tax — wide range because RC truck and SBC pricing both move. Treat this as a ballpark for budgeting, not a quote.

## Open items

- Exact SBC (Pi 5 vs Pi 4) — recommendation given above, final call is yours once you see both in hand.
- Whether a UBEC power setup ever replaces the power bank (§4).
- Pan/tilt mount is out of this list by default — say the word if you want it added.
