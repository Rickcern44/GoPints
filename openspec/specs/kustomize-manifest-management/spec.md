## Purpose

Defines how each app's Kubernetes deployment manifests are managed via kustomize: every app deploy directory carries a `kustomization.yaml` that declares its resources and image reference, and once that file exists, container image tags are changed exclusively through kustomize's structural image transformer rather than direct text edits — ensuring tag bumps are verifiable and fail loudly instead of silently no-op'ing.

## Requirements

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

### Requirement: ArgoCD renders each deploy directory via kustomize
The ArgoCD Application syncing each app's deploy directory SHALL be configured with a `kustomize` source (or no explicit source type, allowing auto-detection), not an explicit `directory` source, so that `kustomization.yaml` is rendered as a build input rather than applied as a literal Kubernetes resource.

#### Scenario: kustomization.yaml present under a directory-mode Application
- **WHEN** an ArgoCD Application's `spec.source` explicitly sets `directory`, and the target path contains a `kustomization.yaml`
- **THEN** ArgoCD attempts to apply `kustomization.yaml` itself as a raw resource and fails sync with an error resembling "could not find kustomize.config.k8s.io/Kustomization ... CRD"; this is a misconfiguration that SHALL be corrected by switching the Application to a `kustomize` source

#### Scenario: kustomize source configured correctly
- **WHEN** an ArgoCD Application's `spec.source` sets `kustomize: {}` (or omits an explicit source type) against a path containing `kustomization.yaml`
- **THEN** ArgoCD renders the directory via `kustomize build` and applies the resulting resources, matching what `kustomize build <dir>` produces locally
