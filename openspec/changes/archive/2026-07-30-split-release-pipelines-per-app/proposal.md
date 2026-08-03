## Why

Today the whole monorepo shares a single release-please-managed version (`.release-please-manifest.json`: `{".": "0.4.0"}`). A change to only `web/` still bumps and releases `gopints/` (and vice versa), and CI (`ci.yml`) runs both the Go and frontend jobs on every PR regardless of which app actually changed. As the two apps' release cadences diverge — the server/agent and the web UI don't need to ship in lockstep — this coupling adds noise (irrelevant CI runs, irrelevant version bumps, irrelevant CHANGELOG entries) and makes it harder to tell what actually shipped in a given release. Additionally, the checked-in Kubernetes manifests (`gopints/infrastructure/deploy/deployment.yaml`, `web/infrastructure/deploy/deployment.yaml`) hardcode an image tag that nothing currently updates — they've been stuck at `0.4.0` since they were added.

## What Changes

- Reconfigure release-please from a single repo-wide package to two independent, path-scoped packages: `gopints/` and `web/`, each with its own manifest entry, semver, and CHANGELOG, via `separate-pull-requests: true` so a change to one app's files never opens or touches the other app's release PR.
- Split `ci.yml` into two path-filtered workflows (backend, frontend) so a PR only runs the checks relevant to what it touched.
- Update `release.yml`'s container-publish job to build/push only the image for whichever package's release actually fired (currently it unconditionally builds both `gopints-web` and `gopints-server` off one shared tag), and gate the GoReleaser (agent binary) job to only run for `gopints/` releases.
- Add a new job (gated the same way, per package) that, after a release publishes, updates that app's `infrastructure/deploy/deployment.yaml` image tag to the newly released version and opens a pull request with the change — no auto-merge, a human reviews and merges it.
- **BREAKING** (to the release process, not the app): existing release tags/CHANGELOG history stay as-is, but going forward the repo produces two independent tag streams (e.g. `gopints-v0.4.1`, `web-v0.5.0`) instead of one shared `vX.Y.Z`. Anything that currently assumes a single repo-wide version (docs, scripts, external references to release tags) needs to be checked.

## Capabilities

### New Capabilities
- `infra-version-bump`: Automatically opening a PR to update an app's Kubernetes manifest image tag to match its latest published release.

### Modified Capabilities
(none — no existing specs in this repo cover CI/release behavior yet; this is new ground for the spec set)

## Impact

- `release-please-config.json`, `.release-please-manifest.json` — restructured to two packages.
- `.github/workflows/ci.yml` — replaced by two path-filtered workflows (or path-filtered jobs within one file; decided in design.md).
- `.github/workflows/release.yml` — container/goreleaser jobs gated by which package's tag fired; new infra-bump job added.
- `gopints/infrastructure/deploy/deployment.yaml`, `web/infrastructure/deploy/deployment.yaml` — target of the new automated version-bump PRs (no manual edits needed going forward).
- No application code changes (`gopints/` Go source, `web/` frontend source untouched) — this is entirely CI/release tooling.
- `docker-compose.yml` and `README.md`'s deployment docs reference `GOPINTS_TAG`/version — worth a quick check that nothing there assumes a single shared tag exists across both images.
