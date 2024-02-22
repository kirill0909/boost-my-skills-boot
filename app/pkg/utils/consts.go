package utils

const (
	// regex layouts
	DirectionNameLayout = "^[a-zA-Z0-9 ]+$"

	// commands
	StartCommand           = "start"
	CreateDirectionCommand = "create_direction"

	// user states
	AwaitingDirectionNameStatus   = 1
	AwaitingParentDirecitonStatus = 2

	// prefix
	AwaitingStatusPrefix  = "awaitingStatusPrefix"
	ParentDirectionPrefix = "parentDirectionPrefix"
)
