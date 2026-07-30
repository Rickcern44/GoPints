## ADDED Requirements

### Requirement: Version-change detection
The kiosk display SHALL capture the server's version at initial load and periodically re-check it, without requiring any backend changes or new endpoints.

#### Scenario: Initial load captures baseline version
- **WHEN** the kiosk display first loads
- **THEN** it fetches `GET /api/health` and stores the returned version as the baseline for comparison

#### Scenario: Periodic re-check
- **WHEN** the kiosk display has been open for an extended period
- **THEN** it re-fetches `GET /api/health` on a recurring interval (approximately every 5 minutes) for as long as the page remains open

#### Scenario: Health check fails
- **WHEN** a periodic re-check fails or does not resolve (e.g. network unavailable)
- **THEN** no error is shown and no update badge appears; the display continues polling on the next interval

### Requirement: Update-available badge
The kiosk display SHALL show a small, tappable badge when the server's version differs from the baseline captured at load, and SHALL NOT show it otherwise.

#### Scenario: Version unchanged
- **WHEN** a periodic re-check returns the same version as the baseline
- **THEN** no update badge is shown

#### Scenario: Version changed
- **WHEN** a periodic re-check returns a version different from the baseline
- **THEN** a tappable "update available" badge appears in the bottom-left corner of the kiosk display, clear of the existing bottom-center pour controls and bottom-right admin link

#### Scenario: Tapping the badge reloads the page
- **WHEN** an admin taps the update-available badge
- **THEN** the browser performs a full page reload, fetching the latest deployed assets

### Requirement: No automatic reload
The kiosk display SHALL NOT reload automatically when a new version is detected — reloading SHALL only happen via an explicit tap on the badge.

#### Scenario: Update detected while idle
- **WHEN** a new version is detected and the badge is shown
- **THEN** the page continues running its current state (carousel, pours, banners) unchanged until the badge is tapped
