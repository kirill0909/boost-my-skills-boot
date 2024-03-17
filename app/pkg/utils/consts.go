package utils

const (
	// regex layouts
	DirectionNameLayout = "^[a-zA-Z0-9 /]+$"
	MessageIDLayout     = "_MessageID_(\\d+)"
	ChatIDLayout        = "_ChatID_(\\d+)"

	// commands
	StartCommand           = "start"
	CreateDirectionCommand = "create_direction"
	AddInfoCommand         = "add_info"
	PrintQuestionsCommand  = "print_questions"

	// user states
	AwaitingDirectionNameStatus    = 1
	AwaitingParentDirecitonStatus  = 2
	AwaitingAddInfoDirectionStatus = 3
	AwaitingQuestionStatus         = 4
	AwaitingAnswerStatus           = 5
	AwaitingPrintQuestionsStatus   = 6
	AwaitingInfoActionsStatus      = 7

	// prefix
	AwaitingStatusPrefix        = "awaitingStatusPrefix"
	ParentDirectionPrefix       = "parentDirectionPrefix"
	ExpirationTimeMessagePrefix = "expirationTimeMessagePreifx"
	DirectionForInfoPrefix      = "directionForInfoPrefix"
	InfoPrefix                  = "infoPrefix"

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
