package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

func (t *TgBot) handleStartCommand(params models.HandleStartCommandParams) error {
	ctx := context.Background()
	if err := t.tgUC.HandleStartCommand(ctx, params); err != nil {
		return err
	}

	return nil
}
