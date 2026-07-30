## Why

The just-shipped per-app release split (independent `gopints`/`web` versioning) fixed the "everything releases together" problem, but doubled the PR overhead: each package release now needs a release-please version-bump PR *and* a separate post-publish infra-manifest-bump PR. For a solo maintainer, that's up to 4 PRs to review and merge for two independent single-package releases. We're moving off release-please (Node-based, opinionated release-PR lifecycle) to GitVersion (git-history-driven version computation) + git-cliff (fast, config-driven changelog generation), and consolidating each package's release-prep PR and infra-tag-bump PR into a single PR — cutting the per-release review burden in half.

## What Changes

- Remove release-please entirely: `release-please-config.json`, `.release-please-manifest.json`, `.github/workflows/release-please.yml`.
- Add GitVersion (computes each package's next semantic version from Conventional Commits on files under that package's path, scoped independently per package) and git-cliff (generates/updates that package's `CHANGELOG.md` from the same commit range) as the new versioning/changelog engine, replacing release-please's manifest-driven approach.
- Replace `release-please.yml` with a new per-package `release-prepare-*.yml` workflow that, on push to `main`: computes the next version via GitVersion, generates the changelog entry via git-cliff, and (for `web`) bumps `package.json`'s version field — then opens **one** PR per package containing that together. No PR opens if a package has no releasable commits since its last tag.
- Merging that single PR is what creates the git tag (`gopints-vX.Y.Z` / `web-vX.Y.Z`) and publishes the GitHub Release for that package — preserving today's "merge = release" mental model, just with one PR per package instead of two.
- The `bump-infra-version` job in `release.yml` is kept, but repositioned: it now runs *after* `build-server`/`build-web` succeed, as a direct commit to `main` (not a PR) updating that package's k8s manifest (`infrastructure/deploy/deployment.yaml`) image tag. This repo has ArgoCD auto-syncing off it, so the manifest must never reference a version that isn't actually in GHCR yet — bumping it pre-merge (as originally planned) would have created exactly that window. The existing `containers`/`goreleaser` jobs and their per-package gating (added in the prior change) are unchanged.
- **BREAKING**: CHANGELOG.md generation moves from release-please's format/config to git-cliff's — history is preserved, but new entries follow git-cliff's section/format conventions going forward.
- Supersedes the infra-bump mechanism proposed in the (open, unarchived) `split-release-pipelines-per-app` change; that change's independent per-package release streams and container/goreleaser gating are preserved and built upon here, not reverted.

## Capabilities

### New Capabilities
- `versioned-changelog-generation`: Per-package next-version computation (GitVersion) and changelog generation (git-cliff) from Conventional Commits, replacing release-please's manifest-driven versioning.
- `consolidated-release-pr`: A single per-package PR bundling the version bump and changelog update; merging it is what triggers that package's tag, GitHub release, and container/binary publish. The infra manifest image-tag bump happens separately, as a direct commit after the image is confirmed pushed.

### Modified Capabilities
(none — no capability from the prior change has been archived into `openspec/specs/` yet, so there is nothing to file a delta against; this proposal's new capabilities supersede that change's still-pending `infra-version-bump` spec in intent)

## Impact

- Removed: `release-please-config.json`, `.release-please-manifest.json`, `.github/workflows/release-please.yml`.
- Added: GitVersion config (per package), git-cliff config (per package), `.github/workflows/release-prepare.yml`.
- Modified: `.github/workflows/release.yml` (`bump-infra-version` job repositioned to run after `build-server`/`build-web` succeed, as a direct commit instead of a PR; `containers`/`goreleaser` jobs and gating unchanged).
- Modified: `gopints/infrastructure/deploy/deployment.yaml`, `web/infrastructure/deploy/deployment.yaml` — still updated automatically, but only after the corresponding image is confirmed pushed, via a direct commit rather than a pre-merge PR.
- Modified: `web/package.json` — version field still auto-bumped, now via GitVersion/git-cliff instead of release-please's node updater.
- No application code changes (`gopints/` Go source, `web/` frontend source untouched) — CI/release tooling only.
