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
				if err := t.sendMessage(update.Message.Chat.ID, "Hello from boost bot"); err != nil {
					log.Println(err)
					continue
				}
				continue
			}
		}
	}
	return nil
}

func (t *TgBot) sendMessage(chatID int64, text ...string) (err error) {
	var msg tgbotapi.MessageConfig
	if len(text) > 0 {
		msg = tgbotapi.NewMessage(chatID, text[0])
	} else {
		msg = tgbotapi.NewMessage(chatID, "")
	}

	if _, err = t.BotAPI.Send(msg); err != nil {
		return err
	}

	return
}
