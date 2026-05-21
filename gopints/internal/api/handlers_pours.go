package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/rickcern44/gopints/pkg/flow"
	"github.com/rickcern44/gopints/pkg/tap"
)

func (s *Server) handleListPours(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r, "limit", 20)
	offset := queryInt(r, "offset", 0)

	if tapIDStr := r.URL.Query().Get("tap_id"); tapIDStr != "" {
		n, err := strconv.ParseUint(tapIDStr, 10, 8)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid tap_id")
			return
		}
		pours, err := s.store.ListPoursByTap(r.Context(), uint8(n), limit, offset)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, pours)
		return
	}

	pours, err := s.store.ListPours(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, pours)
}

func (s *Server) handleDeletePour(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pour id")
		return
	}
	if err := s.store.DeletePour(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleManualPour(w http.ResponseWriter, r *http.Request) {
	tapID, err := parseTapID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tap id")
		return
	}
	var body struct {
		VolumeMl float64 `json:"volume_ml"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.VolumeMl <= 0 {
		writeError(w, http.StatusBadRequest, "volume_ml must be positive")
		return
	}
	now := time.Now().UTC()
	pour, err := s.store.RecordPour(r.Context(), tap.Pour{
		TapID:     tapID,
		VolumeMl:  body.VolumeMl,
		StartedAt: now,
		EndedAt:   now,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.hub.Broadcast(pourEventToMap(flow.PourEvent{
		Type:      flow.PourEnded,
		TapID:     tapID,
		VolumeMl:  body.VolumeMl,
		StartedAt: now,
		EndedAt:   now,
	}))
	writeJSON(w, http.StatusCreated, pour)
}
