package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *Server) handleListTaps(w http.ResponseWriter, r *http.Request) {
	taps, err := s.store.ListTaps(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, taps)
}

func (s *Server) handleGetTap(w http.ResponseWriter, r *http.Request) {
	id, err := parseTapID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tap id")
		return
	}
	t, err := s.store.GetTap(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (s *Server) handleCreateTap(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID uint8 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.ID == 0 {
		writeError(w, http.StatusBadRequest, "id must be 1–255")
		return
	}
	if err := s.store.CreateTap(r.Context(), body.ID); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			writeError(w, http.StatusConflict, "tap already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	t, err := s.store.GetTap(r.Context(), body.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (s *Server) handleDeleteTap(w http.ResponseWriter, r *http.Request) {
	id, err := parseTapID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tap id")
		return
	}
	if err := s.store.DeleteTap(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSetTapKeg(w http.ResponseWriter, r *http.Request) {
	tapID, err := parseTapID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tap id")
		return
	}
	var body struct {
		KegID int64 `json:"keg_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := s.store.SetTapKeg(r.Context(), tapID, body.KegID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleRemoveTapKeg(w http.ResponseWriter, r *http.Request) {
	tapID, err := parseTapID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tap id")
		return
	}
	if err := s.store.RemoveTapKeg(r.Context(), tapID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
