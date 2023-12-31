package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type Usecase interface {
	// commands
	HandleStartCommand(context.Context, models.HandleStartCommandParams) error
}
