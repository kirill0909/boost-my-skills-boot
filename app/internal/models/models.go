package models

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Dependencies struct {
	PgDB        *sqlx.DB
	Redis       *redis.Client
	RedisPubSub *redis.PubSub
}
