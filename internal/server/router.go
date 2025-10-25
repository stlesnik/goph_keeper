package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/handlers"
	"github.com/stlesnik/goph_keeper/internal/middleware"
	"github.com/stlesnik/goph_keeper/internal/store"
	"net/http"
)

// NewRouter creates new chi router and configures it with handlers
func NewRouter(cfg *config.Config, store *store.Store) *chi.Mux {
	r := chi.NewRouter()

	hs := handlers.NewHandlers(cfg, store)
	wrap := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.WithLogging(h)
	}
	authWrap := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.WithAuth(cfg, wrap(h))
	}

	r.Route("/user", func(r chi.Router) {
		r.Post("/register", wrap(hs.RegisterUser))
		r.Post("/login", wrap(hs.LoginUser))
	})

	r.Get("/ping", wrap(hs.Ping))

	return r
}
