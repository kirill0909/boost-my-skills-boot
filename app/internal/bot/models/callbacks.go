package models

type HandleActionsWithMainKeyboardCallbackParams struct {
	ChatID    int64
	MessageID int
}

type HandleAddForUserMainKeyboardActionCallbackParams struct {
	ChatID    int64
	MessageID int
}

type HandleComeBackeToServiceButtonsCallbackParams struct {
	ChatID    int64
	MessageID int
}
