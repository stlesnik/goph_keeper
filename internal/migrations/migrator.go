package migrations

import (
	"errors"
	"github.com/stlesnik/goph_keeper/internal/logger"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationsDir = "./internal/migrations/versions"

// Run applies database migrations using the provided DSN.
func Run(dsn string) {
	m, err := migrate.New("file://"+migrationsDir, dsn)
	if err != nil {
		logger.Logger.Errorw("migrate.New failed", "error", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Logger.Errorw("migrations failed", "error", err)
		os.Exit(1)
	}
	logger.Logger.Infow("migrations applied successfully")
}
