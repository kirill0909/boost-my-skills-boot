package repository

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type botRedisRepo struct {
	db  *redis.Client
	cfg *config.Config
}

func NewBotRedisRepo(redis *redis.Client, cfg *config.Config) bot.RedisRepository {
	return &botRedisRepo{db: redis, cfg: cfg}
}

func (r *botRedisRepo) SetAwaitingStatus(ctx context.Context, params models.SetAwaitingStatusParams) error {
	key := fmt.Sprintf("%d", params.ChatID)
	if _, err := r.db.Set(ctx, key, params.StatusID, 0).Result(); err != nil {
		return errors.Wrapf(err, "botRedisRepo.SetAwaitingStatus.Set(). params(%+v)", params)
	}

	return nil
}
