package utils

const (
	// regex layouts
	DirectionNameLayout = "^[a-zA-Z0-9 ]+$"

	// commands
	StartCommand    = "start"
	CreateDirection = "create_direction"

	// user states
	AwaitingDirectionName   = 1
	AwaitingParentDireciton = 2
)
