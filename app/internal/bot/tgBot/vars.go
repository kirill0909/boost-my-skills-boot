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
	readyMessage                 = "Отлично! Как будете готовы жмите \"ask me\" и бот будет присылать вам вопросы"
	notQuestionsMessage          = "Для этого пока не добавлено ни одного вопроса. Обратитесь к @kirillkorunov."
	addQuestionMessage           = "Для добавления вопроса нажмите на кнопку add_question"
	handleAddQuestionMessage     = "Ввeдите вопрос"
	handleEnteredQuestionMessage = "Отлично! Теперь ввeдите ответ"
	handleEnteredAnswerMessage   = "Ваш вопрос и ответ успешно сохранен"
	unableToGetAnswer            = `Невозможно получить ответ для этого вопроса. 
	Обратитесь к @kirillkorunov для добовления или изменеия ответа к интересующему вас вопросу. Либо довьте вопрос/ответ заново)`
	noOneSubdirectionsFoundMessage = "No subdirections were found for your direction. To add subdirections refer to @kirillkorunov"
	directionQuestionMessage       = "In which direction do you want to add a question?"
	chooseSubdirectionMessage      = "Choose subdirections"

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
	// Ask me callback
	callbackDataAskMe = []string{"1AskMe", "2AskMe", "3AskMe", "4AskMe", "5AskMe", "6AskMe", "7AskMe"}
	// Add question callback
	callbackDataSubdirectionAddQuestion = []string{
		"1_Go_AddQuestion", "2_Computer_Since_AddQuestion", "3_Network_AddQuestion", "4_DB_AddQuestion",
		"5_Algorithm_AddQuestion", "6_Architecture_AddQuestion", "7_General_AddQuestion"}
	callbackDataSubSubdirectionAddQuestion = []string{
		"1_1_Go_Channels_AddQuestion", "1_2_Go_Interface_AddQuestion", "1_3_Go_Sheduler_AddQuestion",
		"1_4_Go_String_AddQuestion", "1_5_Go_sync_AddQuestion", "1_6_Go_Slice_AddQuestion",
		"1_7_Go_Array_AddQuestion", "1_8_Go_Context_AddQuestion", "1_9_Go_GC_AddQuestion",
		"1_10_Go_Goroutines_AddQuestion", "1_11_Go_ErrorGroup_AddQuestion", "1_12_Go_Map_AddQuestion",
		"1_13_Go_OS_AddQuestion", "1_14_Go_buildIn_AddQuestion", "1_15_Go_keywords_AddQuestion",
		"1_16_Go_pointers_AddQuestion", "1_17_Go_net/http_AddQuestion", "1_18_Go_General_AddQuestion",
	}

	// errors
	errUserActivation      = "Ошибка активации аккаута"
	errUUIDAlreadyExists   = "Ошибка активации аккаута. Аккаунт с таким UUID уже существует. Для получения новаго UUID обратитесь к @kirillkorunov"
	errInternalServerError = "Внутренняя ошибка сервера"

	// stat machine
	idle                    = 1
	awaitingQuestion        = 2
	awaitingAnswer          = 3
	awaitingSubdirection    = 4
	awaitingSubSubdirection = 4
)
