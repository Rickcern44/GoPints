## Why

Admins currently can only set a keg's beer image or brewery logo by uploading a file from disk. Many brewery logos and beer labels already exist as URLs (brewery websites, Untappd, distributor catalogs), and forcing a download-then-upload round trip is unnecessary friction. Adding URL import gives admins a faster path while keeping today's upload flow intact.

## What Changes

- Add a new admin-only endpoint pair: `POST /api/kegs/{id}/image/from-url` and `POST /api/kegs/{id}/brewery-image/from-url`, each accepting `{ "url": "..." }` and fetching the image server-side.
- The server downloads the image, validates it (size cap, timeout, no redirects, real-content MIME sniffing), and stores it exactly as today's upload path does (`SetKegImage`/`SetBreweryImage`) — no new storage model, no `image_url` column.
- Gate the new endpoints and the admin UI entry point behind a new feature flag, `feature.remote_image_urls`, following the existing `flow_based_pour` feature-flag convention (`internal/api/handlers_features.go`).
- Add a "paste a URL" option alongside the existing drag-and-drop dropzone in `admin/kegs/[id]` for both the beer image and brewery logo sections, shown only when the feature flag is enabled.
- Reject any fetch that returns a redirect response — admins must supply a direct image URL; no hop-by-hop redirect following.

## Capabilities

### New Capabilities
- `keg-image-url-import`: Server-side fetching of a remotely-hosted image (by admin-supplied URL) into the existing keg/brewery image blob storage, gated behind a feature flag, with SSRF-safe fetch validation.

### Modified Capabilities
(none — no existing specs in this repo yet; the existing upload behavior is unchanged and not being modified by this proposal)

## Impact

- `gopints/internal/api/handlers_kegs.go` — two new handlers (`handleSetKegImageFromURL`, `handleSetBreweryImageFromURL`) plus a shared SSRF-safe fetch helper.
- `gopints/internal/api/server.go` — two new routes registered, both wrapped in `s.requireAuth` and gated by the feature flag check.
- `gopints/internal/api/handlers_features.go` — add `remote_image_urls` to `knownFeatures`.
- This is the **first outbound HTTP call the server makes** (previously, all `net/http` usage was server-side listening only) — new dependency on `net/http.Client` as an outbound client, no new third-party module required.
- `web/src/routes/admin/kegs/[id]/+page.svelte` — add a URL-input mode to both image sections, conditionally shown based on `GET /api/features`.
- No changes to `KegCard.svelte` or the public-facing kiosk display — this change is admin/upload-path only.
