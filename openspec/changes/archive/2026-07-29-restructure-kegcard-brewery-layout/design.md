## Context

`KegCard.svelte` renders each tap on the kiosk carousel as a two-zone flex layout: `.info` (flex: 1 — tap badge, beer image, beer name, inline brewery row, tags) and `.stats` (fixed — level gauge, percent remaining, pours). The brewery logo today is 28px, rendered inline next to the brewery name text inside `.info`, and easy to miss. There are no responsive breakpoints anywhere in this component or its parent page (`+page.svelte`) — sizing is entirely fluid via `clamp()`, and the display only needs to support TV and laptop-sized screens (phone use is limited to the separate `/admin` surface, out of scope here).

This design covers restructuring `.inner` from two zones to three, giving the brewery its own fixed-width "plaque" that mirrors `.stats` on the opposite side.

## Goals / Non-Goals

**Goals:**
- Give the brewery logo/name a fixed-width column of its own, visually differentiated as a distinct block rather than a hairline divider.
- Keep the beer image at its current size (90px) and existing shape variants (circle/square/can) — it only repositions.
- Degrade gracefully when a keg has no brewery name/logo, without leaving dead space.
- Change nothing about the stats zone, the carousel/swipe mechanics, or the backend/API.

**Non-Goals:**
- No responsive/mobile breakpoints — TV and laptop widths only.
- No animation/motion treatment ("dynamic" here means layout structure, not transitions).
- No changes to how brewery/beer images are uploaded or stored (covered by the separate URL-image-input change).

## Decisions

**Three-zone "bookend" flex layout over a two-column or grid restructure.**
`.inner` becomes `[brewery-plaque] [info] [stats]`, all direct flex children, matching the existing flex-based approach rather than switching to CSS grid. Alternative considered: a grid with named template areas — rejected because the current divider/flex-gap approach already handles the fixed/flex/fixed pattern (stats is already `flex-shrink: 0`), so the brewery plaque can reuse the same mechanism (`flex-shrink: 0` + fixed width) without introducing a second layout system into one component.

**Brewery plaque is stacked (logo above name), not logo-only.**
Confirmed with stakeholder: keeping the brewery name as readable text (not just a logo) matters since not every brewery has a recognizable mark. Logo sized 110-140px, name below it in smaller type, both centered in the column.

**Plaque gets a distinct block treatment, not a plain divider.**
The beer-info/stats seam is currently a subtle 1px gradient divider. The brewery plaque is differentiated by giving it its own background tint/border (consistent with the card's existing warm/amber palette) so it reads as a distinct "block" rather than another flowing column — directly addressing the ask that the brewery "take some of the focal point."

**Conditional zone, not conditional visibility.**
When `keg.brewery` is empty, the plaque column is omitted from the flex layout entirely (not just hidden), so `.info` and `.stats` reflow to fill the space — same pattern already used for `hasImage`/`hasBreweryImage` conditionals in the existing template.

**No new state or props.**
This is a template/CSS-only change. `KegCard` continues to receive the same `tap`/`stats` props; `brewery_image_mime_type` and `brewery` are already present on `Keg`.

## Risks / Trade-offs

- **[Risk]** A fixed-width plaque plus a fixed-width stats column could squeeze `.info` uncomfortably at the narrower end of expected widths (smaller laptops). → **Mitigation**: keep plaque width toward the lower end of the 110-140px logo range with modest padding, and verify at a representative laptop width (e.g. 1366px) during implementation, not just at TV widths.
- **[Risk]** Long brewery names could overflow or wrap awkwardly in a narrow fixed column. → **Mitigation**: apply the same `word-break`/line-clamp treatment already used for `.name`, sized down for the plaque.
- **[Trade-off]** Losing the inline brewery byline next to the beer name (it now lives only in the plaque) changes the reading order slightly — brewery identity is now spatially separated from the beer name it used to sit directly beneath. Accepted as the intended effect of making brewery a distinct focal element.

## Migration Plan

Single-component change, no data migration. Land as one PR: restructure `KegCard.svelte` template + styles. No feature flag needed since this is pure presentation with no backend dependency — verify visually against live keg data (with and without brewery images) before merge.

## Open Questions

None outstanding — plaque orientation (stacked), sizing (110-140px), and responsive scope (none needed) were settled during exploration.
