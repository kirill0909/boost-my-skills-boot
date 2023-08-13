package tgbot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
