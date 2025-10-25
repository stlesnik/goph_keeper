package handlers

import (
	"encoding/json"
	"github.com/stlesnik/goph_keeper/internal/logger"
	"github.com/stlesnik/goph_keeper/internal/models"
	"net/http"
)

// ChangePassword handles user password change.
func (h *Handlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var changeReq models.ChangePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&changeReq)
	if err != nil {
		logger.Logger.Error("Error decoding change password request: %w", err)
		http.Error(w, "Got error while decoding change password request", http.StatusBadRequest)
		return
	}

	err = h.service.User.ChangePassword(r.Context(), changeReq)
	if err != nil {
		logger.Logger.Errorw("Change password error: %w", err,
			"error", err.Error(),
			"ip", r.RemoteAddr,
		)
		http.Error(w, "Unable to change password", http.StatusBadRequest)
		return
	}

	logger.Logger.Infow("Change password success",
		"ip", r.RemoteAddr,
	)
	w.WriteHeader(http.StatusOK)
}

// GetUserProfile handles retrieval of user profile information.
func (h *Handlers) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.service.User.GetProfile(r.Context())
	if err != nil {
		logger.Logger.Errorw("Get user profile error: %w", err,
			"error", err.Error(),
			"ip", r.RemoteAddr,
		)
		http.Error(w, "Unable to get user profile", http.StatusBadRequest)
		return
	}

	logger.Logger.Infow("Get user profile success",
		"ip", r.RemoteAddr,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(profile)
	if err != nil {
		logger.Logger.Errorw("JSON encoding error: %w", err)
		http.Error(w, "Unable to get user profile", http.StatusBadRequest)
	}
}
