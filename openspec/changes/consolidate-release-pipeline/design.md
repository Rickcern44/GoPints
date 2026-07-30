## Context

Today's flow (from the just-merged `split-release-pipelines-per-app` change): `release-please.yml` opens a per-package version-bump PR on every push to `main`; merging it makes release-please create a `<package>-vX.Y.Z` tag and publish a GitHub Release; that publish event fires `release.yml`, which (already correctly gated per package) builds/pushes the container image, runs GoReleaser for `gopints` only, and — as a separate job — opens a **second** PR bumping that package's `infrastructure/deploy/deployment.yaml` image tag. Two independent packages × two PRs each = up to 4 PRs per pair of releases, which is the overhead this change removes.

Researched during design (see citations in proposal discussion): GitVersion 6 has no fully "native" monorepo mode yet, but its existing, stable `tag-prefix` (which tags count as this package's release history) and `ignore.paths` (which commits count toward this package's version, filtered by files touched) config keys are sufficient to scope it per package. git-cliff has the equivalent stable `tag_pattern` and `include_paths` config keys for the same purpose. Both are used per-package config files rather than CLI flags, so behavior is explicit and reviewable in the repo rather than hidden in workflow YAML.

## Goals / Non-Goals

**Goals:**
- Replace release-please with GitVersion (version computation) + git-cliff (changelog generation), both scoped independently per package via config, not shared state.
- Collapse each package's release-prep PR and infra-manifest-bump PR into one PR: version bump + CHANGELOG.md + `deployment.yaml` image tag, all in a single diff.
- Preserve "merge = release": merging that one PR is still what creates the tag, publishes the GitHub Release, and (via the unchanged `release.yml`) builds/publishes that package's container image (and agent binaries, for `gopints`).
- Preserve the independent per-package release cadence from the prior change.

**Non-Goals:**
- Not touching `release.yml`'s `prepare`/`build-server`/`build-web`/`goreleaser` jobs or their tag-prefix gating logic — that's already correct and package-scoped from the prior change.
- Not adopting a monorepo build orchestrator (Nx, Turborepo, etc.) — GitVersion/cliff config alone is enough for two packages.
- Not applying k8s manifests to a live cluster — carried over from the prior change's non-goal.
- Not changing application code in `gopints/` or `web/`.

## Decisions

**GitVersion config lives at `gopints/GitVersion.yml` and `web/GitVersion.yml`, each self-contained.**
GitVersion does not merge parent/child configs, so each package needs its own complete file. Each sets:
- `tag-prefix: 'gopints-v'` (or `'web-v'`) so GitVersion only recognizes that package's own tags as release history — it won't be confused by the other package's tags existing in the same repo.
- `ignore: paths: ['^(?!gopints/).+']` (inverse for `web`) — a commit is excluded from version-bump consideration only if *all* files it touched fail to match `gopints/`, so cross-cutting commits still count for both packages, but a `web`-only commit never bumps `gopints`.
- Mode/workflow preset: intend to use GitVersion's trunk-based/mainline preset (closest match to today's "conventional commits directly on `main`" model, no GitFlow branches) — **exact preset name to be confirmed against the currently-installed GitVersion version's docs during implementation** (flagged as a task; GitVersion 6's preset naming has been in flux and guessing wrong here silently produces wrong version numbers, not an error).

