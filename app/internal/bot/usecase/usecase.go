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
	directionCallbackData := splitedCallbackData[len(splitedCallbackData)-2]

	directionID, err := strconv.Atoi(directionCallbackData)
	if err != nil {
		err = errors.Wrapf(err, "BotUC.SetUpDirection.Atoi(%s)", directionCallbackData)
		return
	}
	params.DirectionID = directionID

	if err = u.pgRepo.SetUpDirection(ctx, params); err != nil {
		return
	}

	if err = u.hideKeyboard(params.ChatID, params.MessageID); err != nil {
		return
	}

	var isAdmin bool
	msg := tgbotapi.NewMessage(params.ChatID, readyMessage)
	if params.ChatID == u.cfg.AdminChatID {
		isAdmin = true
		msg.ReplyMarkup = u.createMainMenuKeyboard(isAdmin)
	} else {
		isAdmin = false
		msg.ReplyMarkup = u.createMainMenuKeyboard(isAdmin)
	}

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

func (u *BotUC) HandleGetAnAnswerCallbackData(ctx context.Context, params models.GetAnAnswerParams) (err error) {
	if err = u.hideKeyboard(params.ChatID, params.MessageID); err != nil {
		return
	}

	splitedCallbackData := strings.Split(params.CallbackData, " ")
	questionIDCallbackData := splitedCallbackData[0]
	questionID, err := strconv.Atoi(questionIDCallbackData)
	if err != nil {
		return errors.Wrapf(err, "BotUC.HandleGetAnAnswerCallbackData.Atoi(%s)", questionIDCallbackData)
	}

	result, err := u.pgRepo.GetAnswer(ctx, questionID)
	if err != nil {
		return
	}

	if len(result) == 0 {
		msg := tgbotapi.NewMessage(params.ChatID, "Unable to get answer for this question")
		if _, err = u.BotAPI.Send(msg); err != nil {
			return errors.Wrap(err, "BotUC.HandleGetAnAnswerCallbackData.Send")
		}
		return fmt.Errorf("BotUC.HandleGetAnAnswerCallbackData.len(%s)", result)
	}

	msg := tgbotapi.NewMessage(params.ChatID, result)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return errors.Wrap(err, "BotUC.HandleGetAnAnswerCallbackData.Send")
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
