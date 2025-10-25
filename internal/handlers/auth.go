package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stlesnik/goph_keeper/internal/logger"
	"github.com/stlesnik/goph_keeper/internal/models"
	"net/http"
)

// RegisterUser handles user's initial registration.
func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var regUser models.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&regUser)
	if err != nil {
		logger.Logger.Error("Error decoding registration request: %w", err)
		http.Error(w, "Got error while decoding registration request", http.StatusBadRequest)
	}

	token, err := h.service.Auth.Register(r.Context(), regUser)
	if err != nil {
		logger.Logger.Error("Registration error: %w", err)
		errorMessage := fmt.Sprintf("Unable to register user: %s", err.Error())
		http.Error(w, errorMessage, http.StatusBadRequest)
	}
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusCreated)
}

// LoginUser handles user's log in.
func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	var logUser models.LoginUser
	err := json.NewDecoder(r.Body).Decode(&logUser)
	if err != nil {
		logger.Logger.Error("Error decoding login request: %w", err)
		http.Error(w, "Got error while decoding login request", http.StatusBadRequest)
	}

	token, err := h.service.Auth.Login(r.Context(), logUser)
	if err != nil {
		logger.Logger.Error("Login error: %w", err)
		http.Error(w, "Unable to login", http.StatusBadRequest)
	}
	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

// Ping handles health check requests to verify server responds.
func (h *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		logger.Logger.Error("Error writing response on ping: %w", err)
	}
}
