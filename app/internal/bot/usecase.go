package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type Usecase interface {
	GetUUID(ctx context.Context, params models.GetUUID) (result string, err error)
	UserActivation(ctx context.Context, params models.UserActivation) (err error)
	GetRandomQuestion(ctx context.Context, params models.AksMeCallbackParams) (
		result models.SubdirectionsCallbackResult, err error)
	GetAnswer(ctx context.Context, questionID int) (result string, err error)
	SaveQuestion(ctx context.Context, params models.SaveQuestionParams) (result int, err error)
	SaveAnswer(ctx context.Context, params models.SaveAnswerParams) (err error)
	GetSubdirections(ctx context.Context, params models.GetSubdirectionsParams) (result []string, err error)
	GetSubSubdirections(ctx context.Context, params models.GetSubSubdirectionsParams) (result []string, err error)

	SetUpDirection(ctx context.Context, params models.SetUpDirection) (err error)
	HandleAddInfoCommand(ctx context.Context, params int64) (err error)
	HandleAddInfoSubdirectionCallbackData(ctx context.Context, params models.AddInfoParams) (err error)
	// HandleAddInfoSubSubdirectionCallbackData(ctx context.Context, params models.AddInfoSubdirectionParams) (err error)

	// worker
	SyncDirectionsInfo(ctx context.Context) (err error)
}
