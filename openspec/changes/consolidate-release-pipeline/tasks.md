## 1. Verify assumptions

- [x] 1.1 Confirm GitVersion's exact config key names/casing and the trunk-based/mainline preset name for conventional-commit-driven versioning, against the docs for the GitVersion action version actually pinned in CI
- [x] 1.2 Confirm this repo's actual PR merge strategy (merge commit vs. squash vs. rebase) via repo settings — the merge-detection mechanism in design.md depends on branch names surviving into the merge commit message

## 2. Per-package GitVersion + git-cliff config

- [x] 2.1 Create `gopints/GitVersion.yml`: `tag-prefix: 'gopints-v'`, `ignore.paths` scoped to exclude non-`gopints/` commits, trunk-based/mainline mode
- [x] 2.2 Create `web/GitVersion.yml`: `tag-prefix: 'web-v'`, `ignore.paths` scoped to exclude non-`web/` commits, same mode
- [x] 2.3 Create `gopints/cliff.toml`: `tag_pattern` matching `gopints-v*` only, `include_paths` scoped to `gopints/**`
- [x] 2.4 Create `web/cliff.toml`: `tag_pattern` matching `web-v*` only, `include_paths` scoped to `web/**`

## 3. Backend release-prepare workflow

- [x] 3.1 Create `.github/workflows/release-prepare-backend.yml`: triggers on `push: main` with `paths: ['gopints/**', '!gopints/infrastructure/**']`
- [x] 3.2 Add a step checking `github.event.head_commit.message` for the `release-prep/gopints-v` substring to decide which branch to take
- [x] 3.3 Finalize branch: extract the version from the matched substring, create the `gopints-vX.Y.Z` tag at HEAD, and `gh release create` it
- [x] 3.4 Propose branch: run GitVersion with `gopints/GitVersion.yml`; if the computed version matches the current highest `gopints-v*` tag (no releasable commits), stop without opening a PR
- [x] 3.5 Propose branch: run git-cliff with `gopints/cliff.toml` to update `gopints/CHANGELOG.md`
- [x] 3.6 Propose branch: update `gopints/infrastructure/deploy/deployment.yaml`'s image tag to the computed version (scoped `sed` replace, same approach as the prior change)
- [x] 3.7 Propose branch: open or update a PR via `peter-evans/create-pull-request` on branch `release-prep/gopints-v<version>`, containing the CHANGELOG.md and deployment.yaml changes together

## 4. Web release-prepare workflow

- [x] 4.1 Create `.github/workflows/release-prepare-web.yml`: triggers on `push: main` with `paths: ['web/**', '!web/infrastructure/**']`
- [x] 4.2 Add the same merge-detection step, checking for the `release-prep/web-v` substring
- [x] 4.3 Finalize branch: extract the version, create the `web-vX.Y.Z` tag at HEAD, and `gh release create` it
- [x] 4.4 Propose branch: run GitVersion with `web/GitVersion.yml`; skip opening a PR if there's no version change
- [x] 4.5 Propose branch: run git-cliff with `web/cliff.toml` to update `web/CHANGELOG.md`
- [x] 4.6 Propose branch: update `web/infrastructure/deploy/deployment.yaml`'s image tag, and bump `web/package.json`'s version field, to the computed version
- [x] 4.7 Propose branch: open or update a PR via `peter-evans/create-pull-request` on branch `release-prep/web-v<version>`, containing the CHANGELOG.md, deployment.yaml, and package.json changes together

## 5. Remove release-please

- [x] 5.1 Remove `release-please-config.json`
- [x] 5.2 Remove `.release-please-manifest.json`
- [x] 5.3 Remove `.github/workflows/release-please.yml`

## 6. Simplify release.yml

- [x] 6.1 Remove the `bump-infra-version` job from `.github/workflows/release.yml`; leave `prepare`, `build-server`, `build-web`, and `goreleaser` unchanged

## 7. Docs

- [x] 7.1 Update `README.md`'s CI/CD section: replace `release-please.yml` references with `release-prepare-backend.yml`/`release-prepare-web.yml`, and describe the single consolidated PR per package instead of the two-PR flow

## 8. Verification

- [x] 8.1 Validate all new/edited YAML and TOML files parse correctly
- [x] 8.2 Hand-trace `release-prepare-backend.yml` for both an ordinary `gopints/`-touching push and a push that merges a `release-prep/gopints-v*` branch, confirming the correct branch (propose vs. finalize) is taken in each case
- [x] 8.3 Hand-trace `release-prepare-web.yml` the same way for both cases
- [x] 8.4 Confirm both release-prepare workflows' path filters don't create a trigger loop off the commits their own propose-branch PRs introduce once merged (i.e. merging only touches `CHANGELOG.md`/`deployment.yaml`/version files under the package path, which is expected and handled by the finalize branch, not a new infinite loop)
- [ ] 8.5 Note in the PR/commit description that full end-to-end verification (a real release-prep PR opening, a real merge correctly triggering tag + release, GitVersion computing the expected bump) can only be confirmed by observing real pushes/merges to `main` in the live repo — not reproducible in local sandbox testing
