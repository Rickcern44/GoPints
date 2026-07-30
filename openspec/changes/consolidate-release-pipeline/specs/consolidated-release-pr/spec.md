## ADDED Requirements

### Requirement: Single PR bundles version and changelog
When there are releasable changes for a package, the system SHALL open (or update) a single pull request containing that package's version bump and changelog update together, rather than as separate pull requests.

#### Scenario: Releasable commits since last tag
- **WHEN** commits with a Conventional Commit type warranting a version bump have landed on `main` under `gopints/` since `gopints`'s last release tag
- **THEN** a single pull request is opened (or updated) containing `gopints`'s version bump and its `CHANGELOG.md` entry, in one diff

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
The system SHALL NOT open any additional pull request after a release publishes to update the infrastructure deployment manifest.

#### Scenario: Release publishes
- **WHEN** a package's GitHub Release publishes
- **THEN** no new pull request is opened to update that package's infrastructure deployment manifest

### Requirement: Infra manifest is only updated after the image is confirmed pushed
The system SHALL NOT update a package's infrastructure deployment manifest image tag until that package's container image has been successfully built and pushed to the registry. The manifest update SHALL be a direct commit to `main`, not part of the pre-merge release-prep pull request, so that a GitOps controller watching this repository never observes an image tag that does not yet exist in the registry.

#### Scenario: Image push succeeds
- **WHEN** a package's container image finishes building and pushing to the registry successfully
- **THEN** that package's infrastructure deployment manifest image tag is updated to the new version via a direct commit to `main`

#### Scenario: Image push fails
- **WHEN** a package's container image build or push fails
- **THEN** the infrastructure deployment manifest is not updated, and continues to reference the last successfully published version

#### Scenario: Release-prep pull request never touches the manifest
- **WHEN** a release-prep pull request is opened or updated for either package
- **THEN** it does not include any change to that package's `infrastructure/deploy/deployment.yaml`
