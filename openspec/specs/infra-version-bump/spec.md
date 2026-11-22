## Purpose

Defines independent per-package release and CI behavior for the `gopints/` and `web/` apps in this monorepo — separate versions, changelogs, and release PRs; path-filtered CI; package-scoped container/binary publishing — and records that deployment of the published images is out of this repository's automation scope.

## Requirements

### Requirement: Independent per-package releases
The repository SHALL manage `gopints/` and `web/` as independent release packages, each with its own version, CHANGELOG, and release PR, such that a change scoped to one package's files never bumps, opens a release PR for, or publishes a release for the other package.

#### Scenario: Change scoped to gopints/ only
- **WHEN** a commit touching only files under `gopints/` is merged to `main`
- **THEN** the release pipeline opens or updates a release PR for the `gopints` package only; the `web` package's version and CHANGELOG are untouched

#### Scenario: Change scoped to web/ only
- **WHEN** a commit touching only files under `web/` is merged to `main`
- **THEN** the release pipeline opens or updates a release PR for the `web` package only; the `gopints` package's version and CHANGELOG are untouched

### Requirement: Path-filtered CI
CI SHALL run only the checks relevant to the app(s) a pull request actually modifies.

#### Scenario: PR touches only web/
- **WHEN** a pull request's changes are contained entirely within `web/`
- **THEN** only the frontend CI workflow runs; the backend (Go) CI workflow does not run

#### Scenario: PR touches only gopints/
- **WHEN** a pull request's changes are contained entirely within `gopints/`
- **THEN** only the backend CI workflow runs; the frontend CI workflow does not run

### Requirement: Package-scoped container and binary publishing
Publishing a container image or agent binary SHALL only occur for the package whose release actually fired, not unconditionally for both.

#### Scenario: gopints release publishes
- **WHEN** a `gopints` package release is published
- **THEN** the `gopints-server` container image is built and pushed, and GoReleaser builds and attaches agent binaries; the `gopints-web` container image is not built

#### Scenario: web release publishes
- **WHEN** a `web` package release is published
- **THEN** the `gopints-web` container image is built and pushed; the `gopints-server` container image is not built and GoReleaser does not run

### Requirement: No automatic deployment
The system SHALL NOT deploy published images as part of its release automation; the release pipeline's responsibility ends once the container images and agent binaries are published (e.g. to GHCR and GitHub Release assets). Pulling and running them (e.g. via `docker compose pull && docker compose up -d`) remains a separate action outside this pipeline's scope.

#### Scenario: Release publishes
- **WHEN** a package's GitHub Release publishes and its artifacts are pushed
- **THEN** no process within this repository's CI deploys, or commits deployment configuration for, the new version

### Requirement: Release-prep merge detection is robust to commit message content
The `release-prepare-web.yml` and `release-prepare-backend.yml` workflows' release-prep merge detection SHALL pass the triggering commit message into its shell step via an `env:` variable, not by interpolating `${{ github.event.head_commit.message }}` directly into the `run:` script body, so that special characters in the commit message (quotes, backticks, `$()`, etc.) cannot break the script's shell parsing or be executed as shell syntax.

#### Scenario: Commit message contains a double quote
- **WHEN** a merge commit's message contains a literal `"` character (e.g. a PR title referencing a quoted name)
- **THEN** the "Detect release-prep merge" step still correctly evaluates whether the message references a `release-prep/<package>-v*` branch, without a shell parse error
