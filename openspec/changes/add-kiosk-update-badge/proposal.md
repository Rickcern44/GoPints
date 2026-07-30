## Why

The kiosk display (`/`) runs unattended on a headless touchscreen mounted near the kegerator — nobody is there to manually refresh the browser tab after a new version is deployed. Without some form of self-detection, a deploy silently has no effect until someone physically walks over and reloads the screen. A small, tappable "update available" badge that appears when the server's version changes gives the display a way to pick up new builds without needing a person to intervene, while keeping the reload itself a deliberate, admin-controlled tap rather than an automatic mid-use refresh.

## What Changes

- Add version-change detection to the kiosk display: capture the server's version at initial load (via the existing `GET /api/health`, already used for the admin version badge), then poll periodically for it to change.
- When a version change is detected, show a small tappable badge on the kiosk display. Tapping it does a full page reload (`location.reload()`) to pick up the new build.
- Placed in the bottom-left corner of the kiosk stage, mirroring the existing bottom-right admin-corner icon, and clear of the bottom-center dot-nav/pour-fab controls.
- Poll on a long interval (every 5 minutes) — this is for eventual detection, not real-time, and the display may sit on the same page for days between deploys.
- Scoped to the kiosk display only — the admin UI already shows its version passively (from a prior change) and is manually reloaded by whoever is actively managing it; no polling/badge needed there.
- No backend changes — reuses the existing, already-public `GET /api/health` endpoint.

## Capabilities

### New Capabilities
- `kiosk-update-notification`: Detecting a server version change on the kiosk display and prompting the admin to reload via a tappable badge.

### Modified Capabilities
(none — no existing requirement changes)

## Impact

- `web/src/routes/+page.svelte` — poll `/api/health` on an interval, compare against the version captured at load, render a tappable badge on change.
- Reuses `fetchHealth()` already added to `web/src/lib/api.ts` for the admin version badge — no new API helper needed.
- No Go/API changes; no new endpoints; no new dependencies.
