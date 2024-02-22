package bot

import (
	"boost-my-skills-bot/internal/bot/models"
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
