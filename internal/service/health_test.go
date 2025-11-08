package service

import (
	"context"
	"github.com/stlesnik/goph_keeper/internal/logger"
	"github.com/stlesnik/goph_keeper/internal/store"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDataRepository struct {
	mock.Mock
}

func (m *MockDataRepository) Save(ctx context.Context, item *store.EncryptedDataItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockDataRepository) GetByID(ctx context.Context, id string, userID string) (*store.EncryptedDataItem, error) {
	args := m.Called(ctx, id, userID)
	return args.Get(0).(*store.EncryptedDataItem), args.Error(1)
}

func (m *MockDataRepository) GetAllByUserID(ctx context.Context, userID string, off int) ([]*store.EncryptedDataItem, error) {
	args := m.Called(ctx, userID, off)
	return args.Get(0).([]*store.EncryptedDataItem), args.Error(1)
}

func (m *MockDataRepository) Update(ctx context.Context, item *store.EncryptedDataItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockDataRepository) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockDataRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestHealthService_Check(t *testing.T) {
	mockDataRepo := new(MockDataRepository)
	mockUsersRepo := new(MockUsersRepository)

	healthService := &HealthService{
		DataRepo:  mockDataRepo,
		UsersRepo: mockUsersRepo,
	}

	mockDataRepo.On("Ping", mock.Anything).Return(nil)
	mockUsersRepo.On("Ping", mock.Anything).Return(nil)

	isHealthy := healthService.Check(context.Background())

	assert.True(t, isHealthy)
	mockDataRepo.AssertCalled(t, "Ping", mock.Anything)
	mockUsersRepo.AssertCalled(t, "Ping", mock.Anything)
}

func TestHealthService_Check_DataRepoFails(t *testing.T) {
	err := logger.InitLogger("dev")
	assert.NoError(t, err)
	mockDataRepo := new(MockDataRepository)
	mockUsersRepo := new(MockUsersRepository)

	healthService := &HealthService{
		DataRepo:  mockDataRepo,
		UsersRepo: mockUsersRepo,
	}

	mockDataRepo.On("Ping", mock.Anything).Return(assert.AnError)

	isHealthy := healthService.Check(context.Background())

	assert.False(t, isHealthy)
	mockDataRepo.AssertCalled(t, "Ping", mock.Anything)
	mockUsersRepo.AssertNotCalled(t, "Ping")
}
