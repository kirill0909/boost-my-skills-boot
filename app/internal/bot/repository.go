package bot

import (
	models "boost-my-skills-bot/internal/models/bot"
	"context"
)

type PgRepository interface {
	GetUUID(ctx context.Context) (result string, err error)
	IsAdmin(ctx context.Context, params models.GetUUID) (result bool, err error)
	UserActivation(ctx context.Context, params models.UserActivation) (err error)
	GetRandomQuestion(ctx context.Context, params models.AksMeCallbackParams) (
		result models.SubdirectionsCallbackResult, err error)
	GetAnswer(ctx context.Context, questionID int) (result string, err error)
	SaveQuestion(ctx context.Context, params models.SaveQuestionParams) (result int, err error)
	SaveAnswer(ctx context.Context, params models.SaveAnswerParams) (err error)
	GetSubdirections(ctx context.Context, params models.GetSubdirectionsParams) (result []string, err error)
	GetSubSubdirections(ctx context.Context, params models.GetSubSubdirectionsParams) (result []string, err error)

	GetDirectionIDByChatID(ctx context.Context, param int64) (result int, err error)
	SetUpDirection(ctx context.Context, params models.SetUpDirection) (err error)

	// worker
	GetDirectionsInfo(ctx context.Context) (result []models.DirectionInfo, err error)
	GetSubdirectionsInfo(ctx context.Context) (result []models.SubdirectionInfo, err error)
	GetSubSubdirectionsInfo(ctx context.Context) (result []models.SubSubdirectionInfo, err error)
}
