package bot

import (
	"boost-my-skills-bot/internal/bot/models"
	"context"
)

type Usecase interface {
	// commands
	HandleStartCommand(context.Context, models.HandleStartCommandParams) error
	HandleCreateDirectionCommand(context.Context, models.HandleCreateDirectionCommandParams) error
	HandleAddInfoCommand(context.Context, models.HandleAddInfoCommandParams) error

	GetAwaitingStatus(context.Context, int64) (int, error)
	CreateDirection(context.Context, models.CreateDirectionParams) error
	SetParentDirection(context.Context, models.SetParentDirectionParams) error
	HandleAwaitingQuestion(context.Context, models.HandleAwaitingQuestionParams) error

	// workers
	SyncMainKeyboardWorker()
	ListenExpiredMessageWorker()
}
