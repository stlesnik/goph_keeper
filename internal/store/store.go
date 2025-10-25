package store

import (
	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/goph_keeper/internal/config"
)

// Store represents structure of storage.
type Store struct {
	Users *UsersRepository
	Data  *DataRepository
}

// NewStore creates new Store instance with config
func NewStore(cfg config.Config) (*Store, error) {
	db, err := sqlx.Open("pgx", cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	users := NewUsersRepository(db)
	items := NewDataRepository(db)
	return &Store{
		Users: users,
		Data:  items,
	}, nil
}
