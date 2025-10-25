package service

import (
	"context"
	"github.com/stlesnik/goph_keeper/internal/logger"
	"github.com/stlesnik/goph_keeper/internal/store"
	"go.uber.org/zap"
)

// HealthService handles health logic
type HealthService struct {
	DataRepo  store.DataRepositoryInterface
	UsersRepo store.UsersRepositoryInterface
}

// NewHealthService creates a new health service
func NewHealthService(dataRepo store.DataRepositoryInterface, usersRepo store.UsersRepositoryInterface) *HealthService {
	return &HealthService{
		DataRepo:  dataRepo,
		UsersRepo: usersRepo,
	}
}

func (svc *HealthService) Check(ctx context.Context) bool {
	err := svc.DataRepo.Ping(ctx)
	if err != nil {
		logger.Logger.Error("Data Repository health check failed", zap.Error(err))
		return false
	}
	err := svc.UsersRepo.Ping(ctx)
	if err != nil {
		logger.Logger.Error("Users Repository health check failed", zap.Error(err))
		return false
	}
	return true
}
