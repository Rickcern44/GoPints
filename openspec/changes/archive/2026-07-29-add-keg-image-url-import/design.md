## Context

Keg and brewery images are stored as blobs in SQLite via `pkg/tap.Store` (`SetKegImage`/`GetKegImage`/`DeleteKegImage` and the brewery equivalents), populated today only through `PUT /api/kegs/{id}/image` and `PUT /api/kegs/{id}/brewery-image` — both `requireAuth`-gated, both accepting raw bytes with a client-supplied `Content-Type` header, capped at `maxImageBytes` (10 MB). The server has no outbound HTTP client anywhere in the codebase today; every `net/http` usage is for serving requests (`ListenAndServe`), not making them. This design introduces the server's first outbound network call, which is the primary reason this is a design-worthy (not just task-list-worthy) change: it's new attack surface (SSRF) in an app that has never made an external request before.

There's an existing feature-flag convention (`internal/api/handlers_features.go`): a fixed `knownFeatures` list, flags stored as `feature.<name>` rows via the generic `GetSetting`/`SetSetting` KV store, toggled via `PUT /api/features/{name}` (admin-gated), read via `GET /api/features` (public — the frontend already calls this to check `flow_based_pour`).

## Goals / Non-Goals

**Goals:**
- Let an admin supply a URL instead of a file for either the beer image or brewery logo, with the server fetching and storing it through the exact same blob-storage path as today's upload.
- Make the fetch safe against SSRF: no requests to loopback/private/link-local/metadata addresses, no redirect-following, bounded time and size, and validation of actual image content rather than trusting declared headers.
- Gate the entire feature (both endpoints and the admin UI affordance) behind `feature.remote_image_urls`, off by default.

**Non-Goals:**
- No `image_url` column or "live" (un-downloaded) image references — every stored image, however it arrived, is still a blob served from SQLite via the existing `GET` endpoints. This preserves the current property that the kiosk display works regardless of third-party host uptime.
- No redirect support of any kind (not even single-hop) — out of scope for this iteration; admins must supply a direct image URL.
- No changes to the existing multipart/raw-byte upload path — it continues to work exactly as it does today, flag or no flag.

## Decisions

**Server-side fetch-and-store, not client-side fetch-then-upload or a stored-URL model.**
Considered three shapes (from earlier exploration): (A) browser fetches the URL and re-uploads the blob — rejected, breaks on any remote host without permissive CORS, which is most of them; (B) server fetches and stores as a blob — chosen; (C) store the URL directly and let `<img src>` point at it — rejected, makes the kiosk display's reliability depend on a third party staying up and CORS/hotlink-friendly, which regresses a property the app currently has "for free."

**Feature-flagged via the existing `knownFeatures` convention, not a new mechanism.**
Add `"remote_image_urls"` to `knownFeatures` in `handlers_features.go`. The two new routes check the flag (via `s.store.GetSetting(ctx, "feature.remote_image_urls")`) and return 404 (not just a permission error — the routes should appear not to exist when the flag is off, consistent with disabled features not being discoverable) when disabled. The admin UI reads `GET /api/features` (already fetched today for `flow_based_pour`) to decide whether to show the "paste a URL" option.

**New endpoints, not an overloaded existing endpoint.**
`POST /api/kegs/{id}/image/from-url` and `POST /api/kegs/{id}/brewery-image/from-url`, separate from the existing `PUT .../image`. Dispatch-by-Content-Type on the same route was considered and rejected — keeping upload (raw bytes, `PUT`) and URL-import (`{url}` JSON, `POST`) as distinct routes keeps each handler simple and keeps the existing upload handler completely untouched.

**SSRF-safe fetch helper, shared between both handlers.**
A single helper (e.g. `fetchRemoteImage(ctx, url string) (data []byte, mimeType string, err error)`) used by both new handlers:
1. Parse the URL; reject anything not `http`/`https`.
2. Resolve the hostname; reject if any resolved IP is loopback (`127.0.0.0/8`, `::1`), private (`10/8`, `172.16/12`, `192.168/16`, `fc00::/7`), link-local (`169.254.0.0/16`, `fe80::/10` — this includes the common cloud metadata address `169.254.169.254`), or unspecified.
3. Build an `http.Client` with `Timeout` (e.g. 8s) and `CheckRedirect: func(...) error { return http.ErrUseLastResponse }`-equivalent (return an error) so any 3xx response fails the fetch outright rather than following it.
4. Read the body through `io.LimitReader(resp.Body, maxImageBytes+1)`, reject if it exceeds `maxImageBytes` (reusing the existing constant).
5. Run `http.DetectContentType` on the downloaded bytes; reject if it doesn't sniff as a supported image type. Store the *sniffed* type, not the remote server's `Content-Type` header — this also closes a pre-existing gap where the upload handler trusts the browser's declared header without checking bytes.
6. On success, call the existing `store.SetKegImage`/`store.SetBreweryImage` exactly as the upload handlers do.

**DNS resolution is checked once, before connecting; no re-check on the actual TCP dial.**
A theoretical DNS-rebinding attack (hostname resolves safely at check-time, then to a private IP at connect-time) is a known residual risk of this approach. Accepted for this iteration given the admin-only, single-operator nature of this endpoint (see Risks below) — closing it fully would require a custom `net.Dialer.Control` hook, which adds meaningful complexity for a threat model where the caller is already a trusted authenticated admin, not an anonymous public endpoint.

## Risks / Trade-offs

- **[Risk]** DNS rebinding could bypass the private-IP check between resolution and connection. → **Mitigation**: acceptable residual risk given the endpoint is `requireAuth`-gated (only an authenticated admin can trigger a fetch at all) and flag-gated (off by default); if this app's admin trust model changes (e.g. multi-tenant admins), revisit with a dial-time IP re-check.
- **[Risk]** A slow or hanging remote server could tie up a request goroutine. → **Mitigation**: explicit client `Timeout` (8s) bounds this; no change needed to the server's overall concurrency model.
- **[Risk]** Admins may be confused when a legitimate image URL fails because the host redirects (common for CDN-backed image hosts). → **Mitigation**: surface a clear error message ("This URL redirected; please provide a direct image link") rather than a generic failure, per the earlier decision to keep redirects unsupported.
- **[Trade-off]** No caching/dedup of remote fetches — re-importing the same URL re-downloads it. Acceptable; this is an infrequent admin action, not a hot path.

## Migration Plan

No data migration — this reuses existing blob columns. Deploy as: (1) ship with `feature.remote_image_urls` unset (defaults to off, matching how `flow_based_pour` behaves today), (2) verify the fetch helper against a small set of real brewery-logo URLs in a non-production environment, (3) enable the flag via the existing admin features UI once verified. Rollback is simply disabling the flag — the new routes 404 and the admin UI hides the URL-input affordance again; no upload data is affected either way.

## Open Questions

- Should there be a per-fetch audit log entry (URL + admin + timestamp) given this is the app's first outbound-request feature? Leaning yes for a lightweight log line, but not blocking — can be added as a follow-up task without changing the shape of this design.
