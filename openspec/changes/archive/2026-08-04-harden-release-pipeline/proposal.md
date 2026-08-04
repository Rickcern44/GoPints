## Why

Cutting the first real release through the new kustomize-based pipeline (`2026-08-03-kustomize-image-tag-automation`) surfaced three latent gaps in production, in this order:

1. ArgoCD's `gopints`/`web` Applications were created by hand through the ArgoCD UI (no Application manifest is checked into this repo) with an explicit `directory` source. That source type never invokes kustomize, so ArgoCD tried to apply `kustomization.yaml` itself as a literal Kubernetes resource and failed sync with "could not find kustomize.config.k8s.io/Kustomization ... CRD". This is exactly the follow-up the prior change's proposal flagged and deferred: "if ArgoCD Application manifests pin `directory`/`ksonnet` render mode instead of kustomize, a follow-up change may be needed but is out of scope here."
2. `release-prepare-web.yml` and `release-prepare-backend.yml` both interpolated `${{ github.event.head_commit.message }}` directly into a `run:` script's `MSG="..."` assignment. A merge commit title containing a literal `"` (the `feat(web): redesign frontend as unified "Precision Tap Console"` PR) broke out of the quoted assignment and was parsed as shell syntax, failing the release-prep merge-detection step with `Tap: command not found`.
3. Neither `gopints/infrastructure/deploy/deployment.yaml` nor `web/infrastructure/deploy/deployment.yaml` sets `imagePullPolicy`, which defaults to `IfNotPresent`. Because image tags are mutable — a stray manual `workflow_dispatch` run of `Publish` re-pushed `gopints-web:0.5.2` built from stale `main` before the real redesign release — a node that had already cached that tag never re-pulled the real release's content. The real `web-v0.5.2` release completed successfully end-to-end with no failing step anywhere, yet the deployed pod kept serving the old build.

None of these were caught by the prior change's validation, which had no live cluster or real release to exercise against. They were each diagnosed and fixed directly (PRs #35, #37, plus a live ArgoCD UI edit) without going through a proposal first. Recording them here after the fact so the causes, fixes, and the spec invariants they establish aren't lost.

## What Changes

- **(Already applied, live, not tracked in this repo)** ArgoCD's `gopints` and `web` Application `spec.source` reconfigured from an explicit `directory` source to a `kustomize` source, so `kustomization.yaml` is rendered as a build input instead of applied as a raw manifest.
- `.github/workflows/release-prepare-web.yml`, `.github/workflows/release-prepare-backend.yml`: the "Detect release-prep merge" step now passes the triggering commit message via `env: MSG: ${{ github.event.head_commit.message }}` instead of splicing it into the script body.
- `gopints/infrastructure/deploy/deployment.yaml`, `web/infrastructure/deploy/deployment.yaml`: add `imagePullPolicy: Always` to each Deployment's container.
- **BREAKING**: none — all three are corrections to existing intended behavior; no capability's external contract changes shape.

## Capabilities

### New Capabilities
- `deployment-image-freshness`: Both apps' Deployments always re-pull image content for their tag on pod start, rather than trusting a locally cached image of the same tag string.

### Modified Capabilities
- `kustomize-manifest-management`: adds the requirement that the ArgoCD Application syncing each deploy directory must render via kustomize, not an explicit `directory` source.
- `infra-version-bump`: adds the requirement that release-prep merge detection safely handles arbitrary commit message content.

## Impact

- ArgoCD `gopints`/`web` Applications (cluster-side config, outside this repo) — source type corrected; no repo file changed for this part.
- `.github/workflows/release-prepare-web.yml`, `.github/workflows/release-prepare-backend.yml` — `env:` block added to the "Detect release-prep merge" step (PR #35).
- `gopints/infrastructure/deploy/deployment.yaml`, `web/infrastructure/deploy/deployment.yaml` — `imagePullPolicy: Always` added (PR #37).
- No changes to Go, SvelteKit, or SQLite code paths.
