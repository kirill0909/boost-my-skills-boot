package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type BotUC struct {
	cfg             *config.Config
	pgRepo          bot.PgRepository
	BotAPI          *tgbotapi.BotAPI
	stateUsers      map[int64]models.AddInfoParams
	stateDirections *models.DirectionsData
}

func NewBotUC(
	cfg *config.Config,
	pgRepo bot.PgRepository,
	botAPI *tgbotapi.BotAPI,
	stateUsers map[int64]models.AddInfoParams,
	stateDirections *models.DirectionsData,
) bot.Usecase {
	return &BotUC{
		cfg:             cfg,
		pgRepo:          pgRepo,
		BotAPI:          botAPI,
		stateUsers:      stateUsers,
		stateDirections: stateDirections,
	}
}

func (u *BotUC) GetUUID(ctx context.Context, params models.GetUUID) (result string, err error) {

	isAdmin, err := u.pgRepo.IsAdmin(ctx, params)
	if err != nil {
		return
	}
	if !isAdmin {
		return notAdmin, nil
	}

	result, err = u.pgRepo.GetUUID(ctx)
	if err != nil {
		return
	}

	return u.createTgLink(result), nil
}

func (u *BotUC) createTgLink(param string) string {
	return fmt.Sprintf(u.cfg.TgBot.Prefix, param)
}

func (u *BotUC) UserActivation(ctx context.Context, params models.UserActivation) (err error) {
	return u.pgRepo.UserActivation(ctx, params)
}

func (u *BotUC) SetUpDirection(ctx context.Context, params models.SetUpDirection) (err error) {
	splitedCallbackData := strings.Split(params.CallbackData, " ")
	penultimateElement := splitedCallbackData[len(splitedCallbackData)-2]

	directionID, err := strconv.Atoi(penultimateElement)
	if err != nil {
		err = errors.Wrapf(err, "BotUC.SetUpDirection.Atoi(%s)", penultimateElement)
		return
	}
	params.DirectionID = directionID

	if err = u.pgRepo.SetUpDirection(ctx, params); err != nil {
		return
	}

	if err = u.hideKeyboard(params.ChatID, params.MessageID); err != nil {
		return
	}

	msg := tgbotapi.NewMessage(params.ChatID, readyMessage)
	msg.ReplyMarkup = u.createMainMenuKeyboard()
	if _, err = u.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (u *BotUC) GetRandomQuestion(ctx context.Context, params models.AksMeCallbackParams) (
	result models.AskMeCallbackResult, err error) {
	return u.pgRepo.GetRandomQuestion(ctx, params)
}

func (u *BotUC) GetAnswer(ctx context.Context, questionID int) (result string, err error) {
	return u.pgRepo.GetAnswer(ctx, questionID)
}

func (u *BotUC) SaveQuestion(ctx context.Context, params models.SaveQuestionParams) (result int, err error) {
	return u.pgRepo.SaveQuestion(ctx, params)
}

func (u *BotUC) SaveAnswer(ctx context.Context, params models.SaveAnswerParams) (err error) {
	return u.pgRepo.SaveAnswer(ctx, params)
}

func (u *BotUC) GetSubdirections(ctx context.Context, params models.GetSubdirectionsParams) (result []string, err error) {
	return u.pgRepo.GetSubdirections(ctx, params)
}

func (u *BotUC) GetSubSubdirections(ctx context.Context, params models.GetSubSubdirectionsParams) (result []string, err error) {
	return u.pgRepo.GetSubSubdirections(ctx, params)
}

func (u *BotUC) HandleAddInfoCommand(ctx context.Context, chatID int64) (err error) {

	u.stateUsers[chatID] = models.AddInfoParams{State: awaitingSubdirection}
	directionID, err := u.pgRepo.GetDirectionIDByChatID(ctx, chatID)
	if err != nil {
		return
	}

	subdirections := u.stateDirections.GetSubdirectionsByDirectionID(directionID)
	if len(subdirections) == 0 {
		err = fmt.Errorf("subdirections not found")
		return errors.Wrap(err, "BotUC.HandleAddInfoCommand.len(subdirections)")
	}

	msg := tgbotapi.NewMessage(chatID, "")
	if len(subdirections) == 0 {
		msg.Text = noOneSubdirectionsFoundMessage
		if _, err = u.BotAPI.Send(msg); err != nil {
			return
		}
		u.stateUsers[chatID] = models.AddInfoParams{State: idle}

		return
	}

	msg.Text = subdirectionAddInfoMessage
	msg.ReplyMarkup = u.createSubdirectionsKeyboardAddInfo(subdirections)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

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
		msg := tgbotapi.NewMessage(params.ChatID, notQuestionsMessage)
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

func (u *BotUC) HandleAskMeCommand(ctx context.Context, params models.AskMeParams) (err error) {
	directionID, err := u.pgRepo.GetDirectionIDByChatID(ctx, params.ChatID)
	if err != nil {
		return
	}

	subdirections := u.stateDirections.GetSubdirectionsByDirectionID(directionID)

	msg := tgbotapi.NewMessage(params.ChatID, subdirectionAskMeMessage)
	msg.ReplyMarkup = u.createSubdirectionsKeyboardAskMe(subdirections)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return errors.Wrap(err, "BotUC.HandleAskMeCommand.Send")
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
		msg := tgbotapi.NewMessage(params.ChatID, notQuestionsMessage)
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

func (u *BotUC) SyncDirectionsInfo(ctx context.Context) (err error) {
	directionsInfo, err := u.pgRepo.GetDirectionsInfo(ctx)
	if err != nil {
		return
	}

	subdirectionsInfo, err := u.pgRepo.GetSubdirectionsInfo(ctx)
	if err != nil {
		return
	}

	subSubdirectionsInfo, err := u.pgRepo.GetSubSubdirectionsInfo(ctx)
	if err != nil {
		return
	}

	u.stateDirections.Store(directionsInfo, subdirectionsInfo, subSubdirectionsInfo)

	return
}
