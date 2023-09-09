package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) handleEnteredQuestion(chatID int64, text string, ids ...int) (err error) {
	ctx := context.Background()

	questionParams := models.SaveQuestionParams{ChatID: chatID, Question: text}
	n := len(ids)
	switch {
	case n == 1:
		questionParams.SubdirectionID = ids[0]
	case n == 2:
		questionParams.SubdirectionID = ids[0]
		questionParams.SubSubdirectionID = ids[1]
	default:
		return fmt.Errorf("TgBot.handleEnteredQuestion. Wrong length(%d) of directions ids", n)
	}

	questionID, err := t.tgUC.SaveQuestion(ctx, questionParams)
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, handleEnteredQuestionMessage)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	t.stateUseres[chatID] = models.AddQuestionParams{State: awaitingAnswer, QuestionID: questionID}

	return
}

func (t *TgBot) handleEnteredAnswer(chatID int64, text string) (err error) {
	ctx := context.Background()
	questionParams, ok := t.stateUseres[chatID]
	if !ok {
		return
	}

	if err = t.tgUC.SaveAnswer(ctx, models.SaveAnswerParams{QuestionID: questionParams.QuestionID, Answer: text}); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, handleEnteredAnswerMessage)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	t.stateUseres[chatID] = models.AddQuestionParams{State: idle}

	return
}
