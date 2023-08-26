package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) handleEnteredQuestion(chatID int64, text string, subdirectionID int) (err error) {
	ctx := context.Background()

	questionID, err := t.tgUC.SaveQuestion(ctx, models.SaveQuestionParams{
		ChatID: chatID, Question: text, SubdirectionID: subdirectionID})
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, handleEnteredQuestionMessage)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	t.userStates[chatID] = models.AddQuestionParams{State: awaitingAnswer, QuestionID: questionID}

	return
}

func (t *TgBot) handleEnteredAnswer(chatID int64, text string) (err error) {
	ctx := context.Background()
	questionParams, ok := t.userStates[chatID]
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

	t.userStates[chatID] = models.AddQuestionParams{State: idle}

	return
}
