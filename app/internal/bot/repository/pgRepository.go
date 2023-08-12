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

func (r *BotPGRepo) GetUUID(ctx context.Context) (result string, err error) {
	if err = r.db.GetContext(ctx, &result, queryGetUUID); err != nil {
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

func (r *BotPGRepo) UserActivation(ctx context.Context, params models.UserActivation) (err error) {
	if _, err = r.db.ExecContext(ctx, queryUserActivation, params.TgName, params.ChatID, params.UUID); err != nil {
		return
	}

	return
}
