package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) handleEnteredQuestion(chatID int64, text string) (err error) {
	ctx := context.Background()

	msg := tgbotapi.NewMessage(chatID, handleEnteredQuestionMessage)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	// save question in db and return id
	questionID, err := t.tgUC.SaveQuestion(ctx, models.SaveQuestionParams{ChatID: chatID, Question: text})
	if err != nil {
		return
	}

	t.userStates[chatID] = models.AddQuestionParams{State: awaitingAnswer, QuestionID: questionID}

	return
}

func (t *TgBot) handleEnteredAnswer(chatID int64) (err error) {
	t.userStates[chatID] = models.AddQuestionParams{State: idle}

	msg := tgbotapi.NewMessage(chatID, handleEnteredAnswerMessage)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}
