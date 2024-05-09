package models

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackInfo struct {
	CallbackType int
	DirectionID  int
	InfoID       int
}

type HandleAwaitingPrintAnswerParams struct {
	ChatID    int64
	MessageID int
	InfoID    int
}

type SaveQuestionParams struct {
	Question    string
	DirectionID int
}

type SaveAnswerParams struct {
	Answer string
	InfoID int
}

type GetActiveUsersResult struct {
	ChatID  int64 `db:"tg_chat_id"`
	IsAdmin bool  `db:"is_admin"`
}

type GetUpdatedButtonsResult struct {
	Name         string `db:"name"`
	OnlyForAdmin bool   `db:"only_for_admin"`
}

type CreateDirectionParams struct {
	ChatID            int64
	ParentDirectionID sql.NullInt64
	DirectionName     string
}

type HandleStartCommandParams struct {
	Text   string
	ChatID int64
	TgName string
}

type SetUserActiveParams struct {
	TgName string
	ChatID int64
	UUID   string
}

type HandleCreateDirectionCommandParams struct {
	Text              string
	ChatID            int64
	ParentDirectionID int
}

type HandleAddInfoCommandParams struct {
	ChatID            int64
	ParentDirectionID int
}

type HandlePrintQuestionsCommandParams struct {
	ChatID            int64
	ParentDirectionID int
}

type SetAwaitingStatusParams struct {
	StatusID int
	ChatID   int64
}

type GetUserDirectionParams struct {
	ChatID            int64
	ParentDirectionID int
}

type GetMainButtonsResult struct {
	Name         string `db:"name"`
	OnlyForAdmin bool   `db:"only_for_admin"`
}

type UserDirection struct {
	ID                int    `db:"id"`
	Direction         string `db:"direction"`
	ParentDirectionID int    `db:"parent_directon_id"`
	CreatedAt         int64  `db:"created_at"`
	UpdatedAt         int64  `db:"updated_at"`
}

type Question struct {
	ID   int    `db:"id"`
	Text string `db:"text"`
}

type SetParentDirectionParams struct {
	ChatID            int64
	ParentDirectionID int
}

type HandleAwaitingQuestionParams struct {
	ChatID   int64
	Question string
}

type HandleAwaitingAnswerParams struct {
	ChatID int64
	Answer string
}

type SendMessageParams struct {
	ChatID         int64
	Text           string
	Keyboard       tgbotapi.InlineKeyboardMarkup
	IsNeedToRemove bool
}

type SetDirectionForInfoParams struct {
	ChatID      int64
	DirectionID int
}

type SetInfoIDParams struct {
	ChatID int64
	InfoID int
}
