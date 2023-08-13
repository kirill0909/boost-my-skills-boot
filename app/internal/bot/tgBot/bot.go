package tgbot

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	"context"
	"fmt"
	"log"
	"strings"

	models "boost-my-skills-bot/internal/models/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
				if err := t.handleStartCommand(
					update.Message.Chat.ID,
					models.UserActivation{ChatID: update.Message.Chat.ID, TgName: update.Message.Chat.UserName},
					update.Message.Text,
				); err != nil {
					log.Printf("bot.TgBot.handleStartCommand: %s", err.Error())
					t.sendErrorMessage(context.Background(), update.Message.Chat.ID, errUserActivation)
					continue
				}
			case getUUIDCommand:
				if err := t.handleGetUUIDCommand(
					update.Message.Chat.ID,
					models.GetUUID{ChatID: update.Message.Chat.ID, TgName: update.Message.Chat.UserName}); err != nil {
					log.Printf("bot.TgBot.handleGetUUIDCommand: %s", err.Error())
					continue
				}
			}
		}

		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			chatID := update.CallbackQuery.From.ID
			messageID := update.CallbackQuery.Message.MessageID
			switch callbackData {
			case backendCallbackData:
				if err := t.handleBackendCallbackData(chatID, messageID); err != nil {
					log.Printf("bot.TgBot.handleBackendCallbackData: %s", err.Error())
					continue
				}
			case frontednCallbackData:
				if err := t.handleFrontendCallbackData(chatID, messageID); err != nil {
					log.Printf("bot.TgBot.handleFrontendCallbackData: %s", err.Error())
					continue
				}
			}
		}
	}
	return nil
}

func (t *TgBot) handleFrontendCallbackData(chatID int64, messageID int) (err error) {
	ctx := context.Background()

	if err = t.tgUC.SetUpFrontendDirection(ctx, chatID); err != nil {
		return
	}

	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, readyMessage)
	msg.ReplyMarkup = t.createMainMenuKeyboard()
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleBackendCallbackData(chatID int64, messageID int) (err error) {
	ctx := context.Background()

	if err = t.tgUC.SetUpBackendDirection(ctx, chatID); err != nil {
		return
	}

	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, readyMessage)
	msg.ReplyMarkup = t.createMainMenuKeyboard()
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) hideKeyboard(chatID int64, messageID int) (err error) {
	edit := tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		messageID,
		tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
		},
	)

	if _, err = t.BotAPI.Send(edit); err != nil {
		return
	}

	return
}

func (t *TgBot) handleGetUUIDCommand(chatID int64, params models.GetUUID) (err error) {
	ctx := context.Background()

	result, err := t.tgUC.GetUUID(ctx, params)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, result)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleStartCommand(chatID int64, params models.UserActivation, text string) (err error) {
	ctx := context.Background()

	splitedText := strings.Split(text, " ")
	if len(splitedText) != 2 {
		err = fmt.Errorf("Error invite token extracting")
		return
	}

	params.UUID = splitedText[1]
	if err = t.tgUC.UserActivation(ctx, params); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, wellcomeMessage)
	msg.ReplyMarkup = t.createDirectionsKeyboard()
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) createDirectionsKeyboard() (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backenddButton, backendCallbackData),
			tgbotapi.NewInlineKeyboardButtonData(frontendButton, frontednCallbackData),
		),
	)

	return
}

func (t *TgBot) sendErrorMessage(ctx context.Context, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := t.BotAPI.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func (t *TgBot) createMainMenuKeyboard() (keyboard tgbotapi.ReplyKeyboardMarkup) {

	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(getUUIDButton),
			tgbotapi.NewKeyboardButton(askMeButton),
		),
	)

	keyboard.OneTimeKeyboard = false // Hide keyboard after one use
	keyboard.ResizeKeyboard = true   // Resizes keyboard depending on the user's device

	return
}
