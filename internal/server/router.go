package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/stlesnik/goph_keeper/internal/handlers"
	"github.com/stlesnik/goph_keeper/internal/middleware"
	"github.com/stlesnik/goph_keeper/internal/store"
	"net/http"
)

// NewRouter creates new chi router and configures it with handlers
func NewRouter(store *store.Store) *chi.Mux {
	r := chi.NewRouter()

	hs := handlers.NewHandlers(store)
	wrap := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.WithLogging(h)
	}

	r.Get("/ping", wrap(hs.Ping))

	return r
}
