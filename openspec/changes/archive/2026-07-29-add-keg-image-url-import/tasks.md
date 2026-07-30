## 1. Feature flag

- [x] 1.1 Add `"remote_image_urls"` to `knownFeatures` in `gopints/internal/api/handlers_features.go`

## 2. SSRF-safe fetch helper

- [x] 2.1 Implement `fetchRemoteImage(ctx, url string) (data []byte, mimeType string, err error)` in `gopints/internal/api/handlers_kegs.go` (or a new `handlers_image_fetch.go`)
- [x] 2.2 Reject non-`http`/`https` schemes
- [x] 2.3 Resolve hostname and reject loopback/private/link-local/unspecified addresses (`127.0.0.0/8`, `::1`, `10/8`, `172.16/12`, `192.168/16`, `169.254.0.0/16`, `fc00::/7`, `fe80::/10`)
- [x] 2.4 Configure `http.Client` with a request timeout (e.g. 8s) and a `CheckRedirect` that errors on any redirect
- [x] 2.5 Enforce `maxImageBytes` via `io.LimitReader` on the response body
- [x] 2.6 Sniff real content type via `http.DetectContentType`; reject unsupported types
- [x] 2.7 Unit test the helper: rejects private IPs, rejects redirects, rejects oversized bodies, rejects non-image content, succeeds on a valid small image (use `httptest.Server` for the success/failure cases)

## 3. API handlers and routes

- [x] 3.1 Implement `handleSetKegImageFromURL` and `handleSetBreweryImageFromURL`, each: check feature flag (404 if disabled) → parse `{url}` body → call `fetchRemoteImage` → call existing `store.SetKegImage`/`SetBreweryImage` on success → map fetch errors to appropriate HTTP status codes
- [x] 3.2 Register `POST /api/kegs/{id}/image/from-url` and `POST /api/kegs/{id}/brewery-image/from-url` in `server.go`, both wrapped in `s.requireAuth`
- [x] 3.3 Add handler tests in `handlers_test.go` using a mock store, covering: flag disabled → 404, unauthenticated → 401/403, successful import stores expected bytes/mime type, fetch failure → error response with no store call

## 4. Admin UI

- [x] 4.1 In `admin/kegs/[id]/+page.svelte`, read `remote_image_urls` from the existing features fetch and conditionally show a "paste a URL" toggle/tab next to each dropzone (beer image, brewery logo)
- [x] 4.2 Add URL input field + submit action calling the new `from-url` endpoints via `adminFetch`, reusing existing `saving`/`error`/`success` state
- [x] 4.3 On success, refresh `hasImage`/`hasBreweryImage` the same way `uploadImage`/`uploadBreweryImage` do today
- [x] 4.4 Surface fetch-specific error messages (e.g. "This URL redirected; please provide a direct image link", "Image exceeds size limit", "URL could not be reached")

## 5. Verification

- [x] 5.1 Enable the flag locally and manually import a real brewery logo URL and a real beer label URL end to end
- [x] 5.2 Manually verify a private-IP/localhost URL is rejected and a redirecting URL is rejected, with clear admin-facing error messages
- [x] 5.3 Run `go test -race ./...` in `gopints/`
- [x] 5.4 Run `npm run check` and `npm run lint` in `web/`
