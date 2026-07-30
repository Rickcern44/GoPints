## ADDED Requirements

### Requirement: Single PR bundles version, changelog, and infra tag
When there are releasable changes for a package, the system SHALL open (or update) a single pull request containing that package's version bump, changelog update, and infrastructure deployment manifest image-tag update together, rather than as separate pull requests.

#### Scenario: Releasable commits since last tag
- **WHEN** commits with a Conventional Commit type warranting a version bump have landed on `main` under `gopints/` since `gopints`'s last release tag
- **THEN** a single pull request is opened (or updated) containing `gopints`'s version bump, its `CHANGELOG.md` entry, and `gopints/infrastructure/deploy/deployment.yaml`'s image tag update, all in one diff

#### Scenario: No releasable changes
- **WHEN** no commits since a package's last release tag warrant a version bump
- **THEN** no pull request is opened or updated for that package

### Requirement: Merging the consolidated PR triggers the release
Merging a package's consolidated release pull request SHALL create that package's git tag and publish its GitHub Release, without requiring any separate manual tagging step.

#### Scenario: Consolidated PR merged
- **WHEN** a `gopints` consolidated release pull request is merged into `main`
- **THEN** a `gopints-vX.Y.Z` tag is created at the merge commit and a GitHub Release is published for that tag, with no additional manual action required

#### Scenario: Ordinary push is not mistaken for a release merge
- **WHEN** a push to `main` touching `gopints/**` is not the merge of a `gopints` consolidated release pull request
- **THEN** the system does not create a tag or publish a release from that push; it only evaluates whether to open or update a consolidated release pull request

### Requirement: No separate post-publish infra-bump pull request
The system SHALL NOT open any additional pull request after a release publishes to update the infrastructure deployment manifest — that update SHALL already be present in the merged consolidated pull request.

#### Scenario: Release publishes
- **WHEN** a package's GitHub Release publishes
- **THEN** no new pull request is opened to update that package's infrastructure deployment manifest, since it was already updated and merged as part of the consolidated release pull request
