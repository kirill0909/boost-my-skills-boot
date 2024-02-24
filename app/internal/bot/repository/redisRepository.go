package repository

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	"boost-my-skills-bot/internal/bot/models"
	"boost-my-skills-bot/pkg/utils"
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
	key := fmt.Sprintf("%s_%d", utils.AwaitingStatusPrefix, params.ChatID)
	delay := time.Duration(time.Second * time.Duration(r.cfg.AwaitingDirectionNameDelay))
	if _, err := r.db.Set(ctx, key, params.StatusID, delay).Result(); err != nil {
		return errors.Wrapf(err, "botRedisRepo.SetAwaitingStatus.Set.Result(). params(%+v)", params)
	}

	return nil
}

func (r *botRedisRepo) ResetAwaitingStatus(ctx context.Context, chatID int64) error {
	key := fmt.Sprintf("%s_%d", utils.AwaitingStatusPrefix, chatID)
	if _, err := r.db.Del(ctx, key).Result(); err != nil {
		return errors.Wrapf(err, "botRedisRepo.ResetAwaitingStatus.Del.Result(). chatID: %d", chatID)
	}

	return nil
}

func (r *botRedisRepo) GetAwaitingStatus(ctx context.Context, chatID int64) (string, error) {
	key := fmt.Sprintf("%s_%d", utils.AwaitingStatusPrefix, chatID)
	value, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return "", errors.Wrapf(err, "botRedisRepo.GetAwaitingStatus.Get.Result(). chatID: %d", chatID)
	}

	return value, nil
}

func (r *botRedisRepo) SetParentDirection(ctx context.Context, params models.SetParentDirectionParams) error {
	key := fmt.Sprintf("%s_%d", utils.ParentDirectionPrefix, params.ChatID)
	delay := time.Duration(time.Second * time.Duration(r.cfg.AwaitingParentDirectionDelay))
	if _, err := r.db.Set(ctx, key, params.ParentDirectionID, delay).Result(); err != nil {
		return errors.Wrapf(err, "botRedisRepo.SetParentDirection.Result(). params(%+v)", params)
	}

	return nil
}

func (r *botRedisRepo) GetParentDirection(ctx context.Context, chatID int64) (string, error) {
	key := fmt.Sprintf("%s_%d", utils.ParentDirectionPrefix, chatID)
	parentDirectionID, err := r.db.Get(ctx, key).Result()
	if err != nil {
		return "", errors.Wrapf(err, "botRedisRepo.GetParentDirection.Result(). chatID: %d", chatID)
	}

	return parentDirectionID, nil
}

func (r *botRedisRepo) ResetParentDirection(ctx context.Context, chatID int64) error {
	key := fmt.Sprintf("%s_%d", utils.ParentDirectionPrefix, chatID)
	if _, err := r.db.Del(ctx, key).Result(); err != nil {
		return errors.Wrapf(err, "botRedisRepo.ResetParentDirection.Result(). chatID: %d", chatID)
	}

	return nil
}

func (r *botRedisRepo) SetExpirationTimeForMessage(ctx context.Context, messageID int, chatID int64) error {
	key := fmt.Sprintf("%s_MessageID_%d_ChatID_%d", utils.ExpirationTimeMessagePrefix, messageID, chatID)
	delay := time.Duration(time.Second * time.Duration(r.cfg.AwaitingParentDirectionDelay))
	if _, err := r.db.Set(ctx, key, "", delay).Result(); err != nil {
		return errors.Wrapf(err, "botRedisRepo.SetExpirationTimeForMessage.Result(). messageID: %d", messageID)
	}

	return nil
}
