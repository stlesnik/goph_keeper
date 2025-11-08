package store

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
	"github.com/stlesnik/goph_keeper/internal/config"
	"github.com/stlesnik/goph_keeper/internal/migrations"
)

// Store represents structure of storage.
type Store struct {
	Users *UsersRepository
	Data  *DataRepository
}

// NewStore creates new Store instance with config
func NewStore(cfg config.ServerConfig) (*Store, error) {
	db, err := sqlx.Open("pgx", cfg.PostgresDSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	migrations.Run(cfg.PostgresDSN)

	users := NewUsersRepository(db)
	items := NewDataRepository(db)
	return &Store{
		Users: users,
		Data:  items,
	}, nil
}

func (s *Store) Close() error {
	err := s.Data.Close()
	if err != nil {
		return err
	}
	err = s.Users.Close()
	if err != nil {
		return err
	}
	return nil
}
