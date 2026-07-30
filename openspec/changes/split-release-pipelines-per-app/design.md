## Context

Three workflows currently govern releases: `release-please.yml` (opens/updates a release PR on every push to `main`, driven by `release-please-config.json` + `.release-please-manifest.json`, both currently single-package `"."`), `release.yml` (fires on `release: published`, unconditionally builds+pushes both `gopints-web` and `gopints-server` Docker images off the one shared tag, plus runs GoReleaser for agent binaries), and `ci.yml` (runs a `backend` and a `frontend` job on every PR, no path filtering beyond excluding `**/infrastructure/**` entirely). Both apps already have their own `infrastructure/deploy/deployment.yaml` with a hardcoded image tag that nothing updates automatically. `docker-compose.yml` and `README.md` currently assume one shared `GOPINTS_TAG` applies to both images.

Decided with the user before writing this design: independent per-package versions (not just split CI triggers), infra-bump lands as an auto-generated PR (not a direct commit), and scope is limited to keeping the manifest's image tag current in git — not applying it to a live cluster (no GitOps controller exists in this repo today; the k8s manifests are checked in but nothing currently applies them automatically).

## Goals / Non-Goals

**Goals:**
- `gopints/` and `web/` release independently: a change to one never bumps, opens a release PR for, or publishes a container for the other.
- CI only runs the checks relevant to what a PR actually touched.
- Each app's k8s manifest image tag automatically tracks that app's latest published release, via a reviewable PR.

**Non-Goals:**
- Not applying manifests to a real cluster — no cluster credentials are being added to CI as part of this change.
- Not changing the agent's install script (`scripts/setup-agent.sh`) or how it discovers "latest" releases via GitHub's release API — worth a follow-up check but not touching it here unless the per-package tag format breaks its assumptions (see Risks).
- Not changing application code in either `gopints/` or `web/` — this is CI/release tooling only.

## Decisions

**Two release-please packages, path-scoped, with `separate-pull-requests: true`.**
```json
{
  "release-type": "simple",
  "separate-pull-requests": true,
  "packages": {
    "gopints": {},
    "web": { "release-type": "node" }
  }
}
```
`web/` uses release-please's `node` release-type so `web/package.json`'s `"version"` field gets bumped automatically as part of the release (a free correctness win — that field exists and currently just sits at a stale manually-set `0.0.1`). `gopints/` uses the generic `simple` type — there's no analogous checked-in version field to bump there (the Go binary's version comes from `-ldflags -X main.version=...` at build time, not a source file), so `simple` (bump nothing but the manifest + CHANGELOG) is the right fit, not `go`.

`.release-please-manifest.json` becomes `{"gopints": "0.4.0", "web": "0.4.0"}` — both start from the current shared version, then diverge independently from here.

**Per-package tags, verified during implementation rather than assumed.**
Release-please's default tag format for named (non-`"."`) packages is `<package>-v<version>` (e.g. `gopints-v0.4.1`, `web-v0.5.0`) — but this is exactly the kind of detail worth confirming against the actual `googleapis/release-please-action@v4` behavior/docs during implementation rather than trusting from memory, since every downstream conditional (container job, goreleaser job, infra-bump job) depends on matching this prefix correctly. Flagged as an explicit early task.

**`release.yml` gates every job on the tag prefix.**
The `containers` job's matrix currently always builds both `web` and `server`. Add a per-matrix-entry (or per-job) condition checking `startsWith(github.event.release.tag_name, 'gopints-')` / `'web-'`, so a `web-vX.Y.Z` release only builds the web image, and a `gopints-vX.Y.Z` release only builds the server image. The `goreleaser` job (agent binaries) gets `if: startsWith(github.event.release.tag_name, 'gopints-')` — no reason to cut new agent binaries on a web-only release. The version passed to `docker/metadata-action` and `--build-arg VERSION=` needs the package prefix stripped before use (e.g. via a small `steps` output extracting everything after the first `-v`), since the Docker image tags and the `main.version` ldflag should be the bare semver, not the prefixed release-please tag.

