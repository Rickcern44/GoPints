## 1. Independent release-please packages

- [x] 1.1 Verify release-please's default tag format for named (non-`"."`) packages against `googleapis/release-please-action@v4` docs — confirm the `<package>-v<version>` assumption before wiring conditionals to it
- [x] 1.2 Update `release-please-config.json`: two packages (`gopints` as `simple`, `web` as `node`), `separate-pull-requests: true`
- [x] 1.3 Update `.release-please-manifest.json`: `{"gopints": "0.4.0", "web": "0.4.0"}`

## 2. Split CI by app

- [x] 2.1 Create `ci-backend.yml`: the existing `backend` job from `ci.yml`, triggered on `pull_request` with `paths: ['gopints/**']`
- [x] 2.2 Create `ci-frontend.yml`: the existing `frontend` job from `ci.yml`, triggered on `pull_request` with `paths: ['web/**']`
- [x] 2.3 Remove `ci.yml`

## 3. Gate release.yml by package

- [x] 3.1 Add a step (or reusable expression) extracting the bare semver and package name from `github.event.release.tag_name` (e.g. strip the `gopints-`/`web-` prefix)
- [x] 3.2 Gate the `containers` job's `web` matrix entry on the tag being a `web-` release, and the `server` matrix entry on the tag being a `gopints-` release
- [x] 3.3 Gate the `goreleaser` job on the tag being a `gopints-` release
- [x] 3.4 Update `docker/metadata-action` tag patterns and the `VERSION` build-arg to use the stripped bare semver, not the raw prefixed tag

## 4. Infra version-bump job

- [x] 4.1 Add a `bump-infra-version` job to `release.yml`, gated per package the same way as the containers job
- [x] 4.2 Update the `image:` line in the corresponding `infrastructure/deploy/deployment.yaml` to the new bare semver (scoped sed/yq replace, not a full file rewrite)
- [x] 4.3 Open a PR with the change via `peter-evans/create-pull-request` (no auto-merge)

## 5. Docker Compose / docs

- [x] 5.1 Rename `GOPINTS_TAG` to `GOPINTS_SERVER_TAG` and `GOPINTS_WEB_TAG` in `docker-compose.yml` (each defaulting to `latest`)
- [x] 5.2 Update `README.md`'s deployment section to reflect the two independent tag variables

## 6. Agent install script audit

- [x] 6.1 Read `scripts/setup-agent.sh`'s "latest release" logic; determine whether it needs to filter to `gopints-v*` tags specifically now that the repo has two independent tag streams
- [x] 6.2 Update the script if needed

## 7. Verification

- [x] 7.1 Validate all edited/new YAML files parse correctly (e.g. `python3 -c "import yaml; yaml.safe_load(open(...))"` or equivalent) for every workflow file touched
- [x] 7.2 Trace through `release.yml`'s conditionals by hand for both a hypothetical `gopints-v0.4.1` tag and a `web-v0.5.0` tag, confirming exactly one image (and goreleaser only for gopints) would build in each case
- [x] 7.3 Confirm `ci-backend.yml`/`ci-frontend.yml` path filters correctly scope to `gopints/**` and `web/**` respectively, with no gap or unwanted overlap
- [x] 7.4 Note in the PR/commit description that full end-to-end verification (an actual release-please PR opening correctly per package, an actual tagged release firing the gated jobs correctly) can only be confirmed by observing real merges to `main` in the live repo — not something reproducible in local sandbox testing
