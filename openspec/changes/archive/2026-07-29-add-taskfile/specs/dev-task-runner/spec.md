## ADDED Requirements

### Requirement: Unified task entrypoint
The repository SHALL provide a root `Taskfile.yml` that aggregates module-level Taskfiles for `gopints/` and `web/`, discoverable via `task --list`.

#### Scenario: Listing available tasks from repo root
- **WHEN** a developer runs `task --list` from the repository root
- **THEN** tasks from both the `server`/`agent` (Go) and `web` (frontend) namespaces are listed with descriptions

### Requirement: Manual-testing run tasks
The Taskfile SHALL provide tasks to run the server in a mode usable without physical GPIO hardware, run the agent simulator, and run the web dev server, individually and concurrently.

#### Scenario: Running the server for manual testing
- **WHEN** a developer runs the server task intended for local testing
- **THEN** the server starts with simulate mode enabled (`KEGERATOR_SIMULATE=true`) and exposes its `/dev` endpoints

#### Scenario: Running the web dev server
- **WHEN** a developer runs the web dev task
- **THEN** the SvelteKit dev server starts on its default port, proxying `/api` to the local server

#### Scenario: Running both concurrently
- **WHEN** a developer runs the composite `dev` task
- **THEN** both the server (simulate mode) and the web dev server start concurrently in a single terminal session, with output from both visible

### Requirement: Test and lint tasks mirror CI
The Taskfile SHALL provide tasks that run the same checks as `.github/workflows/ci.yml`, so a developer can reproduce CI results locally with one command.

#### Scenario: Running backend tests locally
- **WHEN** a developer runs the backend test task
- **THEN** `go test -race ./...` runs in the `gopints/` module, matching the CI "Test" step

#### Scenario: Running frontend checks locally
- **WHEN** a developer runs the frontend check task
- **THEN** `npm run check` and the Biome lint check both run in the `web/` module, matching the CI "Type check" and "Lint" steps

### Requirement: Build tasks wrap documented build commands
The Taskfile SHALL provide tasks that wrap the existing documented build commands for the server and agent binaries.

#### Scenario: Building the server binary via Task
- **WHEN** a developer runs the server build task
- **THEN** it produces the same `kegerator-server` binary as running `go build -o kegerator-server ./cmd/server` directly

#### Scenario: Building the agent binary via Task
- **WHEN** a developer runs the agent build task
- **THEN** it produces the same `kegerator-agent` binary as running `go build -o kegerator-agent ./cmd/agent` directly
