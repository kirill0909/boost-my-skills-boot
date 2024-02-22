package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/kirill0909/logger"
	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	PgDB   *sqlx.DB
	Redis  *redis.Client
	Logger *logger.Logger
}
