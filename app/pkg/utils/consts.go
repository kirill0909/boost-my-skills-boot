package utils

const (
	// regex layouts
	DirectionNameLayout = "^[a-zA-Z0-9 /-]+$"
	MessageIDLayout     = "_MessageID_(\\d+)"
	ChatIDLayout        = "_ChatID_(\\d+)"

	// commands
	StartCommand           = "start"
	CreateDirectionCommand = "create_direction"
	AddInfoCommand         = "add_info"
	PrintQuestionsCommand  = "print_questions"

	// user states
	AwaitingDirectionNameStatus = 1
	AwaitingQuestionStatus      = 2
	AwaitingAnswerStatus        = 3

	// prefix
	AwaitingStatusPrefix        = "awaitingStatusPrefix"
	ParentDirectionPrefix       = "parentDirectionPrefix"
	ExpirationTimeMessagePrefix = "expirationTimeMessagePreifx"
	DirectionForInfoPrefix      = "directionForInfoPrefix"
	InfoPrefix                  = "infoPrefix"

	// callback types
	AwaitingParentDirectionCallbackType  = 1
	AwaitingAddInfoDirectionCallbackType = 2
	AwaitingPrintQuestionsCallbackType   = 3
	AwaitingInfoActionsCallbackType      = 4

	// snipets
	// go
	SnipetGoBegin = "snipet go begin"
	// sql
	SnipetSQLBegin = "snipet sql begin"
	// bash
	SnipetBashBegin = "snipet bash begin"
	// rust
	SnipetRustBegin = "snipet rust begin"
	SnipetEnd       = "snipet end"
)
