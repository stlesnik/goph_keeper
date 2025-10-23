package store

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (users *UsersRepository) GetById(ctx context.Context, id string) (User, error) {
	var user User
	err := users.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return User{}, fmt.Errorf("error getting user by id %s: %w", id, err)
	}
	return user, nil
}
