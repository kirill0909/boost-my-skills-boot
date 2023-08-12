package repository

import (
	"aifory-pay-admin-bot/internal/bot"

	"github.com/jmoiron/sqlx"
)

// it need later
type BotRepo struct {
	db *sqlx.DB
}

func NewBotRepository(db *sqlx.DB) bot.PgRepository {
	return &BotRepo{db: db}
}
