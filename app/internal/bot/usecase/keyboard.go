package usecase

import (
	"boost-my-skills-bot/app/internal/bot/models"
	"boost-my-skills-bot/app/pkg/utils"
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *botUC) createMainMenuKeyboard(isAdmin bool) (keyboard tgbotapi.ReplyKeyboardMarkup) {

	ctx := context.Background()
	buttons, err := u.pgRepo.GetMainButtons(ctx)
	if err != nil {
		u.log.Error("botUC.createMainMenuKeyboard()", "error", err.Error())
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

func (u *botUC) createDirectionsKeyboard(directions []models.UserDirection, callbackType int) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton
	callbackData := `{"directionID": %d, "callbackType": %d}`

	for i := 0; i < len(directions); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(
			directions[i].Direction,
			fmt.Sprintf(callbackData, directions[i].ID, callbackType))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(directions)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}
	}

	return
}

func (u *botUC) createInfoKeyboard(questionID int, callbackType int) tgbotapi.InlineKeyboardMarkup {

	var keyboard tgbotapi.InlineKeyboardMarkup
	var rows []tgbotapi.InlineKeyboardButton
	callbackData := `{"infoID": %d, "callbackType": %d}`

	button := tgbotapi.NewInlineKeyboardButtonData("get an answer", fmt.Sprintf(callbackData, questionID, callbackType))
	rows = append(rows, button)

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))

	return keyboard
}

func (u *botUC) createServiceKeyboard() tgbotapi.InlineKeyboardMarkup {

	var keyboard tgbotapi.InlineKeyboardMarkup
	actionsWithMainKeyboardCallbackData := `{"callbackType": %d}`
	actionsWithUsersCallbackData := `{"callbackType": %d}`

	keyboard.InlineKeyboard = append(
		keyboard.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow( // first row
			tgbotapi.NewInlineKeyboardButtonData("keyboard", fmt.Sprintf(actionsWithMainKeyboardCallbackData, utils.ActionsWithMainKeyboardCallbackType)),
			tgbotapi.NewInlineKeyboardButtonData("users", fmt.Sprintf(actionsWithUsersCallbackData, utils.ActionsWithUsersCallbackType)),
		),
	)

	return keyboard
}

func (u *botUC) hideKeyboard(chatID int64, messageID int) (err error) {
	edit := tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		messageID,
		tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
		},
	)

	if _, err = u.BotAPI.Send(edit); err != nil {
		u.log.Error("botUC.hideKeyboard.Send()", "error", err.Error())
		return
	}

	return
}
