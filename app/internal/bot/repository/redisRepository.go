package repository

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	"boost-my-skills-bot/internal/bot/models"
	"context"
	"fmt"
	"time"

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
	delay := time.Duration(time.Second * time.Duration(r.cfg.AwaitingDirectionNameDelay))
	if _, err := r.db.Set(ctx, key, params.StatusID, delay).Result(); err != nil {
		return errors.Wrapf(err, "botRedisRepo.SetAwaitingStatus.Set.Result(). params(%+v)", params)
	}

	return nil
}

func (r *botRedisRepo) ResetAwaitingStatus(ctx context.Context, chatID int64) error {
	key := fmt.Sprintf("%d", chatID)
	if _, err := r.db.Del(ctx, key).Result(); err != nil {
		return errors.Wrapf(err, "botRedisRepo.ResetAwaitingStatus.Del.Result(). chatID: %d", chatID)
	}

	return nil
}

func (r *botRedisRepo) GetAwaitingStatus(ctx context.Context, chatID int64) (string, error) {
	key := fmt.Sprintf("%d", chatID)
	value, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return "", errors.Wrapf(err, "botRedisRepo.GetAwaitingStatus.Get.Result(). chatID: %d", chatID)
	}

	return value, nil
}
