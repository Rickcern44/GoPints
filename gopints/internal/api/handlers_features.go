package api

import (
	"encoding/json"
	"net/http"
)

const featurePrefix = "feature."

var knownFeatures = []string{"flow_based_pour"}

func (s *Server) handleListFeatures(w http.ResponseWriter, r *http.Request) {
	features := make(map[string]bool, len(knownFeatures))
	for _, name := range knownFeatures {
		features[name] = false
		val, err := s.store.GetSetting(r.Context(), featurePrefix+name)
		if err == nil && val == "true" {
			features[name] = true
		}
	}
	writeJSON(w, http.StatusOK, features)
}

func (s *Server) handleSetFeature(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	known := false
	for _, f := range knownFeatures {
		if f == name {
			known = true
			break
		}
	}
	if !known {
		writeError(w, http.StatusNotFound, "unknown feature")
		return
	}
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	val := "false"
	if body.Enabled {
		val = "true"
	}
	if err := s.store.SetSetting(r.Context(), featurePrefix+name, val); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
