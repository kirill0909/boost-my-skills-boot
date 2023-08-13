package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type PgRepository interface {
	GetUUID(ctx context.Context) (result string, err error)
	IsAdmin(ctx context.Context, params models.GetUUID) (result bool, err error)
	UserActivation(ctx context.Context, params models.UserActivation) (err error)
	SetUpBackendDirection(ctx context.Context, chatID int64) (err error)
}
