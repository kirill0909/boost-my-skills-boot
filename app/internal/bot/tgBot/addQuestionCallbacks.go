package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"regexp"
	"strconv"

	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
	t.stateUsers[chatID] = models.AddInfoParams{
		State: awaitingSubSubdirection, SubdirectionID: subdirectionID}

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
