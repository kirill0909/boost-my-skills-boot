package usecase

import (
	"context"

	"boost-my-skills-bot/app/internal/bot/models"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (u *botUC) HandleStartCommand(ctx context.Context, params models.HandleStartCommandParams) error {
	splitedText := strings.Split(params.Text, " ")

	if len(splitedText) != 2 {
		return fmt.Errorf("botUC.handleStartCommand. wrong len of splited text: %d != 2. params(%+v)", len(splitedText), params)
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

func (u *botUC) HandleServiceCommand(ctx context.Context, chatID int64) error {
	userInfo, err := u.pgRepo.GetUserInfo(ctx, chatID)
	if err != nil {
		return err
	}

	if !userInfo.IsAdmin || !userInfo.IsActive {
		return nil
	}

	sendMessageParams := models.SendMessageParams{
		ChatID: userInfo.TgChatID, Text: "choose service action", InlineKeyboard: u.createServiceKeyboard()}
	u.sendMessage(sendMessageParams)

	return nil
}
