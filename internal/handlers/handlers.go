package handlers

import (
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/service"
	"github.com/stlesnik/goph_keeper/internal/store"
)

// Handlers is a structure representing handlers layer.
type Handlers struct {
	service *service.Service
	store   *store.Store
}

// NewHandlers creates a new Handler with the provided service.
func NewHandlers(cfg *config.Config, store *store.Store) *Handlers {
	svc := service.NewService(cfg, store)
	return &Handlers{
		service: svc,
		store:   store,
	}
}