**git-cliff config lives at `gopints/cliff.toml` and `web/cliff.toml`, run from repo root with `--config <path>`.**
Each sets `tag_pattern` (matching only that package's tag prefix, e.g. `"gopints-v[0-9].*"`) and `include_paths` (e.g. `["gopints/**"]`) so changelog entries are scoped the same way GitVersion's bump decision is. Running with an explicit `--config` flag from repo root (rather than relying on git-cliff's cwd-based auto-detection) keeps CI invocation unambiguous regardless of working directory.

**One workflow per package (`release-prepare-backend.yml`, `release-prepare-web.yml`), not one shared file.**
Mirrors the existing `ci-backend.yml`/`ci-frontend.yml` split (precedent already in this repo) — each triggers on `push: main` with `paths: ['gopints/**', '!gopints/infrastructure/**']` (inverse for web), same exclusion rationale as CI: don't re-trigger off the manifest-bump change this same workflow makes.

**A single workflow handles both "propose the release" and "finalize the release" via merge-commit detection — no separate finalize workflow.**
This is the part release-please used to handle internally, and now has to be hand-built. The mechanism:
1. Release-prep PRs are opened on a branch named `release-prep/<package>-v<version>` (e.g. `release-prep/gopints-v0.5.0`).
2. This repo merges PRs via GitHub's default "Merge pull request" button (evidenced by existing history — `git log` shows `Merge pull request #NN from <branch>` commits, not squash commits), which means the resulting merge commit's message contains the source branch name.
3. On every push to `main` matching the package's path filter, the workflow's first step checks `github.event.head_commit.message` for the substring `release-prep/gopints-v` (or `web-v`):
   - **Match found** → this push *is* the merge of a release-prep PR. Extract the version from the matched substring, create the git tag `<package>-v<version>` at the current `main` HEAD, and `gh release create` it (publishing the release, which fires the existing `release.yml`). Do not compute a new version or open a new PR.
   - **No match** → this is an ordinary code push. Run GitVersion to compute the next version; if it's unchanged from the current highest matching tag (no releasable commits), stop — don't open an empty PR. Otherwise run git-cliff to generate the changelog entry, update `infrastructure/deploy/deployment.yaml`'s image tag (same scoped `sed` approach as the prior change) and, for `web`, `package.json`'s version field, and open/update the consolidated PR via `peter-evans/create-pull-request`.

**`bump-infra-version` is removed from `release.yml`.**
Its work now happens pre-merge, inside the same PR/branch the version-prep step produces — there's no longer a post-publish step needed.

## Risks / Trade-offs

- **[Risk]** The merge-detection mechanism assumes this repo's PRs are always merged via a standard two-parent merge commit (not squash or rebase merge), since it relies on the branch name surviving into the merge commit message. → **Mitigation**: verify the repo's actual configured/used merge strategy as an explicit early implementation task (existing `git log` history strongly suggests merge commits, but repo settings could allow squash too — worth confirming before relying on it). If squash-merge is possible, fall back to checking the squash commit's message/title directly (which becomes the PR title) instead of a branch-name substring.
- **[Risk]** GitVersion 6's monorepo-adjacent config (`tag-prefix`, `ignore.paths`) is being applied in a pattern this repo hasn't used before; exact YAML key casing/nesting should be verified against the GitVersion version actually pinned in CI, not trusted from secondary sources. → **Mitigation**: explicit verification task before wiring the real config, mirroring how the previous change verified release-please's tag format before depending on it.
- **[Risk]** Hand-rolling "propose vs. finalize" branching logic moves real complexity that release-please used to own into this repo's own workflow YAML/shell — more surface area to get subtly wrong (e.g., a false-positive substring match, a version-extraction regex bug). → **Mitigation**: task 7-style hand-tracing of both branches (ordinary push vs. merge-of-release-prep-branch) for both packages before considering this done, same verification discipline as the prior change.
- **[Trade-off]** The window between a release-prep PR merging (manifest already bumped) and the corresponding image actually landing in GHCR is now slightly longer than before (merge → tag-detection step → `release.yml`'s `published` trigger → build), not shorter — accepted, same non-goal (no live-cluster auto-apply) as previously.
- **[Trade-off]** Losing release-please's Node-ecosystem maturity/battle-testing in exchange for lighter, more explicit-but-custom tooling — accepted per explicit user request to reduce per-release PR count.

## Migration Plan

1. Verify GitVersion's exact monorepo-scoping config syntax and confirm this repo's PR merge strategy (merge-commit vs. squash) against live settings.
2. Add `gopints/GitVersion.yml`, `gopints/cliff.toml`, `web/GitVersion.yml`, `web/cliff.toml`.
3. Add `release-prepare-backend.yml` / `release-prepare-web.yml` implementing the propose/finalize branch described above.
4. Remove `release-please-config.json`, `.release-please-manifest.json`, `.github/workflows/release-please.yml`.
5. Remove the `bump-infra-version` job from `release.yml`; leave `prepare`/`build-server`/`build-web`/`goreleaser` untouched.
6. Update `README.md`'s CI/CD section to describe the new single-PR flow.
7. Hand-trace both workflow branches (ordinary push, merge-of-release-prep-branch) for both packages; validate all new/edited YAML and TOML parse correctly.
8. Note in the PR/commit description that full end-to-end behavior (a real release-prep PR opening, a real merge correctly triggering tag+release, GitVersion computing the expected bump) can only be confirmed by observing real pushes/merges on `main` — not reproducible in local sandbox testing.

Rollback: revert the added/removed files. Unlike the release-please manifest (a stateful file tracking "current version" that would need reconciling on rollback), GitVersion computes purely from git tag/commit history, so rollback has no extra state-reconciliation step — re-adding `release-please-config.json`/`.release-please-manifest.json` with the last-known versions is sufficient if reverting.

## Open Questions

- Exact GitVersion 6 preset/mode name for "trunk-based, conventional-commit-driven, no branch-based GitFlow" — verification task, not a blocking unknown.
- Whether this repo's merge strategy is fixed at merge-commit (confirmed) or configurable per-PR (i.e., could someone squash-merge a release-prep PR by mistake, breaking detection) — worth a guard: the finalize-detection step should fail loudly (not silently skip tagging) if it can't find the expected substring pattern *and* the diff clearly looks like a release-prep merge (e.g., branch name still visible via `gh pr list --search`), rather than assuming "no match = ordinary push" unconditionally.
