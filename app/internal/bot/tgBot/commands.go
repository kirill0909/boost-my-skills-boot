package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"

	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) handleGetUUIDCommand(chatID int64, params models.GetUUID) (err error) {
	ctx := context.Background()

	result, err := t.tgUC.GetUUID(ctx, params)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, result)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleStartCommand(chatID int64, params models.UserActivation, text string) (err error) {
	ctx := context.Background()

	splitedText := strings.Split(text, " ")
	if len(splitedText) != 2 {
		err = fmt.Errorf("Error invite token extracting")
		return
	}

	params.UUID = splitedText[1]
	if err = t.tgUC.UserActivation(ctx, params); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, wellcomeMessage)
	msg.ReplyMarkup = t.createDirectionsKeyboard()
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleAskMeCommand(chatID int64, params models.AskMeParams) (err error) {
	ctx := context.Background()

	result, err := t.tgUC.GetRandomQuestion(ctx, params)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(params.ChatID, result.Question)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = t.createAnswerKeyboard(fmt.Sprintf("%d", result.QuestionID))
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

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
