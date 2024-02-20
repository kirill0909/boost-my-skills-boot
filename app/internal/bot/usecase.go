package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type Usecase interface {
	// commands
	HandleStartCommand(context.Context, models.HandleStartCommandParams) error
	HandleCreateDirectionCommand(context.Context, models.HandleCreateDirectionCommandParams) error

	GetAwaitingStatus(context.Context, int64) (int, error)
	CreateDirection(context.Context, models.CreateDirectionParams) error
	// workers
	SyncMainKeyboardWorker()
}
