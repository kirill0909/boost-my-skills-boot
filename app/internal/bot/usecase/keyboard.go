package usecase

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *BotUC) createMainMenuKeyboard(isAdmin bool) (keyboard tgbotapi.ReplyKeyboardMarkup) {

	var keyboardButtons []tgbotapi.KeyboardButton
	if isAdmin { // buttons available only to the administrator
		keyboardButtons = append(
			keyboardButtons,
			tgbotapi.NewKeyboardButton(getUUIDButton))
	}

	keyboardButtons = append(
		keyboardButtons,
		tgbotapi.NewKeyboardButton(askMeButton),
		tgbotapi.NewKeyboardButton(addInfoButton),
		tgbotapi.NewKeyboardButton(printInfoButton),
	)

	middle := len(keyboardButtons) / 2

	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			keyboardButtons[:middle]...,
		),
		tgbotapi.NewKeyboardButtonRow(
			keyboardButtons[middle:]...,
		),
	)

	keyboard.OneTimeKeyboard = false // Hide keyboard after one use
	keyboard.ResizeKeyboard = true   // Resizes keyboard depending on the user's device

	return
}
