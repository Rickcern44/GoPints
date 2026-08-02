## ADDED Requirements

### Requirement: Admin interactive controls use brand tokens, not default Tailwind colors
Buttons, links, form inputs, selects, and focus states on admin content pages (`kegs`, `kegs/[id]`, `taps`, `pours`, `banner`, `features`, `login`) SHALL use the brand accent token in place of Tailwind's default `blue-*` palette.

#### Scenario: Primary action button
- **WHEN** an admin page renders a primary action button (e.g. "+ Register Keg", "Sign In", "Register")
- **THEN** it is styled with the brand accent token (`.btn-console`, `bg-accent`) rather than `bg-blue-600`/`hover:bg-blue-700`

#### Scenario: Text link and focus ring
- **WHEN** an admin page renders an inline text link (e.g. "Edit") or a focused form input
- **THEN** the link color and input focus ring use the brand accent token rather than `text-blue-600` or `focus:ring-blue-500`/`focus:border-blue-500`

### Requirement: Admin content surface matches the kiosk's dark instrument-panel surface
The admin content area (`main-content`) SHALL render on the same dark void/panel surface as the kiosk display and the admin shell chrome, rather than a separate light theme or Tailwind's default `gray-50`/`white`/`gray-100`.

#### Scenario: Page background
- **WHEN** any authenticated admin content page renders
- **THEN** its background is the shared dark `--color-void` token, not a light neutral or a literal `bg-gray-100`/`bg-white`/`#fafaf9`

#### Scenario: Card/list surface on the dark background
- **WHEN** a panel, table, or log-row container renders on an admin content page
- **THEN** its border and surface colors are drawn from the token set (`--color-panel`/`--color-panel-raised`/`--color-line`) rather than literal `border-gray-200`/`bg-white` values, reading as a slightly raised dark panel against the darker page background

### Requirement: Admin status states match the brand
Error, empty, and loading states on admin content pages SHALL be restyled to the token palette instead of Tailwind's stock red/gray defaults, while remaining clearly legible as their respective state.

#### Scenario: Error message
- **WHEN** an admin page shows a form or fetch error (e.g. "Failed to load kegs")
- **THEN** the error surface uses the token-derived critical-signal red (`--color-error`/`--color-error-bg`) rather than the literal `bg-red-50`/`text-red-700` pair, and remains visually distinct from the success/neutral states around it

#### Scenario: Empty and loading states
- **WHEN** an admin page has no data to show (e.g. "No kegs registered yet") or is still loading
- **THEN** the text color is the token-derived muted foreground (`--color-fg-muted`) rather than the literal `text-gray-500`
