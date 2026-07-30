## Why

Running this project today means remembering separate commands across two modules (`gopints/` Go builds/tests, `web/` npm scripts) plus Docker Compose, none of which are unified anywhere. Onboarding and day-to-day manual testing (spin up the server in `--simulate`-equivalent mode, run the web dev server, run the agent simulator, run backend/frontend tests) requires knowing the right flags and working directories by memory or by re-reading `CLAUDE.md`. Adopting [Task](https://taskfile.dev) (`go-task`) with a `Taskfile.yml` gives a single, discoverable entrypoint (`task --list`, `task dev`, `task test`) for exercising the app during manual and automated testing.

## What Changes

- Add a root `Taskfile.yml` that composes two module-level Taskfiles via `includes:`: `gopints/Taskfile.yml` (Go build/run/test tasks) and `web/Taskfile.yml` (npm script wrappers).
- Add tasks for the common manual-testing loop: running the server with simulate mode enabled, running the agent simulator, running the web dev server, and a composite `task dev` that runs the server and web dev server concurrently in one terminal.
- Add tasks that mirror the existing CI checks (`go test -race ./...`, `golangci-lint`, `npm run check`, `npm run lint`) so the same checks CI runs can be run locally with one command (`task test`, `task lint`).
- Add a `task build` (or per-module `task server:build` / `task agent:build` / `task web:build`) wrapping the existing build commands documented in `CLAUDE.md`.
- Document the new `task` entrypoints briefly in `CLAUDE.md`'s Commands section, alongside (not replacing) the existing raw `go`/`npm` commands.
- No changes to application behavior, APIs, or CI itself — this is developer tooling only.

## Capabilities

### New Capabilities
- `dev-task-runner`: The Task-based developer entrypoints (root + module Taskfiles) for building, running, and testing the server, agent, and web app locally.

### Modified Capabilities
(none — no application behavior changes)

## Impact

- New files: `Taskfile.yml` (root), `gopints/Taskfile.yml`, `web/Taskfile.yml`.
- `CLAUDE.md` — Commands section gains a short pointer to the new `task` entrypoints.
- New tooling dependency: `go-task` CLI (developer machines and optionally CI runners), no change to application runtime dependencies.
- No changes to `gopints/` or `web/` application source, Dockerfiles, `docker-compose.yml`, or `.github/workflows/ci.yml`.
