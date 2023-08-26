package tgbot

import (
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
