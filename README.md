# GoPints

An open-source kegerator tap monitoring and management system. A Raspberry Pi GPIO agent sends pulse data over UDP to a Go server, which persists pours to SQLite and streams events to a SvelteKit web interface.

---

## Overview

```
Raspberry Pi (GPIO flow sensor)
  → gopints-agent  (UDP pulses)
  → gopints-server (REST + WebSocket + SQLite)
  → web UI         (live pour display + admin panel)
```

---

## Deployment

### Docker Compose (recommended)

Pull the pre-built images from GHCR and start:

```bash
docker compose pull
docker compose up -d
```

The server and web images version independently, so each has its own tag variable. Override them to pin specific releases:

```bash
GOPINTS_SERVER_TAG=v1.2.3 GOPINTS_WEB_TAG=v0.5.0 docker compose up -d
```

The web UI is available at `http://localhost:8081` by default. To use port 80 in production, set `WEB_PORT=80`. Configure the server via environment variables in `docker-compose.yml`:

| Variable | Default | Description |
|---|---|---|
| `KEGERATOR_DB_PATH` | `/app/data/kegerator.db` | SQLite database path |
| `KEGERATOR_HTTP_ADDR` | `:8080` | HTTP listen address |
| `KEGERATOR_UDP_ADDR` | `:9876` | UDP listen address for agent pulses |

Database data is persisted in the `db_data` Docker volume.

---

### Agent (Raspberry Pi)

The agent binary runs on the Pi and sends GPIO pulse events to the server over UDP.

**Quick install:**

```bash
curl -fsSL https://raw.githubusercontent.com/rickcern44/gopints/main/scripts/setup-agent.sh | sudo bash
```

The installer will:
1. Detect your Pi's architecture (`arm64` / `armv7` / `amd64`)
2. Download the correct binary from the latest GitHub Release
3. Create a config template at `/etc/gopints/config.json`
4. Register and start a `systemd` service (`gopints-agent`)

**Install a specific version:**

```bash
sudo bash setup-agent.sh --version gopints-v1.2.3 --server 192.168.1.10:9876
```

**After install, edit the config:**

```bash
sudo nano /etc/gopints/config.json
sudo systemctl restart gopints-agent
```

Config reference:

```json
{
  "server_addr": "192.168.1.10:9876",
  "taps": [
    { "id": 1, "gpio_pin": 17 }
  ],
  "pulses_per_liter": 450,
  "flow_timeout_ms": 2000
}
```

**Useful commands:**

```bash
sudo systemctl status  gopints-agent      # check status
sudo journalctl -u     gopints-agent -f   # stream logs
sudo systemctl restart gopints-agent      # apply config changes
```

---

## Local development

```bash
# Backend — simulate pours without GPIO hardware
cd gopints
go run ./cmd/server --simulate            # server on :8080, enables /api/dev/* endpoints

# In a second terminal — fire synthetic pour events
go run ./cmd/agent simulate --tap 1 --hz 10 --pulses 450 --interval 30s

# Frontend
cd web
npm run dev                               # dev server on :5173, proxies /api → :8080
```

---

## Repository structure

```
project-torrent/
├── gopints/                  Go backend
│   ├── cmd/agent/            Pi GPIO agent binary
│   ├── cmd/server/           REST + WebSocket server binary
│   ├── internal/api/         HTTP handlers + WebSocket hub
│   ├── internal/udp/         UDP sender/receiver
│   ├── pkg/agent/            GPIO abstraction (real + simulator)
│   ├── pkg/flow/             Pulse → volume metering, PourEvent emission
│   ├── pkg/tap/              Store interface + SQLiteStore
│   ├── pkg/protocol/         10-byte UDP wire format
│   ├── pkg/config/           Config loader chain
│   └── .goreleaser.yaml      Agent binary release config
├── web/                      SvelteKit frontend
│   ├── src/routes/           Pages (public display + admin panel)
│   ├── src/lib/              API client, WebSocket store, components
│   ├── Dockerfile            Multi-stage nginx build
│   └── nginx.conf            Nginx config with /api proxy + WebSocket
├── scripts/
│   └── setup-agent.sh        Pi agent installer
├── docker-compose.yml        Production compose file
├── release-please-config.json
└── .release-please-manifest.json
```

---

## API reference

All endpoints are served by the Go server on `:8080`.

### Taps

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/taps` | — | List all taps (keg populated) |
| `GET` | `/api/taps/{id}` | — | Single tap |
| `POST` | `/api/taps` | ✓ | Create tap |
| `DELETE` | `/api/taps/{id}` | ✓ | Delete tap |
| `PUT` | `/api/taps/{id}/keg` | ✓ | Assign keg to tap |
| `DELETE` | `/api/taps/{id}/keg` | ✓ | Remove keg from tap |
| `POST` | `/api/taps/{id}/pour` | — | Record a manual pour |

### Kegs

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/kegs` | — | List all kegs |
| `POST` | `/api/kegs` | ✓ | Create keg |
| `GET` | `/api/kegs/{id}` | — | Single keg |
| `PATCH` | `/api/kegs/{id}` | ✓ | Partial update |
| `DELETE` | `/api/kegs/{id}` | ✓ | Delete keg |
| `GET` | `/api/kegs/{id}/stats` | — | Pour count, volume poured, % remaining |
| `PUT` | `/api/kegs/{id}/image` | ✓ | Upload beer label image (stored as BLOB) |
| `GET` | `/api/kegs/{id}/image` | — | Fetch beer label image |
| `DELETE` | `/api/kegs/{id}/image` | ✓ | Remove beer label image |
| `PUT` | `/api/kegs/{id}/brewery-image` | ✓ | Upload brewery logo |
| `GET` | `/api/kegs/{id}/brewery-image` | — | Fetch brewery logo |
| `DELETE` | `/api/kegs/{id}/brewery-image` | ✓ | Remove brewery logo |

