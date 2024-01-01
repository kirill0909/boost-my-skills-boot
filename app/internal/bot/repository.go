package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type PgRepository interface {
	SetStatusActive(context.Context, models.SetStatusActiveParams) error
	GetMainButtons(context.Context) ([]models.GetMainButtonsResult, error)
	GetActiveUsers(context.Context) ([]models.GetActiveUsersResult, error)
	GetUpdatedButtons(context.Context, int64) ([]models.GetUpdatedButtonsResult, error)
}
