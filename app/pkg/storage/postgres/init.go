package postgres

import (
	"aifory-pay-admin-bot/config"
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
)

func InitPsqlDB(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	_, span := otel.Tracer("").Start(ctx, "storage.InitPsqlDB")
	defer span.End()

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
		return nil, err
	}

	database.SetMaxOpenConns(cfg.Postgres.Settings.MaxOpenConns)
	database.SetConnMaxLifetime(cfg.Postgres.Settings.ConnMaxLifetime * time.Second)
	database.SetMaxIdleConns(cfg.Postgres.Settings.MaxIdleConns)
	database.SetConnMaxIdleTime(cfg.Postgres.Settings.ConnMaxIdleTime * time.Second)

	if err = database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}