### Pours

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/pours` | — | List pours (paginated, `?tap_id=` filter) |
| `DELETE` | `/api/pours/{id}` | ✓ | Delete a pour record |

### Banner

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/banner` | — | Get custom banner message |
| `PUT` | `/api/banner` | ✓ | Set custom banner message |
| `DELETE` | `/api/banner` | ✓ | Clear custom banner |

### Features

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/features` | — | List feature flags |
| `PUT` | `/api/features/{name}` | ✓ | Enable or disable a feature |

Available feature flags:

| Flag | Default | Description |
|---|---|---|
| `flow_based_pour` | `false` | When enabled, hides the manual pour button (use once hardware flow meters are installed) |

### Auth

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/admin/status` | Whether an admin password has been set |
| `POST` | `/api/admin/setup` | Set the admin password (first-run only) |
| `POST` | `/api/admin/login` | Login, returns session cookie |
| `POST` | `/api/admin/logout` | Invalidate session |

### WebSocket

| Path | Description |
|---|---|
| `GET /api/ws` | Live pour events: `PourStarted`, `PourUpdated`, `PourEnded` |

### Dev (simulate mode only)

| Method | Path | Description |
|---|---|---|
| `POST` | `/api/dev/pour/{id}` | Trigger a simulated pour on tap `{id}` |

### System

| Method | Path | Description |
|---|---|---|
| `GET` | `/api/health` | Version + status |

---

## Web UI

### Public display (`/`)

Full-screen keg carousel with one card per active tap. Live pour banners appear over the carousel via WebSocket events. Guests can log manual pours directly from the public page (no login required).

### Admin panel (`/admin`)

Password-protected. First visit prompts for an admin password setup.

| Page | Path |
|---|---|
| Kegs | `/admin/kegs` |
| Taps | `/admin/taps` |
| Banner | `/admin/banner` |
| Pour History | `/admin/pours` |
| Features | `/admin/features` |

---

## CI / CD

`gopints/` (backend) and `web/` (frontend) release independently, each with its own version, CHANGELOG, and tag stream.

| Workflow | Trigger | What it does |
|---|---|---|
| `ci-backend.yml` | Pull request → `main` touching `gopints/**` | Go lint + test + build |
| `ci-frontend.yml` | Pull request → `main` touching `web/**` | Frontend lint + type-check + build |
| `release-please.yml` | Push → `main` | Maintains independent Release PRs for `gopints` and `web`; merging one bumps that package's version, updates its `CHANGELOG.md`, and publishes a GitHub Release |
| `release.yml` | GitHub Release published | Builds + pushes only the image for the package that released (`gopints-server` or `gopints-web`) to GHCR; runs GoReleaser to attach agent binaries only for `gopints` releases; opens a PR bumping that package's k8s manifest image tag |

### Versioning

Versioning is driven by [Conventional Commits](https://www.conventionalcommits.org/), scoped per package by which files a commit touches:

| Commit prefix | Version bump |
|---|---|
| `fix:` | patch |
| `feat:` | minor |
| `feat!:` or `BREAKING CHANGE:` | major |

When commits land on `main`, release-please opens or updates a Release PR for whichever package(s) changed. Merging a Release PR:
1. Bumps that package's version in `.release-please-manifest.json`
2. Updates that package's `CHANGELOG.md`
3. Creates a `gopints-vX.Y.Z` or `web-vX.Y.Z` tag and GitHub Release

The release workflow then fires and publishes only that package's artifacts, and opens a follow-up PR (no auto-merge) updating that package's Kubernetes manifest image tag.

### Released artifacts

| Artifact | Triggered by | Registry / location |
|---|---|---|
| `ghcr.io/rickcern44/gopints-server:X.Y.Z` | `gopints-vX.Y.Z` release | GHCR |
| `ghcr.io/rickcern44/gopints-web:X.Y.Z` | `web-vX.Y.Z` release | GHCR |
| `gopints-agent-linux-amd64` | `gopints-vX.Y.Z` release | GitHub Release assets |
| `gopints-agent-linux-arm64` | `gopints-vX.Y.Z` release | GitHub Release assets |
| `gopints-agent-linux-armv7` | `gopints-vX.Y.Z` release | GitHub Release assets |
| `checksums.txt` | `gopints-vX.Y.Z` release | GitHub Release assets |

Both container images are built for `linux/amd64` and `linux/arm64`.

---

## Tests

```bash
cd gopints
go test -race ./...
```

| Package | Tests |
|---|---|
| `pkg/protocol` | Encode/decode round-trip, boundary values, error cases |
| `pkg/flow` | PourStarted/Updated/Ended, volume math, sequential pours, full-channel safety |
| `pkg/agent` | SimulatorRequester registration, pulse delivery, close/cleanup |
| `pkg/tap` | Full CRUD, migrations, image blob, stats, pagination |
| `pkg/config` | Default, FileLoader, EnvLoader, StaticLoader |
| `internal/api` | All HTTP handler routes, feature flags, manual pour, mock store |

---

## Architecture notes

- **Agent binary** compiles only on Linux (`go-gpiocdev` uses Linux GPIO character device). Use `--simulate` on other platforms.
- **Pure-Go SQLite** (`modernc.org/sqlite`) — `CGO_ENABLED=0` throughout; no C toolchain needed in Docker.
- **Tap ID** is `uint8` (1–255) at every layer.
- **Volume** is always stored in milliliters. Default sensor calibration: 450 pulses/liter.
- **SQLite timestamps** are milliseconds since epoch; UDP protocol uses nanoseconds.
- **Session tokens** are in-memory only (24-hour TTL). Restarting the server invalidates all sessions.


