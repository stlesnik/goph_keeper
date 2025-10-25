package store

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// UsersRepository is a structure with database connection
type UsersRepository struct {
	db *sqlx.DB
}

// NewUsersRepository creates new instance of UserRepository
func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// Save creates user in db
func (users *UsersRepository) Save(ctx context.Context, username, email, hashedPassword string) error {
	_, err := users.db.ExecContext(ctx, `
		INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3)
	`, username, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	return nil
}

// GetByEmail returns user from db if it exists
func (users *UsersRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := users.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return User{}, fmt.Errorf("error getting user by email %s: %w", email, err)
	}
	return user, nil
}
