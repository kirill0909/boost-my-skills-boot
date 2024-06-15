package bot

import (
	"boost-my-skills-bot/app/internal/bot/models"
	"context"
)

type Usecase interface {
	// commands
	HandleStartCommand(context.Context, models.HandleStartCommandParams) error
	HandleCreateDirection(context.Context, models.HandleCreateDirectionParams) error
	HandleAddInfo(context.Context, models.HandleAddInfoParams) error
	HandlePrintQuestions(context.Context, models.HandlePrintQuestionsParams) error
	HandleGetInviteLink(context.Context, int64) error
	// new commands
	HandleServiceCommand(context.Context, int64) error

	GetAwaitingStatus(context.Context, int64) (int, error)
	CreateDirection(context.Context, models.CreateDirectionParams) error
	SetParentDirection(context.Context, models.SetParentDirectionParams) error
	HandleAwaitingQuestion(context.Context, models.HandleAwaitingQuestionParams) error
	HandleAwaitingAnswer(context.Context, models.HandleAwaitingAnswerParams) error
	HandleAwaitingPrintAnswer(context.Context, models.HandleAwaitingPrintAnswerParams) error
	// new callbacks
	HandleActionsWithMainKeyboardCallback(context.Context, models.HandleActionsWithMainKeyboardCallback) error

	// workers
	SyncMainKeyboardWorker()
	ListenExpiredMessageWorker()
}
