package migrations

import (
	"database/sql"
	"BlockchainCurrency/pkg/logger"

	"github.com/pressly/goose/v3"
)

const (
	migrationsDir = "./migrations"
	dialect       = "postgres"
)

func Run(db *sql.DB, logger logger.Logger) error {
	logger.Info("Running migrations...")

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return err
	}

	logger.Info("Migrations completed successfully")
	return nil
}
