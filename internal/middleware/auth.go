package middleware

import (
	"context"
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/models"
	"github.com/stlesnik/goph_keeper/internal/util"
	"net/http"
)

type contextKey string

const UserContextKey = contextKey("user")

// WithAuth is a middleware that checks if the user is authenticated.
func WithAuth(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "request does not contain an access token", http.StatusUnauthorized)
			return
		}
		claims, err := util.ParseToken(tokenString, cfg.JWTSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, &models.UserContext{
			UserID: claims.UserID,
			Email:  claims.Email,
		})
		next(w, r.WithContext(ctx))
	}
}
