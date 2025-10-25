package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/handlers"
	"github.com/stlesnik/goph_keeper/internal/middleware"
	"github.com/stlesnik/goph_keeper/internal/store"
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
		r.Put("/password", authWrap(hs.ChangePassword))
		r.Get("/profile", authWrap(hs.GetUserProfile))
	})

	r.Route("/api/data", func(r chi.Router) {
		r.Post("/", authWrap(hs.CreateData))
		r.Get("/", authWrap(hs.GetAllData))

		r.Get("/{id}", authWrap(hs.GetDataByID))
		r.Put("/{id}", authWrap(hs.UpdateData))
		r.Delete("/{id}", authWrap(hs.DeleteData))
	})

	r.Get("/ping", wrap(hs.Ping))
	r.Get("/health", wrap(hs.Health))
	r.Get("/version", wrap(hs.Version))

	return r
}
