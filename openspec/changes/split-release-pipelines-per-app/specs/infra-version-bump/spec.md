## ADDED Requirements

### Requirement: Independent per-package releases
The repository SHALL manage `gopints/` and `web/` as independent release-please packages, each with its own version, CHANGELOG, and release PR, such that a change scoped to one package's files never bumps, opens a release PR for, or publishes a release for the other package.

#### Scenario: Change scoped to gopints/ only
- **WHEN** a commit touching only files under `gopints/` is merged to `main`
- **THEN** release-please opens or updates a release PR for the `gopints` package only; the `web` package's version and CHANGELOG are untouched

#### Scenario: Change scoped to web/ only
- **WHEN** a commit touching only files under `web/` is merged to `main`
- **THEN** release-please opens or updates a release PR for the `web` package only; the `gopints` package's version and CHANGELOG are untouched

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

### Requirement: Automated infra manifest version bump
After a package's release is published, the system SHALL open a pull request updating that package's Kubernetes deployment manifest image tag to the newly released version. It SHALL NOT apply the manifest to a live cluster and SHALL NOT auto-merge the pull request.

#### Scenario: gopints release triggers an infra bump PR
- **WHEN** a `gopints` package release is published with version `X.Y.Z`
- **THEN** a pull request is opened updating `gopints/infrastructure/deploy/deployment.yaml`'s image tag to `X.Y.Z`, requiring manual review and merge

#### Scenario: web release triggers an infra bump PR
- **WHEN** a `web` package release is published with version `X.Y.Z`
- **THEN** a pull request is opened updating `web/infrastructure/deploy/deployment.yaml`'s image tag to `X.Y.Z`, requiring manual review and merge

#### Scenario: No automatic cluster deployment
- **WHEN** an infra-bump pull request is merged
- **THEN** no automated process applies the updated manifest to a live cluster; applying it remains a manual, separate action
