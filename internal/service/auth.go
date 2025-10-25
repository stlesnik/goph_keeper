package service

import (
	"context"
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/models"
	"github.com/stlesnik/goph_keeper/internal/store"
	"github.com/stlesnik/goph_keeper/internal/util"
)

// AuthService is a structure representing authentication service
type AuthService struct {
	cfg  *config.Config
	repo store.UsersRepositoryInterface
}

// NewAuthService creates new authentication service
func NewAuthService(cfg *config.Config, repo *store.UsersRepository) *AuthService {
	return &AuthService{
		cfg:  cfg,
		repo: repo,
	}
}

// Register validates registration request and creates db entity and authorization token
func (svc *AuthService) Register(ctx context.Context, regUser models.RegisterUserRequest) (string, error) {
	err := util.ValidatePassword(regUser.Password)
	if err != nil {
		return "", err
	}
	hashedPassword, err := util.HashPassword(regUser.Password)
	if err != nil {
		return "", err
	}
	err = svc.repo.Save(ctx, regUser.Username, regUser.Email, hashedPassword)
	if err != nil {
		return "", err
	}
	token, err := util.GenerateJWT(regUser.Username, regUser.Email, svc.cfg.JWTSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Login confirms that user is registered and generates authorization token
func (svc *AuthService) Login(ctx context.Context, loginUser models.LoginUserRequest) (string, error) {
	user, err := svc.repo.GetByEmail(ctx, loginUser.Email)
	if err != nil {
		return "", err
	}
	err = util.CheckPassword(loginUser.Password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	token, err := util.GenerateJWT(user.Username, user.Email, svc.cfg.JWTSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}
