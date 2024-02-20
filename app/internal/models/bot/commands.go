package bot

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

type GetMainButtonsResult struct {
	Name         string `db:"name"`
	OnlyForAdmin bool   `db:"only_for_admin"`
}

type GetUserDirectionsResult struct {
	ID                int    `db:"id"`
	Direction         string `db:"direction"`
	ParentDirectionID int    `db:"parent_directon_id"`
	CreatedAt         int64  `db:"created_at"`
	UpdatedAt         int64  `db:"updated_at"`
}
