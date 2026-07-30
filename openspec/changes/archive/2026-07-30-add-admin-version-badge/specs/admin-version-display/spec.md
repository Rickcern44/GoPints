## ADDED Requirements

### Requirement: Version badge in the admin shell
The authenticated admin UI SHALL display the running server's version, sourced from `GET /api/health`, as a small non-interactive badge fixed to the lower-right corner of the viewport on desktop widths.

#### Scenario: Health fetch succeeds
- **WHEN** the admin shell loads and `GET /api/health` returns a version string
- **THEN** a badge reading `v{version}` renders fixed to the lower-right corner of the viewport, at desktop widths (≥768px)

#### Scenario: Health fetch fails
- **WHEN** `GET /api/health` fails or does not resolve
- **THEN** no badge is rendered and no error is shown to the admin

#### Scenario: Mobile width
- **WHEN** the admin shell is viewed below the existing 768px desktop breakpoint
- **THEN** the version badge is not shown, avoiding collision with the mobile bottom navigation bar

#### Scenario: Login page
- **WHEN** the admin is on `/admin/login` (not yet authenticated)
- **THEN** no version badge is shown, consistent with the rest of the admin shell chrome only rendering once authenticated
