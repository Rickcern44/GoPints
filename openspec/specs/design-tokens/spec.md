## Purpose

A shared CSS custom-property layer (color, type) defined once in `layout.css` and consumed by both the kiosk display and the admin UI, replacing hardcoded per-component color literals.

## Requirements

### Requirement: Centralized design tokens
The frontend SHALL define its color, typography, spacing, and radius values as a single set of design tokens declared once in `web/src/routes/layout.css` via Tailwind's `@theme` directive, rather than as literal values repeated per component.

#### Scenario: Token available as a CSS custom property
- **WHEN** a hand-written `<style>` block (e.g. in `KegCard.svelte` or `+page.svelte`) needs the brand accent color
- **THEN** it references it as `var(--color-*)` rather than a literal `rgba()`/hex value

#### Scenario: Token available as a Tailwind utility class
- **WHEN** a Tailwind-utility-authored page (an admin content page) needs the brand accent color
- **THEN** it uses a generated utility class (e.g. `bg-accent`, `text-accent`, `border-accent/30`) rather than a stock Tailwind palette class (e.g. `bg-blue-600`) or an arbitrary-value class (e.g. `bg-[#c8821a]`)

### Requirement: Shared components consume tokens for state color
Components whose color communicates state (e.g. `LevelGauge`'s healthy/low/critical fill) SHALL derive that color from the design tokens rather than Tailwind's default palette classes, while preserving a distinguishable multi-tier signal.

#### Scenario: Gauge fill reflects remaining percentage using token colors
- **WHEN** `LevelGauge` renders at a "healthy," "low," or "critical" remaining percentage
- **THEN** each tier renders a token-derived color distinct from the other two tiers, and none of the three is Tailwind's stock `green-500`, `amber-400`, or `red-500`

### Requirement: Brand fonts load once, globally
Chakra Petch (display/headings) and JetBrains Mono (data/readouts) SHALL be loaded via a single set of font imports declared globally (in `layout.css`), not via per-component `@import url(...)` statements.

#### Scenario: Font import removed from component style blocks
- **WHEN** `KegCard.svelte` or `admin/+layout.svelte` render
- **THEN** neither component's `<style>` block contains its own `@import url(...)` for a Google Font; both fonts are already available from the global import

#### Scenario: Admin page titles adopt the shared display face
- **WHEN** an admin content page renders its top-level page heading (e.g. "Keg Stock", "Tap Manifest")
- **THEN** the heading uses the same display font family (`var(--font-display)`/`var(--font-brand)`, Chakra Petch) as the admin sidebar wordmark, rather than the default system font stack

#### Scenario: Data values render in the mono readout face
- **WHEN** a page renders a data value (a tap ID, a volume, a timestamp, a form field's typed value)
- **THEN** it renders in `var(--font-mono)` (JetBrains Mono), distinguishing telemetry/data text from heading/label text
