package handlers

import (
	"encoding/json"
	"github.com/stlesnik/goph_keeper/internal/logger"
	"github.com/stlesnik/goph_keeper/internal/models"
	"net/http"
)

// Health handles health check requests to verify server status.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	isHealthy := h.service.Health.Check(r.Context())
	if !isHealthy {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	logger.Logger.Infow("Health check success",
		"ip", r.RemoteAddr,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	if err != nil {
		logger.Logger.Errorw("Error encoding health status", "error", err)
		http.Error(w, "Error encoding health status", http.StatusInternalServerError)
	}
}

// Version handles version information requests.
func (h *Handlers) Version(w http.ResponseWriter, r *http.Request) {
	version := models.VersionResponse{
		Version: "1.1.0",
		Build:   "dev",
		Date:    "26-10-2025",
	}

	logger.Logger.Infow("Version request success",
		"version", version.Version,
		"ip", r.RemoteAddr,
	)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(version)
	if err != nil {
		logger.Logger.Errorw("Error encoding version", "error", err)
		http.Error(w, "Error encoding version", http.StatusInternalServerError)
	}
}
