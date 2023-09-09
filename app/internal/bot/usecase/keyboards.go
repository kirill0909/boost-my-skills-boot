package usecase

import (
	models "boost-my-skills-bot/internal/models/bot"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *BotUC) createSubdirectionsKeyboardAddInfo(subdirections []models.SubdirectionInfo) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(
			subdirections[i].SubdirectionName, fmt.Sprintf("%d %d", subdirections[i].SubdirectionID, t.cfg.CallbackType.SubdirectionAddInfo))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}

func (t *BotUC) hideKeyboard(chatID int64, messageID int) (err error) {
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

func (t *BotUC) createMainMenuKeyboard() (keyboard tgbotapi.ReplyKeyboardMarkup) {

	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(getUUIDButton),
			tgbotapi.NewKeyboardButton(askMeButton),
			tgbotapi.NewKeyboardButton(addInfoButton),
		),
	)

	keyboard.OneTimeKeyboard = false // Hide keyboard after one use
	keyboard.ResizeKeyboard = true   // Resizes keyboard depending on the user's device

	return
}
