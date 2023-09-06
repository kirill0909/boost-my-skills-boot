package tgbot

import (
	"context"

	models "boost-my-skills-bot/internal/models/bot"
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) handleSubdirectionsCallbackAddQuestion(chatID int64, subdirection string, messageID int) (err error) {
	ctx := context.Background()

	subdirectionID, err := t.extractSubDirectionID(subdirection)
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
		if err = t.handleSubSubdirectionsExistsCaseAddQuestion(chatID, subSubdirections, subdirectionID); err != nil {
			return
		}
	default:
		if err = t.handleAddQuestionCase(chatID, subdirectionID); err != nil {
			return
		}
	}

	return
}

func (t *TgBot) extractSubDirectionID(callbackData string) (result int, err error) {
	layout := "^\\d+"
	re := regexp.MustCompile(layout)

	subdirection := re.FindString(callbackData)

	result, err = strconv.Atoi(subdirection)
	if err != nil {
		err = errors.Wrap(err, "TgBot.extractSubDirectionID")
		return
	}

	return
}

func (t *TgBot) handleSubSubdirectionsExistsCaseAddQuestion(chatID int64, subSubdirections []string, subdirectionID int) (err error) {
	msg := tgbotapi.NewMessage(chatID, "Choose sub sub direction")
	msg.ReplyMarkup = t.createSubSubdirectionsKeyboardAddQuestion(subSubdirections)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}
	t.userStates[chatID] = models.AddQuestionParams{
		State: awaitingSubSubdirection, SubdirectionID: subdirectionID}

	return
}

func (t *TgBot) handleAddQuestionCase(chatID int64, ids ...int) (err error) {
	n := len(ids)
	switch {
	case n == 1:
		t.userStates[chatID] = models.AddQuestionParams{State: awaitingQuestion, SubdirectionID: ids[0]}
	case n == 2:
		t.userStates[chatID] = models.AddQuestionParams{State: awaitingQuestion, SubdirectionID: ids[0], SubSubdirectionID: ids[1]}
	default:
		return fmt.Errorf("TgBot.handleEnteredQuestion. Wrong length(%d) of directions ids", n)
	}

	msg := tgbotapi.NewMessage(chatID, "Alright, enter your question")
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleSubSubdirectionsCallbackAddQuestion(chatID int64, callbackData string, messageID int) (err error) {

	if err = t.hideKeyboard(chatID, messageID); err != nil {
		return
	}

	ids, err := t.extractDirectionsIDs(callbackData)
	if err != nil {
		return
	}

	if err = t.handleAddQuestionCase(chatID, ids[0], ids[1]); err != nil {
		return
	}

	return
}

func (t *TgBot) extractDirectionsIDs(callbackData string) (result []int, err error) {
	layout := "\\d+"
	re := regexp.MustCompile(layout)
	directionsIDs := re.FindAllString(callbackData, 2)

	for _, value := range directionsIDs {
		id, err := strconv.Atoi(value)
		if err != nil {
			return []int{}, errors.Wrap(err, "tgBot.extractDirectionsIDs")
		}

		result = append(result, id)
	}

	return
}
