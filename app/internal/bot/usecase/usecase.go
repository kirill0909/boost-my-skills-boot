package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type BotUC struct {
	cfg    *config.Config
	pgRepo bot.PgRepository
	BotAPI *tgbotapi.BotAPI
}

func NewBotUC(
	cfg *config.Config,
	pgRepo bot.PgRepository,
	botAPI *tgbotapi.BotAPI,
) bot.Usecase {
	return &BotUC{
		cfg:    cfg,
		pgRepo: pgRepo,
		BotAPI: botAPI,
	}
}

func (u *BotUC) HandleStartCommand(ctx context.Context, params models.HandleStartCommandParams) error {
	splitedText := strings.Split(params.Text, " ")

	if len(splitedText) != 2 {
		return fmt.Errorf("TgBot.handleStartCommand. wrong len of splited text: %d != 2. params(%+v)", len(splitedText), params)
	}

	uuid := splitedText[1]
	result, err := u.pgRepo.CompareUUID(ctx, models.CompareUUIDParams{
		UUID:   uuid,
		ChatID: params.ChatID})
	if err != nil {
		return err
	}

	if !result {
		return fmt.Errorf("wrong uuid(%s) comparison for user(%d)", uuid, params.ChatID)
	}

	if err := u.pgRepo.SetStatusActive(ctx, params.ChatID); err != nil {
		return err
	}

	return nil
}
