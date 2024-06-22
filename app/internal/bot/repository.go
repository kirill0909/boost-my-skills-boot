package bot

import (
	"boost-my-skills-bot/app/internal/bot/models"
	"context"
)

type PgRepository interface {
	SetUserActive(context.Context, models.SetUserActiveParams) error
	GetMainButtons(context.Context) ([]models.GetMainButtonsResult, error) // TODO: Remove
	GetActiveUsers(context.Context) ([]models.GetActiveUsersResult, error)
	GetUpdatedButtons(context.Context, int64) ([]models.GetUpdatedButtonsResult, error)
	GetUserDirection(context.Context, models.GetUserDirectionParams) ([]models.UserDirection, error)
	CreateDirection(context.Context, models.CreateDirectionParams) (string, error)
	SaveQuestion(context.Context, models.SaveQuestionParams) (int, error)
	SaveAnswer(context.Context, models.SaveAnswerParams) error
	GetQuestionsByDirectionID(context.Context, int) ([]models.Question, error)
	GetAnswerByInfoID(context.Context, int) (string, error)
	CreateInActiveUser(context.Context) (string, error)

	AddNewButtonToMainKeyboard(context.Context, models.AddNewButtonToMainKeyboardParams) error
	GetMainKeyboard(context.Context) ([]models.GetMainKeyboardResult, error)
	GetUserInfo(context.Context, int64) (models.UserInfo, error)
}

type RedisRepository interface {
	SetAwaitingStatus(context.Context, models.SetAwaitingStatusParams) error
	ResetAwaitingStatus(context.Context, int64) error
	ResetParentDirection(context.Context, int64) error
	GetAwaitingStatus(context.Context, int64) (string, error)
	SetParentDirection(context.Context, models.SetParentDirectionParams) error
	GetParentDirection(context.Context, int64) (string, error)
	SetExpirationTimeForMessage(context.Context, int, int64) error
	SetDirectionForInfo(context.Context, models.SetDirectionForInfoParams) error
	GetDirectionForInfo(context.Context, int64) (string, error)
	SetInfoID(context.Context, models.SetInfoIDParams) error
	GetInfoID(context.Context, int64) (string, error)
	ResetInfoID(context.Context, int64) error
}
