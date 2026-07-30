## 1. Go module Taskfile

- [x] 1.1 Create `gopints/Taskfile.yml` with `version: '3'`
- [x] 1.2 Add `build` (server), `build-agent` tasks wrapping `go build -o kegerator-server ./cmd/server` / `go build -o kegerator-agent ./cmd/agent`
- [x] 1.3 Add `simulate` task running the server with `KEGERATOR_SIMULATE=true go run ./cmd/server`
- [x] 1.4 Add `agent-simulate` task wrapping `go run ./cmd/agent simulate --tap 1 --pulses 450 --hz 10 --interval 30s`
- [x] 1.5 Add `test` task wrapping `go test -race ./...`
- [x] 1.6 Add `lint` task wrapping `golangci-lint run` (matching the CI lint step)

## 2. Web module Taskfile

- [x] 2.1 Create `web/Taskfile.yml` with `version: '3'`
- [x] 2.2 Add `dev`, `build`, `check`, `lint` tasks wrapping the corresponding `npm run` scripts

## 3. Root Taskfile

- [x] 3.1 Create root `Taskfile.yml` with `version: '3'` and `includes:` for `gopints: ./gopints` (namespace `server`/`gopints`) and `web: ./web` (namespace `web`)
- [x] 3.2 Add composite `test` task depending on `gopints:test` and `web:check` (and `web:lint`)
- [x] 3.3 Add composite `lint` task depending on `gopints:lint` and `web:lint`
- [x] 3.4 Add composite `dev` task with parallel `deps:` on the server simulate task and the web dev task
- [x] 3.5 Add a header comment linking to https://taskfile.dev/installation for setup

## 4. Documentation

- [x] 4.1 Add a short "Task runner" note to `CLAUDE.md`'s Commands section pointing to `task --list`, without removing the existing raw `go`/`npm` commands
- [x] 4.2 Fix the stale `go run ./cmd/server --simulate` example in `CLAUDE.md` to reflect the actual `KEGERATOR_SIMULATE=true` env var

## 5. Verification

- [x] 5.1 Run `task --list` from repo root and confirm all expected tasks appear
- [x] 5.2 Run `task gopints:test` / `task test` and confirm it matches `go test -race ./...` output run directly
- [x] 5.3 Run `task web:check` and `task web:lint` and confirm they match running the npm scripts directly
- [x] 5.4 Run `task dev`, confirm both the server (simulate mode) and web dev server start and are reachable, then stop cleanly (Ctrl-C)
