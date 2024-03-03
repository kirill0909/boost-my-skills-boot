package postgres

import (
	"boost-my-skills-bot/app/config"
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func InitPgDB(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {

	connectionURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	database, err := sqlx.Open("pgx", connectionURL)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open DB connection")
	}

	if err = database.Ping(); err != nil {
		return nil, errors.Wrap(err, "unable to ping db")
	}

	return database, nil
}
