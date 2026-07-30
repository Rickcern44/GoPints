## 1. Version-change detection

- [ ] 1.1 In `web/src/routes/+page.svelte`, fetch health via `fetchHealth()` on mount and store the version as a baseline
- [ ] 1.2 Add a `setInterval` (5 minutes) that re-fetches health and compares against the baseline
- [ ] 1.3 Use `cache: 'no-store'` on the health fetch (or an equivalent option on `fetchHealth()`) to avoid the PWA's `NetworkFirst` `/api/` caching rule masking a version change
- [ ] 1.4 Clear the interval on component teardown

## 2. Update-available badge

- [ ] 2.1 Add `updateAvailable` state, set to true when a poll returns a version different from the baseline
- [ ] 2.2 Render a small tappable badge in the bottom-left corner (mirroring the existing bottom-right `admin-corner` styling/positioning pattern) when `updateAvailable` is true
- [ ] 2.3 Wire the badge's click handler to `location.reload()`

## 3. Verification

- [ ] 3.1 Run the app locally; confirm no badge appears when the version hasn't changed
- [ ] 3.2 Simulate a version change (e.g. restart the local server with a different `-ldflags -X main.version=...` value, or temporarily shorten the poll interval) and confirm the badge appears and tapping it reloads the page
- [ ] 3.3 Confirm the badge doesn't overlap the dot-nav, pour-fab, or admin-corner elements
- [ ] 3.4 Run `npm run check` and `npm run lint` in `web/`
