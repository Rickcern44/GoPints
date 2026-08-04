## Purpose

Defines pull-policy behavior for both apps' container images given that image tags are mutable (a version tag can be re-pushed with different content, e.g. via a manual `workflow_dispatch` re-run) — ensuring a node never serves stale cached content under a tag string it has seen before.

## Requirements

### Requirement: Always pull on deploy
Both the `gopints-server` and `gopints-web` Deployments SHALL set `imagePullPolicy: Always` on their container, so that every pod start (rollout, restart, reschedule) re-checks the registry for the image content at the given tag rather than trusting a locally cached image of the same tag.

#### Scenario: Same tag re-pushed with different content
- **WHEN** a node already has an image cached locally under `ghcr.io/rickcern44/gopints-web:<tag>` (or `gopints-server:<tag>`), and a new build pushes different content to that same tag
- **THEN** the next pod start on that node pulls the new image content from the registry rather than reusing the stale local cache
