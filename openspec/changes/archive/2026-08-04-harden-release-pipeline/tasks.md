## 1. ArgoCD render mode

- [x] 1.1 Diagnose ArgoCD sync error "could not find kustomize.config.k8s.io/Kustomization ... CRD" as the Application's `spec.source` being pinned to `directory` instead of `kustomize`
- [x] 1.2 Reconfigure the `gopints` and `web` ArgoCD Applications to a `kustomize` source, via the ArgoCD UI (no repo change — no Application manifest is checked in)
- [x] 1.3 Confirm sync succeeds against the rendered kustomize output

## 2. Release-prep commit message handling

- [x] 2.1 Diagnose `release-prepare-web.yml`'s "Detect release-prep merge" step failing (`Tap: command not found`) as `${{ github.event.head_commit.message }}` being spliced directly into the `run:` script, breaking on a merge commit title containing `"`
- [x] 2.2 Fix `release-prepare-web.yml` to pass the message via `env: MSG: ...` instead
- [x] 2.3 Apply the same fix to `release-prepare-backend.yml`, which had the identical pattern and would have hit the same failure on a future gopints release
- [x] 2.4 Merge (PR #35) and confirm a subsequent push with a quote-containing commit message correctly resolves `is_finalize`

## 3. Deployment image freshness

- [x] 3.1 Diagnose the deployed `gopints-web` pod serving pre-redesign content despite a successful `web-v0.5.2` release, tracing it to a stray manual `workflow_dispatch` run that had already pushed a `0.5.2`-tagged image built from stale `main`, combined with `imagePullPolicy` defaulting to `IfNotPresent`
- [x] 3.2 Add `imagePullPolicy: Always` to `gopints/infrastructure/deploy/deployment.yaml` and `web/infrastructure/deploy/deployment.yaml`
- [x] 3.3 Verify `kustomize build` for both directories renders `imagePullPolicy: Always`
- [x] 3.4 Merge (PR #37)

## 4. Record spec coverage

- [x] 4.1 Add an ArgoCD render-mode requirement to `kustomize-manifest-management`
- [x] 4.2 Add a commit-message-handling requirement to `infra-version-bump`
- [x] 4.3 Add a new `deployment-image-freshness` capability covering `imagePullPolicy: Always`
