package migrator

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/haranton/go-graphql-blog/internals/config"
)

func RunMigrations(cfg *config.Config, logger *slog.Logger) error {

	if cfg.Database.User == "" ||
		cfg.Database.Password == "" ||
		cfg.Database.Port == 0 ||
		cfg.Database.Name == "" ||
		cfg.Database.Host == "" {
		return fmt.Errorf(
			"incomplete DB configuration: user=%q, name=%q, host=%q, port=%d",
			cfg.Database.User, cfg.Database.Name, cfg.Database.Host, cfg.Database.Port,
		)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
	)

	logger.Info("Starting database migrations")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			logger.Warn("failed to close DB connection", slog.String("error", cerr.Error()))
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create driver: %w", err)
	}

	migrationPath := cfg.Migrations.Path

	logger.Info("Using migrations path", slog.String("path", migrationPath))

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("No new migrations to apply")
		} else {
			return fmt.Errorf("run migrations: %w", err)
		}
	}

	logger.Info("Database migrations ran successfully")
	return nil
}

func MustRunMigrations(cfg *config.Config, logger *slog.Logger) {
	err := RunMigrations(cfg, logger)
	if err != nil {
		panic(err)
	}
}
