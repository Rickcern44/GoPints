## Purpose

Defines the display behavior of the `KegCard` component, which renders per-tap keg information (brewery, beer, and pour stats) on the kegerator monitoring frontend.

## Requirements

### Requirement: Three-zone bookend layout
The `KegCard` component SHALL render each tap's content in up to three horizontally arranged zones within `.inner`: a fixed-width brewery plaque, a flexible beer-info zone, and a fixed-width stats zone.

#### Scenario: Keg with brewery information renders all three zones
- **WHEN** a tap's keg has a non-empty `brewery` name
- **THEN** the card renders the brewery plaque (left), beer info (middle), and stats (right) as three distinct zones

#### Scenario: Keg without brewery information collapses to two zones
- **WHEN** a tap's keg has an empty `brewery` field
- **THEN** the brewery plaque zone is omitted from the layout entirely (not merely hidden), and the beer-info zone expands to fill the freed space

### Requirement: Brewery plaque presentation
The brewery plaque SHALL display the brewery logo (when present) stacked above the brewery name, contained in a fixed-width column with a distinct background/border treatment that differentiates it from the beer-info and stats zones.

#### Scenario: Brewery has a logo image
- **WHEN** `keg.brewery_image_mime_type` is set
- **THEN** the plaque renders the logo image (110-140px) above the brewery name text, both centered within the plaque column

#### Scenario: Brewery has no logo image
- **WHEN** `keg.brewery_image_mime_type` is empty but `keg.brewery` is non-empty
- **THEN** the plaque renders the brewery name text alone, without a broken image or empty image placeholder

#### Scenario: Long brewery name in a narrow column
- **WHEN** `keg.brewery` exceeds the plaque column's display width
- **THEN** the brewery name wraps or truncates without breaking the column's fixed width or overflowing into adjacent zones

### Requirement: Beer image retains size and position within beer-info zone
The beer image SHALL continue to render at its existing size (90px) and existing shape variants (circle/square/can), positioned within the beer-info zone independent of whether the brewery plaque is present.

#### Scenario: Beer image renders regardless of brewery plaque presence
- **WHEN** a keg has a beer image (`image_mime_type` set) and either has or lacks brewery information
- **THEN** the beer image renders at its existing size and selected `image_style` shape in both cases

### Requirement: No new responsive breakpoints
The three-zone layout SHALL rely on the existing fluid (`clamp()`-based) sizing approach and SHALL NOT introduce viewport-width media queries, consistent with the display's TV/laptop-only usage.

#### Scenario: Layout at TV and laptop widths
- **WHEN** the card is viewed at typical TV or laptop screen widths
- **THEN** all three zones (when present) remain legible and proportionally balanced without relying on breakpoint-specific layout rules
