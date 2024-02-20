package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kirill0909/logger"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type BotUC struct {
	cfg                  *config.Config
	pgRepo               bot.PgRepository
	rdb                  *redis.Client
	BotAPI               *tgbotapi.BotAPI
	log                  *logger.Logger
	lastKeyboardChecking int64
}

func NewBotUC(
	cfg *config.Config,
	pgRepo bot.PgRepository,
	rdb *redis.Client,
	botAPI *tgbotapi.BotAPI,
	log *logger.Logger,
) bot.Usecase {
	return &BotUC{
		cfg:                  cfg,
		pgRepo:               pgRepo,
		rdb:                  rdb,
		BotAPI:               botAPI,
		log:                  log,
		lastKeyboardChecking: time.Now().Unix(),
	}
}

func (u *BotUC) HandleStartCommand(ctx context.Context, params models.HandleStartCommandParams) error {
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

func (u *BotUC) HandleCreateDirectionCommand(ctx context.Context, params models.HandleCreateDirectionCommandParams) error {
	directions, err := u.pgRepo.GetUserDirection(ctx, params.ChatID)
	if err != nil {
		return err
	}

	switch len(directions) {
	case 0:
		// set status awaiting direction name
		// create first direction
		u.sendMessage(params.ChatID, "enter name of your FIRST direction")
	default:
		// set status awaiting parent dirction
		// create another direction
		u.sendMessage(params.ChatID, "enter name of your  direction")
	}

	return nil
}

func (u *BotUC) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := u.BotAPI.Send(msg); err != nil {
		u.log.Errorf(err.Error())
	}
}
