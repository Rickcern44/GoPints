package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rickcern44/gopints/pkg/tap"
)

func (s *Server) handleListKegs(w http.ResponseWriter, r *http.Request) {
	kegs, err := s.store.ListKegs(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, kegs)
}

func (s *Server) handleGetKeg(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	keg, err := s.store.GetKeg(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, keg)
}

func (s *Server) handleCreateKeg(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BeerName   string  `json:"beer_name"`
		Style      string  `json:"style"`
		ABV        float64 `json:"abv"`
		Brewery    string  `json:"brewery"`
		Notes      string  `json:"notes"`
		CapacityMl float64 `json:"capacity_ml"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	keg, err := s.store.CreateKeg(r.Context(), tap.Keg{
		BeerName:   body.BeerName,
		Style:      body.Style,
		ABV:        body.ABV,
		Brewery:    body.Brewery,
		Notes:      body.Notes,
		CapacityMl: body.CapacityMl,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, keg)
}

func (s *Server) handleUpdateKeg(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}

	existing, err := s.store.GetKeg(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	var patch struct {
		BeerName   *string  `json:"beer_name"`
		Style      *string  `json:"style"`
		ABV        *float64 `json:"abv"`
		Brewery    *string  `json:"brewery"`
		Notes      *string  `json:"notes"`
		CapacityMl *float64 `json:"capacity_ml"`
		ImageStyle *string  `json:"image_style"`
	}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if patch.BeerName != nil {
		existing.BeerName = *patch.BeerName
	}
	if patch.Style != nil {
		existing.Style = *patch.Style
	}
	if patch.ABV != nil {
		existing.ABV = *patch.ABV
	}
	if patch.Brewery != nil {
		existing.Brewery = *patch.Brewery
	}
	if patch.Notes != nil {
		existing.Notes = *patch.Notes
	}
	if patch.CapacityMl != nil {
		existing.CapacityMl = *patch.CapacityMl
	}
	if patch.ImageStyle != nil {
		existing.ImageStyle = *patch.ImageStyle
	}

	updated, err := s.store.UpdateKeg(r.Context(), existing)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (s *Server) handleDeleteKeg(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	if err := s.store.DeleteKeg(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetKegStats(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	stats, err := s.store.GetKegStats(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// handleSetKegImage accepts the raw image as the request body.
// The caller must set Content-Type to the image MIME type (e.g. image/jpeg).
func (s *Server) handleSetKegImage(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	mimeType := r.Header.Get("Content-Type")
	if mimeType == "" {
		writeError(w, http.StatusBadRequest, "Content-Type header required")
		return
	}
	data, err := io.ReadAll(io.LimitReader(r.Body, maxImageBytes+1))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read image")
		return
	}
	if len(data) > maxImageBytes {
		writeError(w, http.StatusRequestEntityTooLarge, "image exceeds 10 MB limit")
		return
	}
	if len(data) == 0 {
		writeError(w, http.StatusBadRequest, "empty image body")
		return
	}
	if err := s.store.SetKegImage(r.Context(), id, data, mimeType); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSetKegImageFromURL(w http.ResponseWriter, r *http.Request) {
	if !s.featureEnabled(r.Context(), "remote_image_urls") {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.URL == "" {
		writeError(w, http.StatusBadRequest, "url required")
		return
	}
	data, mimeType, err := fetchRemoteImage(r.Context(), body.URL)
	if err != nil {
		writeFetchError(w, err)
		return
	}
	if err := s.store.SetKegImage(r.Context(), id, data, mimeType); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetKegImage(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	data, mimeType, err := s.store.GetKegImage(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Cache-Control", "max-age=3600")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) handleDeleteKegImage(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	if err := s.store.DeleteKegImage(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSetBreweryImage(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	mimeType := r.Header.Get("Content-Type")
	if mimeType == "" {
		writeError(w, http.StatusBadRequest, "Content-Type header required")
		return
	}
	data, err := io.ReadAll(io.LimitReader(r.Body, maxImageBytes+1))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read image")
		return
	}
	if len(data) > maxImageBytes {
		writeError(w, http.StatusRequestEntityTooLarge, "image exceeds 10 MB limit")
		return
	}
	if len(data) == 0 {
		writeError(w, http.StatusBadRequest, "empty image body")
		return
	}
	if err := s.store.SetBreweryImage(r.Context(), id, data, mimeType); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleSetBreweryImageFromURL(w http.ResponseWriter, r *http.Request) {
	if !s.featureEnabled(r.Context(), "remote_image_urls") {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	var body struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.URL == "" {
		writeError(w, http.StatusBadRequest, "url required")
		return
	}
	data, mimeType, err := fetchRemoteImage(r.Context(), body.URL)
	if err != nil {
		writeFetchError(w, err)
		return
	}
	if err := s.store.SetBreweryImage(r.Context(), id, data, mimeType); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetBreweryImage(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	data, mimeType, err := s.store.GetBreweryImage(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Cache-Control", "max-age=3600")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) handleDeleteBreweryImage(w http.ResponseWriter, r *http.Request) {
	id, err := parseKegID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid keg id")
		return
	}
	if err := s.store.DeleteBreweryImage(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
