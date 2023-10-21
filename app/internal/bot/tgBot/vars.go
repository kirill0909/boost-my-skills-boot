package tgbot

var (
	// Messages
	wellcomeMessage = `
	Поздравляю! Вы успешно активировали свой аккаунт в BoostMySkillsBot. 
	Это бот который поможет прокочать вам навыки программирования. 
	
	⚠️Как он будет это делать❔
	Очень просто. Бот будет присылать вам вопрос по программированию.
	A вы должны ясно и четко ответить на него желательно в слух, а еще лучше если вы будете расказывать
	об этом кому нибудь, после нажимаете "Get an answer" и сравниваете с тем что сказали вы

	⚠️Как часто он будет присылать мне вопросы❔
	Сейчас бот будет присылать вaм вопросы тогда когда вы его об этом попросите.
	То есть появилась свободная минута, нажимаете на кнопоку "ask me", бот прислал вам вопрос,
	Едете в общественном транспорте, нажимаете на кнопку "ask me"... ну вы поняли.
	Отличная альтернатива социальным сетям.

	⚠️Ok. Что дальше❔
	Внизу под этим сообщением нужно выбрать, вопросы из какой ниши программирования вы хотите получать.
	Frontend или Backend. После этого вы можете добавить свои вопросы нажав на кнопку add_question
	и получать их в случайном порядке нажав на кнопку ask_me. Удачи)
	`
	addQuestionMessage           = "To add info press on the button add_question"
	handleEnteredQuestionMessage = "Cool! Now enter your answer"
	handleEnteredAnswerMessage   = "Your info successfully stored"

	// Buttons
	getUUIDButton  = "/get_uuid"
	askMeButton    = "/ask_me"
	addInfoButton  = "/add_info"
	backenddButton = "I ❤️  backend"
	frontendButton = "I ❤️  frontend"

	// Commands name
	startCommand   = "start"
	getUUIDCommand = "get_uuid"
	askMeCommend   = "ask_me"
	addInfo        = "add_info"
	printInfo      = "print_info"

	// errors
	errUserActivation      = "Error account activation"
	errUUIDAlreadyExists   = "Error account activation. Account with some UUID already exists. For get new UUID connect to @kirillkorunov"
	errInternalServerError = "Internal server error"
)
