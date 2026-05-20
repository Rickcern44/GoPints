# project-torrent

An open-source kegerator tap monitoring and management system. Tracks pours in real time using GPIO flow sensors on a Raspberry Pi, persists data in SQLite, and serves a REST + WebSocket API for a display frontend.

---

## Repository structure

```
project-torrent/
├── gopints/          Go backend (agent + server binaries + shared packages)
└── web/              SvelteKit display frontend
```

---

## What's been built

### `gopints/` — Go backend

#### Shared packages (externally distributable)
| Package | Description |
|---|---|
| `pkg/agent` | GPIO abstraction — `Observer` reads pulse events via `LineRequester` interface; `SimulatorRequester` for hardware-free dev |
| `pkg/flow` | Flow metering — `Meter` converts pulse counts to pour volume, emits `PourStarted` / `PourUpdated` / `PourEnded` events |
| `pkg/tap` | Persistence — `Store` interface + `SQLiteStore` implementation; models for Tap, Keg, Pour, KegStats |
| `pkg/protocol` | UDP wire format — 10-byte fixed message (type + tap ID + timestamp) |
| `pkg/config` | Config loading — `Loader` interface with `FileLoader` (JSON), `EnvLoader` (env vars), `StaticLoader` (tests), `Default()` |

#### Binaries
| Binary | Description |
|---|---|
| `cmd/agent` | Runs on Pi — reads GPIO pulses, sends UDP datagrams to server. `simulate` subcommand fires synthetic pulses over UDP for development |
| `cmd/server` | Runs anywhere — receives UDP, tracks pours, serves REST API + WebSocket, persists to SQLite |

#### Internal packages
| Package | Description |
|---|---|
| `internal/udp` | UDP sender/receiver wrappers |
| `internal/api` | HTTP server — all REST endpoints + WebSocket hub |

#### API endpoints (server)
| Method | Path | Description |
|---|---|---|
| GET | `/api/health` | Version + status |
| GET | `/api/taps` | List all taps (with current keg populated) |
| GET | `/api/taps/{id}` | Single tap |
| PUT | `/api/taps/{id}/keg` | Assign keg to tap |
| DELETE | `/api/taps/{id}/keg` | Remove keg from tap |
| GET | `/api/kegs` | List all kegs |
| POST | `/api/kegs` | Create keg |
| GET | `/api/kegs/{id}` | Single keg |
| PATCH | `/api/kegs/{id}` | Partial update keg |
| DELETE | `/api/kegs/{id}` | Delete keg |
| GET | `/api/kegs/{id}/stats` | Pour count, poured ml, remaining ml, % remaining |
| PUT | `/api/kegs/{id}/image` | Upload keg image (stored as SQLite BLOB) |
| GET | `/api/kegs/{id}/image` | Fetch keg image |
| DELETE | `/api/kegs/{id}/image` | Remove keg image |
| GET | `/api/pours` | List pours (paginated, optional `?tap_id=` filter) |
| DELETE | `/api/pours/{id}` | Delete a pour record |
| GET | `/api/ws` | WebSocket — broadcasts `PourStarted`, `PourUpdated`, `PourEnded` events |
| POST | `/api/dev/pour/{id}` | Simulate a pour (only when `--simulate` flag set) |

#### SQLite schema (v3)
- `kegs` — beer metadata + image BLOB
- `taps` — maps tap ID to current keg
- `pours` — pour records with volume + timestamps (millisecond precision)

#### GoReleaser
- Agent: `linux/arm64` + `linux/armv7` only (Pi targets)
- Server: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`
- `CGO_ENABLED=0` throughout (modernc.org/sqlite is pure Go)

#### Tests
Full test suite, stdlib only, passing with `-race`:

| Package | Tests |
|---|---|
| `pkg/protocol` | 5 — encode/decode round-trip, boundary values, error cases |
| `pkg/flow` | 7 — PourStarted/Updated/Ended, volume math, sequential pours, full-channel safety |
| `pkg/agent` | 5 — SimulatorRequester registration, pulse delivery, close/cleanup |
| `pkg/tap` | 18 — full CRUD, migrations, image blob, stats, pagination |
| `pkg/config` | 10 — Default, FileLoader, EnvLoader, StaticLoader |
| `internal/api` | 27 — all HTTP handler routes, mockStore with func fields |

---

### `web/` — SvelteKit frontend

**Current state:** Clean scaffold only — SvelteKit v2 + Svelte 5 (runes mode) + Tailwind v4 + TypeScript. No pages implemented yet.

---

## Planned next: display frontend

A read-only keg display UI for a tablet or wall-mounted screen.

### UX
- Full-screen per-keg pages, swipeable left/right (one page per active tap)
- CSS scroll snap — no JS swipe library needed
- No admin controls (admin is a separate tool/computer)

### Per-keg page shows
- Beer name, brewery, style, ABV
- Animated vertical fill gauge (green → amber → red as keg empties)
- Remaining volume + percentage
- Total pour count

### Three banner types (fixed overlay at top of screen)
| Banner | Trigger | Color |
|---|---|---|
| Active pour | WebSocket `PourStarted`/`PourUpdated` event | Amber, pulsing |
| Low keg | Any keg below 20% remaining | Orange |
| Custom message | Admin sets via `PUT /api/banner` | Indigo |

### Required backend additions
- SQLite migration v4: `settings` key-value table
- `GetBannerMessage` / `SetBannerMessage` on Store + SQLiteStore
- New routes: `GET /api/banner`, `PUT /api/banner`, `DELETE /api/banner`

### Frontend files to create
| File | Purpose |
|---|---|
| `web/vite.config.ts` | Dev proxy `/api` → `localhost:8080` |
| `web/src/lib/api.ts` | TypeScript types + fetch helpers |
| `web/src/lib/ws.ts` | Svelte 5 `$state` WebSocket store for live pour events |
| `web/src/routes/+page.ts` | Client-side load (SSR off) |
| `web/src/routes/+page.svelte` | Carousel + banner wiring |
| `web/src/lib/components/LevelGauge.svelte` | Animated vertical fill bar |
| `web/src/lib/components/KegCard.svelte` | Full-screen keg card |
| `web/src/lib/components/BannerStack.svelte` | Stacked banner overlay |

### UI library
None — raw Tailwind utilities only. Tailwind v4 + Svelte 5 runes mode breaks most component libraries; the component needs here are simple enough without one.

---

## Running locally

```bash
# Backend
cd gopints
go run ./cmd/server --simulate       # server on :8080

# In another terminal — simulate pours
cd gopints
go run ./cmd/agent simulate --tap 1 --hz 10 --pulses 450 --interval 30s

# Frontend (once implemented)
cd web
npm run dev                          # dev server on :5173
```

## Running tests

```bash
cd gopints
go test -race ./...
```
