package repository

import (
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type BotPGRepo struct {
	db *sqlx.DB
}

func NewBotPGRepo(db *sqlx.DB) bot.PgRepository {
	return &BotPGRepo{db: db}
}

func (r *BotPGRepo) SetStatusActive(ctx context.Context, params models.SetStatusActiveParams) error {
	res, err := r.db.ExecContext(ctx, querySetStatusActive, params.ChatID, params.UUID, params.TgName)
	if err != nil {
		return errors.Wrapf(err, "BotPGRepo.SetStatusActive.querySetStatusActive. params(%+v)", params)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "BotPGRepo.SetStatusActive.RowsAffected. params(%+v)", params)

	}

	if rowsAffected != 1 {
		return fmt.Errorf("BotPGRepo.SetStatusActive rowsAffected != 1. params(%+v)", params)
	}

	return nil
}
