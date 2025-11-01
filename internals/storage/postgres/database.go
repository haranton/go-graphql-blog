package postgres

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/haranton/go-graphql-blog/internals/config"
	"github.com/jmoiron/sqlx"
)

func GetDBConnect(config *config.Config, logger *slog.Logger) *sqlx.DB {

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
	)

	logger.Info("Connecting to database", slog.String("dsn", dsn))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error("failed to connect database")
		os.Exit(1)
	}

	return db

}
