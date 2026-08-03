## Why

The release pipeline currently bumps the container image tag in each app's deployment manifest (`gopints/infrastructure/deploy/deployment.yaml`, `web/infrastructure/deploy/deployment.yaml`) with a `sed -i -E` regex against the raw YAML (`.github/workflows/release.yml`, `bump-infra-version` job). This works today but is brittle: it depends on the manifest's exact formatting, silently no-ops if the `image:` line ever shifts shape, and gives no structural validation that the result is still well-formed YAML. Kustomize ships a purpose-built, structural image-tag transformer (`kustomize edit set image`) that is the standard tool for exactly this GitOps write-back pattern and is what ArgoCD already understands natively.

## What Changes

- Add a `kustomization.yaml` to each app's deploy directory (`gopints/infrastructure/deploy/`, `web/infrastructure/deploy/`) declaring its existing resources (`deployment.yaml`, `service.yaml`, `pvc.yaml` for gopints; `deployment.yaml`, `service.yaml`, `ingress.yaml`, `namespace.yaml` for web) and an `images:` entry pinning the current tag.
- Replace the `sed -i -E` step in the `bump-infra-version` job of `.github/workflows/release.yml` with `kustomize edit set image <image>=<image>:<version>`, run from each app's deploy directory.
- No change to the ArgoCD application targets is required for this proposal — ArgoCD's kustomize support renders `kustomization.yaml` automatically when present in the sync path; if ArgoCD Application manifests pin `directory`/`ksonnet` render mode instead of kustomize, a follow-up change may be needed but is out of scope here.
- **BREAKING**: none — output of `kustomize build` for each dir is byte-for-byte equivalent to the current raw manifests (same resources, same fields), so this is a structural/tooling change only.

## Capabilities

### New Capabilities
- `kustomize-manifest-management`: Deploy manifests for each app are organized as a kustomize base (`kustomization.yaml` + resources); image tags are read and written exclusively through kustomize's image transformer rather than text-editing YAML directly.

### Modified Capabilities
- `infra-version-bump`: The "Manifest updated" scenario's mechanism changes from a regex-based `sed` substitution to `kustomize edit set image`, which fails loudly (non-zero exit) if the target image reference isn't declared in `kustomization.yaml`, instead of silently no-op'ing on a mismatched regex.

## Impact

- `.github/workflows/release.yml` — `bump-infra-version` job: swap the "Update image tag and commit" step's `sed` command for a `kustomize edit set image` invocation; requires the `kustomize` CLI available on the runner (via `pip`/GitHub Action or downloading the release binary).
- `gopints/infrastructure/deploy/` — new `kustomization.yaml`.
- `web/infrastructure/deploy/` — new `kustomization.yaml`.
- No changes to Go, SvelteKit, or SQLite code paths; this is deploy-tooling only.
