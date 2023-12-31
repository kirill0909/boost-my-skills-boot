package repository

import (
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type BotPGRepo struct {
	db *sqlx.DB
}

func NewBotPGRepo(db *sqlx.DB) bot.PgRepository {
	return &BotPGRepo{db: db}
}

func (r *BotPGRepo) CompareUUID(ctx context.Context, params models.CompareUUIDParams) (bool, error) {
	var result bool
	if err := r.db.GetContext(ctx, &result, queryCompareUUID, params.ChatID, params.UUID); err != nil {
		err = errors.Wrapf(err, "BotPGRepo.CompareUUID.queryCompareUUID. params(%+v)", params)
		return false, err
	}

	return result, nil
}

func (r *BotPGRepo) SetStatusActive(ctx context.Context, chatID int64) error {
	_, err := r.db.ExecContext(ctx, querySetStatusActive, chatID)
	if err != nil {
		err = errors.Wrapf(err, "BotPGRepo.SetStatusActive.querySetStatusActive. chatID: %d", chatID)
		return err
	}

	return nil
}
