package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
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

func (t *TgBot) handleStartCommand(params models.UserActivation) (err error) {
	ctx := context.Background()

	splitedText := strings.Split(params.Text, " ")
	if len(splitedText) != 2 {
		err = fmt.Errorf("Error invite token extracting")
		return
	}

	params.UUID = splitedText[1]
	if err = t.tgUC.UserActivation(ctx, params); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(params.ChatID, wellcomeMessage)
	msg.ReplyMarkup = t.createDirectionsKeyboard(t.stateDirections.DirectionInfo)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleAskMeCommand(chatID int64, params models.AskMeParams) (err error) {
	ctx := context.Background()

	if err = t.tgUC.HandleAskMeCommand(ctx, params); err != nil {
		return
	}

	return
}

func (t *TgBot) handleAddInfoCommand(chatID int64) (err error) {
	ctx := context.Background()

	if err = t.tgUC.HandleAddInfoCommand(ctx, chatID); err != nil {
		return
	}

	return
}

func (t *TgBot) handlePrintQuestionsCommand(chatID int64) (err error) {
	ctx := context.Background()

	if err = t.tgUC.HandlePrintQuestions(ctx, models.PrintQuestionsParams{ChatID: chatID}); err != nil {
		return err
	}

	return
}
