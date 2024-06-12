package usecase

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (u *botUC) HandleServiceCommand(ctx context.Context, chatID int64) error {
	userInfo, err := u.pgRepo.GetUserInfo(ctx, chatID)
	if err != nil {
		return err
	}

	if !userInfo.IsAdmin || !userInfo.IsActive {
		return nil
	}

	msg := tgbotapi.NewMessage(userInfo.TgChatID, "choose service action")
	msg.ReplyMarkup = u.createServiceKeyboard()
	if _, err := u.BotAPI.Send(msg); err != nil {
		return errors.Wrapf(err, "BotUC.HandleServiceCommand.Send")
	}

	return nil
}
