package database

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/thanhnamdk2710/auth-service/internal/config"
	"github.com/thanhnamdk2710/auth-service/internal/shared/logger"
)

func RunMigrations(cfg *config.PostgresConfig) error {
	migrationPath := "migrations"
	absPath, _ := filepath.Abs(migrationPath)

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", absPath),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	logger.Info("Database migration completed successfully")
	return nil
}
