package usecase

import (
	"boost-my-skills-bot/config"
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotUC struct {
	cfg    *config.Config
	pgRepo bot.PgRepository
	BotAPI *tgbotapi.BotAPI
}

func NewBotUC(cfg *config.Config, pgRepo bot.PgRepository, botAPI *tgbotapi.BotAPI) bot.Usecase {
	return &BotUC{cfg: cfg, pgRepo: pgRepo, BotAPI: botAPI}
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

func (u *BotUC) SetUpBackendDirection(ctx context.Context, chatID int64) (err error) {
	return u.pgRepo.SetUpBackendDirection(ctx, chatID)
}

func (u *BotUC) SetUpFrontendDirection(ctx context.Context, chatID int64) (err error) {
	return u.pgRepo.SetUpFrontendDirection(ctx, chatID)
}

func (u *BotUC) GetRandomQuestion(ctx context.Context, params models.AskMeParams) (result models.AskMeResult, err error) {
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
