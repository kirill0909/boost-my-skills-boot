package usecase

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

func (u *BotUC) HandleAskMeCommand(ctx context.Context, params models.AskMeParams) (err error) {
	directionID, err := u.pgRepo.GetDirectionIDByChatID(ctx, params.ChatID)
	if err != nil {
		return
	}

	subdirections := u.stateDirections.GetSubdirectionsByDirectionID(directionID)

	msg := tgbotapi.NewMessage(params.ChatID, subdirectionAskMeMessage)
	msg.ReplyMarkup = u.createSubdirectionsKeyboardAskMe(subdirections)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return errors.Wrap(err, "BotUC.HandleAskMeCommand.Send")
	}

	return
}
