package service

import (
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/store"
)

// Service is a structure representing service layer
type Service struct {
	Auth *AuthService
	User *UserService
	Data *DataService
}

// NewService creates new Service
func NewService(cfg *config.Config, store *store.Store) *Service {
	return &Service{
		Auth: NewAuthService(cfg, store.Users),
		User: NewUserService(store.Users),
		Data: NewDataService(store.Data),
	}
}
