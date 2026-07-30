## Context

`GET /api/health` returns `{"status": "ok", "version": s.version}`, already consumed by the admin UI (`fetchHealth()` in `web/src/lib/api.ts`, added for the admin version badge) to show a passive, non-interactive version indicator. The kiosk display (`web/src/routes/+page.svelte`) is a different situation: it's a long-lived, unattended browser tab on a physical touchscreen, with no admin present to notice a stale build or manually reload. The frontend is also a PWA (`vite-plugin-pwa`, `vite.config.ts`) with a workbox `NetworkFirst` runtime-caching rule for `/api/` (5s network timeout before falling back to cache) — relevant here because a version-check fetch needs to reliably hit the network when the network is actually up, not silently serve a stale cached response.

The kiosk page already has an established pattern for small, low-emphasis fixed-position UI: the `admin-corner` link (bottom-right, `position: absolute`, low-opacity icon) and the `dot-nav`/`pour-anchor` (bottom-center). This gives a clear, uncontested spot for a new corner element: bottom-left.

## Goals / Non-Goals

**Goals:**
- Detect when the server's deployed version changes while the kiosk page is already open, without any backend changes.
- Give the display a way to self-recover from a stale build via a deliberate tap, not an automatic reload (a screen doing a background pour or showing a live pour animation shouldn't reload out from under whoever's standing at the tap).
- Keep the polling cheap and infrequent — this runs indefinitely on a device nobody is actively watching.

**Non-Goals:**
- Not adding this to the admin UI — it already shows its version (from the prior `add-admin-version-badge` change) and is reloaded manually by whoever's actively managing it there.
- Not auto-reloading on version change — the tap is the whole point; an unattended auto-reload could interrupt an in-progress pour animation or banner.
- Not using the PWA service worker's own update-detection APIs (`registration.waiting`, `updatefound`) — those detect a new *service worker* script, which is a different (and more complex) signal than "the server reports a different version." The existing `/api/health` version is simpler, already wired up, and directly reflects what actually changed (the deployed backend build, which in this monorepo's release process ships alongside the frontend build it's paired with).

## Decisions

**Poll `/api/health` on a `setInterval`, comparing against the version captured at initial load.**
On mount, fetch health once and store `initialVersion`. Every 5 minutes thereafter, fetch again and compare; if different (and non-empty), set `updateAvailable = true`. Alternative considered: `setTimeout` recursion instead of `setInterval` — no meaningful difference here since the fetch is cheap and infrequent; `setInterval` is simpler and the existing codebase doesn't already have a recursive-poll pattern to mirror.

**5-minute interval.**
This is background housekeeping on a screen nobody is watching for the sole purpose of an eventual update — sub-minute polling would just be extra unnecessary load for essentially zero user-facing benefit (nobody needs a new build applied within seconds of deploy on an ambient kiosk display). Matches the "this is for eventual detection, not real-time" framing from the proposal.

**Explicit `cache: 'no-store'` on the health fetch.**
Belt-and-suspenders against the PWA's `NetworkFirst` workbox rule for `/api/` — `NetworkFirst` already tries the network first when online (which is the case that matters here), but adding `no-store` removes any ambiguity about whether a version-check fetch could ever resolve from a cached response while online.

**Tapping the badge does `location.reload()` — a hard reload, not a client-side route change.**
A real deploy ships new JS/CSS bundles, not just new data; a SvelteKit client-side navigation wouldn't pick those up. `location.reload()` forces the browser to refetch the document and its assets.

**Badge placed bottom-left, mirroring the existing bottom-right `admin-corner`.**
Keeps it out of the way of the bottom-center `dot-nav`/`pour-anchor` controls and out of the top banner-stack region, consistent with how the existing kiosk chrome is already laid out at the corners.

## Risks / Trade-offs

- **[Risk]** The kiosk stays open across a version bump that then rolls back (e.g. bad deploy reverted) — the badge would appear, then technically still be "valid" to tap since reloading just picks up whatever's live at tap-time, which is the desired forgiving behavior regardless of the specific version history. No special handling needed.
- **[Trade-off]** Up to a 5-minute detection lag after a deploy. Accepted per the design goal — this isn't a real-time signal.
- **[Trade-off]** If the kiosk loses network entirely, polling fetches will fail silently (same pattern as `fetchHealth()` already swallowing errors) and the badge simply never appears — acceptable, matches how the admin version badge already handles fetch failure.
