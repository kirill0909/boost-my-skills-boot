package repository

import (
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"

	"github.com/jmoiron/sqlx"
)

type BotPGRepo struct {
	db *sqlx.DB
}

func NewBotPGRepo(db *sqlx.DB) bot.PgRepository {
	return &BotPGRepo{db: db}
}

func (r *BotPGRepo) GetUUID(ctx context.Context, params models.GetUUID) (result string, err error) {
	if err = r.db.GetContext(ctx, &result, queryGetUUID, params.TgName, params.ChatID); err != nil {
		return
	}

	return
}

func (r *BotPGRepo) IsAdmin(ctx context.Context, params models.GetUUID) (result bool, err error) {
	if err = r.db.GetContext(ctx, &result, queryIsAdmin, params.TgName, params.ChatID); err != nil {
		return
	}

	return
}
