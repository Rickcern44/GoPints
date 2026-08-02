## Why

The kiosk display had a dark, copper-and-oak "late-night brewpub" identity (italic serif beer names, warm amber accents), but that identity stopped at the admin sidebar — every admin content page (kegs, taps, pours, banner, features, login) fell back to generic default-Tailwind blue-on-white forms and tables, copy-pasted per page, and `LevelGauge` used stock Tailwind green/amber/red utility colors unrelated to the brand. There was also no shared token layer — every component hardcoded its own `rgba(200, 130, 26, …)` literals. An initial re-skin pass unified the admin onto the kiosk's warm palette, but a follow-up design review concluded the two-toned "dark kiosk / light admin" split, and the brewpub-ledger aesthetic itself, didn't read as a deliberate design so much as decoration on top of an unchanged scaffold. This change replaces both with one considered identity, end to end: kiosk and admin as a single **Precision Tap Console** — a dark, technical instrument-panel aesthetic (matching what the app actually is: a kegerator telemetry/monitoring tool).

## What Changes

- Replace the token layer in `web/src/routes/layout.css`: cool graphite void/panel/border surfaces (not warm oak), copper (`#c8821a`) retained but used sparingly as a signal color, plus a phosphor-green/copper/red telemetry palette (`--color-signal-good/warn/critical`) for status. Fonts switch to Chakra Petch (geometric display headers) and JetBrains Mono (all data/readouts) — Playfair Display, Bebas Neue, and Special Elite are dropped entirely.
- Unify kiosk and admin onto one dark surface — the admin content area (previously a light "paper" surface) now matches the kiosk's dark instrument-panel background; there is no light admin theme.
- Rebuild `LevelGauge` as a circular donut gauge (conic-gradient ring with the percentage centered inside) replacing the vertical fill bar.
- Restyle `KegCard`/`+page.svelte` (kiosk): drop the italic-serif "menu" treatment for the beer name in favor of the geometric display face, add a status LED to the tap badge, sharpen chrome (pour flyout, dot-nav, admin-corner) to rectangular/technical shapes, add a faint blueprint-grid background texture.
- Rebuild every admin content page (`kegs`, `kegs/[id]`, `taps`, `pours`, `banner`, `features`, `login`) around a shared console component set: rectangular connector-plate ID readouts (replacing table row numbers), a system-log list for pour history (replacing a plain table), an illuminated rocker switch for feature toggles, and a "system access" login panel.
- Retint `BannerStack` (kiosk overlay banners) off stock Tailwind amber/orange/indigo onto the signal-color palette.
- **BREAKING**: none — purely visual/presentational; no route, API, or data shape changes.

## Capabilities

### New Capabilities
- `design-tokens`: A shared CSS custom-property layer (color, type) defined once in `layout.css` and consumed by both the kiosk display and the admin UI, replacing hardcoded per-component color literals.
- `admin-visual-theme`: Admin content pages (forms, tables/lists, buttons, status states) render as part of one unified dark instrument-panel product with the kiosk display, instead of a separate light admin theme or generic Tailwind defaults.

### Modified Capabilities
(none — `keg-card-display`'s layout/behavior requirements are unchanged; this change only affects the visual treatment components render with, not the zones, sizing, or breakpoint rules already specified)

## Impact

- `web/src/routes/layout.css` — replaced token definitions and shared "console" component classes (`.panel`, `.data-table`, `.field`, `.log-row`, `.btn-console`, `.readout-badge`, `.led`, etc.).
- `web/src/lib/components/LevelGauge.svelte` — rebuilt as a circular gauge.
- `web/src/lib/components/KegCard.svelte`, `web/src/routes/+page.svelte`, `web/src/lib/components/BannerStack.svelte` — retinted/restyled to the console palette.
- `web/src/routes/admin/+layout.svelte`, `admin/kegs/+page.svelte`, `admin/kegs/[id]/+page.svelte`, `admin/taps/+page.svelte`, `admin/pours/+page.svelte`, `admin/banner/+page.svelte`, `admin/features/+page.svelte`, `admin/login/+page.svelte` — rebuilt around the shared console component set; admin content surface moved from light to dark.
- No Go/API/database changes; no new runtime dependencies (two Google Fonts swapped for two others).
