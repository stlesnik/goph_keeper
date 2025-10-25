package store

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// UsersRepositoryInterface is an interface for structure UsersRepository for mocks
type UsersRepositoryInterface interface {
	Save(ctx context.Context, username, email, hashedPassword string) error
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
func (r *UsersRepository) Save(ctx context.Context, username, email, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3)
	`, username, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

// GetByEmail returns user from db if it exists
func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return User{}, fmt.Errorf("error getting user by email %s: %w", email, err)
	}
	return user, nil
}

// Ping returns table status
func (r *UsersRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
