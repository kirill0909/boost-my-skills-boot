package tgbot

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TgBot struct {
	BotAPI *tgbotapi.BotAPI
	cfg    *config.Config
	tgUC   bot.Usecase
}

func NewTgBot(
	cfg *config.Config,
	usecase bot.Usecase,
	botAPI *tgbotapi.BotAPI,
) *TgBot {
	return &TgBot{
		cfg:    cfg,
		BotAPI: botAPI,
		tgUC:   usecase,
	}
}

func (t *TgBot) Run() error {
	log.Printf("Authorized on account %s", t.BotAPI.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := t.BotAPI.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case startCommand:
				log.Println("start command")
			}
		}
	}

	return nil
}

func (t *TgBot) sendErrorMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
