package tgbot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) createSubdirectionsKeyboardAddQuestion(subdirections []string) (keyboard tgbotapi.InlineKeyboardMarkup) {

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subdirections[0], GoCallbackDataAddQuestion),
			tgbotapi.NewInlineKeyboardButtonData(subdirections[1], ComputerSinceCallbackDataAddQuestion),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subdirections[2], NetworkCallbackDataAddQuestion),
			tgbotapi.NewInlineKeyboardButtonData(subdirections[3], DBCallbackDataAddQuestion),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subdirections[4], AlgorithmsCallbackDataAddQuestion),
			tgbotapi.NewInlineKeyboardButtonData(subdirections[5], ArchitectureCallbackDataAddQuestion),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subdirections[6], GeneralCallbackDataAddQuestion),
		),
	)

	return
}

func (t *TgBot) createSubdirectionsKeyboardAskMe(subdirections []string) (keyboard tgbotapi.InlineKeyboardMarkup) {

	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subdirections[0], GoCallbackDataAskMe),
			tgbotapi.NewInlineKeyboardButtonData(subdirections[1], ComputerSinceCallbackDataAskMe),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subdirections[2], NetworkCallbackDataAskMe),
			tgbotapi.NewInlineKeyboardButtonData(subdirections[3], DBCallbackDataAskMe),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subdirections[4], AlgorithmsCallbackDataAskMe),
			tgbotapi.NewInlineKeyboardButtonData(subdirections[5], ArchitectureCallbackDataAskMe),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(subdirections[6], GeneralCallbackDataAskMe),
		),
	)

	return
}

func (t *TgBot) createAnswerKeyboard(questionID string) (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(getAnswerButton, fmt.Sprintf("%s %s", getAnswerCallbackData, questionID)),
		),
	)

	return
}

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

func (t *TgBot) createDirectionsKeyboard() (keyboard tgbotapi.InlineKeyboardMarkup) {
	keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backenddButton, backendCallbackData),
			tgbotapi.NewInlineKeyboardButtonData(frontendButton, frontednCallbackData),
		),
	)

	return
}

func (t *TgBot) createMainMenuKeyboard() (keyboard tgbotapi.ReplyKeyboardMarkup) {

	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(getUUIDButton),
			tgbotapi.NewKeyboardButton(askMeButton),
			tgbotapi.NewKeyboardButton(addQuestionButton),
		),
	)

	keyboard.OneTimeKeyboard = false // Hide keyboard after one use
	keyboard.ResizeKeyboard = true   // Resizes keyboard depending on the user's device

	return
}
