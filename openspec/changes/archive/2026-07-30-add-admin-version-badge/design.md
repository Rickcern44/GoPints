## Context

`GET /api/health` already returns `{"status": "ok", "version": s.version}` (`gopints/internal/api/handlers.go`), and `version` is already injected at Docker build time via `-ldflags -X main.version=${VERSION}`, tied to the release-please-managed monorepo version (`.release-please-manifest.json`). Nothing in the frontend calls this endpoint today. `admin/+layout.svelte` is the single shared shell every authenticated admin page renders inside, with an existing `768px` breakpoint that already distinguishes desktop sidebar layout from mobile top-bar/bottom-nav layout.

## Goals / Non-Goals

**Goals:**
- Surface the running server's build version somewhere in the admin UI, sourced from the existing `/api/health` endpoint (no new backend surface).
- Keep it small and out of the way — an indicator for admins who go looking, not a prominent UI element.

**Non-Goals:**
- Not adding a backend endpoint or changing how `version` is set — purely consuming what already exists.
- Not showing the frontend's own `package.json` version — the server binary's version is the one that matters for "which build is deployed," and is already the authoritative, release-please-tracked value.
- Not showing on mobile — bottom-right is already occupied by the mobile bottom nav bar; not worth the layout complexity for a low-priority indicator.

## Decisions

**Fetch `/api/health` from `admin/+layout.svelte`, not from `+layout.ts`'s load function.**
The endpoint is a fire-and-forget nicety, not blocking data the page needs to render — an `onMount` fetch (same pattern as the main site's `fetchFeatures()` call in `admin/kegs/[id]/+page.svelte`) keeps it off the critical render path. If it fails, the badge just doesn't show.

**Fixed-position badge, shown only inside the authenticated shell (not on `/admin/login`).**
Placed as a sibling to the existing `.admin-shell` markup, `position: fixed; bottom; right;`, reusing the sidebar's existing `@media (min-width: 768px)` breakpoint to hide it on mobile. This avoids colliding with the mobile bottom nav without adding a new breakpoint.

**Plain text, not a link or interactive element.**
`v{version}`, low-opacity, non-interactive — informational only, matching the low-emphasis treatment of other incidental chrome in this layout (e.g. `brand-sub`, `mobile-brand-sub`).

## Risks / Trade-offs

- **[Risk]** `/api/health` fails or is slow. → **Mitigation**: fetch is fire-and-forget; badge simply doesn't render if the fetch fails, no error state needed for something this low-stakes.
- **[Trade-off]** Dev builds show `v dev` (the `main.version` default) rather than a real version number when running via `go run`/`task dev` instead of the Docker-built binary — expected and fine, since the whole point is to distinguish deployed builds.
