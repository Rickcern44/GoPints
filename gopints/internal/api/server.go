package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/rickcern44/gopints/pkg/flow"
	"github.com/rickcern44/gopints/pkg/tap"
)

// Server wires together the HTTP router, WebSocket hub, tap store, and flow meters.
type Server struct {
	store    tap.Store
	meters   map[uint8]*flow.Meter
	hub      *Hub
	simulate bool
	version  string
	http     *http.Server
}

// NewServer creates a Server. Call ListenAndServe to start accepting requests.
//
//   - addr is the TCP address to listen on (e.g. ":8080")
//   - store is the tap/keg/pour persistence layer
//   - meters maps tap IDs to their FlowMeters (used by the dev simulate endpoint)
//   - simulate enables the POST /api/dev/pour/{id} endpoint
//   - version is embedded in the health response
func NewServer(addr string, store tap.Store, meters map[uint8]*flow.Meter, simulate bool, version string) *Server {
	s := &Server{
		store:    store,
		meters:   meters,
		hub:      newHub(),
		simulate: simulate,
		version:  version,
	}

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("GET /api/health", s.handleHealth)

	// Taps
	mux.HandleFunc("GET /api/taps", s.handleListTaps)
	mux.HandleFunc("GET /api/taps/{id}", s.handleGetTap)
	mux.HandleFunc("PUT /api/taps/{id}/keg", s.handleSetTapKeg)
	mux.HandleFunc("DELETE /api/taps/{id}/keg", s.handleRemoveTapKeg)

	// Kegs
	mux.HandleFunc("GET /api/kegs", s.handleListKegs)
	mux.HandleFunc("POST /api/kegs", s.handleCreateKeg)
	mux.HandleFunc("GET /api/kegs/{id}", s.handleGetKeg)
	mux.HandleFunc("PATCH /api/kegs/{id}", s.handleUpdateKeg)
	mux.HandleFunc("DELETE /api/kegs/{id}", s.handleDeleteKeg)
	mux.HandleFunc("GET /api/kegs/{id}/stats", s.handleGetKegStats)

	// Keg images
	mux.HandleFunc("PUT /api/kegs/{id}/image", s.handleSetKegImage)
	mux.HandleFunc("GET /api/kegs/{id}/image", s.handleGetKegImage)
	mux.HandleFunc("DELETE /api/kegs/{id}/image", s.handleDeleteKegImage)

	// Pours
	mux.HandleFunc("GET /api/pours", s.handleListPours)
	mux.HandleFunc("DELETE /api/pours/{id}", s.handleDeletePour)

	// WebSocket
	mux.HandleFunc("GET /api/ws", s.hub.ServeWS)

	// Dev only
	if simulate {
		mux.HandleFunc("POST /api/dev/pour/{id}", s.handleSimulatePour)
	}

	s.http = &http.Server{Addr: addr, Handler: mux}
	return s
}

// Handler returns the HTTP handler, for use with httptest.NewServer in tests.
func (s *Server) Handler() http.Handler {
	return s.http.Handler
}

// Broadcast sends a PourEvent to all connected WebSocket clients.
func (s *Server) Broadcast(e flow.PourEvent) {
	s.hub.Broadcast(pourEventToMap(e))
}

// ListenAndServe starts the HTTP server. It blocks until the server is shut down.
func (s *Server) ListenAndServe() error {
	slog.Info("HTTP server listening", "addr", s.http.Addr)
	return s.http.ListenAndServe()
}

// Shutdown gracefully stops the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
