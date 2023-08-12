package tgbot

import (
	"aifory-pay-admin-bot/config"
	"aifory-pay-admin-bot/internal/bot"
	"log"

	"aifory-pay-admin-bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	BotAPI *tgbotapi.BotAPI
	cfg    *config.Config
	tgRepo bot.PgRepository
	tgUc   bot.Usecase
}

func NewTgBot(
	cfg *config.Config,
	repo bot.PgRepository,
	usecase bot.Usecase,
	botAPI *tgbotapi.BotAPI,
) *TgBot {
	return &TgBot{
		cfg:    cfg,
		tgRepo: repo,
		BotAPI: botAPI,
		tgUc:   usecase,
	}
}

func (t *TgBot) Run() error {
	log.Printf("Authorized on account %s", t.BotAPI.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = utils.UpdateConfigTime

	updates := t.BotAPI.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			if update.Message.Command() == "start" {
				continue
			}
		}
	}
	return nil
}
