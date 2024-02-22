package models

import "database/sql"

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
	Text   string
	TgName string
	ChatID int64
}

type SetAwaitingStatusParams struct {
	StatusID int
	ChatID   int64
}

type GetUserDirectionParams struct {
	ChatID            int64
	ParentDirectionID sql.NullInt64
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