**CI split into two path-filtered workflows, not just path-filtered jobs in one file.**
Separate `ci-backend.yml` (triggers on `pull_request` with `paths: ['gopints/**']`) and `ci-frontend.yml` (`paths: ['web/**']`), replacing today's single `ci.yml`. Two files (rather than one file with per-job `paths-ignore`) map more directly to "split the CI to individual apps," and make it simpler to see status of each app's CI independently in the PR checks list. Alternative considered: keep one `ci.yml` with path-filtered jobs — rejected only because it's marginally less legible as "one CI per app" than literally separate files, not because of any functional difference.

**Infra-bump job added to `release.yml`, gated identically to the containers job.**
A new `bump-infra-version` job runs `if: startsWith(github.event.release.tag_name, <package prefix>)`, checks out `main`, updates the `image:` line in that package's `infrastructure/deploy/deployment.yaml` to the new version (a scoped `sed`/`yq` replace — the line already has a well-known exact format: `image: ghcr.io/rickcern44/gopints-<name>:<version>`), and opens a PR via `peter-evans/create-pull-request` (a widely-used, already-trusted action for exactly this "bot commits a small change, opens a PR" pattern — avoids hand-rolling git branch/commit/push logic in the workflow). No auto-merge; a human reviews and merges.

**`docker-compose.yml` / `README.md` move from one shared `GOPINTS_TAG` to two independent tag variables.**
Since the two images now version independently, a single `GOPINTS_TAG` pin no longer makes sense for both. Rename to `GOPINTS_SERVER_TAG` / `GOPINTS_WEB_TAG` (each still defaulting to `latest`), and update the README's pinning example accordingly. This is a breaking change to the compose interface, called out in the proposal.

## Risks / Trade-offs

- **[Risk]** `scripts/setup-agent.sh` downloads "the latest GitHub Release" for the agent binary — if it queries the repo's single "latest" release via GitHub's API, that API concept doesn't cleanly map to two independent tag streams (GitHub's "latest release" is a single repo-wide flag, typically whichever released most recently by timestamp, which could now sometimes be a `web-vX` release with no agent binary attached at all). → **Mitigation**: audit the script during implementation; if it relies on "latest" generically, it likely needs to filter to tags matching `gopints-v*` specifically (GitHub's release API supports listing releases and filtering client-side). Flagged as a task, not silently ignored.
- **[Risk]** Existing consumers/bookmarks referencing a bare `vX.Y.Z` tag scheme will need to adjust to `gopints-vX.Y.Z` / `web-vX.Y.Z`. → **Mitigation**: this is a solo/small-scale project per prior context in this session; acceptable one-time adjustment, called out as **BREAKING** in the proposal.
- **[Trade-off]** Two release PRs instead of one means two CHANGELOGs to check when reviewing "what shipped recently" instead of one unified view. Accepted as the direct, intended consequence of independent release cadence.

## Migration Plan

1. Update `release-please-config.json` / `.release-please-manifest.json` to the two-package structure (both starting at `0.4.0`).
2. Split `ci.yml` into `ci-backend.yml` / `ci-frontend.yml`.
3. Update `release.yml`: gate `containers` matrix entries, gate `goreleaser`, strip tag prefix for version extraction, add `bump-infra-version` job.
4. Update `docker-compose.yml` / `README.md` for the two-tag-variable scheme.
5. Audit `scripts/setup-agent.sh` for "latest release" assumptions.
6. Merge to `main`; next push should produce the new per-package release-please PR behavior — verify by observing the next real release PR(s) rather than assuming success from workflow YAML alone.

Rollback: revert the config/manifest/workflow files; release-please and the existing single-package history remain intact since no past tags/releases are modified by this change.

## Open Questions

- Exact default tag format for release-please's named-package mode — treated as a verification task (see Migration Plan step 6 / tasks.md), not a blocking unknown, since it's cheap to confirm against the live action's behavior or docs during implementation.
