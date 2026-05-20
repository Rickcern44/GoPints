package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rickcern44/gopints/pkg/flow"
	"github.com/rickcern44/gopints/pkg/tap"
	"golang.org/x/crypto/bcrypt"
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

// --- Health ---

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": s.version,
	})
}

// --- Taps ---

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

// --- Kegs ---

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

	// Read current state then merge in the patch fields.
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

// --- Keg image ---

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

// --- Pours ---

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
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
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

// --- Banner ---

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

// --- Admin auth ---

const adminPasswordKey = "admin_password"

func (s *Server) handleAdminStatus(w http.ResponseWriter, r *http.Request) {
	_, err := s.store.GetSetting(r.Context(), adminPasswordKey)
	writeJSON(w, http.StatusOK, map[string]bool{"password_set": err == nil})
}

func (s *Server) handleAdminSetup(w http.ResponseWriter, r *http.Request) {
	if _, err := s.store.GetSetting(r.Context(), adminPasswordKey); err == nil {
		writeError(w, http.StatusConflict, "password already set")
		return
	}
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Password == "" {
		writeError(w, http.StatusBadRequest, "password required")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}
	if err := s.store.SetSetting(r.Context(), adminPasswordKey, string(hash)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	token := uuid.New().String()
	s.sessions.Store(token, time.Now().Add(sessionTTL))
	writeJSON(w, http.StatusCreated, map[string]string{"token": token})
}

func (s *Server) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	hash, err := s.store.GetSetting(r.Context(), adminPasswordKey)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	token := uuid.New().String()
	s.sessions.Store(token, time.Now().Add(sessionTTL))
	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (s *Server) handleAdminLogout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	s.sessions.Delete(token)
	w.WriteHeader(http.StatusNoContent)
}

// --- Dev simulate ---

func (s *Server) handleSimulatePour(w http.ResponseWriter, r *http.Request) {
	if !s.simulate {
		writeError(w, http.StatusForbidden, "not in simulate mode")
		return
	}
	tapID, err := parseTapID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid tap id")
		return
	}
	pulses := queryInt(r, "pulses", 200)

	meter, ok := s.meters[tapID]
	if !ok {
		writeError(w, http.StatusNotFound, "tap not found")
		return
	}

	go func() {
		for i := 0; i < pulses; i++ {
			meter.HandlePulse()
			time.Sleep(10 * time.Millisecond)
		}
	}()

	writeJSON(w, http.StatusAccepted, map[string]any{
		"tap_id": tapID,
		"pulses": pulses,
	})
}

// --- Helpers ---

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
