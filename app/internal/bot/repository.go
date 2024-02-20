package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type PgRepository interface {
	SetUserActive(context.Context, models.SetUserActiveParams) error
	GetMainButtons(context.Context) ([]models.GetMainButtonsResult, error)
	GetActiveUsers(context.Context) ([]models.GetActiveUsersResult, error)
	GetUpdatedButtons(context.Context, int64) ([]models.GetUpdatedButtonsResult, error)
	GetUserDirection(context.Context, int64) ([]models.GetUserDirectionsResult, error)
	CreateDirection(context.Context, models.CreateDirectionParams) error
}

type RedisRepository interface {
	SetAwaitingStatus(context.Context, models.SetAwaitingStatusParams) error
	ResetAwaitingStatus(context.Context, int64) error
	GetAwaitingStatus(context.Context, int64) (string, error)
}
