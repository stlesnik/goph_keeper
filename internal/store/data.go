package store

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// DataRepositoryInterface defines interface for data items repository
type DataRepositoryInterface interface {
	Save(ctx context.Context, item *EncryptedDataItem) error
	GetByID(ctx context.Context, id string, userID string) (*EncryptedDataItem, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*EncryptedDataItem, error)
	Update(ctx context.Context, item *EncryptedDataItem) error
	Delete(ctx context.Context, id string, userID string) error
	Ping(ctx context.Context) error
}

// DataRepository handles data items storage operations
type DataRepository struct {
	db *sqlx.DB
}

// NewDataRepository creates a new data repository
func NewDataRepository(db *sqlx.DB) *DataRepository {
	return &DataRepository{
		db: db,
	}
}

// Save stores a new encrypted data item
func (r *DataRepository) Save(ctx context.Context, item *EncryptedDataItem) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO data (id, user_id, type, title, encrypted_data, iv, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
		item.ID, item.UserID, item.Type, item.Title,
		item.EncryptedData, item.IV, item.Metadata,
		item.CreatedAt, item.UpdatedAt)

	return err
}

// GetByID retrieves a data item by ID for specific user
func (r *DataRepository) GetByID(ctx context.Context, id string, userID string) (*EncryptedDataItem, error) {
	var item EncryptedDataItem
	err := r.db.GetContext(ctx, &item, `
		SELECT id, user_id, type, title, encrypted_data, iv, metadata, created_at, updated_at
		FROM data 
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// GetAllByUserID retrieves all data items for a specific user
func (r *DataRepository) GetAllByUserID(ctx context.Context, userID string) ([]*EncryptedDataItem, error) {
	var items []*EncryptedDataItem
	err := r.db.SelectContext(ctx, &items, `
		SELECT id, user_id, type, title, encrypted_data, iv, metadata, created_at, updated_at
		FROM data 
		WHERE user_id = $1
		ORDER BY created_at DESC
		`, userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Update updates an existing data item
func (r *DataRepository) Update(ctx context.Context, item *EncryptedDataItem) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE data 
		SET type = $2, title = $3, encrypted_data = $4, iv = $5, metadata = $6, updated_at = $7
		WHERE id = $1 AND user_id = $8
		`,
		item.ID, item.Type, item.Title, item.EncryptedData,
		item.IV, item.Metadata, item.UpdatedAt, item.UserID)

	return err
}

// Delete removes a data item
func (r *DataRepository) Delete(ctx context.Context, id string, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM data WHERE id = $1 AND user_id = $2
		`, id, userID)
	return err
}

// Ping returns table status
func (r *DataRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
