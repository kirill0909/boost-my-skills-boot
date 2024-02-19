package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

func (t *TgBot) handleStartCommand(ctx context.Context, params models.HandleStartCommandParams) error {
	if err := t.tgUC.HandleStartCommand(ctx, params); err != nil {
		return err
	}

	return nil
}

func (t *TgBot) handleCreateDirectionCommand(ctx context.Context, params models.HandleCreateDirectionCommandParams) error {
	if err := t.tgUC.HandleCreateDirectionCommand(ctx, params); err != nil {
		return err
	}

	return nil
}
