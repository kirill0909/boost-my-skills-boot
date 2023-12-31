package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"fmt"
)

func (t *TgBot) handleStartCommand(params models.HandleStartCommandParams) error {
	return fmt.Errorf("This is my test error with args: %d", 3)
}
