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
	GetInviteLinkCommand   = "get_invite_link"
	// new commands
	ServiceCommand = "service"

	// user states
	AwaitingDirectionNameStatus = 1
	AwaitingQuestionStatus      = 2
	AwaitingAnswerStatus        = 3
	// new status
	AwaitingNewMainButtonNameForUserStatus = 4

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
	AwaitingPrintAnswerCallbackType      = 4
	// new callbacks
	ActionsWithMainKeyboardCallbackType       = 5
	ActionsWithUsersCallbackType              = 6
	RenameMainKeyboardActionCallbackType      = 7
	AddForUserMainKeyboardActionCallbackType  = 8
	AddForAdminMainKeyboardActionCallbackType = 9
	RemoveMainKeyboardActionCallbackType      = 10
	ComeBackToServiceButtonsCallbackType      = 11

	// TODO: Remove
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
