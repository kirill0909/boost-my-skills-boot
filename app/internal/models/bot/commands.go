package bot

type HandleStartCommandParams struct {
	Text   string
	ChatID int64
	TgName string
}

// type SetStatusActiveParams struct {
// 	TgName string
// 	ChatID int64
// 	UUID   string
// }

type UserActivationParams struct {
	TgName string
	ChatID int64
	UUID   string
}

type GetMainButtonsResult struct {
	Name         string `db:"name"`
	OnlyForAdmin bool   `db:"only_for_admin"`
}
