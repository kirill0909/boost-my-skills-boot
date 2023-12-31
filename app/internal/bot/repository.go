package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type PgRepository interface {
	// commands
	SetStatusActive(context.Context, models.SetStatusActiveParams) error
}
