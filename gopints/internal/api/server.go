package api

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rickcern44/gopints/pkg/flow"
	"github.com/rickcern44/gopints/pkg/tap"
)

const sessionTTL = 24 * time.Hour

// Server wires together the HTTP router, WebSocket hub, tap store, and flow meters.
type Server struct {
	store    tap.Store
	meters   map[uint8]*flow.Meter
	hub      *Hub
	simulate bool
	version  string
	http     *http.Server
	sessions sync.Map // token string → expiry time.Time
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
	mux.HandleFunc("POST /api/taps", s.requireAuth(s.handleCreateTap))
	mux.HandleFunc("GET /api/taps/{id}", s.handleGetTap)
	mux.HandleFunc("DELETE /api/taps/{id}", s.requireAuth(s.handleDeleteTap))
	mux.HandleFunc("PUT /api/taps/{id}/keg", s.requireAuth(s.handleSetTapKeg))
	mux.HandleFunc("DELETE /api/taps/{id}/keg", s.requireAuth(s.handleRemoveTapKeg))

	// Kegs
	mux.HandleFunc("GET /api/kegs", s.handleListKegs)
	mux.HandleFunc("POST /api/kegs", s.requireAuth(s.handleCreateKeg))
	mux.HandleFunc("GET /api/kegs/{id}", s.handleGetKeg)
	mux.HandleFunc("PATCH /api/kegs/{id}", s.requireAuth(s.handleUpdateKeg))
	mux.HandleFunc("DELETE /api/kegs/{id}", s.requireAuth(s.handleDeleteKeg))
	mux.HandleFunc("GET /api/kegs/{id}/stats", s.handleGetKegStats)

	// Keg images
	mux.HandleFunc("PUT /api/kegs/{id}/image", s.requireAuth(s.handleSetKegImage))
	mux.HandleFunc("GET /api/kegs/{id}/image", s.handleGetKegImage)
	mux.HandleFunc("DELETE /api/kegs/{id}/image", s.requireAuth(s.handleDeleteKegImage))
	mux.HandleFunc("PUT /api/kegs/{id}/brewery-image", s.requireAuth(s.handleSetBreweryImage))
	mux.HandleFunc("GET /api/kegs/{id}/brewery-image", s.handleGetBreweryImage)
	mux.HandleFunc("DELETE /api/kegs/{id}/brewery-image", s.requireAuth(s.handleDeleteBreweryImage))

	// Pours
	mux.HandleFunc("GET /api/pours", s.handleListPours)
	mux.HandleFunc("DELETE /api/pours/{id}", s.requireAuth(s.handleDeletePour))
	mux.HandleFunc("POST /api/taps/{id}/pour", s.handleManualPour)

	// Features
	mux.HandleFunc("GET /api/features", s.handleListFeatures)
	mux.HandleFunc("PUT /api/features/{name}", s.requireAuth(s.handleSetFeature))

	// Banner
	mux.HandleFunc("GET /api/banner", s.handleGetBanner)
	mux.HandleFunc("PUT /api/banner", s.requireAuth(s.handleSetBanner))
	mux.HandleFunc("DELETE /api/banner", s.requireAuth(s.handleDeleteBanner))

	// Admin auth
	mux.HandleFunc("GET /api/admin/status", s.handleAdminStatus)
	mux.HandleFunc("POST /api/admin/setup", s.handleAdminSetup)
	mux.HandleFunc("POST /api/admin/login", s.handleAdminLogin)
	mux.HandleFunc("POST /api/admin/logout", s.requireAuth(s.handleAdminLogout))

	// WebSocket
	mux.HandleFunc("GET /api/ws", s.hub.ServeWS)

	// Dev only
	if simulate {
		mux.HandleFunc("POST /api/dev/pour/{id}", s.handleSimulatePour)
	}

	s.http = &http.Server{Addr: addr, Handler: requestLogger(mux)}
	return s
}

// statusResponseWriter captures the HTTP status code written by a handler.
// It also forwards http.Hijacker so WebSocket upgrades work through the middleware.
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *statusResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *statusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("underlying ResponseWriter does not implement http.Hijacker")
	}
	return h.Hijack()
}

// requestLogger wraps a handler and emits a structured log line per request.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusResponseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"ms", time.Since(start).Milliseconds(),
		)
	})
}

// requireAuth wraps a handler to require a valid session token.
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !s.validSession(token) {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next(w, r)
	}
}

func (s *Server) validSession(token string) bool {
	if token == "" {
		return false
	}
	v, ok := s.sessions.Load(token)
	if !ok {
		return false
	}
	if time.Now().After(v.(time.Time)) {
		s.sessions.Delete(token)
		return false
	}
	return true
}

// SeedSession inserts a token directly into the session store. For testing only.
func (s *Server) SeedSession(token string) {
	s.sessions.Store(token, time.Now().Add(sessionTTL))
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
