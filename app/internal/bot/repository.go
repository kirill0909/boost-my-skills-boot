package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type PgRepository interface {
	GetUUID(ctx context.Context, params models.GetUUID) (result string, err error)
	IsAdmin(ctx context.Context, params models.GetUUID) (result bool, err error)
}
