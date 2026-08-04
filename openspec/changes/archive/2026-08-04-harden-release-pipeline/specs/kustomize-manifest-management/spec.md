## ADDED Requirements

### Requirement: ArgoCD renders each deploy directory via kustomize
The ArgoCD Application syncing each app's deploy directory SHALL be configured with a `kustomize` source (or no explicit source type, allowing auto-detection), not an explicit `directory` source, so that `kustomization.yaml` is rendered as a build input rather than applied as a literal Kubernetes resource.

#### Scenario: kustomization.yaml present under a directory-mode Application
- **WHEN** an ArgoCD Application's `spec.source` explicitly sets `directory`, and the target path contains a `kustomization.yaml`
- **THEN** ArgoCD attempts to apply `kustomization.yaml` itself as a raw resource and fails sync with an error resembling "could not find kustomize.config.k8s.io/Kustomization ... CRD"; this is a misconfiguration that SHALL be corrected by switching the Application to a `kustomize` source

#### Scenario: kustomize source configured correctly
- **WHEN** an ArgoCD Application's `spec.source` sets `kustomize: {}` (or omits an explicit source type) against a path containing `kustomization.yaml`
- **THEN** ArgoCD renders the directory via `kustomize build` and applies the resulting resources, matching what `kustomize build <dir>` produces locally
