package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"log"
)

func (t *TgBot) handleStartCommand(params models.HandleStartCommandParams) error {
	ctx := context.Background()
	if err := t.tgUC.HandleStartCommand(ctx, params); err != nil {
		return err
	}

	return nil
}

func (t *TgBot) handleCreateDirectionCommand() error {
	log.Println("hello from create direction case")
	return nil
}
