## Context

The frontend (SvelteKit 5 + Tailwind v4) went through two iterations in this change:

1. An initial re-skin kept the kiosk's warm copper/oak "brewpub" identity and admin's light "paper" surface, just retinting admin's generic-Tailwind buttons/tables/forms onto that palette.
2. A design review concluded the light-admin/dark-kiosk split — and the brewpub metaphor itself — read as decoration rather than a considered design. This change supersedes that pass: kiosk and admin are now unified into one dark, technical "instrument panel" identity (**Precision Tap Console**), since a kegerator monitoring/telemetry tool is naturally suited to that register.

There is no Playwright suite yet (despite it being in the project's tech stack) — verification for this change was manual, via the dev server against an isolated scratch database.

## Goals / Non-Goals

**Goals:**
- One dark surface end to end — kiosk and admin share the same void/panel/border tokens, so there is no seam between them.
- One token layer (`web/src/routes/layout.css`, via Tailwind v4 `@theme`) that both hand-written `<style>` blocks and Tailwind utility classes consume.
- A small shared "console" component vocabulary (`.panel`, `.data-table`, `.field`, `.log-row`, `.btn-console`/`.btn-console-ghost`, `.readout-badge`, `.led`, `.console-heading`) used directly by every admin route template, rather than each page reinventing card/table/form patterns.
- `LevelGauge` as a genuine dashboard element (circular gauge) rather than a retinted generic progress bar.
- Preserve the `keg-card-display` spec's locked-in zone/sizing rules for `KegCard` — the console pass changes color/type/shape of chrome within those zones, not the zones themselves.

**Non-Goals:**
- No new pages, routes, or admin features; no API/data changes.
- Not adding a Playwright/visual-regression harness (real gap, flagged as a follow-up, not this change's job).
- Not preserving the admin-is-light-paper decision from the superseded pass — explicitly reversed here.

## Decisions

**1. Cool graphite surfaces, not warm oak — copper becomes a signal color, not a wash.**
The prior palette used copper at many opacities as a decorative brand wash (borders, backgrounds, text) throughout. The console palette keeps copper (`#c8821a`) as *the* accent but narrows its role: primary actions, focus rings, active nav state, and the "low" gauge/LED tier — not a tint applied to every surface. Void (`#0a0d0f`)/panel (`#12171a`)/panel-raised (`#1a2226`) surfaces are a cool near-black rather than the previous warm browns, and a dedicated telemetry palette (`--color-signal-good` phosphor-green, `--color-signal-warn` = accent, `--color-signal-critical` red) drives status (LEDs, gauge tiers, banners) independently of the brand accent.

**2. Admin content moves from light to dark — reversing the prior pass's explicit decision.**
The earlier design argued a light admin surface serves data-entry legibility better than full dark. This pass overrides that: a unified dark surface is the entire point of "one console instead of two apps," and dashboard/ops tooling (Grafana, Linear, most IoT consoles) is conventionally dark without a legibility penalty, provided body text contrast is kept high (`--color-fg` `#e8edf0` on `--color-void` `#0a0d0f` is comfortably AA+). Accepted trade-off: a bright-daylight kitchen counter may have more glare on a dark admin panel than a light one would — judged acceptable since the admin is typically used briefly, not for extended data entry.

**3. Typography: Chakra Petch (display) + JetBrains Mono (data), replacing all three previous faces.**
Playfair Display (italic serif), Bebas Neue (condensed poster face), and Special Elite (typewriter) all carried "warm/vintage brewpub" connotations incompatible with a technical-console register. Chakra Petch is a geometric sans with slightly cut corners (HUD/technical character) for headings and brand type; JetBrains Mono renders every data value (tap IDs, volumes, timestamps, form inputs) so the interface reads as instrument telemetry rather than prose. Both load once via a single `@import` in `layout.css` (carried over from the prior pass's fix for the original per-component duplicate-`@import` problem).

**4. `LevelGauge` becomes a circular gauge (conic-gradient ring), not a retinted bar.**
A vertical fill bar is a generic progress-bar pattern; a circular instrument gauge is the idiomatic "console" element and doubles as the keg's hero stat — the big percentage now renders inside the ring itself, so `KegCard` drops the separate large-number block it used to render beside the gauge (was duplicated information, one visual element now carries it). The three color tiers (`--color-signal-good`/`--color-accent`/`--color-signal-critical`) preserve the existing at-a-glance healthy/low/critical read.

**5. Readouts over decoration: tap IDs, pour history, and toggles get purpose-built console controls.**
Tap IDs render as a rectangular connector-plate badge (`.readout-badge`, clipped corner) instead of a plain number or (in the superseded pass) a rubber-stamp circle. Pour history is a system-log list (`.log-row`) instead of a table or a receipt-roll metaphor. Feature flags use an illuminated rocker switch (`.console-switch`, half-track thumb that glows green when on) instead of an iOS pill toggle or a tap-handle lever. Each of these is a small, purpose-specific component living in the shared stylesheet, reused by every page that needs it, rather than bespoke per-page markup.

## Risks / Trade-offs

- **[Risk]** This is the second full visual pass in one change — reviewers diffing against the *original* pre-change code see two superimposed rewrites. **Mitigation**: `proposal.md`/`design.md` describe only the final state; git history retains the intermediate pass for anyone who needs it.
- **[Risk]** No automated visual regression coverage exists. **Mitigation**: manual pass through every kiosk/admin route in the dev server against an isolated scratch database before calling this done (performed); Playwright snapshots remain a flagged follow-up.
- **[Trade-off]** Reversing Decision 2 from the superseded design doc means the "light admin for legibility" argument made there no longer applies — recorded here so a future reader doesn't wonder why the codebase contradicts an archived doc elsewhere.

## Migration Plan

Purely additive/CSS — no data migration.
1. Replace the `@theme` token block and shared component classes in `layout.css` (surfaces, fonts, `.panel`/`.data-table`/`.field`/`.log-row`/`.btn-console*`/`.readout-badge`/`.led`/`.console-heading`).
2. Rebuild `LevelGauge` (circular gauge) and retint `KegCard`/kiosk `+page.svelte`/`BannerStack`.
3. Rebuild admin shell (`+layout.svelte`) and each content page (`login` → `kegs` → `kegs/[id]` → `taps` → `pours` → `banner` → `features`) around the shared console classes.
4. No feature flag or rollback beyond normal git revert.

## Open Questions

- None blocking. A future Playwright pass would want visual snapshots of the console baseline established here.
