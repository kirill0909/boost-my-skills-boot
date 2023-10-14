package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) createDirectionsKeyboard(directions []models.DirectionInfo) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(directions); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(
			directions[i].DirectionName,
			fmt.Sprintf("%d %d", directions[i].DirectionID, t.cfg.CallbackType.Direction))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(directions)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}
