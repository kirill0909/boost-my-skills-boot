package tgbot

import (
	"context"
	"fmt"
	"strconv"

	models "boost-my-skills-bot/internal/models/bot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) handleFrontendCallbackData(chatID int64, messageID int) (err error) {
	ctx := context.Background()

	if err = t.tgUC.SetUpFrontendDirection(ctx, chatID); err != nil {
		return
	}

	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, readyMessage)
	msg.ReplyMarkup = t.createMainMenuKeyboard()
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleBackendCallbackData(chatID int64, messageID int) (err error) {
	ctx := context.Background()

	if err = t.tgUC.SetUpBackendDirection(ctx, chatID); err != nil {
		return
	}

	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, readyMessage)
	msg.ReplyMarkup = t.createMainMenuKeyboard()
	if _, err = t.BotAPI.Send(msg); err != nil {
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

func (t *TgBot) handleSubdirectionsCallbackAskMe(chatID int64, subdirection string, messageID int) (err error) {
	ctx := context.Background()

	// hide subdirections keyboard
	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	subdirectionID, err := strconv.Atoi(subdirection)
	if err != nil {
		return
	}

	result, err := t.tgUC.GetRandomQuestion(ctx, models.SubdirectionsCallbackParams{
		ChatID:         chatID,
		SubdirectionID: subdirectionID})
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

func (t *TgBot) handleSubdirectionsCallbackAddQuestion(chatID int64, subdirection string, messageID int) (err error) {
	ctx := context.Background()
	subdirectionID, err := strconv.Atoi(subdirection)
	if err != nil {
		return
	}

	// hide subdirections keyboard
	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	subSubdirections, err := t.tgUC.GetSubSubdirections(ctx, models.GetSubSubdirectionsParams{
		ChatID:         chatID,
		SubdirectionID: subdirectionID,
	})
	if err != nil {
		return
	}

	n := len(subSubdirections)
	switch {
	case n > 0:
		if err = t.handleSubSubdirectionsExistsCase(chatID, subSubdirections, subdirectionID); err != nil {
			return
		}
	default:
		if err = t.handleDefaultSubdirectionsCase(chatID, subdirectionID); err != nil {
			return
		}
	}

	return
}

func (t *TgBot) handleSubSubdirectionsExistsCase(chatID int64, subSubdirections []string, subdirectionID int) (err error) {
	msg := tgbotapi.NewMessage(chatID, "Choose sub sub direction")
	msg.ReplyMarkup = t.createSubSubdirectionsKeyboardAskMe(subSubdirections)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}
	t.userStates[chatID] = models.AddQuestionParams{
		State: awaitingSubSubdirection, SubdirectionID: subdirectionID}

	return
}

func (t *TgBot) handleDefaultSubdirectionsCase(chatID int64, subdirectionID int) (err error) {
	t.userStates[chatID] = models.AddQuestionParams{State: awaitingQuestion, SubdirectionID: subdirectionID}

	msg := tgbotapi.NewMessage(chatID, "Alright, enter your question")
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}
