package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	"boost-my-skills-bot/internal/bot/models"
	"boost-my-skills-bot/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type botUC struct {
	cfg                  *config.Config
	pgRepo               bot.PgRepository
	rdb                  bot.RedisRepository
	redisPubSub          *redis.PubSub
	BotAPI               *tgbotapi.BotAPI
	log                  *logger.Logger
	lastKeyboardChecking int64
}

func NewBotUC(
	cfg *config.Config,
	pgRepo bot.PgRepository,
	rdb bot.RedisRepository,
	redisPubSub *redis.PubSub,
	botAPI *tgbotapi.BotAPI,
	log *logger.Logger,
) bot.Usecase {
	return &botUC{
		cfg:                  cfg,
		pgRepo:               pgRepo,
		rdb:                  rdb,
		redisPubSub:          redisPubSub,
		BotAPI:               botAPI,
		log:                  log,
		lastKeyboardChecking: time.Now().Unix(),
	}
}

func (u *botUC) HandleStartCommand(ctx context.Context, params models.HandleStartCommandParams) error {
	splitedText := strings.Split(params.Text, " ")

	if len(splitedText) != 2 {
		return fmt.Errorf("TgBot.handleStartCommand. wrong len of splited text: %d != 2. params(%+v)", len(splitedText), params)
	}

	uuid := splitedText[1]
	setUserActiveParams := models.SetUserActiveParams{
		TgName: params.TgName,
		ChatID: params.ChatID,
		UUID:   uuid}

	if err := u.pgRepo.SetUserActive(ctx, setUserActiveParams); err != nil {
		return err
	}

	var isAdmin bool
	msg := tgbotapi.NewMessage(params.ChatID, "your account has been successfully activated")
	if params.ChatID == u.cfg.AdminChatID {
		isAdmin = true
		msg.ReplyMarkup = u.createMainMenuKeyboard(isAdmin)
	} else {
		isAdmin = false
		msg.ReplyMarkup = u.createMainMenuKeyboard(isAdmin)
	}

	if _, err := u.BotAPI.Send(msg); err != nil {
		return errors.Wrapf(err, "BotUC.HandleStartCommand.Send")
	}

	return nil
}

func (u *botUC) HandleCreateDirectionCommand(ctx context.Context, params models.HandleCreateDirectionCommandParams) error {
	var err error
	var parentDirectionID int
	var getUserDirectionParms models.GetUserDirectionParams
	if params.CallbackData != "" {
		parentDirectionID, err = strconv.Atoi(params.CallbackData)
		if err != nil {
			return errors.Wrapf(err, "botUC.HandleCreateDirectionCommand.Atoi(). params(%+v)", params)
		}
		getUserDirectionParms.ParentDirectionID = sql.NullInt64{Int64: int64(parentDirectionID), Valid: true}
	}
	getUserDirectionParms.ChatID = params.ChatID

	directions, err := u.pgRepo.GetUserDirection(ctx, getUserDirectionParms)
	if err != nil {
		return err
	}

	var statusID int
	if len(directions) == 0 {
		statusID = utils.AwaitingDirectionNameStatus
		sendMessageParams := models.SendMessageParams{
			ChatID: params.ChatID,
			Text:   "enter name of your direction"}
		u.sendMessage(sendMessageParams)
	} else {
		statusID = utils.AwaitingParentDirecitonStatus
		sendMessageParams := models.SendMessageParams{
			ChatID:         params.ChatID,
			Text:           "choose parent direciton",
			Keyboard:       u.createDirectionsKeyboard(directions),
			IsNeedToRemove: true}
		u.sendMessage(sendMessageParams)
	}

	setAwaitingStatusParams := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: statusID}
	if err := u.rdb.SetAwaitingStatus(ctx, setAwaitingStatusParams); err != nil {
		return err
	}

	return nil
}

func (u *botUC) GetAwaitingStatus(ctx context.Context, chatID int64) (int, error) {
	value, err := u.rdb.GetAwaitingStatus(ctx, chatID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil
		}

		return 0, err
	}

	statusID, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.Wrapf(err, "botUC.GetAwaitingStatus.Atoi(). uanble convert value:%s to int for chatID:%d", value, chatID)
	}

	return statusID, nil
}

