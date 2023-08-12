package usecase

import (
	"aifory-pay-admin-bot/config"
	"aifory-pay-admin-bot/internal/bot"
	"context"
	"fmt"
	"log"
	"strings"

	pb "gitlab.axarea.ru/main/aiforypay/package/admin-bot-proto"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotUsecase struct {
	botAPI *tgbotapi.BotAPI
	tgRepo bot.PgRepository
	cfg    *config.Config
}

func NewBotUsecase(
	repo bot.PgRepository,
	cfg *config.Config,
	botAPI *tgbotapi.BotAPI,
) bot.Usecase {
	return &BotUsecase{
		tgRepo: repo,
		cfg:    cfg,
		botAPI: botAPI,
	}
}

func (u *BotUsecase) NotifyGeneralAppealPercent(
	ctx context.Context, req *pb.NotifyGeneralAppealPercentRequest) (err error) {
	tagMe := strings.Join(u.cfg.AiforyPayAdminLog.TagMe, " ")
	info := fmt.Sprintf(generalAppealPercent, tagMe, req.GetPercent())
	msg := tgbotapi.NewMessage(u.cfg.AiforyPayAdminLog.ChatID, info)
	msg.DisableNotification = req.GetDisableNotification()

	if _, err = u.botAPI.Send(msg); err != nil {
		log.Println(err)
		return
	}

	return nil
}

func (u *BotUsecase) NotifyTraderAppealPercent(
	ctx context.Context, req *pb.NotifyTraderAppealPercentRequest) (err error) {
	tagMe := strings.Join(u.cfg.AiforyPayAdminLog.TagMe, " ")
	info := fmt.Sprintf(traderAppealPercent, tagMe, req.GetTraderID(), req.GetPercent())
	msg := tgbotapi.NewMessage(u.cfg.AiforyPayAdminLog.ChatID, info)

	if _, err = u.botAPI.Send(msg); err != nil {
		log.Println(err)
		return
	}

	return nil
}
