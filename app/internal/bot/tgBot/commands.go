package tgbot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"fmt"
)

func (t *TgBot) handleStartCommand(params models.HandleStartCommandParams) error {
	t.log.Infof("This is my log info with args: (%s)", "my args")
	return fmt.Errorf("This is my test error with args: %d", 3)
}
