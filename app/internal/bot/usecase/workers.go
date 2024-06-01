package usecase

import (
	"boost-my-skills-bot/app/internal/bot/models"
	"boost-my-skills-bot/app/pkg/utils"
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (u *botUC) SyncMainKeyboardWorker() {
	u.log.Info("botUC.SyncMainKeyboardWorker() is started", "info", "SyncMainKeyboardWorker is started")
	ctx := context.Background()
	ticker := time.NewTicker(time.Second * 10)

	for ; true; <-ticker.C {
		getUpdatedButtonsResult, err := u.pgRepo.GetUpdatedButtons(ctx, u.lastKeyboardChecking)
		if err != nil {
			u.log.Error("botUC.SyncMainKeyboardWorker.GetUpdatedButtons()", "error", err.Error())
			continue
		}
		u.lastKeyboardChecking = time.Now().Unix()

		if len(getUpdatedButtonsResult) == 0 {
			continue
		}

		users, err := u.pgRepo.GetActiveUsers(ctx)
		if err != nil {
			u.log.Error("botUC.SyncMainKeyboardWorker.GetActiveUsers()", "error", err.Error())
			continue
		}

		for _, updatedButton := range getUpdatedButtonsResult {
			switch updatedButton.OnlyForAdmin {
			case true:
				if err := u.handleOnlyForAdminCaes(users); err != nil {
					u.log.Error("botUC.SyncMainKeyboardWorker.handleOnlyForAdminCaes()", "error", err.Error())
				}
			case false:
				if err := u.handleOnlyForUserCase(users); err != nil {
					u.log.Error("botUC.SyncMainKeyboardWorker.handleOnlyForUserCase()", "error", err.Error())
				}
			}
		}
	}
}

func (u *botUC) handleOnlyForAdminCaes(users []models.GetActiveUsersResult) error {
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

func (u *botUC) handleOnlyForUserCase(users []models.GetActiveUsersResult) error {
	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ChatID, "main kyboard was updated")
		msg.ReplyMarkup = u.createMainMenuKeyboard(user.IsAdmin)
		if _, err := u.BotAPI.Send(msg); err != nil {
			return errors.Wrap(err, "BotUC.handleOnlyForUserCase")
		}
	}

	return nil
}

func (u *botUC) ListenExpiredMessageWorker() {
	u.log.Info("botUC.ListenExpirationMessageWorker()", "info", "ListenExpirationMessageWorker is started")
	for msg := range u.redisPubSub.Channel() {
		if strings.HasPrefix(msg.Payload, utils.ExpirationTimeMessagePrefix) {
			if err := u.removeMessage(msg.Payload); err != nil {
				u.log.Info("botUC.ListenExpirationMessageWorker()", "error", err.Error())
			}
		}
	}
}

func (u *botUC) removeMessage(payload string) error {
	messageIDRegex := regexp.MustCompile(utils.MessageIDLayout)
	submatchResult := messageIDRegex.FindStringSubmatch(payload)
	messageID, err := strconv.Atoi(submatchResult[1])
	if err != nil {
		return errors.Wrapf(err, "botUC.removeMessage.Atoi(). payload: %s", payload)
	}

	chatIDRegex := regexp.MustCompile(utils.ChatIDLayout)
	submatchResult = chatIDRegex.FindStringSubmatch(payload)
	chatID, err := strconv.Atoi(submatchResult[1])
	if err != nil {
		return errors.Wrapf(err, "botUC.removeMessage.Atoi(). payload: %s", payload)
	}

	deleteMsg := tgbotapi.NewDeleteMessage(int64(chatID), messageID)
	u.BotAPI.Send(deleteMsg)

	return nil
}
