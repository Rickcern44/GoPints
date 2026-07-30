## Context

The repo has two independently-tooled modules: `gopints/` (Go, built/tested via `go build`/`go test -race ./...`, linted via `golangci-lint`) and `web/` (Node/SvelteKit, via `npm run dev|build|check|lint`). `.github/workflows/ci.yml` already runs the canonical commands for both as separate jobs. There's also a root `docker-compose.yml` for containerized runs and a `scripts/setup-agent.sh` helper. None of this is unified — a developer has to know the right working directory and command for whichever module and mode (build, test, run/simulate) they need.

## Goals / Non-Goals

**Goals:**
- One discoverable entrypoint (`task --list`) for the commands already documented in `CLAUDE.md` and used in CI.
- A single command to bring up both the server (in a mode usable for manual testing without real GPIO hardware) and the web dev server together, for interactive testing.
- Task definitions that are thin wrappers around existing commands — no new build logic, no duplication of what CI already does differently.

**Non-Goals:**
- Not replacing CI's own step definitions in `ci.yml` — Task is for local developer use; CI continues to run its own explicit steps (kept in sync manually, not by having CI shell out to `task`, to avoid coupling CI behavior to a tool not all contributors may have installed in exactly the same version).
- Not wrapping `docker-compose` deployment flows — out of scope for this "run for testing" change; can be added later if wanted.
- Not changing any application code, Dockerfiles, or CI workflow files.

## Decisions

**Root Taskfile with per-module `includes:`, not one flat file.**
`Taskfile.yml` at the repo root uses Task's `includes:` to pull in `gopints/Taskfile.yml` (namespaced `server:`, `agent:`) and `web/Taskfile.yml` (namespaced `web:`). This mirrors the existing CI job split (backend/frontend) and keeps each module's tasks colocated with the code they operate on — a Go contributor can `cd gopints && task test` without needing the root file at all. Alternative considered: one large root Taskfile with all tasks flat — rejected as it would grow unwieldy and duplicates the module boundary the repo already has.

**Server "simulate" mode via the `KEGERATOR_SIMULATE=true` env var, not a `--simulate` CLI flag.**
`CLAUDE.md` documents `go run ./cmd/server --simulate`, but the actual flag registered in `cmd/server/main.go` is only `--config`; simulate mode is controlled by `cfg.Server.Simulate`, set via the `KEGERATOR_SIMULATE` env var (confirmed in `pkg/config/loader.go`). The `server:simulate` task sets this env var rather than passing a nonexistent flag. (This also means `CLAUDE.md`'s existing command example is stale — worth a one-line fix while touching this area, called out as a task.)

**`task dev` runs server + web concurrently via Task's parallel `deps:`.**
Task runs a task's `deps:` list in parallel by default. A composite `dev` task with `deps: [server:simulate, web:dev]` starts both long-running dev processes in one terminal with interleaved output — matching the "easily run these applications for testing" ask directly, rather than requiring two manual terminal tabs.

**Task installation is documented, not vendored.**
No `go install` wrapper or version-pinning script is added — `go-task` is installed by the developer per the [official install docs](https://taskfile.dev/installation) (Homebrew, `go install`, etc.), same as any other dev-machine tool (`go`, `node`) already assumed present. A `# Requires: https://taskfile.dev/installation` comment at the top of the root Taskfile points contributors there.

## Risks / Trade-offs

- **[Risk]** Task commands could drift from what CI actually runs (e.g. someone updates `ci.yml`'s lint flags but forgets the Taskfile). → **Mitigation**: keep task bodies as literal one-line wrappers around the same commands used in `ci.yml` (no extra flags/logic), so future maintainers can diff them side by side easily. Accepted as a documentation-discipline risk rather than something to automate away in this change.
- **[Risk]** Contributors without `go-task` installed hit an unfamiliar `command not found`. → **Mitigation**: existing raw `go`/`npm` commands in `CLAUDE.md` are kept, not removed — Task is additive, not a required path.
- **[Trade-off]** `task dev`'s parallel `deps:` interleaves both processes' stdout in one terminal, which is harder to read than two separate terminals for heavy debugging. Accepted as the right default for quick manual testing; nothing prevents running `task server:simulate` and `task web:dev` separately in two terminals when that's preferred.

## Migration Plan

Purely additive — no existing files are removed or behaviorally changed (aside from the one-line `CLAUDE.md` fix noted above). Land as a single PR: add the three Taskfiles, verify each task runs, update `CLAUDE.md`. No rollback concerns beyond deleting the new files if unwanted.

## Open Questions

None outstanding.
