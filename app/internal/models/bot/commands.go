package bot

type HandleStartCommandParams struct {
	Text   string
	ChatID int64
}

type CompareUUIDParams struct {
	UUID   string
	ChatID int64
}
