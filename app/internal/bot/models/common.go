package models

import (
	"database/sql"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GetMainKeyboardResult struct {
	ID           int          `db:"id"`
	Name         string       `db:"name"`
	OnlyForAdmin bool         `db:"only_for_admin"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}

type SendMessageParams struct {
	ChatID         int64
	Text           string
	InlineKeyboard tgbotapi.InlineKeyboardMarkup
	ReplyKeyboard  tgbotapi.ReplyKeyboardMarkup
	IsNeedToRemove bool
}

type EditMessageParams struct {
	ChatID    int64
	MessageID int
	Text      string
	Keyboard  tgbotapi.InlineKeyboardMarkup
}

type HandleAwaitingNewMainButtonNameParams struct {
	ChatID       int64
	ButtonName   string
	OnlyForAdmin bool
}

type AddNewButtonToMainKeyboardParams struct {
	ButtonName   string
	OnlyForAdmin bool
}
