## Why

The brewery identity on the tap display is currently an afterthought — a 28px logo squeezed inline next to the brewery name inside the beer-info column. On a kiosk display meant to run unattended on a TV, the brand that supplied the beer deserves equal visual weight to the beer itself and to the pour stats. The current layout also has no room to grow the brewery treatment without cannibalizing the beer name/tags it currently shares space with.

## What Changes

- Restructure `KegCard` from a two-zone layout (`.info` | `.stats`) to a three-zone "bookend" layout: a new fixed-width brewery plaque (left) | beer info (middle, flex) | stats (right, unchanged).
- Brewery plaque is a stacked layout: large logo (110-140px) above the brewery name, visually differentiated as its own block (background/border treatment) rather than a plain hairline divider like the existing beer/stats seam.
- Remove the existing inline `.brewery-row` (28px logo + name) from the beer-info column — that content now lives exclusively in the new plaque.
- Beer image keeps its current size (90px) and shape variants (circle/square/can); it only repositions within the beer-info column now that the brewery row has moved out.
- No responsive breakpoints are introduced — the display targets TVs and laptop-sized screens only (phone use is limited to the separate `/admin` surface); existing fluid `clamp()` sizing continues to apply.
- Cards with no brewery name/logo (`keg.brewery` empty) render without the plaque zone, collapsing back to a two-zone layout so empty space isn't reserved for nothing.

## Capabilities

### New Capabilities
- `keg-card-display`: Presentation and layout of the per-tap `KegCard` shown on the kiosk carousel — the brewery plaque, beer info, and stats zones, and how they resize/collapse based on available keg data.

### Modified Capabilities
(none — no existing specs in this repo yet)

## Impact

- `web/src/lib/components/KegCard.svelte` — template restructure (new brewery plaque markup, removal of `.brewery-row` from `.info`) and corresponding style changes (new fixed-width column, block treatment, updated flex layout for `.inner`).
- No backend or API changes — this is presentation-only, consuming the same `Tap`/`Keg` fields (`brewery`, `brewery_image_mime_type`) already returned by the existing `/api/kegs/{id}/brewery-image` endpoint.
- No changes to `admin/kegs/[id]` upload UI (out of scope for this change; covered separately by the URL-image-input change).
