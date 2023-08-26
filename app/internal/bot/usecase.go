package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type Usecase interface {
	GetUUID(ctx context.Context, params models.GetUUID) (result string, err error)
	UserActivation(ctx context.Context, params models.UserActivation) (err error)
	SetUpBackendDirection(ctx context.Context, chatID int64) (err error)
	SetUpFrontendDirection(ctx context.Context, chatID int64) (err error)
	GetRandomQuestion(ctx context.Context, params models.AskMeParams) (result models.AskMeResult, err error)
	GetAnswer(ctx context.Context, questionID int) (result string, err error)
	SaveQuestion(ctx context.Context, params models.SaveQuestionParams) (result int, err error)
	SaveAnswer(ctx context.Context, params models.SaveAnswerParams) (err error)
	GetSubdirections(ctx context.Context, params models.GetSubdirectionsParams) (result []string, err error)
}
