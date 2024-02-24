package utils

const (
	// regex layouts
	DirectionNameLayout = "^[a-zA-Z0-9 ]+$"
	MessageIDLayout     = "_MessageID_(\\d+)"
	ChatIDLayout        = "_ChatID_(\\d+)"

	// commands
	StartCommand           = "start"
	CreateDirectionCommand = "create_direction"
	AddInfoCommand         = "add_info"

	// user states
	AwaitingDirectionNameStatus   = 1
	AwaitingParentDirecitonStatus = 2
	AwaitingAddInfoDirection      = 3
	AwaitingQuestion              = 4

	// prefix
	AwaitingStatusPrefix        = "awaitingStatusPrefix"
	ParentDirectionPrefix       = "parentDirectionPrefix"
	ExpirationTimeMessagePrefix = "ExpirationTimeMessagePreifx"
)
