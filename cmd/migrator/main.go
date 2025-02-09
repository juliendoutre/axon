package main

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/juliendoutre/axon/internal/config"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals
var (
	GoVersion string
	Os        string
	Arch      string
)

func main() {
	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}

	defer func() { _ = logger.Sync() }()

	migrationURL := config.MigrationsURL()

	migrator, err := migrate.New(migrationURL.String(), config.PostgresURL().String())
	if err != nil {
		logger.Panic("Creating migrator", zap.Error(err))
	}
	defer migrator.Close()

	version, isDirty, err := migrator.Version()
	if err != nil && !errors.Is(err, migrate.ErrInvalidVersion) {
		logger.Panic("Getting migrations version", zap.Error(err))
	}

	logger.Info("Current migrations version", zap.Uint("version", version), zap.Bool("is_dirty", isDirty))

	logger.Info("Running migration...")

	if err := migrator.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			logger.Panic("Running migrations", zap.Error(err))
		}
	}

	version, isDirty, err = migrator.Version()
	if err != nil && !errors.Is(err, migrate.ErrInvalidVersion) {
		logger.Panic("Getting migrations version", zap.Error(err))
	}

	logger.Info("New migrations version", zap.Uint("version", version), zap.Bool("is_dirty", isDirty))
}
