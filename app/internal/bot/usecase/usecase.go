package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"boost-my-skills-bot/pkg/utils"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
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
	directions, err := u.pgRepo.GetUserDirection(ctx, params.ChatID)
	if err != nil {
		return err
	}

	switch len(directions) {
	case 0:
		// set status awaiting direction name
		params := models.SetAwaitingStatusParams{ChatID: params.ChatID, StatusID: utils.AwaitingDirectionName}
		if err := u.rdb.SetAwaitingStatus(ctx, params); err != nil {
			return err
		}
		// create first direction
		u.sendMessage(params.ChatID, "enter name of your FIRST direction")
	default:
		// set status awaiting parent dirction
		// create another direction
		u.sendMessage(params.ChatID, "enter name of your  direction")
	}

	return nil
}

func (u *botUC) GetAwaitingStatus(ctx context.Context, chatID int64) (int, error) {
	value, err := u.rdb.GetAwaitingStatus(ctx, chatID)
	if err != nil {
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

	return u.pgRepo.CreateDirection(ctx, params)
}

func (u *botUC) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := u.BotAPI.Send(msg); err != nil {
		u.log.Errorf(err.Error())
	}
}
