## 1. API helper

- [x] 1.1 Add `Health` interface and `fetchHealth()` to `web/src/lib/api.ts`, calling `GET /api/health`

## 2. Admin shell badge

- [x] 2.1 In `admin/+layout.svelte`, call `fetchHealth()` on mount (only when authenticated) and store the version in state
- [x] 2.2 Render a fixed-position `v{version}` badge in the lower-right corner, only when a version is loaded
- [x] 2.3 Hide the badge below the existing 768px desktop breakpoint and on `/admin/login`

## 3. Verification

- [x] 3.1 Run the app locally, confirm the badge appears on desktop width admin pages and shows `vdev` (the local `go run` default)
- [x] 3.2 Confirm the badge is hidden on mobile width and on `/admin/login`
- [x] 3.3 Run `npm run check` and `npm run lint` in `web/`
