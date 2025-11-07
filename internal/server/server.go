package server

import (
	"context"
	"crypto/tls"
	"github.com/go-chi/chi/v5"
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/logger"
	"github.com/stlesnik/goph_keeper/internal/store"
	"net/http"
	"time"
)

// Server represents the HTTP server.
type Server struct {
	cfg    config.ServerConfig
	store  *store.Store
	router *chi.Mux
	http   *http.Server
}

// NewServer creates a new Server instance with given storage and config
func NewServer(cfg *config.ServerConfig, store *store.Store) (*Server, error) {
	r := NewRouter(cfg, store)

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	s := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      r,
		TLSConfig:    tlsConfig,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
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
	if s.cfg.EnableHTTPS {
		logger.Logger.Infof("Starting HTTPS server on %s", s.cfg.ServerAddress)
		return s.http.ListenAndServeTLS(s.cfg.TLSCertFile, s.cfg.TLSKeyFile)
	}

	logger.Logger.Infof("Starting HTTP server on %s", s.cfg.ServerAddress)
	return s.http.ListenAndServe()
}

// Stop method stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