func (u *botUC) CreateDirection(ctx context.Context, params models.CreateDirectionParams) error {
	re := regexp.MustCompile(utils.DirectionNameLayout)
	if !re.MatchString(params.DirectionName) {
		return fmt.Errorf("direction name contains unacceptable symbols. params(%+v)", params)
	}

	getParentDirectionResult, err := u.rdb.GetParentDirection(ctx, params.ChatID)
	if errors.Is(err, redis.Nil) {
		err = nil
	} else if err != nil {
		return err
	}

	if getParentDirectionResult != "" {
		parentDirectionID, err := strconv.Atoi(getParentDirectionResult)
		if err != nil {
			return err
		}
		params.ParentDirectionID = sql.NullInt64{Int64: int64(parentDirectionID), Valid: true}
	}

	direction, err := u.pgRepo.CreateDirection(ctx, params)
	if err != nil {
		return err
	}

	if err := u.rdb.ResetAwaitingStatus(ctx, params.ChatID); err != nil {
		return err
	}

	if err := u.rdb.ResetParentDirection(ctx, params.ChatID); err != nil {
		return err
	}

	sendMessageParams := models.SendMessageParams{
		ChatID: params.ChatID,
		Text:   fmt.Sprintf("successfully created \"%s\" direction", direction)}
	u.sendMessage(sendMessageParams)

	return nil
}

func (u *botUC) SetParentDirection(ctx context.Context, params models.SetParentDirectionParams) error {
	parendDirectionID, err := strconv.Atoi(params.CallbackData)
	if err != nil {
		return errors.Wrapf(err, "botUC.SetParentDirection.Atoi(). params(%+v)", params)
	}

	params.ParentDirectionID = parendDirectionID

	if err := u.rdb.SetParentDirection(ctx, params); err != nil {
		return err
	}

	return nil
}

func (u *botUC) HandleAddInfoCommand(ctx context.Context, params models.HandleAddInfoCommandParams) error {
	var err error
	var parentDirectionID int
	var getUserDirectionParams models.GetUserDirectionParams
	if params.CallbackData != "" {
		parentDirectionID, err = strconv.Atoi(params.CallbackData)
		if err != nil {
			return errors.Wrapf(err, "botUC.HandleCreateDirectionCommand.Atoi(). params(%+v)", params)
		}
		getUserDirectionParams.ParentDirectionID = sql.NullInt64{Int64: int64(parentDirectionID), Valid: true}
	}
	getUserDirectionParams.ChatID = params.ChatID

	directions, err := u.pgRepo.GetUserDirection(ctx, getUserDirectionParams)
	if err != nil {
		return err
	}

	if len(directions) == 0 && params.CallbackData == "" { // executed when the user does not have one direction
		sendMessageParms := models.SendMessageParams{ChatID: params.ChatID, Text: "To add information, create at least one direction"}
		u.sendMessage(sendMessageParms)
		return nil
	} else if len(directions) == 0 && params.CallbackData != "" { // executed when the user has directions but has reached the lowest level
		setAwaitingStatusParams := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: utils.AwaitingQuestion}
		if err := u.rdb.SetAwaitingStatus(ctx, setAwaitingStatusParams); err != nil {
			return err
		}

		setDirectionForInfoParams := models.SetDirectionForInfoParams{
			ChatID: params.ChatID, DirectionID: int(getUserDirectionParams.ParentDirectionID.Int64)}
		if err := u.rdb.SetDirectionForInfo(ctx, setDirectionForInfoParams); err != nil {
			return err
		}
		sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: "enter your question"}
		u.sendMessage(sendMessageParams)
		return nil
	}

	setAwaitingStatusParams := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: utils.AwaitingAddInfoDirection}
	if err := u.rdb.SetAwaitingStatus(ctx, setAwaitingStatusParams); err != nil {
		return err
	}

	sendMessageParams := models.SendMessageParams{
		ChatID:         params.ChatID,
		Text:           "choose direction for add info",
		Keyboard:       u.createDirectionsKeyboard(directions),
		IsNeedToRemove: true}
	u.sendMessage(sendMessageParams)

	return nil
}

