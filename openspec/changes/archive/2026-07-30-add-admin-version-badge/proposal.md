## Why

There's currently no way to tell which build of the server an admin is looking at from the admin UI itself — the only place the version exists today is the `version` field returned by `GET /api/health`, which nothing in the frontend calls. When debugging a deployed instance or confirming a release rolled out, an admin has to shell into the server or check Docker directly. A small, unobtrusive version indicator in the admin UI closes that gap.

## What Changes

- Add a `fetchHealth()` helper to `web/src/lib/api.ts` that calls the existing (already public, unauthenticated) `GET /api/health` endpoint.
- In `admin/+layout.svelte` (the shared shell wrapping every admin page), fetch health on mount and render the version as a small, low-emphasis badge fixed to the lower-right corner of the viewport.
- Shown on desktop widths only (matching the existing `768px` sidebar breakpoint already in this layout) — hidden on mobile, where the bottom nav already occupies that corner.
- Not shown on the `/admin/login` page (mirrors how the sidebar/nav chrome already only renders once authenticated).
- No backend changes — the version is already embedded in the server binary at build time (`-ldflags -X main.version=...` in `gopints/Dockerfile`) and already served by `/api/health`.

## Capabilities

### New Capabilities
- `admin-version-display`: Displaying the running server's version as a small badge within the admin UI shell.

### Modified Capabilities
(none — no existing requirement changes)

## Impact

- `web/src/lib/api.ts` — new `fetchHealth()` export.
- `web/src/routes/admin/+layout.svelte` — fetch health on mount, render a fixed-position version badge (desktop only).
- No Go/API changes; no new endpoints; no new dependencies.
