## ADDED Requirements

### Requirement: Kustomize base per app deploy directory
Each app's Kubernetes deploy directory (`gopints/infrastructure/deploy/`, `web/infrastructure/deploy/`) SHALL contain a `kustomization.yaml` declaring that directory's manifests as `resources` and its container image as an `images` entry, such that `kustomize build <dir>` renders a manifest set equivalent to applying the directory's raw YAML files directly.

#### Scenario: Building the gopints deploy directory
- **WHEN** `kustomize build gopints/infrastructure/deploy/` is run
- **THEN** the rendered output includes the `gopints-server` Deployment, Service, and PVC with the same fields as the directory's raw manifest files, with the image reference sourced from the `kustomization.yaml` `images` entry

#### Scenario: Building the web deploy directory
- **WHEN** `kustomize build web/infrastructure/deploy/` is run
- **THEN** the rendered output includes the `gopints-web` Deployment, Service, Ingress, and Namespace with the same fields as the directory's raw manifest files, with the image reference sourced from the `kustomization.yaml` `images` entry

### Requirement: Image tags are set only through kustomize's image transformer
Once a `kustomization.yaml` exists for a deploy directory, that directory's container image tag SHALL be changed exclusively via `kustomize edit set image` (or an equivalent structural edit to the `images` entry), not via direct text edits to the `deployment.yaml`'s `image:` field.

#### Scenario: Setting a new image tag
- **WHEN** `kustomize edit set image ghcr.io/rickcern44/gopints-server=ghcr.io/rickcern44/gopints-server:0.6.0` is run from `gopints/infrastructure/deploy/`
- **THEN** the `kustomization.yaml`'s `images` entry is updated to the new tag, and `kustomize build` subsequently renders the Deployment with that tag, without any edit to `deployment.yaml` itself

#### Scenario: Target image not declared in kustomization
- **WHEN** `kustomize edit set image` is run against an image reference that has no corresponding entry in the deploy directory's `kustomization.yaml`
- **THEN** the CI step that verifies rendered output (e.g. `kustomize build | grep <expected-tag>`) fails, so the tag bump is not silently dropped
