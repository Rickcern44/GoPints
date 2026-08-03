## ADDED Requirements

### Requirement: Independent per-package version computation
The system SHALL compute each package's (`gopints`, `web`) next semantic version independently, considering only commits that touched files under that package's own path and only git tags matching that package's own tag prefix (`gopints-v*` / `web-v*`) as prior release history.

#### Scenario: Commit touches only web/
- **WHEN** a commit touching only files under `web/` lands on `main`
- **THEN** `web`'s next computed version reflects that commit's Conventional Commit type; `gopints`'s next computed version is unaffected

#### Scenario: Commit touches only gopints/
- **WHEN** a commit touching only files under `gopints/` lands on `main`
- **THEN** `gopints`'s next computed version reflects that commit's Conventional Commit type; `web`'s next computed version is unaffected

#### Scenario: Cross-cutting commit touches both packages
- **WHEN** a single commit touches files under both `gopints/` and `web/`
- **THEN** both packages' next-version computation considers that commit

#### Scenario: No releasable commits since last tag
- **WHEN** no commits since a package's last matching tag contain a Conventional Commit type that warrants a version bump (`feat`, `fix`, or a breaking change)
- **THEN** that package's computed next version is unchanged from its last released version

### Requirement: Independent per-package changelog generation
The system SHALL generate each package's `CHANGELOG.md` entries from only the commits scoped to that package's path, bounded by that package's own tag history.

#### Scenario: gopints changelog excludes web-only commits
- **WHEN** a release-prep run generates `gopints/CHANGELOG.md`
- **THEN** the generated entries include only commits that touched files under `gopints/`; commits that touched only `web/` do not appear

#### Scenario: web changelog excludes gopints-only commits
- **WHEN** a release-prep run generates `web/CHANGELOG.md`
- **THEN** the generated entries include only commits that touched files under `web/`; commits that touched only `gopints/` do not appear
