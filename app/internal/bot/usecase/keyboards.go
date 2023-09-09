package usecase

import (
	models "boost-my-skills-bot/internal/models/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *BotUC) createSubdirectionsKeyboardAddInfo(subdirections []models.SubdirectionInfo) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(subdirections[i].SubdirectionName, subdirections[i].SubdirectionName)
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}
