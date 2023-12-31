package bot

type HandleStartCommandParams struct {
	Text   string
	ChatID int64
	TgName string
}

type SetStatusActiveParams struct {
	TgName string
	ChatID int64
	UUID   string
}
