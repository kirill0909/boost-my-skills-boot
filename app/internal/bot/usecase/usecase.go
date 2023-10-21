package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"strconv"
	"strings"
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

func (u *BotUC) HandlePrintQuestions(ctx context.Context, params models.PrintQuestionsParams) (err error) {
	directionID, err := u.pgRepo.GetDirectionIDByChatID(ctx, params.ChatID)
	if err != nil {
		return
	}

	subdirections := u.stateDirections.GetSubdirectionsByDirectionID(directionID)
	if len(subdirections) == 0 {
		err = fmt.Errorf("subdirections not found")
		return errors.Wrap(err, "BotUC.HandlePrintQuestions.len(subdirections)")
	}

	msg := tgbotapi.NewMessage(params.ChatID, "")
	if len(subdirections) == 0 {
		msg.Text = noOneSubdirectionsFoundMessage
		if _, err = u.BotAPI.Send(msg); err != nil {
			return
		}

		return
	}

	msg.Text = subdirectionPrintQuestions
	msg.ReplyMarkup = u.createSubdirectionsKeyboardPrintQuestions(subdirections)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return
	}

	return
}

func (u *BotUC) HandlePrintQuestionsSubdirectionCallbackData(ctx context.Context, params models.PrintQuestionsParams) (err error) {
	if err = u.hideKeyboard(params.ChatID, params.MessageID); err != nil {
		return
	}

	splitedCallbackData := strings.Split(params.CallbackData, " ")
	subdirectionIDCallbackData := splitedCallbackData[len(splitedCallbackData)-2]

	subdirectionID, err := strconv.Atoi(subdirectionIDCallbackData)
	if err != nil {
		err = errors.Wrapf(err, "BotUC.HandlePrintQuestionsSubdirectionCallbackData.Atoi(%s)", subdirectionIDCallbackData)
		return
	}
	params.SubdirectionID = subdirectionID

	subSubdirections := u.stateDirections.GetSubSubdirectionsBySubdirectionID(params.SubdirectionID)

	n := len(subSubdirections)
	switch {
	case n > 0:
		if err = u.handlePrintQuestionsSubSubdirectionsCase(ctx, params); err != nil {
			return
		}
	default:
		if err = u.handlePrintQuestionsSubSubdirectionsDefaultCase(ctx, params); err != nil {
			return
		}
	}

	return
}

func (u *BotUC) handlePrintQuestionsSubSubdirectionsCase(ctx context.Context, params models.PrintQuestionsParams) (err error) {
	subSubdirections := u.stateDirections.GetSubSubdirectionsBySubdirectionID(params.SubdirectionID)
	u.stateUsers[params.ChatID] = models.AddInfoParams{
		SubdirectionID: params.SubdirectionID,
	}

	msg := tgbotapi.NewMessage(params.ChatID, subSubdirectionAskMeMessage)
	msg.ReplyMarkup = u.createSubSubdirectionsKeyboardPrintQuestions(subSubdirections)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return errors.Wrap(err, "BotUC.handlePrintQuestionsSubSubdirectionscase.Send()")
	}

	return
}

func (u *BotUC) handlePrintQuestionsSubSubdirectionsDefaultCase(ctx context.Context, params models.PrintQuestionsParams) (err error) {
	msg := tgbotapi.NewMessage(params.ChatID, notSubSubDirectionMessage)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return errors.Wrap(err, "BotUC.handlePrintQuestionsSubSubdirectionsDefaultCase.Send()")
	}

	return
}

func (u *BotUC) HandlePrintQuestionsSubSubdirectionCallbackData(ctx context.Context, params models.PrintQuestionsParams) (err error) {
	if err = u.hideKeyboard(params.ChatID, params.MessageID); err != nil {
		return
	}

	splitedCallbackData := strings.Split(params.CallbackData, " ")
	subSubdirectionIDCallbackData := splitedCallbackData[len(splitedCallbackData)-2]

	subSubdirectionID, err := strconv.Atoi(subSubdirectionIDCallbackData)
	if err != nil {
		err = errors.Wrapf(err, "BotUC.HandlePrintQuestionsSubSubdirectionCallbackData.Atoi(%s)", subSubdirectionIDCallbackData)
		return
	}

	params.SubSubdirectionID = subSubdirectionID
	params.SubdirectionID = u.stateUsers[params.ChatID].SubdirectionID

	result, err := u.pgRepo.PrintQuestions(models.PrintQuestionsParams{
		ChatID:            params.ChatID,
		SubdirectionID:    params.SubdirectionID,
		SubSubdirectionID: params.SubSubdirectionID,
	})
	if err != nil {
		return
	}

	if len(result) == 0 {
		msg := tgbotapi.NewMessage(params.ChatID, notOneQuestion)
		if _, err = u.BotAPI.Send(msg); err != nil {
			return errors.Wrap(err, "BotUC.HandlePrintQuestionsSubSubdirectionCallbackData.len(result).Send")
		}

		return
	}

	for _, question := range result {
		msg := tgbotapi.NewMessage(params.ChatID, question.Question)
		if _, err = u.BotAPI.Send(msg); err != nil {
			return errors.Wrap(err, "BotUC.HandlePrintQuestionsSubSubdirectionCallbackData..Send")
		}
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
