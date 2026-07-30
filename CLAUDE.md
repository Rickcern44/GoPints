# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GPIO-based kegerator monitoring system: a Raspberry Pi GPIO agent sends pulse data over UDP to a REST/WebSocket server, which persists pours to SQLite and streams events to a SvelteKit frontend.

## Commands

### Task runner

A [Task](https://taskfile.dev) `Taskfile.yml` wraps the commands below for convenience — run `task --list` from the repo root to see all available tasks (`task test`, `task lint`, `task dev` to run the server + web dev server together, etc.). The raw commands below always work too.

### Backend (Go — from `gopints/`)

```bash
go build -o kegerator-agent ./cmd/agent     # Pi-only GPIO binary
go build -o kegerator-server ./cmd/server   # REST/WS server binary
KEGERATOR_SIMULATE=true go run ./cmd/server # Server on :8080 with /dev endpoints
go run ./cmd/agent simulate --tap 1 --pulses 450 --hz 10 --interval 30s
```

**Tests:**
```bash
go test -race ./...                                            # All tests
go test -race ./pkg/tap -v                                     # Single package
go test -race -run TestSQLiteStore_Migration ./pkg/tap         # Single test
```

### Frontend (from `web/`)

```bash
npm run dev       # Dev server on :5173
npm run build     # Production build
npm run check     # Type check + Svelte validation
npm run lint      # Prettier + ESLint
npm run format    # Auto-format
```

## Architecture

### Data Flow

```
Raspberry Pi GPIO pin
  → pkg/agent (Observer + LineRequester)
  → internal/udp sender (10-byte UDP datagram)
  → internal/udp receiver
  → pkg/flow Meter (pulses → ml volume)
  → PourEvent (Started/Updated/Ended)
  → eventLoop: WebSocket broadcast + SQLite persist (PourEnded only)
  → SvelteKit frontend via /api/ws
```

### Key Packages

- **`pkg/protocol`** — 10-byte UDP wire format: `type(1) + tap_id(1) + timestamp_ns(8)`
- **`pkg/agent`** — GPIO abstraction; `LineRequester` interface has real (`requester_linux.go`) and simulator (`requester_sim.go`) implementations
- **`pkg/flow`** — `Meter` accumulates pulses per tap; goroutine-safe via sync.Mutex; emits `PourEvent` on buffered channel (non-blocking sends to avoid GPIO stall)
- **`pkg/tap`** — `Store` interface + `SQLiteStore` implementation; schema auto-migrates to v3 on first connect; `SetMaxOpenConns(1)` for writer serialization
- **`pkg/config`** — `Loader` chain: FileLoader → EnvLoader → StaticLoader → defaults
- **`internal/api`** — HTTP routing via Go's built-in `http.ServeMux` with pattern matching; WebSocket hub broadcasts all events to all clients

### Important Invariants

- Tap ID is `uint8` (1–255) throughout all layers
- Volume is always in milliliters: `pulses / PulsesPerLiter * 1000` (default 450 pulses/liter)
- SQLite timestamps are milliseconds since epoch; UDP protocol uses nanoseconds
- The agent binary only compiles on Linux (`requester_linux.go` uses `go-gpiocdev`); use `agent simulate` on other platforms
- `/api/dev/pour/{id}` endpoint only registers when the server starts with `KEGERATOR_SIMULATE=true`

### Test Patterns

- Store tests use in-memory SQLite (`:memory:`)
- API handler tests use a mock store with function fields (no real DB)
- All tests use `-race` detector

### Frontend

SvelteKit v2 with Svelte 5 runes mode (enforced in `svelte.config.js`), Tailwind v4, TypeScript. Currently a scaffold — no pages implemented yet.
