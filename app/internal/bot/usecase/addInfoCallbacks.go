package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	models "boost-my-skills-bot/internal/models/bot"
)

func (u *BotUC) HandleAddInfoSubdirectionCallbackData(ctx context.Context, params models.AddInfoSubdirectionParams) (err error) {
	if err = u.hideKeyboard(params.ChatID, params.MessageID); err != nil {
		return
	}

	splitedCallbackData := strings.Split(params.CallbackData, " ")
	subdirectionIDCallbackData := splitedCallbackData[len(splitedCallbackData)-2]

	subdirectionID, err := strconv.Atoi(subdirectionIDCallbackData)
	if err != nil {
		err = errors.Wrapf(err, "BotUC.HandleAddInfoSubdirectionCallbackData.Atoi(%s)", subdirectionIDCallbackData)
		return
	}
	params.SubdirectionID = subdirectionID
	u.stateUsers[params.ChatID] = models.AddInfoParams{State: awaitingSubSubdirection, SubdirectionID: subdirectionID}

	subSubdirections := u.stateDirections.GetSubSubdirectionsBySubdirectionID(params.SubdirectionID)

	n := len(subSubdirections)
	switch {
	case n > 0:
		if err = u.handleAddInfoSubdirectionsCase(ctx, params); err != nil {
			return
		}
	default:
		if err = u.hanleAddInfoSubdirectionsDefaultCase(ctx, params); err != nil {
			return
		}
	}

	return
}

func (u *BotUC) handleAddInfoSubdirectionsCase(ctx context.Context, params models.AddInfoSubdirectionParams) (err error) {
	u.stateUsers[params.ChatID] = models.AddInfoParams{State: awaitingSubSubdirection, SubdirectionID: params.SubdirectionID}

	subSubdirections := u.stateDirections.GetSubSubdirectionsBySubdirectionID(params.SubdirectionID)
	if len(subSubdirections) == 0 {
		err = fmt.Errorf("sub subdirections not found")
		return errors.Wrap(err, "handleAddInfoSubdirectionsCase.len(subSubdirections)")
	}

	msg := tgbotapi.NewMessage(params.ChatID, subSubdirectionAddInfoMessage)
	msg.ReplyMarkup = u.createSubSubdirectionsKeyboardAddInfo(subSubdirections)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return errors.Wrap(err, "BotUC.handleAddInfoSubdirectionsCase.Send")
	}

	return
}

func (u *BotUC) hanleAddInfoSubdirectionsDefaultCase(ctx context.Context, params models.AddInfoSubdirectionParams) (err error) {
	u.stateUsers[params.ChatID] = models.AddInfoParams{State: awaitingQuestion, SubdirectionID: params.SubdirectionID}

	msg := tgbotapi.NewMessage(params.ChatID, enterQuestionMessage)
	if _, err = u.BotAPI.Send(msg); err != nil {
		err = errors.Wrap(err, "BotUC.hanleAddInfoSubdirectionsDefaultCase.Send")
		return
	}

	return
}

func (u *BotUC) HandleAddInfoSubSubdirectionCallbackData(ctx context.Context, params models.AddInfoSubSubdirectionParams) (err error) {
	if err = u.hideKeyboard(params.ChatID, params.MessageID); err != nil {
		return
	}

	splitedCallbackData := strings.Split(params.CallbackData, " ")
	subSubdirectionIDCallbackData := splitedCallbackData[len(splitedCallbackData)-2]

	subSubdirectionID, err := strconv.Atoi(subSubdirectionIDCallbackData)
	if err != nil {
		err = errors.Wrapf(err, "BotUC.HandleAddInfoSubSubdirectionCallbackData.Atoi(%s)", subSubdirectionIDCallbackData)
		return
	}

	u.stateUsers[params.ChatID] = models.AddInfoParams{
		State:             awaitingQuestion,
		SubdirectionID:    params.SubdirectionID,
		SubSubdirectionID: subSubdirectionID,
	}

	msg := tgbotapi.NewMessage(params.ChatID, enterQuestionMessage)
	if _, err = u.BotAPI.Send(msg); err != nil {
		err = errors.Wrap(err, "BotUC.HandleAddInfoSubSubdirectionCallbackData.Send")
		return
	}

	return
}
