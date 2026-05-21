package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rickcern44/gopints/pkg/flow"
)

const maxImageBytes = 10 << 20 // 10 MB

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseTapID(r *http.Request) (uint8, error) {
	n, err := strconv.ParseUint(r.PathValue("id"), 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(n), nil
}

func parseKegID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func queryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}

var pourEventTypeNames = map[flow.EventType]string{
	flow.PourStarted: "PourStarted",
	flow.PourUpdated: "PourUpdated",
	flow.PourEnded:   "PourEnded",
}

func pourEventToMap(e flow.PourEvent) map[string]any {
	m := map[string]any{
		"type":       pourEventTypeNames[e.Type],
		"tap_id":     e.TapID,
		"volume_ml":  e.VolumeMl,
		"started_at": e.StartedAt,
	}
	if !e.EndedAt.IsZero() {
		m["ended_at"] = e.EndedAt
	}
	return m
}

// --- Health ---

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": s.version,
	})
}
