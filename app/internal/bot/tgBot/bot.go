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

			case startCommand:
				if err := t.handleStartCommand(update.Message.Chat.ID); err != nil {
					log.Printf("bot.TgBot.handleStartCommand: %s", err.Error())
					continue
				}
			case getUUIDCommand:
				if err := t.handleGetUUIDButton(update.Message.Chat.ID); err != nil {
					log.Printf("bot.TgBot.handleGetUUIDButton: %s", err.Error())
					continue
				}
			}
		}
	}
	return nil
}

func (t *TgBot) handleGetUUIDButton(chatID int64) (err error) {
	msg := tgbotapi.NewMessage(chatID, "You tap on the get uuid buttons")
	msg.ReplyMarkup = t.createMainMenuKeyboard()
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
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
			tgbotapi.NewKeyboardButton(getUUIDButton),
		),
	)

	keyboard.OneTimeKeyboard = false // Hide keyboard after one use
	keyboard.ResizeKeyboard = true   // Resizes keyboard depending on the user's device

	return
}
