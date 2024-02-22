package utils

const (
	// regex layouts
	DirectionNameLayout = "^[a-zA-Z0-9 ]+$"

	// commands
	StartCommand           = "start"
	CreateDirectionCommand = "create_direction"

	// user states
	AwaitingDirectionName   = 1
	AwaitingParentDireciton = 2
)
