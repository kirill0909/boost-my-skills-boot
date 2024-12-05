package usecase

import (
	"boost-my-skills-bot/app/config"
	"boost-my-skills-bot/app/internal/bot"
	"boost-my-skills-bot/app/internal/bot/models"
	"boost-my-skills-bot/app/pkg/utils"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type botUC struct {
	cfg                  *config.Config
	pgRepo               bot.PgRepository
	rdb                  bot.RedisRepository
	redisPubSub          *redis.PubSub
	BotAPI               *tgbotapi.BotAPI
	log                  *slog.Logger
	lastKeyboardChecking int64
}

func NewBotUC(
	cfg *config.Config,
	pgRepo bot.PgRepository,
	rdb bot.RedisRepository,
	redisPubSub *redis.PubSub,
	botAPI *tgbotapi.BotAPI,
	log *slog.Logger,
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

func (u *botUC) HandleCreateDirection(ctx context.Context, params models.HandleCreateDirectionParams) error {
	setAwaitingStatusParams := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: utils.AwaitingDirectionNameStatus}
	if err := u.rdb.SetAwaitingStatus(ctx, setAwaitingStatusParams); err != nil {
		return err
	}

	getUserDirectionParms := models.GetUserDirectionParams{
		ChatID:            params.ChatID,
		ParentDirectionID: params.ParentDirectionID,
	}

	directions, err := u.pgRepo.GetUserDirection(ctx, getUserDirectionParms)
	if err != nil {
		return err
	}

	if len(directions) == 0 {
		sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: "enter name of your direction"}
		u.sendMessage(sendMessageParams)
	} else {
		sendMessageParams := models.SendMessageParams{
			ChatID:         params.ChatID,
			Text:           "choose parent direciton",
			InlineKeyboard: u.createDirectionsKeyboard(directions, utils.AwaitingParentDirectionCallbackType),
			IsNeedToRemove: true}
		u.sendMessage(sendMessageParams)
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
		Text:   fmt.Sprintf("successfully created \"%s\" direction", utils.FormatBadCharacters(direction))}
	u.sendMessage(sendMessageParams)

	return nil
}

func (u *botUC) SetParentDirection(ctx context.Context, params models.SetParentDirectionParams) error {
	if err := u.rdb.SetParentDirection(ctx, params); err != nil {
		return err
	}

	return nil
}

func (u *botUC) HandleAddInfo(ctx context.Context, params models.HandleAddInfoParams) error {
	getUserDirectionParams := models.GetUserDirectionParams{ChatID: params.ChatID, ParentDirectionID: params.ParentDirectionID}

	directions, err := u.pgRepo.GetUserDirection(ctx, getUserDirectionParams)
	if err != nil {
		return err
	}

	// TODO: move every case to separete function
	switch {
	// executed when the user does not have one direction
	case len(directions) == 0 && params.ParentDirectionID == 0:
		sendMessageParms := models.SendMessageParams{ChatID: params.ChatID, Text: "To add information, create at least one direction"}
		u.sendMessage(sendMessageParms)
	// executed when the user has directions and has parrent directions OR has directions but has NOT parent directions
	case (len(directions) != 0 && params.ParentDirectionID == 0) || (len(directions) != 0 && params.ParentDirectionID != 0):
		sendMessageParams := models.SendMessageParams{
			ChatID:         params.ChatID,
			Text:           "choose direction for add info",
			InlineKeyboard: u.createDirectionsKeyboard(directions, utils.AwaitingAddInfoDirectionCallbackType),
			IsNeedToRemove: true}
		u.sendMessage(sendMessageParams)
	// executed when user reached low level
	default:
		setAwaitingStatusParams := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: utils.AwaitingQuestionStatus}
		if err := u.rdb.SetAwaitingStatus(ctx, setAwaitingStatusParams); err != nil {
			return err
		}

		setDirectionForInfoParams := models.SetDirectionForInfoParams{
			ChatID: params.ChatID, DirectionID: getUserDirectionParams.ParentDirectionID}
		if err := u.rdb.SetDirectionForInfo(ctx, setDirectionForInfoParams); err != nil {
			return err
		}
		sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: "enter your question"}
		u.sendMessage(sendMessageParams)
	}

	return nil
}

func (u *botUC) HandleGetInviteLink(ctx context.Context, chatID int64) error {
	uuid, err := u.pgRepo.CreateInActiveUser(ctx)
	if err != nil {
		return err
	}

	link := fmt.Sprintf(u.cfg.TgBot.Prefix, uuid)

	sendMssageParams := models.SendMessageParams{ChatID: chatID, Text: utils.FormatBadCharacters(link)}
	u.sendMessage(sendMssageParams)

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

	saveQuestionParams := models.SaveQuestionParams{Question: utils.FormatSnipets(params.Question), DirectionID: directionID}
	infoID, err := u.pgRepo.SaveQuestion(ctx, saveQuestionParams)
	if err != nil {
		return err
	}

	setInfoIDParams := models.SetInfoIDParams{ChatID: params.ChatID, InfoID: infoID}
	if err := u.rdb.SetInfoID(ctx, setInfoIDParams); err != nil {
		return err
	}

	setAwaitingStatusParams := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: utils.AwaitingAnswerStatus}
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

	saveAnswerParams := models.SaveAnswerParams{Answer: utils.FormatSnipets(params.Answer), InfoID: infoID}
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

