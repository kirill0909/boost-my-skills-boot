package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"log"

	// "log"

	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func (t *TgBot) handleStartCommand(chatID int64, params models.UserActivation, text string) (err error) {
	ctx := context.Background()

	splitedText := strings.Split(text, " ")
	if len(splitedText) != 2 {
		err = fmt.Errorf("Error invite token extracting")
		return
	}

	params.UUID = splitedText[1]
	if err = t.tgUC.UserActivation(ctx, params); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(chatID, wellcomeMessage)
	msg.ReplyMarkup = t.createDirectionsKeyboard()
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleAskMeCommand(chatID int64, params models.AskMeParams) (err error) {
	ctx := context.Background()

	subdirections, err := t.tgUC.GetSubdirections(ctx, models.GetSubdirectionsParams{ChatID: chatID})
	if err != nil {
		return
	}

	msg := tgbotapi.NewMessage(params.ChatID, chooseSubdirectionMessage)
	msg.ReplyMarkup = t.createSubdirectionsKeyboardAskMe(subdirections)
	if _, err = t.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (t *TgBot) handleAddQuestionCommand(chatID int64) (err error) {
	// ctx := context.Background()

	subdirections := t.stateDirections.GetSubdirectionsByDirectionID(1)
	if len(subdirections) == 0 {
		log.Println("Not found")
		return
	}

	log.Println(subdirections)

	// subdirectionsData, ok := value.(models.DirectionsData)
	// if !ok {
	// 	log.Println("Unable to coast any to models.SubdirectionInfo")
	// 	return
	// }
	//
	// log.Println(subdirectionsData.SubSubdirectionInfo)

	// t.stateUser[chatID] = models.AddQuestionParams{State: awaitingSubdirection}
	// subdirections, err := t.tgUC.GetSubdirections(ctx, models.GetSubdirectionsParams{ChatID: chatID})
	// if err != nil {
	// 	return
	// }
	//
	// msg := tgbotapi.NewMessage(chatID, "")
	// if len(subdirections) == 0 {
	// 	msg.Text = noOneSubdirectionsFoundMessage
	// 	if _, err = t.BotAPI.Send(msg); err != nil {
	// 		return
	// 	}
	// 	t.stateUser[chatID] = models.AddQuestionParams{State: idle}
	//
	// 	return
	// }
	//
	// msg.Text = directionQuestionMessage
	// msg.ReplyMarkup = t.createSubdirectionsKeyboardAddQuestion(subdirections)
	// if _, err = t.BotAPI.Send(msg); err != nil {
	// 	return
	// }

	return
}
