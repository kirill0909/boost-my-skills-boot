package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
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
	BotAPI               *tgbotapi.BotAPI
	log                  *logger.Logger
	lastKeyboardChecking int64
}

func NewBotUC(
	cfg *config.Config,
	pgRepo bot.PgRepository,
	rdb bot.RedisRepository,
	botAPI *tgbotapi.BotAPI,
	log *logger.Logger,
) bot.Usecase {
	return &botUC{
		cfg:                  cfg,
		pgRepo:               pgRepo,
		rdb:                  rdb,
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
	getUserDirectionParms := models.GetUserDirectionParams{ChatID: params.ChatID}
	directions, err := u.pgRepo.GetUserDirection(ctx, getUserDirectionParms)
	if err != nil {
		return err
	}

	var statusID int
	if len(directions) == 0 {
		statusID = utils.AwaitingDirectionName
		u.sendMessage(params.ChatID, "enter name of your FIRST direction")
	} else {
		statusID = utils.AwaitingParentDireciton
		u.sendMessage(params.ChatID, "choose parent direciton", u.createDirectionsKeyboard(directions))
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

	direction, err := u.pgRepo.CreateDirection(ctx, params)
	if err != nil {
		return err
	}

	if err := u.rdb.ResetAwaitingStatus(ctx, params.ChatID); err != nil {
		return err
	}

	text := fmt.Sprintf("successfully created \"%s\" direction", direction)
	u.sendMessage(params.ChatID, text)

	return nil
}

func (u *botUC) sendMessage(chatID int64, text string, keyboard ...tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	if len(keyboard) > 0 {
		msg.ReplyMarkup = keyboard[0]
	}
	if _, err := u.BotAPI.Send(msg); err != nil {
		u.log.Errorf(err.Error())
	}
}
