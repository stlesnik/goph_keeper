package service

import (
	"context"
	"github.com/stlesnik/goph_keeper/internal/middleware"
	"github.com/stlesnik/goph_keeper/internal/models"
	"github.com/stlesnik/goph_keeper/internal/store"
	"github.com/stlesnik/goph_keeper/internal/util"
)

// UserService is a structure representing user preferences service
type UserService struct {
	repo store.UsersRepositoryInterface
}

// NewUserService creates new user preferences service
func NewUserService(repo *store.UsersRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// ChangePassword validates password and updates it
func (svc *UserService) ChangePassword(ctx context.Context, changeReq models.ChangePasswordRequest) error {
	userContext := ctx.Value(middleware.UserContextKey).(*models.UserContext)
	user, err := svc.repo.GetByEmail(ctx, userContext.Email)
	if err != nil {
		return err
	}
	err = util.CheckPassword(changeReq.CurrentPassword, user.PasswordHash)
	if err != nil {
		return err
	}
	err = util.ValidatePassword(changeReq.NewPassword)
	if err != nil {
		return err
	}
	user.PasswordHash, err = util.HashPassword(changeReq.NewPassword)
	if err != nil {
		return err
	}
	err = svc.repo.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}

// GetProfile gather user's profile info
func (svc *UserService) GetProfile(ctx context.Context) (models.UserProfile, error) {
	userContext := ctx.Value(middleware.UserContextKey).(*models.UserContext)
	user, err := svc.repo.GetByEmail(ctx, userContext.Email)
	if err != nil {
		return models.UserProfile{}, err
	}
	return models.UserProfile{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	}, nil
}
