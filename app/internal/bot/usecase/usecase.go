package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotUC struct {
	cfg    *config.Config
	pgRepo bot.PgRepository
	BotAPI *tgbotapi.BotAPI
}

func NewBotUC(
	cfg *config.Config,
	pgRepo bot.PgRepository,
	botAPI *tgbotapi.BotAPI,
) bot.Usecase {
	return &BotUC{
		cfg:    cfg,
		pgRepo: pgRepo,
		BotAPI: botAPI,
	}
}
