# Project Torrent

A monorepo for a GPIO monitoring and management system, featuring a Go backend and a SvelteKit frontend.

## Project Structure

- `gopints/`: Go backend project for GPIO interaction.
  - `cmd/agent/`: The agent responsible for monitoring hardware GPIO pins.
  - `cmd/server/`: The central server for coordinating agents and serving data.
  - `pkg/agent/`: Shared logic for GPIO observation and pulse handling.
- `web/`: SvelteKit frontend for visualization and management.
  - `src/routes/`: Application pages and layout.
  - `src/lib/`: Reusable components and utilities.

## Technologies

### Backend (`gopints`)
- **Language:** Go 1.26
- **GPIO Library:** `github.com/warthog618/go-gpiocdev`
- **Architecture:** Standard Go layout with `cmd/` for entry points and `pkg/` for reusable logic.

### Frontend (`web`)
- **Framework:** SvelteKit 5
- **Styling:** Tailwind CSS 4
- **Language:** TypeScript
- **Build Tool:** Vite

## Getting Started

### Backend

To build the backend components:

```bash
cd gopints
# Build agent
go build -o agent ./cmd/agent
# Build server
go build -o server ./cmd/server
```

To run tests:

```bash
go test ./...
```

### Frontend

To set up and run the frontend:

```bash
cd web
npm install
npm run dev
```

To run checks and linting:

```bash
npm run check
npm run lint
```

## Development Conventions

- **Go:** Follow standard Go idioms and project structure. Ensure code is formatted with `go fmt`.
- **Frontend:** Use Svelte 5 runes and patterns. Adhere to TypeScript for type safety. Formatting and linting are managed via Prettier and ESLint.
- **Communication:** (TODO) Define the API/protocol between `gopints` and `web`.
