package tgbot

import (
	"context"
	"fmt"

	models "boost-my-skills-bot/internal/models/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) handleSubSubdirectonsCallbackAskMe(chatID int64, callbackData string, messageID int) (err error) {
	ctx := context.Background()

	// hide sub sub directions keyboard
	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	ids, err := t.extractDirectionsIDs(callbackData)
	if err != nil {
		return
	}

	if err = t.handleAskMeCase(ctx, chatID, ids[0], ids[1]); err != nil {
		return err
	}

	return
}

func (t *TgBot) handleSubdirectionsCallbackAskMe(chatID int64, callbackData string, messageID int) (err error) {
	ctx := context.Background()

	subdirectionID, err := t.extractSubDirectionID(callbackData)
	if err != nil {
		return
	}

	// hide subdirections keyboard
	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	subSubdirections, err := t.tgUC.GetSubSubdirections(
		ctx, models.GetSubSubdirectionsParams{ChatID: chatID, SubdirectionID: subdirectionID})
	if err != nil {
		return
	}

	n := len(subSubdirections)
	switch {
	case n > 0:
		if err = t.handleSubSubdirectionsExistsAskMeCase(chatID, subSubdirections, subdirectionID); err != nil {
			return
		}
	default:
		if err = t.handleAskMeCase(ctx, chatID, subdirectionID); err != nil {
			return
		}
	}

	return
}

func (t *TgBot) handleAskMeCase(ctx context.Context, chatID int64, ids ...int) (err error) {

	askMeCallbackParams := models.AksMeCallbackParams{ChatID: chatID}
	n := len(ids)
	switch {
	case n == 1:
		askMeCallbackParams.SubdirectionID = ids[0]
	case n == 2:
		askMeCallbackParams.SubdirectionID = ids[0]
		askMeCallbackParams.SubSubdirectionID = ids[1]
	}

	result, err := t.tgUC.GetRandomQuestion(ctx, askMeCallbackParams)
	if err != nil {
		return
	}

	if len(result.Question) == 0 {
		msg := tgbotapi.NewMessage(chatID, notQuestionsMessage)
		if _, err = t.BotAPI.Send(msg); err != nil {
			return
		}
		return
	}

	msg := tgbotapi.NewMessage(chatID, result.Question)
	msg.ReplyMarkup = t.createAnswerKeyboard(fmt.Sprintf("%d", result.QuestionID))
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleSubSubdirectionsExistsAskMeCase(chatID int64, subSubdirections []string, subdirectionID int) (err error) {
	msg := tgbotapi.NewMessage(chatID, "Choose sub sub direction")
	msg.ReplyMarkup = t.createSubSubdirectionsKeyboardAskMe(subSubdirections)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}
