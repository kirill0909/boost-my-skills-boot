package usecase

import (
	models "boost-my-skills-bot/internal/models/bot"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (t *BotUC) createSubdirectionsKeyboardAddInfo(subdirections []models.SubdirectionInfo) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(
			subdirections[i].SubdirectionName,
			fmt.Sprintf("%d %d", subdirections[i].SubdirectionID, t.cfg.CallbackType.SubdirectionAddInfo))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}

func (t *BotUC) createSubSubdirectionsKeyboardAddInfo(subSubdirections []models.SubSubdirectionInfo) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subSubdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(
			subSubdirections[i].SubSubdirectionName,
			fmt.Sprintf("%d %d", subSubdirections[i].SubSubdirectionID, t.cfg.CallbackType.SubSubdirectionAddInfo))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subSubdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}

func (t *BotUC) createSubdirectionsKeyboardAskMe(subdirections []models.SubdirectionInfo) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(
			subdirections[i].SubdirectionName,
			fmt.Sprintf("%d %d", subdirections[i].SubdirectionID, t.cfg.CallbackType.SubdirectionAskMe))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}

func (t *BotUC) createSubSubdirectionsKeyboardAskMe(subSubdirections []models.SubSubdirectionInfo) (keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subSubdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(
			subSubdirections[i].SubSubdirectionName,
			fmt.Sprintf("%d %d", subSubdirections[i].SubSubdirectionID, t.cfg.CallbackType.SubSubdirectionAskMe))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subSubdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}

func (t *BotUC) createSubdirectionsKeyboardPrintQuestions(subdirections []models.SubdirectionInfo) (
	keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(
			subdirections[i].SubdirectionName,
			fmt.Sprintf("%d %d", subdirections[i].SubdirectionID, t.cfg.CallbackType.SubdirectionPrintQuestions))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}

func (t *BotUC) createSubSubdirectionsKeyboardPrintQuestions(subSubdirections []models.SubSubdirectionInfo) (
	keyboard tgbotapi.InlineKeyboardMarkup) {

	var rows []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(subSubdirections); i++ {
		buttons := tgbotapi.NewInlineKeyboardButtonData(
			subSubdirections[i].SubSubdirectionName,
			fmt.Sprintf("%d %d", subSubdirections[i].SubSubdirectionID, t.cfg.CallbackType.SubSubdirectionPrintQuestions))
		rows = append(rows, buttons)

		if (i+1)%2 == 0 || i == len(subSubdirections)-1 {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(rows...))
			rows = rows[:0]
		}

	}

	return
}

func (t *BotUC) createAnswerKeyboard(questionID int) (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(getAnswerButton, fmt.Sprintf("%d %d", questionID, t.cfg.CallbackType.GetAnAnswer)),
		),
	)

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
		err = errors.Wrap(err, "BotUC.hideKeyboard.Send")
		return
	}

	return
}

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

	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			keyboardButtons...,
		),
	)

	keyboard.OneTimeKeyboard = false // Hide keyboard after one use
	keyboard.ResizeKeyboard = true   // Resizes keyboard depending on the user's device

	return
}
