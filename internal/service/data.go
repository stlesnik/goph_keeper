package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stlesnik/goph_keeper/internal/middleware"
	"time"

	"github.com/stlesnik/goph_keeper/internal/models"
	"github.com/stlesnik/goph_keeper/internal/store"
	"github.com/stlesnik/goph_keeper/internal/util"
)

// DataService handles data items business logic
type DataService struct {
	repo    store.DataRepositoryInterface
	encrypt *util.EncryptionService
}

// NewDataService creates a new data service
func NewDataService(repo store.DataRepositoryInterface) *DataService {
	return &DataService{
		repo:    repo,
		encrypt: util.NewEncryptionService(),
	}
}

// Create creates a new data item
func (svc *DataService) Create(ctx context.Context, req models.CreateDataRequest) error {
	userClaims := ctx.Value(middleware.UserContextKey).(*util.Claims)
	userID := userClaims.UserID

	userKey := svc.encrypt.DeriveUserKey(userClaims.Email, []byte(userClaims.Email))
	encryptedData, iv, err := svc.encrypt.EncryptData(req.Data, userKey)
	if err != nil {
		return err
	}

	now := time.Now()
	item := &store.EncryptedDataItem{
		ID:            uuid.New().String(),
		UserID:        userID,
		Type:          req.Type,
		Title:         req.Title,
		EncryptedData: encryptedData,
		IV:            iv,
		Metadata:      req.Metadata,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = svc.repo.Save(ctx, item)
	if err != nil {
		return err
	}

	return nil
}

// GetAll retrieves all data item's titles for the user
func (svc *DataService) GetAll(ctx context.Context) ([]models.DataItemResponse, error) {
	userClaims := ctx.Value(middleware.UserContextKey).(*util.Claims)
	userID := userClaims.UserID

	items, err := svc.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []models.DataItemResponse
	for _, item := range items {
		responses = append(responses, models.DataItemResponse{
			ID:        item.ID,
			Type:      item.Type,
			Title:     item.Title,
			Data:      "",
			Metadata:  item.Metadata,
			CreatedAt: item.CreatedAt.Format(time.RFC3339),
			UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
		})
	}
	return responses, nil
}

// GetByID retrieves a data item by ID
func (svc *DataService) GetByID(ctx context.Context, id string) (models.DataItemResponse, error) {
	userClaims := ctx.Value(middleware.UserContextKey).(*util.Claims)
	userID := userClaims.UserID

	item, err := svc.repo.GetByID(ctx, id, userID)
	if err != nil {
		return models.DataItemResponse{}, err
	}

	userKey := svc.encrypt.DeriveUserKey(userClaims.Email, []byte(userClaims.Email))
	decryptedData, err := svc.encrypt.DecryptData(item.EncryptedData, item.IV, userKey)
	if err != nil {
		return models.DataItemResponse{}, err
	}

	return models.DataItemResponse{
		ID:        item.ID,
		Type:      item.Type,
		Title:     item.Title,
		Data:      decryptedData,
		Metadata:  item.Metadata,
		CreatedAt: item.CreatedAt.Format(time.RFC3339),
		UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// Update updates an existing data item
func (svc *DataService) Update(ctx context.Context, id string, req models.UpdateDataRequest) error {
	userClaims := ctx.Value(middleware.UserContextKey).(*util.Claims)
	userID := userClaims.UserID

	existingItem, err := svc.repo.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	userKey := svc.encrypt.DeriveUserKey(userClaims.Email, []byte(userClaims.Email))
	encryptedData, iv, err := svc.encrypt.EncryptData(req.Data, userKey)
	if err != nil {
		return err
	}

	existingItem.Type = req.Type
	existingItem.Title = req.Title
	existingItem.EncryptedData = encryptedData
	existingItem.IV = iv
	existingItem.Metadata = req.Metadata
	existingItem.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, existingItem)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a data item
func (svc *DataService) Delete(ctx context.Context, id string) error {
	userClaims := ctx.Value(middleware.UserContextKey).(*util.Claims)
	userID := userClaims.UserID

	return svc.repo.Delete(ctx, id, userID)
}
