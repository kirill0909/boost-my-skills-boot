package repository

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	redis *redis.Client
	cfg   *config.Config
}

func NewBotRedisRepo(redis *redis.Client, cfg *config.Config) bot.RedisRepository {
	return &RedisRepo{redis: redis, cfg: cfg}
}
