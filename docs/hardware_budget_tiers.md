# RC Car Cam — Hardware Budget & Tiers

**Status**: Draft — Recommendations and Cost Tiers
**Date**: 2026-07-24

---

This document breaks down the required hardware for the RC Car FPV/remote driving project into three budget tiers: **Budget** (barely working), **Standard** (working comfortably with room for upgrades), and **Ultra** (no $ cap, go nuts). 

This serves as the bill of materials for the physical build, complementing `docs/foundation.md` (the why) and `CLAUDE.md` (the software architecture).

*Note: Prices below are estimated ballpark figures (USD) as of mid-2026. Verify current stock and pricing before purchasing.*

---

## 🟢 Tier 1: Budget (Barely Working)
**Goal:** Get a moving vehicle with a live feed onto the network for the absolute lowest cost. Expect low framerates, reduced range, and improvised mounting.

| Category | Item | Est. Cost | Why this was chosen |
| :--- | :--- | :--- | :--- |
| **RC Car** | WLtoys 104001 (1/10 Scale Buggy) | ~$120 | One of the cheapest 1/10 scale "hobby-grade" RC cars available. It provides enough physical space to mount electronics, unlike smaller 1/16/1/18 scales, but keeps costs rock-bottom. |
| **SBC** | Raspberry Pi Zero 2 W | ~$15 | Extremely cheap and low-power. It barely meets the requirement for serving a basic MJPEG feed, but has very little CPU headroom for anything else (recording, heavy Go processing). |
| **Camera** | Generic Pi Camera Module V1/V2 clone | ~$15 | Basic 5MP/8MP fixed-focus camera. It gets the job done for cheap, though the field of view (FOV) is relatively narrow for FPV driving. |
| **Power** | 5V/3A UBEC (Battery Tap) | ~$10 | Tapping directly into the RC car's battery via a UBEC is much cheaper than buying a dedicated power bank. *Warning:* Motor electrical noise may occasionally brown out the Pi. |
| **Mounting** | Zip ties, double-sided tape, cardboard | ~$5 | No custom brackets. Just tape the Pi to the chassis and zip-tie the camera to the front shock tower. |
| **Total** | **~ $165** | | |

---

## 🟡 Tier 2: Standard (Working Comfortably)
**Goal:** The recommended baseline. Solid performance, smooth video streaming, easy to assemble, and plenty of headroom for future software upgrades.

| Category | Item | Est. Cost | Why this was chosen |
| :--- | :--- | :--- | :--- |
| **RC Car** | Traxxas Stampede 2WD RTR | ~$250 | Long-running, extremely durable "basher." Excellent for a first hardware project where crashes are expected. Has a flat, roomy chassis top deck for mounting electronics. |
| **SBC** | Raspberry Pi 5 (4GB) + Active Cooler | ~$120 | The Cortex-A76 CPU has plenty of margin for encoding/serving video plus running the Go backend. Active cooler prevents thermal throttling inside an enclosed RC car body. |
| **Camera** | Raspberry Pi Camera Module 3 (Wide) | ~$38 | The "Wide" variant provides a 120° FOV, massively improving situational awareness while FPV driving. High quality with autofocus. |
| **Power** | 10,000mAh USB-C PD Power Bank | ~$35 | Clean, isolated power. Electrically isolates the Pi from the RC motor's noise/voltage spikes, and easy to recharge independently. |
| **Mounting** | 3D-printed/Acrylic Tray & Velcro | ~$15 | Secures the Pi, camera, and power bank as a single removable module. Protects boards from basic impacts. |
| **Total** | **~ $458** | | |

---

## 🔴 Tier 3: Ultra (No Cap, Go Nuts)
**Goal:** The ultimate FPV exploration rig. Maximum stability, flawless high-resolution onboard recording alongside streaming, and top-tier build quality.

| Category | Item | Est. Cost | Why this was chosen |
| :--- | :--- | :--- | :--- |
| **RC Car** | Traxxas TRX-4 (e.g., Ford Bronco) | ~$500 | A high-end rock crawler rather than a basher. Crawlers are inherently slower, incredibly stable, and have massive suspension travel—resulting in a buttery smooth camera feed. Portal axles mean it rarely gets stuck. |
| **SBC** | Pi 5 (8GB) + NVMe Base + 256GB SSD | ~$180 | Maxed-out Pi. The NVMe SSD provides insanely fast I/O for saving 4K raw video locally *while* streaming the live feed, without hitting SD card bottlenecks. |
| **Camera** | Pi High Quality (HQ) Cam + Wide CS Lens | ~$100 | Massive image sensor compared to standard modules. Provides drastically better low-light performance, dynamic range, and perfectly crisp optics for FPV. |
| **Power** | Premium 100W+ USB-C PD Power Bank | ~$100 | Over-spec'd power bank (e.g., Anker Prime) ensures absolutely perfectly stable 5V/5A power under extreme load, plus hours of run-time. |
| **Mounting** | Custom CNC Aluminum Deck + Dampeners | ~$70 | Custom-machined mounting plate with rubber anti-vibration dampeners (like drone gimbals) to eliminate rolling shutter "jello" effect in the video feed. |
| **Total** | **~ $950** | | |

---

## General Open Items / Next Steps
- **Networking Limitation:** Regardless of tier, range is limited by standard Wi-Fi. If you plan to drive out of your driveway/yard, you will eventually need a travel router or external high-gain antennas. 
- **Exact SBC Call:** A Pi 4 (4GB) could replace the Pi 5 in the Standard tier if you want to rely on the Pi 4's dedicated hardware H.264 encoder instead of the Pi 5's brute-force CPU approach.
- **Pan/Tilt Camera:** None of the tiers currently feature a pan/tilt servo mechanism for the camera. This can be added to the Standard or Ultra tier later for around $25-$50.