func (u *botUC) HandlePrintQuestions(ctx context.Context, params models.HandlePrintQuestionsParams) error {
	getUserDirectionParams := models.GetUserDirectionParams{ChatID: params.ChatID, ParentDirectionID: params.ParentDirectionID}
	directions, err := u.pgRepo.GetUserDirection(ctx, getUserDirectionParams)
	if err != nil {
		return err
	}

	// TODO: move every case to separete function
	switch {
	// executed when user has not one direction
	case len(directions) == 0 && params.ParentDirectionID == 0:
		sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: "Create at least one direction and add info to it"}
		u.sendMessage(sendMessageParams)
		// executed when user has directions but not reached low level
	case (len(directions) != 0 && params.ParentDirectionID == 0) || (len(directions) != 0 && params.ParentDirectionID != 0):
		sendMessageParams := models.SendMessageParams{
			ChatID:         params.ChatID,
			Text:           "choose direction for print questions",
			InlineKeyboard: u.createDirectionsKeyboard(directions, utils.AwaitingPrintQuestionsCallbackType),
			IsNeedToRemove: true}
		u.sendMessage(sendMessageParams)
		// executed when user has directions and reached low level
	default:
		questions, err := u.pgRepo.GetQuestionsByDirectionID(ctx, getUserDirectionParams.ParentDirectionID)
		if err != nil {
			return err
		}

		if len(questions) == 0 {
			sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: "To print info create at least one direction and add info to it"}
			u.sendMessage(sendMessageParams)
			return nil
		}

		for _, v := range questions {
			sendMessageParams := models.SendMessageParams{
				ChatID: params.ChatID,
				Text:   utils.FormatBadCharacters(v.Text), InlineKeyboard: u.createInfoKeyboard(v.ID, utils.AwaitingPrintAnswerCallbackType)}
			u.sendMessage(sendMessageParams)
			time.Sleep(time.Millisecond * 50)
		}

	}

	return nil
}

func (u *botUC) HandleAwaitingPrintAnswer(ctx context.Context, params models.HandleAwaitingPrintAnswerParams) error {
	answer, err := u.pgRepo.GetAnswerByInfoID(ctx, params.InfoID)
	if err != nil {
		return err
	}

	sendMessageParams := models.SendMessageParams{ChatID: params.ChatID, Text: utils.FormatBadCharacters(answer)}
	u.sendMessage(sendMessageParams)
	u.hideKeyboard(params.ChatID, params.MessageID)

	return nil
}

func (u *botUC) HandleAwaitingNewMainButtonName(ctx context.Context, params models.HandleAwaitingNewMainButtonNameParams) error {
	if len(params.ButtonName) == 0 {
		err := fmt.Errorf("name of main button should be more then 0")
		return errors.Wrapf(err, "botUC.HandleAwaitingNewMainButtonName.len(). params(%+v)", params)
	}

	buttonName := fmt.Sprintf("/%s", params.ButtonName)
	addNewButtonToMainKeyboardParams := models.AddNewButtonToMainKeyboardParams{ButtonName: buttonName, OnlyForAdmin: params.OnlyForAdmin}
	if err := u.pgRepo.AddNewButtonToMainKeyboard(ctx, addNewButtonToMainKeyboardParams); err != nil {
		return err
	}

	if err := u.rdb.ResetAwaitingStatus(ctx, params.ChatID); err != nil {
		return err
	}

	activeUsers, err := u.pgRepo.GetActiveUsers(ctx)
	if err != nil {
		return err
	}

	for _, user := range activeUsers {
		sendMessageParams := models.SendMessageParams{
			ChatID: user.ChatID, Text: "main keyboard was updated", ReplyKeyboard: u.createMainMenuKeyboard(user.IsAdmin)}
		u.sendMessage(sendMessageParams)
	}

	return nil
}

func (u *botUC) sendMessage(params models.SendMessageParams) {
	msg := tgbotapi.NewMessage(params.ChatID, params.Text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	if params.InlineKeyboard.InlineKeyboard != nil {
		msg.ReplyMarkup = params.InlineKeyboard
	}

	if params.ReplyKeyboard.Keyboard != nil {
		msg.ReplyMarkup = params.ReplyKeyboard
	}

	sendedMsg, err := u.BotAPI.Send(msg)
	if err != nil {
		u.log.Error("botUC.sendMessage.Send()", "error", err.Error())
		return
	}

	if params.IsNeedToRemove {
		if err := u.rdb.SetExpirationTimeForMessage(context.Background(), sendedMsg.MessageID, params.ChatID); err != nil {
			u.log.Error("botUC.sendMessage.SetExpirationTimeForMessage()", "error", err.Error())
		}
	}
}

func (u *botUC) editMessage(params models.EditMessageParams) {
	if params.Keyboard.InlineKeyboard == nil {
		params.Keyboard = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{})
	}

	msg := tgbotapi.NewEditMessageTextAndMarkup(params.ChatID, params.MessageID, params.Text, params.Keyboard)
	_, err := u.BotAPI.Send(msg)
	if err != nil {
		u.log.Error("botUC.editMessage.Send()", "error", err.Error())
		return
	}
}
