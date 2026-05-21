package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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
