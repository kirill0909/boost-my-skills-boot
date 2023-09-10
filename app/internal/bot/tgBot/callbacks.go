package tgbot

import (
	"context"
	"log"
	"strconv"

	models "boost-my-skills-bot/internal/models/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) handleDirectionCallbackData(chatID int64, messageID int, callbackData string) (err error) {
	ctx := context.Background()

	if err = t.tgUC.SetUpDirection(ctx, models.SetUpDirection{
		ChatID:       chatID,
		MessageID:    messageID,
		CallbackData: callbackData}); err != nil {
		return
	}

	return
}

func (t *TgBot) handleAddInfoSubdirectionCallbackData(chatID int64, messageID int, callbackData string) (err error) {
	ctx := context.Background()

	if err = t.tgUC.HandleAddInfoSubdirectionCallbackData(ctx, models.AddInfoSubdirectionParams{
		ChatID:       chatID,
		MessageID:    messageID,
		CallbackData: callbackData,
	}); err != nil {
		return
	}

	return
}

func (t *TgBot) handleAddInfoSubSubdirectionCallbackData(chatID int64, messageID int, callbackData string) (err error) {
	ctx := context.Background()

	log.Println(t.stateUsers[chatID].SubdirectionID)
	if err = t.tgUC.HandleAddInfoSubSubdirectionCallbackData(ctx, models.AddInfoSubSubdirectionParams{
		ChatID:         chatID,
		MessageID:      messageID,
		CallbackData:   callbackData,
		SubdirectionID: t.stateUsers[chatID].SubdirectionID,
	}); err != nil {
		return
	}

	return
}

func (t *TgBot) handleAskMeSubdirectionCallbackData(chatID int64, messageID int, callbackData string) (err error) {
	ctx := context.Background()

	if err = t.tgUC.HandleAskMeSubdirectionCallbackData(ctx, models.AskMeParams{
		ChatID:       chatID,
		MessageID:    messageID,
		CallbackData: callbackData,
	}); err != nil {
		return
	}

	return
}

func (t *TgBot) handleGetAnswerCallbackData(chatID int64, questionID string, messageID int) (err error) {
	ctx := context.Background()

	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	id, err := strconv.Atoi(questionID)
	if err != nil {
		return
	}

	answer, err := t.tgUC.GetAnswer(ctx, id)
	if err != nil {
		return
	}

	if len(answer) == 0 {
		msg := tgbotapi.NewMessage(chatID, unableToGetAnswer)
		if _, err = t.BotAPI.Send(msg); err != nil {
			return
		}
		return
	}

	msg := tgbotapi.NewMessage(chatID, answer)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}
