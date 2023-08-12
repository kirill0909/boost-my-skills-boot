package bot

type GetUUID struct {
	ChatID int64
	TgName string
}

type UserActivation struct {
	ChatID int64
	TgName string
	UUID   string
}
