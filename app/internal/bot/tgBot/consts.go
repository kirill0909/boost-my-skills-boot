package tgbot

const (
	// Messages
	wellcomeMessage = `
	Поздравляю! Вы успешно активировали свой аккаунт в BoostMySkillsBot. 
	Это бот который поможет прокочать вам навыки програмировани. 
	
	⚠️Как он будет это делать❔
	Очень просто. Бот будет присылать вам вопрос по програмированию.
	A вы должны ясно и четко ответить на него желательно в слух, а еще лучше если вы будете расказывать
	об этом кому нибудь, после нажимаете "Get an answer" и сравниваете с тем что сказали вы

	⚠️Как часто он будет присылать мне вопросы❔
	Сейчас бот будет присылать вaм вопросы тогда когда вы его об этом попросите.
	То есть появилась свободная минута, нажимаете на кнопоку "ask me", бот прислал вам вопрос,
	Едете в общественном транспорте, нажимаете на кнопку "ask me"... ну вы поняли.
	Отличная альтернатива социальным сетям.

	⚠️Ok. Что дальше❔
	Внизу под этим сообщением нужно выбрать, вопросы из какой ниши программирования вы хотите получать.
	Frontend или Backend. После этого бот будет готов к тому что бы мучить вас вопросами.
	`
	readyMessage                 = "Отлично! Как будете готовы жмите \"ask me\" и бот будет присылать вам вопросы"
	notQuestionsMessage          = "Для вашего напровления пока не добавлено ни одного вопроса. Обратитесь к @kirillkorunov."
	addQuestionMessage           = "Для добавления вопроса нажмите на кнопку add_question"
	handleAddQuestionMessage     = "Ввeдите вопрос"
	handleEnteredQuestionMessage = "Отлично! Теперь ввeдите ответ"
	handleEnteredAnswerMessage   = "Ваш вопрос и ответ успешно сохранен"
	unableToGetAnswer            = `Невозможно получить ответ для этого вопроса. 
	Обратитесь к @kirillkorunov для добовления или изменеия ответа к интересующему вас вопросу. Либо довьте вопрос/ответ заново)`

	// Buttons
	getUUIDButton     = "/get_uuid"
	askMeButton       = "/ask_me"
	addQuestionButton = "/add_question"
	backenddButton    = "I ❤️  backend"
	frontendButton    = "I ❤️  frontend"
	getAnswerButton   = "Get an answer"

	// Commands name
	startCommand   = "start"
	getUUIDCommand = "get_uuid"
	askMeCommend   = "ask_me"
	addQuestion    = "add_question"

	// Callback data
	backendCallbackData   = "backend"
	frontednCallbackData  = "frontend"
	getAnswerCallbackData = "Get_an_answer"

	// errors
	errUserActivation      = "Ошибка активации аккаута"
	errInternalServerError = "Внутренняя ошибка сервера"

	// stat machine
	idle             = 1
	awaitingQuestion = 2
	awaitingAnswer   = 3
)
