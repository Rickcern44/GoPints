## MODIFIED Requirements

### Requirement: No automatic cluster deployment
The system SHALL NOT apply the Kubernetes deployment manifest to a live cluster as part of its release automation; applying it (e.g. via ArgoCD sync) remains a separate action outside this pipeline's scope. Where a package's deploy directory has a `kustomization.yaml`, the release pipeline SHALL update that package's infrastructure deployment manifest image tag structurally via `kustomize edit set image`, rather than a text-based edit (e.g. `sed`) of the manifest file, so that a mismatched image reference or manifest shape causes the pipeline step to fail rather than silently no-op.

#### Scenario: Manifest updated
- **WHEN** a package's infrastructure deployment manifest image tag is updated by the release pipeline
- **THEN** the update is performed via `kustomize edit set image` against that package's `kustomization.yaml`, and no process within this repository's CI applies the resulting manifest to a live cluster

#### Scenario: Image tag update fails to apply
- **WHEN** the release pipeline runs `kustomize edit set image` for a package whose deploy directory's `kustomization.yaml` does not declare the target image, or whose rendered output does not contain the expected new tag
- **THEN** the `bump-infra-version` job fails, and no commit updating the manifest image tag is created
