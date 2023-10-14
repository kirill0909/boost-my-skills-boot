package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
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

func (t *TgBot) handleGetAnAnswerCallbackData(chatID int64, messageID int, callbackData string) (err error) {
	ctx := context.Background()

	if err = t.tgUC.HandleGetAnAnswerCallbackData(ctx, models.GetAnAnswerParams{
		ChatID:       chatID,
		MessageID:    messageID,
		CallbackData: callbackData,
	}); err != nil {
		return
	}

	return
}

func (t *TgBot) handleAskMeSubSubdirectionCallbackData(chatID int64, messageID int, callbackData string) (err error) {
	ctx := context.Background()

	if err = t.tgUC.HandleAskMeSubSubdirectionCallbackData(ctx, models.AskMeParams{
		ChatID:       chatID,
		MessageID:    messageID,
		CallbackData: callbackData,
	}); err != nil {

	}

	return
}
