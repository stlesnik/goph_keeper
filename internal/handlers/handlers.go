package handlers

import (
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/service"
	"github.com/stlesnik/goph_keeper/internal/store"
	"net/http"
)

// Handlers is a structure representing handlers layer.
type Handlers struct {
	service service.Service
	store   *store.Store
}

// NewHandlers creates a new Handler with the provided service.
func NewHandlers(store *store.Store) *Handlers {
	svc := service.NewService(store)

	return &Handlers{
		service: svc,
		store:   store,
	}

}

// Ping handles health check requests to verify server responds.
func (h *Handlers) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
