package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type Usecase interface {
	GetUUID(ctx context.Context, params models.GetUUID) (result string, err error)
}
