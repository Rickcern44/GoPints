package api

import (
	"encoding/json"
	"net/http"
)

const bannerKey = "banner"

func (s *Server) handleGetBanner(w http.ResponseWriter, r *http.Request) {
	msg, err := s.store.GetSetting(r.Context(), bannerKey)
	if err != nil {
		writeError(w, http.StatusNotFound, "no banner set")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": msg})
}

func (s *Server) handleSetBanner(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if body.Message == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}
	if err := s.store.SetSetting(r.Context(), bannerKey, body.Message); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDeleteBanner(w http.ResponseWriter, r *http.Request) {
	if err := s.store.DeleteSetting(r.Context(), bannerKey); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
