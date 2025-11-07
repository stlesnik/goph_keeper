package service

import (
	"context"
	"github.com/stlesnik/goph_keeper/internal/util"
	"testing"

	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/models"
	"github.com/stlesnik/goph_keeper/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUsersRepository is a mock implementation of UsersRepositoryInterface
type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) Save(ctx context.Context, item *store.User) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockUsersRepository) GetByEmail(ctx context.Context, email string) (store.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(store.User), args.Error(1)
}

func (m *MockUsersRepository) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	cfg := &config.ServerConfig{
		JWTSecret: "test-secret-that-is-long-enough-for-validation",
	}

	authService := &AuthService{
		cfg:  cfg,
		repo: mockRepo,
	}

	regUser := models.RegisterUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "ValidPass123",
	}

	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*store.User")).Return(nil)

	token, err := authService.Register(context.Background(), regUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertCalled(t, "Save", mock.Anything, mock.AnythingOfType("*store.User"))
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	cfg := &config.ServerConfig{
		JWTSecret: "test-secret-that-is-long-enough-for-validation",
	}

	authService := &AuthService{
		cfg:  cfg,
		repo: mockRepo,
	}

	loginUser := models.LoginUserRequest{
		Email:    "test@example.com",
		Password: "ValidPass123",
	}

	hashedPassword, _ := util.HashPassword("ValidPass123")
	mockUser := store.User{
		ID:           "test-user-id",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)

	token, err := authService.Login(context.Background(), loginUser)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertCalled(t, "GetByEmail", mock.Anything, "test@example.com")
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(MockUsersRepository)
	cfg := &config.ServerConfig{
		JWTSecret: "test-secret-that-is-long-enough-for-validation",
	}

	authService := &AuthService{
		cfg:  cfg,
		repo: mockRepo,
	}

	loginUser := models.LoginUserRequest{
		Email:    "test@example.com",
		Password: "WrongPassword123",
	}

	hashedPassword, _ := util.HashPassword("ValidPass123")
	mockUser := store.User{
		ID:           "test-user-id",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
	}

	mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(mockUser, nil)

	token, err := authService.Login(context.Background(), loginUser)

	assert.Error(t, err)
	assert.Empty(t, token)
}
