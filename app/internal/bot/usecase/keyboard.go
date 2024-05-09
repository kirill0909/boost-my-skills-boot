package usecase

import (
	"boost-my-skills-bot/app/internal/bot/models"
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

	buttons := tgbotapi.NewInlineKeyboardButtonData("get an answer", fmt.Sprintf(callbackData, questionID, callbackType))
	rows = append(rows, buttons)

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))

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
		u.log.Errorf("botUC.hideKeyboard.Send(). %s", err.Error())
		return
	}

	return
}
