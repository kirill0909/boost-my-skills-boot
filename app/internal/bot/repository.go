package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type PgRepository interface {
	// commands
	CompareUUID(context.Context, models.CompareUUIDParams) (bool, error)
	SetStatusActive(context.Context, int64) error
}
