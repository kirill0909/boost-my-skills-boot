package bot

type GetActiveUsersResult struct {
	ChatID  int64 `db:"tg_chat_id"`
	IsAdmin bool  `db:"is_admin"`
}

type GetUpdatedButtonsResult struct {
	Name         string `db:"name"`
	OnlyForAdmin bool   `db:"only_for_admin"`
}
