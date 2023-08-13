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
	SetUpFrontendDirection(ctx context.Context, chatID int64) (err error)
	GetRandomQuestion(ctx context.Context, params models.AskMeParams) (result models.AskMeResult, err error)
	GetAnswer(ctx context.Context, questionID int) (result string, err error)
}
