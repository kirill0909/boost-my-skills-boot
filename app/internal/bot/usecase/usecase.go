package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"boost-my-skills-bot/pkg/logger"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"strings"
)

type BotUC struct {
	cfg    *config.Config
	pgRepo bot.PgRepository
	BotAPI *tgbotapi.BotAPI
	log    *logger.Logger
}

func NewBotUC(
	cfg *config.Config,
	pgRepo bot.PgRepository,
	botAPI *tgbotapi.BotAPI,
	log *logger.Logger,
) bot.Usecase {
	return &BotUC{
		cfg:    cfg,
		pgRepo: pgRepo,
		BotAPI: botAPI,
		log:    log,
	}
}

func (u *BotUC) HandleStartCommand(ctx context.Context, params models.HandleStartCommandParams) error {
	splitedText := strings.Split(params.Text, " ")

	if len(splitedText) != 2 {
		return fmt.Errorf("TgBot.handleStartCommand. wrong len of splited text: %d != 2. params(%+v)", len(splitedText), params)
	}

	uuid := splitedText[1]
	if err := u.pgRepo.SetStatusActive(ctx, models.SetStatusActiveParams{
		TgName: params.TgName,
		ChatID: params.ChatID,
		UUID:   uuid}); err != nil {
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
