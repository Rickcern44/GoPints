## 1. Design tokens

- [x] 1.1 Define the `@theme` token block in `web/src/routes/layout.css`: color tokens (surfaces, accent + accent variants, cream/cream-muted text, warm-neutral light surface, error/warning tiers for `LevelGauge`), font-family tokens, radius scale, spacing scale.
- [x] 1.2 Move the Playfair Display (italic) and Bebas Neue `@import url(...)` statements out of `KegCard.svelte` and `admin/+layout.svelte` into the single global import in `layout.css`.
- [x] 1.3 Confirm `npm run build` and `npm run check` still pass with the new `@theme` block in place (no consumers changed yet).

## 2. Shared/kiosk components

- [x] 2.1 Recolor `LevelGauge.svelte`'s healthy/low/critical fill from stock `bg-green-500`/`bg-amber-400`/`bg-red-500` to token-derived colors, keeping the 3-tier distinction visually obvious.
- [x] 2.2 Update `KegCard.svelte`'s hardcoded `rgba()`/hex literals to reference the new tokens (no layout/behavior change — verify against `keg-card-display` spec scenarios still hold).
- [x] 2.3 Update `+page.svelte` (kiosk stage, pour flyout, dot-nav, empty state, admin-corner) to reference tokens in place of literals.
- [x] 2.4 Manually verify the kiosk display in the dev server: empty state, single tap, multi-tap carousel + dot-nav, pour flyout open/close/log, admin-corner link — confirm no visual regression from the pre-change look.

## 3. Admin shell

- [x] 3.1 Update `admin/+layout.svelte` (sidebar, mobile topbar, mobile bottomnav, version badge) to reference tokens in place of its existing literals (visual parity expected — this file was already on-brand).
- [x] 3.2 Retint `main-content`'s background from `#fafaf9` to the new warm-neutral light token.

## 4. Admin content pages

- [x] 4.1 Restyle `admin/login/+page.svelte`: card surface, primary button, input focus states, error banner — off default Tailwind blue/gray onto tokens.
- [x] 4.2 Restyle `admin/kegs/+page.svelte`: primary button, create-form inputs/selects, table, row hover, edit/delete actions, empty/loading/error states.
- [x] 4.3 Restyle `admin/kegs/[id]/+page.svelte` to match the same control/token treatment established in 4.2.
- [x] 4.4 Restyle `admin/taps/+page.svelte` to match.
- [x] 4.5 Restyle `admin/pours/+page.svelte` to match.
- [x] 4.6 Restyle `admin/banner/+page.svelte` to match.
- [x] 4.7 Restyle `admin/features/+page.svelte` to match.

## 5. Verification (initial re-skin pass)

- [x] 5.1 Walk every admin route in the dev server (`login`, `kegs`, `kegs/[id]`, `taps`, `pours`, `banner`, `features`) and confirm buttons, links, inputs, tables, and empty/loading/error states all read as one brand-consistent surface with no leftover default-Tailwind blue/gray.
- [x] 5.2 Confirm `npm run check` and `npm run lint` pass after all restyling.
- [x] 5.3 Confirm no data/behavior regression: keg CRUD, tap assignment, pour logging/history, banner message, and feature toggles all still function exactly as before across both mobile and desktop admin layouts.

## 6. Design direction supersession — Precision Tap Console

The initial re-skin (sections 1-5) kept the kiosk's warm copper/oak identity and a light admin surface, only retinting admin's generic-Tailwind controls onto that palette. A design review concluded this read as decoration rather than a considered redesign; the following tasks replaced it with a unified dark "instrument panel" identity across kiosk and admin. See `design.md` for the full rationale (Decisions 1-5) — this supersedes Decision 2 ("keep admin light") from the original design doc.

- [x] 6.1 Replace the `@theme` token block in `layout.css`: cool graphite void/panel/border surfaces, copper retained as a sparing signal color, new `--color-signal-good/warn/critical` telemetry palette, Chakra Petch + JetBrains Mono fonts (replacing Playfair Display/Bebas Neue/Special Elite). Add shared console component classes (`.panel`, `.data-table`, `.field`, `.log-row`, `.btn-console`/`.btn-console-ghost`, `.readout-badge`, `.led`, `.console-heading`, `.console-switch`).
- [x] 6.2 Rebuild `LevelGauge.svelte` as a circular conic-gradient gauge with the percentage centered inside; update `KegCard.svelte` to drop the now-redundant separate big-number block and retint to console tokens.
- [x] 6.3 Restyle kiosk `+page.svelte` (stage grid texture, dot-nav, pour flyout/fab, admin-corner) and `BannerStack.svelte` (signal-color banners) to the console palette.
- [x] 6.4 Rebuild admin shell (`admin/+layout.svelte`): dark `main-content` (was light), active-nav notch, sidebar retint.
- [x] 6.5 Rebuild `admin/login/+page.svelte` as a "system access" panel (connector-plate logo chip, eyebrow label).
- [x] 6.6 Rebuild `admin/kegs/+page.svelte` and `admin/kegs/[id]/+page.svelte` around `.panel`/`.data-table`/`.field`; sharpen the beer-image dropzone/style-picker controls to match.
- [x] 6.7 Rebuild `admin/taps/+page.svelte` with `.readout-badge` tap IDs and console-styled pour flyout.
- [x] 6.8 Rebuild `admin/pours/+page.svelte` as a system-log list (`.log-row`) instead of a table.
- [x] 6.9 Rebuild `admin/banner/+page.svelte` ("System Broadcast") and `admin/features/+page.svelte` ("System Modules", illuminated rocker `.console-switch` replacing the pill toggle) around the shared console classes.
- [x] 6.10 Confirm `npm run build`, `npm run check`, and `npm run lint` all pass; sweep the full source tree for zero leftover references to the superseded palette/class names (`--color-ink-*`, `--color-cream*`, `--color-paper*`, `ledger-*`, `btn-stamp*`, `tap-stamp`, `ticket-row`).
- [x] 6.11 Manually verify kiosk and every admin route in the dev server against an isolated scratch database (not the real `kegerator.db`); confirmed with the user directly.
