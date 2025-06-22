package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sg3t41/api/pkg/config"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Migrator struct {
	db     *sql.DB
	config *config.Config
	logger *zap.Logger
}

func NewMigrator(db *sql.DB, cfg *config.Config, logger *zap.Logger) *Migrator {
	return &Migrator{
		db:     db,
		config: cfg,
		logger: logger,
	}
}

func (m *Migrator) Migrate() error {
	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		m.logger.Error("Failed to create source driver", zap.Error(err))
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		m.config.DBUser, m.config.DBPassword, m.config.DBHost, m.config.DBPort, m.config.DBName)
	migrator, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		m.logger.Error("Failed to create migrator", zap.Error(err))
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	defer func() {
		if sourceErr, dbErr := migrator.Close(); sourceErr != nil || dbErr != nil {
			m.logger.Error("Failed to close migrator", zap.Error(sourceErr), zap.Error(dbErr))
		}
	}()

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			m.logger.Info("No new migrations to apply")
			return nil
		}
		m.logger.Error("Failed to run migrations", zap.Error(err))
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	m.logger.Info("Migrations applied successfully")
	return nil
}

func (m *Migrator) Version() (uint, bool, error) {
	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return 0, false, fmt.Errorf("failed to create source driver: %w", err)
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		m.config.DBUser, m.config.DBPassword, m.config.DBHost, m.config.DBPort, m.config.DBName)
	migrator, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}

	defer func() {
		if sourceErr, dbErr := migrator.Close(); sourceErr != nil || dbErr != nil {
			m.logger.Error("Failed to close migrator", zap.Error(sourceErr), zap.Error(dbErr))
		}
	}()

	version, dirty, err := migrator.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}