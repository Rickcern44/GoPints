## 1. Verify ArgoCD render path

- [x] 1.1 Confirm (against the live ArgoCD Application specs, outside this repo) that the Applications syncing `gopints/infrastructure/deploy/` and `web/infrastructure/deploy/` already render via kustomize, or will once `kustomization.yaml` is present — resolve the design's open question before proceeding

## 2. Add kustomize bases

- [x] 2.1 Add `gopints/infrastructure/deploy/kustomization.yaml` listing `deployment.yaml`, `service.yaml`, `pvc.yaml` as `resources` and an `images` entry pinning `ghcr.io/rickcern44/gopints-server` at its current tag
- [x] 2.2 Add `web/infrastructure/deploy/kustomization.yaml` listing `deployment.yaml`, `service.yaml`, `ingress.yaml`, `namespace.yaml` as `resources` and an `images` entry pinning `ghcr.io/rickcern44/gopints-web` at its current tag
- [x] 2.3 Run `kustomize build` against both directories and diff the output against `kubectl apply --dry-run=client -f <dir>` of the raw manifests to confirm equivalence
- Note: `kubectl` was not available locally; equivalence was confirmed by running `kustomize build` on both directories and comparing the rendered resources/fields directly against the raw manifest files (same kinds, names, namespaces, and fields, with the image tag sourced from `kustomization.yaml`).

## 3. Update release workflow

- [x] 3.1 Add a `kustomize` CLI setup step to the `bump-infra-version` job in `.github/workflows/release.yml`
- [x] 3.2 Replace the `sed -i -E` image-tag substitution in the "Update image tag and commit" step with `kustomize edit set image <image>=<image>:<version>`, run from the target package's deploy directory
- [x] 3.3 Add a verification step that runs `kustomize build <dir>` and asserts the new tag appears in the rendered output, failing the job if not
- [x] 3.4 Update the "already at version; nothing to commit" no-op check to diff `kustomization.yaml` instead of `deployment.yaml`

## 4. Validate end to end

- [x] 4.1 Dry-run the updated `bump-infra-version` job (e.g. via `act` or a scratch branch) against a fake version bump for one package and confirm the commit only touches `kustomization.yaml`
- Note: `act` was not used; the exact shell logic from the "Update image tag and commit" step was run locally against `gopints/infrastructure/deploy/` with a throwaway version (`0.6.0-test`), confirmed only `kustomization.yaml` changed and the rendered image tag matched, then reverted to the real pinned tag (`0.5.1-10`).
- [x] 4.2 Confirm `openspec/specs/infra-version-bump/spec.md`'s "No automatic cluster deployment" requirement still holds — no step in the updated job applies manifests to a cluster
- Note: verified via `grep` for `kubectl apply`/`argocd`/`helm` across the updated `release.yml` — no matches.
