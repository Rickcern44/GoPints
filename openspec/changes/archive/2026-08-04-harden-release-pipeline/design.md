## Context

This change documents three fixes discovered and applied in immediate succession while cutting the first release through the `2026-08-03-kustomize-image-tag-automation` pipeline against a real cluster. Each was diagnosed from a live failure (an ArgoCD sync error, a failed GitHub Actions run, and a deployed pod not matching its release), not from planning ahead of time.

## Goals / Non-Goals

**Goals:**
- Capture the causes and fixes so they aren't rediscovered from scratch on the next incident.
- Add spec coverage for invariants that were previously implicit (ArgoCD's render mode, commit-message handling, image pull freshness).

**Non-Goals:**
- Fixing the in-progress-pour data-loss gap identified during the same investigation (an in-flight pour only lives in the server's in-memory `Meter` until it ends, and is lost on any restart). That's tracked separately and intentionally out of scope here — this change is deploy-pipeline hardening only.
- Checking an ArgoCD Application manifest into this repo. The Applications remain UI-created/cluster-side; this change only corrects their source-type configuration.

## Decisions

**Fix `imagePullPolicy` rather than making tags immutable (e.g. digest-pinning).**
Switching to immutable, digest-pinned references would also solve the staleness problem and is arguably more correct, but it's a bigger change to the release pipeline's tag-based `kustomize edit set image` flow. `imagePullPolicy: Always` is a one-line fix that closes the actual failure mode (a node silently trusting a stale local cache) without touching how images are tagged or referenced elsewhere in the pipeline.

**Document the ArgoCD source-type fix here even though it isn't a repo change.**
The Application's `spec.source` lives only in the cluster, so there's no file to diff — but the failure mode and fix are exactly the scenario the prior change's proposal explicitly deferred, so recording it against `kustomize-manifest-management` keeps that capability's spec honest about what's actually required for kustomize rendering to work, not just what's required inside this repo.
