package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) createSubSubdirectionsKeyboardAddQuestion(subSubdirections []string) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subSubdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(subSubdirections[i], callbackDataSubSubdirectionAddQuestion[i])
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subSubdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}

func (t *TgBot) createSubdirectionsKeyboardAskMe(subdirections []string) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(subdirections[i], callbackDataAskMe[i])
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}
	}

	return
}

func (t *TgBot) createSubSubdirectionsKeyboardAskMe(subdirections []string) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(subdirections[i], callbackDataSubSubdirectionAskMe[i])
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}
	}

	return
}

// func (t *TgBot) createAnswerKeyboard(questionID string) (keyboard tgbotapi.InlineKeyboardMarkup) {
// 	keyboard = tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(
// 			tgbotapi.NewInlineKeyboardButtonData(getAnswerButton, fmt.Sprintf("%s %s", getAnswerCallbackData, questionID)),
// 		),
// 	)
//
// 	return
// }

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
