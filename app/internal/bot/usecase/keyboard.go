package usecase

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *botUC) createMainMenuKeyboard(isAdmin bool) (keyboard tgbotapi.ReplyKeyboardMarkup) {

	ctx := context.Background()
	buttons, err := u.pgRepo.GetMainButtons(ctx)
	if err != nil {
		u.log.Errorf(err.Error())
	}

	var keyboardButtons []tgbotapi.KeyboardButton
	for _, button := range buttons {
		if isAdmin { // all buttons are available to the admin
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(button.Name))
		} else if !button.OnlyForAdmin && !isAdmin {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(button.Name))
		}
	}

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

func (t *botUC) createDirectionsKeyboard(directions []models.UserDirection) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(directions); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(directions[i].Direction, fmt.Sprintf("%d", directions[i].ID))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(directions)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}
	}

	return
}
