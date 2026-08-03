## Context

`gopints/infrastructure/deploy/` and `web/infrastructure/deploy/` each hold plain Kubernetes manifests (no `kustomization.yaml`). ArgoCD Applications sync directly against these directories. The `bump-infra-version` job in `.github/workflows/release.yml` runs after a package's image is confirmed pushed to GHCR, and mutates the relevant `deployment.yaml`'s `image:` line with:

```bash
sed -i -E "s#(image: ghcr.io/rickcern44/gopints-server:).*#\1${VERSION}#" gopints/infrastructure/deploy/deployment.yaml
```

then commits the change directly to `main`. This works but has no structural awareness of the YAML — a reformatted manifest, a renamed container, or a second container in the pod silently breaks the regex with no error.

## Goals / Non-Goals

**Goals:**
- Give each app's deploy directory a `kustomization.yaml` so image tags can be set structurally via `kustomize edit set image`.
- Make the CI bump step fail loudly (non-zero exit) if the image reference it's targeting doesn't exist in the kustomization, instead of silently no-op'ing.
- Keep `kustomize build` output for both apps byte-equivalent (module resources, fields, ordering of keys aside) to what's applied today, so this is a no-behavior-change swap from ArgoCD's perspective.

**Non-Goals:**
- Introducing kustomize overlays (dev/staging/prod) — there's only one environment today; this only adds the base.
- Changing how ArgoCD Applications are defined or how sync/rollback works.
- Automating deploys to the cluster — `infra-version-bump`'s existing "no automatic cluster deployment" invariant is unchanged; this only touches how the manifest *file* is edited pre-sync.

## Decisions

**Use the `kustomize` CLI's `edit set image` subcommand, not the Go library (`sigs.k8s.io/kustomize/api`) embedded in a custom tool.**
The CLI is a single static binary, needs no Go toolchain step in the release workflow (which otherwise only needs Go for GoReleaser on backend releases, not on web releases), and `kustomize edit set image` is the exact command this GitOps write-back pattern is built for. Writing a Go program against the kustomize API would add a maintenance surface for no behavioral gain, since the requirement is purely "update one image tag field."

**One `kustomization.yaml` per app deploy directory, not a shared root kustomization.**
`gopints/infrastructure/deploy/` and `web/infrastructure/deploy/` are already independent per-package deploy roots (mirroring the independent per-package release model in `infra-version-bump`). A shared root kustomization would couple the two apps' manifest builds together, contradicting the "package-scoped" principle already established for releases and publishing.

**Pin the `images:` transformer using the same `ghcr.io/rickcern44/<image>` reference already in `deployment.yaml`, with `newTag` matching the current tag at migration time.**
This keeps `kustomize build` output identical to the current raw manifest on day one — no manifest content changes ship as part of this change, only the addition of `kustomization.yaml` alongside the existing files.

## Risks / Trade-offs

- [Risk] CI runner doesn't have `kustomize` pre-installed → build step fails. **Mitigation**: pin a `kustomize` release binary download step (or `azure/setup-kustomize`-style action) in the `bump-infra-version` job, matching the version used for any local verification.
- [Risk] ArgoCD Application's sync path doesn't recognize `kustomization.yaml` and instead applies raw manifests directly (skipping the kustomize render). **Mitigation**: out of scope for this change per the proposal, but flagged as a follow-up check — if ArgoCD isn't already kustomize-aware for this path, the bumped tag in `kustomization.yaml` would never reach the rendered `deployment.yaml`'s effective image. This should be verified against the live ArgoCD Application spec before merging the workflow change.
- [Risk] `kustomize edit set image` requires the target image to already be declared in `kustomization.yaml`'s `images:` list, or it silently adds a new entry rather than erroring — differs from "fails loudly" framing in the proposal. **Mitigation**: task list should include a CI assertion step (`kustomize build | grep <expected-tag>`) that verifies the rendered output actually contains the new tag before committing, closing this gap explicitly rather than trusting the command's own exit code.

## Migration Plan

1. Add `kustomization.yaml` to `gopints/infrastructure/deploy/` and `web/infrastructure/deploy/`, each listing existing resources and pinning the current live image tag — verify `kustomize build <dir>` output matches today's raw-`kubectl apply` manifests.
2. Swap the `sed` step in `bump-infra-version` for `kustomize edit set image`, plus a rendered-output assertion (see Risks).
3. Ship both directory additions and the workflow change in the same PR so the workflow is never pointed at a directory without a `kustomization.yaml`.
4. Rollback: revert the workflow-file commit; the `kustomization.yaml` files are additive and inert if the `sed` step is reverted alongside them (they don't change `kubectl apply -f` behavior against the individual files).

## Open Questions

- Does the current ArgoCD Application spec already build via kustomize for these paths, or apply raw manifests? Needs to be checked against the live ArgoCD config (not in this repo) before this ships.
