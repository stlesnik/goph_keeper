package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/store"
	"net/http"
)

// Server represents the HTTP server.
type Server struct {
	cfg    config.Config
	store  *store.Store
	router *chi.Mux
	http   *http.Server
}

// NewServer creates a new Server instance with given storage and config
func NewServer(cfg *config.Config, store *store.Store) (*Server, error) {
	r := NewRouter(cfg, store)
	s := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	server := &Server{
		cfg:    *cfg,
		store:  store,
		router: r,
		http:   s,
	}
	return server, nil
}

// Start method starts the HTTP server
func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

// Stop method stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
