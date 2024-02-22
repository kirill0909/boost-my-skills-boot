package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *TgBot) handleAwaitingDirectionName(ctx context.Context, params models.HandleAwaitingDirectionNameParams) error {
	if err := t.tgUC.CreateDirection(ctx, models.CreateDirectionParams{ChatID: params.ChatID, DirectionName: params.DirectionName}); err != nil {
		return err
	}

	t.sendMessage(params.ChatID, "new direction has been successfully created")

	return nil
}
