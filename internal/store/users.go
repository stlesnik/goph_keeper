package store

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// UsersRepositoryInterface is an interface for structure UsersRepository for mocks
type UsersRepositoryInterface interface {
	Save(ctx context.Context, item *User) error
	Update(ctx context.Context, item *User) error
	GetByEmail(ctx context.Context, email string) (User, error)
	Ping(ctx context.Context) error
}

// UsersRepository is a structure with database connection
type UsersRepository struct {
	db *sqlx.DB
}

// NewUsersRepository creates new instance of UserRepository
func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// Save creates user in db
func (r *UsersRepository) Save(ctx context.Context, item *User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, username, email, password_hash, salt, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, item.ID, item.Username, item.Email,
		item.PasswordHash, item.Salt, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

// Update updates user in db
func (r *UsersRepository) Update(ctx context.Context, item *User) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users 
		SET username = $2,
		    email = $3,
		    password_hash = $4,
		    salt = $5,
		    created_at = $6,
		    updated_at = $7
		WHERE id = $1
		`, item.ID, item.Username, item.Email,
		item.PasswordHash, item.Salt, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

// GetByEmail returns user from db if it exists
func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := r.db.GetContext(ctx, &user, `
		SELECT id, username, email, password_hash, salt, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`, email)
	if err != nil {
		return User{}, fmt.Errorf("error getting user by email %s: %w", email, err)
	}
	return user, nil
}

// Ping returns table status
func (r *UsersRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *UsersRepository) Close() error {
	return r.db.Close()
}
