package middleware

import (
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/util"
	"net/http"
)

// WithAuth is a middleware that checks if the user is authenticated.
func WithAuth(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "request does not contain an access token", http.StatusUnauthorized)
			return
		}
		err := util.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
