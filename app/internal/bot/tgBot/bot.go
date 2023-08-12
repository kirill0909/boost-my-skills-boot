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
			switch update.Message.Command() {
			case "start":
				if err := t.handleStartCommand(update.Message.Chat.ID); err != nil {
					log.Printf("bot.TgBot.handleStartCommand: %s", err.Error())
				}
			}
		}
	}
	return nil
}

func (t *TgBot) handleStartCommand(chatID int64) (err error) {
	msg := tgbotapi.NewMessage(chatID, wellcomeMessage)
	msg.ReplyMarkup = t.createMainMenuKeyboard()
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) createMainMenuKeyboard() (keyboard tgbotapi.ReplyKeyboardMarkup) {

	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/option1"),
			tgbotapi.NewKeyboardButton("/option1.1"),
			tgbotapi.NewKeyboardButton("/option1.2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/option2"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/option3"),
		),
	)

	keyboard.OneTimeKeyboard = false // Set this to true if you want the keyboard to hide after one use
	keyboard.ResizeKeyboard = true   // Resizes keyboard depending on the user's device

	return
}
