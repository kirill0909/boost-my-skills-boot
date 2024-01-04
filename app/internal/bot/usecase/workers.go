package usecase

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"time"
)

func (u *BotUC) SyncMainKeyboardWorker() {
	u.log.Infof("SyncMainKeyboardWorker is started")
	ctx := context.Background()
	ticker := time.NewTicker(time.Second * 10)

	for ; true; <-ticker.C {
		getUpdatedButtonsResult, err := u.pgRepo.GetUpdatedButtons(ctx, u.lastKeyboardChecking)
		if err != nil {
			u.log.Errorf(err.Error())
			continue
		}
		u.lastKeyboardChecking = time.Now().Unix()

		if len(getUpdatedButtonsResult) == 0 {
			continue
		}

		users, err := u.pgRepo.GetActiveUsers(ctx)
		if err != nil {
			u.log.Errorf(err.Error())
			continue
		}

		for _, updatedButton := range getUpdatedButtonsResult {
			switch updatedButton.OnlyForAdmin {
			case true:
				if err := u.handleOnlyForAdminCaes(users); err != nil {
					u.log.Errorf(err.Error())
				}
			case false:
				if err := u.handleOnlyForUserCase(users); err != nil {
					u.log.Errorf(err.Error())
				}
			}
		}
	}
}

func (u *BotUC) handleOnlyForAdminCaes(users []models.GetActiveUsersResult) error {
	for _, user := range users {
		if user.IsAdmin {
			msg := tgbotapi.NewMessage(user.ChatID, "main kyboard was updated")
			msg.ReplyMarkup = u.createMainMenuKeyboard(user.IsAdmin)
			if _, err := u.BotAPI.Send(msg); err != nil {
				return errors.Wrap(err, "BotUC.handleOnlyForAdminCaes")
			}
		}
	}

	return nil
}

func (u *BotUC) handleOnlyForUserCase(users []models.GetActiveUsersResult) error {
	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ChatID, "main kyboard was updated")
		msg.ReplyMarkup = u.createMainMenuKeyboard(user.IsAdmin)
		if _, err := u.BotAPI.Send(msg); err != nil {
			return errors.Wrap(err, "BotUC.handleOnlyForUserCase")
		}
	}

	return nil
}