func (u *botUC) HandleAwaitingQuestion(ctx context.Context, params models.HandleAwaitingQuestionParams) error {
	getDirectionForInfoResult, err := u.rdb.GetDirectionForInfo(ctx, params.ChatID)
	if err != nil {
		return err
	}

	directionID, err := strconv.Atoi(getDirectionForInfoResult)
	if err != nil {
		return errors.Wrapf(err, "botUC.HandleAwaitingQuestion.Atoi(%s).", getDirectionForInfoResult)
	}

	saveQuestionParams := models.SaveQuestionParams{Question: params.Question, DirectionID: directionID}
	infoID, err := u.pgRepo.SaveQuestion(ctx, saveQuestionParams)
	if err != nil {
		return err
	}

	setInfoIDParams := models.SetInfoIDParams{ChatID: params.ChatID, InfoID: infoID}
	if err := u.rdb.SetInfoID(ctx, setInfoIDParams); err != nil {
		return err
	}

	setAwaitingStatusParams := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: utils.AwaitingAnswer}
	if err := u.rdb.SetAwaitingStatus(ctx, setAwaitingStatusParams); err != nil {
		return err
	}

	sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: "alright, enter your answer"}
	u.sendMessage(sendMessageParams)

	return nil
}

func (u *botUC) HandleAwaitingAnswer(ctx context.Context, params models.HandleAwaitingAnswerParams) error {
	getInfoIDResult, err := u.rdb.GetInfoID(ctx, params.ChatID)
	if err != nil {
		return err
	}

	infoID, err := strconv.Atoi(getInfoIDResult)
	if err != nil {
		return errors.Wrapf(err, "botUC.HandleAwaitingAnswer.Atoi(%s)", getInfoIDResult)
	}

	saveAnswerParams := models.SaveAnswerParams{Answer: params.Answer, InfoID: infoID}
	if err := u.pgRepo.SaveAnswer(ctx, saveAnswerParams); err != nil {
		return err
	}

	sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: "your answer successfully stored"}
	u.sendMessage(sendMessageParams)

	if err := u.rdb.ResetAwaitingStatus(ctx, params.ChatID); err != nil {
		return err
	}

	if err := u.rdb.ResetInfoID(ctx, params.ChatID); err != nil {
		return err
	}

	return nil
}

func (u *botUC) HandlePrintQuestionsCommand(ctx context.Context, params models.HandlePrintQuestionsCommandParams) error {
	var err error
	var parentDirectionID int
	var getUserDirectionParams models.GetUserDirectionParams
	if params.CallbackData != "" {
		parentDirectionID, err = strconv.Atoi(params.CallbackData)
		if err != nil {
			return errors.Wrapf(err, "botUC.HandlePrintInfoCommand.Atoi(). params(%+v)", params)
		}
		getUserDirectionParams.ParentDirectionID = sql.NullInt64{Int64: int64(parentDirectionID), Valid: true}
	}
	getUserDirectionParams.ChatID = params.ChatID

	directions, err := u.pgRepo.GetUserDirection(ctx, getUserDirectionParams)
	if err != nil {
		return err
	}

	if len(directions) == 0 && params.CallbackData == "" {
		sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: "To print info create at least one direction and add info to it"}
		u.sendMessage(sendMessageParams)
		return nil
	} else if len(directions) == 0 && params.CallbackData != "" {
		questions, err := u.pgRepo.GetQuestionsByDirectionID(ctx, int(getUserDirectionParams.ParentDirectionID.Int64))
		if err != nil {
			return err
		}
		for _, v := range questions {
			sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: v.Text, Keyboard: u.createInfoKeyboard(v.ID)}
			u.sendMessage(sendMessageParams)
			time.Sleep(time.Millisecond * 500)
		}

		return nil
	}

	setAwaitingStatusParams := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: utils.AwaitingPrintQuestions}
	if err := u.rdb.SetAwaitingStatus(ctx, setAwaitingStatusParams); err != nil {
		return err
	}

	sendMessageParams := models.SendMessageParams{
		ChatID:         params.ChatID,
		Text:           "choose direction for print questions",
		Keyboard:       u.createDirectionsKeyboard(directions),
		IsNeedToRemove: true}
	u.sendMessage(sendMessageParams)

	return nil
}

func (u *botUC) sendMessage(params models.SendMessageParams) {
	msg := tgbotapi.NewMessage(params.ChatID, params.Text)
	if params.Keyboard.InlineKeyboard != nil {
		msg.ReplyMarkup = params.Keyboard
	}

	sendedMsg, err := u.BotAPI.Send(msg)
	if err != nil {
		u.log.Errorf(err.Error())
	}

	if params.IsNeedToRemove {
		if err := u.rdb.SetExpirationTimeForMessage(context.Background(), sendedMsg.MessageID, params.ChatID); err != nil {
			u.log.Errorf(err.Error())
		}
	}
}
