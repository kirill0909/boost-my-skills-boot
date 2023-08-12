package tgbot

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	"log"

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
	updateConfig.Timeout = 30

	updates := t.BotAPI.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			if update.Message.Command() == "start" {
				log.Println("hello from boost my skills bot")
				continue
			}
		}
	}
	return nil
}
