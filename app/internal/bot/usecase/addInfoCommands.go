package usecase

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	models "boost-my-skills-bot/internal/models/bot"
)

func (u *BotUC) HandleAddInfoCommand(ctx context.Context, chatID int64) (err error) {

	u.stateUsers[chatID] = models.AddInfoParams{State: awaitingSubdirection}
	directionID, err := u.pgRepo.GetDirectionIDByChatID(ctx, chatID)
	if err != nil {
		return
	}

	subdirections := u.stateDirections.GetSubdirectionsByDirectionID(directionID)
	if len(subdirections) == 0 {
		err = fmt.Errorf("subdirections not found")
		return errors.Wrap(err, "BotUC.HandleAddInfoCommand.len(subdirections)")
	}

	msg := tgbotapi.NewMessage(chatID, "")
	if len(subdirections) == 0 {
		msg.Text = noOneSubdirectionsFoundMessage
		if _, err = u.BotAPI.Send(msg); err != nil {
			return
		}
		u.stateUsers[chatID] = models.AddInfoParams{State: idle}

		return
	}

	msg.Text = subdirectionAddInfoMessage
	msg.ReplyMarkup = u.createSubdirectionsKeyboardAddInfo(subdirections)
	if _, err = u.BotAPI.Send(msg); err != nil {
		return
	}

	return
}
