package usecase

import (
	"context"
	"strconv"
	"strings"

	models "boost-my-skills-bot/internal/models/bot"
	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (u *BotUC) HandleAskMeSubdirectionCallbackData(ctx context.Context, params models.AskMeParams) (err error) {
	if err = u.hideKeyboard(params.ChatID, params.MessageID); err != nil {
		return
	}

	splitedCallbackData := strings.Split(params.CallbackData, " ")
	subdirectionIDCallbackData := splitedCallbackData[len(splitedCallbackData)-2]

	subdirectionID, err := strconv.Atoi(subdirectionIDCallbackData)
	if err != nil {
		err = errors.Wrapf(err, "BotUC.HandleAskMeSubdirectionCallbackData.Atoi(%s)", subdirectionIDCallbackData)
		return
	}
	params.SubdirectionID = subdirectionID

	subSubdirections := u.stateDirections.GetSubSubdirectionsBySubdirectionID(params.SubdirectionID)

	n := len(subSubdirections)
	switch {
	case n > 0:
		if err = u.handleAskMeSubSubdirectionsCase(ctx, params); err != nil {
			return
		}
	default:
		if err = u.handleAskMeSubdirectionsDefaultCase(ctx, params); err != nil {
			return
		}
	}

	return
}

func (u *BotUC) handleAskMeSubSubdirectionsCase(ctx context.Context, params models.AskMeParams) (err error) {

	subSubdirections := u.stateDirections.GetSubSubdirectionsBySubdirectionID(params.SubdirectionID)
	u.stateUsers[params.ChatID] = models.AddInfoParams{
		SubdirectionID: params.SubdirectionID,
	}

	msg := tgbotapi.NewMessage(params.ChatID, subSubdirectionAskMeMessage)
	msg.ReplyMarkup = u.createSubSubdirectionsKeyboardAskMe(subSubdirections)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return errors.Wrap(err, "BotUC.handleAskMeSubSubdirectionsCase.Send()")
	}

	return
}

func (u *BotUC) handleAskMeSubdirectionsDefaultCase(ctx context.Context, params models.AskMeParams) (err error) {

	result, err := u.pgRepo.GetRandomQuestion(ctx, models.AksMeCallbackParams{
		ChatID:         params.ChatID,
		SubdirectionID: params.SubdirectionID})
	if err != nil {
		return
	}

	if len(result.Question) == 0 {
		msg := tgbotapi.NewMessage(params.ChatID, notOneQuestion)
		if _, err = u.BotAPI.Send(msg); err != nil {
			return errors.Wrap(err, "BotUC.hanleAskMeSubdirectionsDefaultCase.len(result.Question).Send")
		}

		return
	}

	msg := tgbotapi.NewMessage(params.ChatID, result.Question)
	msg.ReplyMarkup = u.createAnswerKeyboard(result.QuestionID)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return errors.Wrap(err, "BotUC.hanleAskMeSubdirectionsDefaultCase.Send")
	}

	return
}

func (u *BotUC) HandleAskMeSubSubdirectionCallbackData(ctx context.Context, params models.AskMeParams) (err error) {
	if err = u.hideKeyboard(params.ChatID, params.MessageID); err != nil {
		return
	}

	splitedCallbackData := strings.Split(params.CallbackData, " ")
	subSubdirectionIDCallbackData := splitedCallbackData[len(splitedCallbackData)-2]

	subSubdirectionID, err := strconv.Atoi(subSubdirectionIDCallbackData)
	if err != nil {
		err = errors.Wrapf(err, "BotUC.HandleAskMeSubdirectionCallbackData.Atoi(%s)", subSubdirectionIDCallbackData)
		return
	}
	params.SubSubdirectionID = subSubdirectionID
	params.SubdirectionID = u.stateUsers[params.ChatID].SubdirectionID

	result, err := u.pgRepo.GetRandomQuestion(ctx, models.AksMeCallbackParams{
		ChatID:            params.ChatID,
		SubdirectionID:    params.SubdirectionID,
		SubSubdirectionID: params.SubSubdirectionID,
	})
	if err != nil {
		return
	}

	if len(result.Question) == 0 {
		msg := tgbotapi.NewMessage(params.ChatID, notOneQuestion)
		if _, err = u.BotAPI.Send(msg); err != nil {
			return errors.Wrap(err, "BotUC.HandleAskMeSubSubdirectionCallbackData.len(result.Question).Send")
		}

		return
	}

	msg := tgbotapi.NewMessage(params.ChatID, result.Question)
	msg.ReplyMarkup = u.createAnswerKeyboard(result.QuestionID)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return errors.Wrap(err, "BotUC.HandleAskMeSubSubdirectionCallbackData.Send")
	}

	return
}
