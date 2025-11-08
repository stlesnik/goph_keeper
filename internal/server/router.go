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
func NewRouter(cfg *config.ServerConfig, store *store.Store) *chi.Mux {
	r := chi.NewRouter()

	hs := handlers.NewHandlers(cfg, store)

	r.Use(middleware.WithLogging)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", hs.RegisterUser)
		r.Post("/login", hs.LoginUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return middleware.WithAuth(cfg, next)
		})

		r.Route("/user", func(r chi.Router) {
			r.Put("/password", hs.ChangePassword)
			r.Get("/profile", hs.GetUserProfile)
		})

		r.Route("/data", func(r chi.Router) {
			r.Get("/{offset}", hs.GetAllData)

			r.Route("/item", func(r chi.Router) {
				r.Post("/", hs.CreateData)
				r.Get("/{id}", hs.GetDataByID)
				r.Put("/{id}", hs.UpdateData)
				r.Delete("/{id}", hs.DeleteData)
			})
		})
	})

	r.Get("/ping", hs.Ping)
	r.Get("/health", hs.Health)
	r.Get("/version", hs.Version)

	return r
}
