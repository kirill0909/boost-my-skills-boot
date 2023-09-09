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
	getUUIDButton   = "/get_uuid"
	askMeButton     = "/ask_me"
	addInfoButton   = "/add_info"
	backenddButton  = "I ❤️  backend"
	frontendButton  = "I ❤️  frontend"
	getAnswerButton = "Get an answer"

	// Commands name
	startCommand   = "start"
	getUUIDCommand = "get_uuid"
	askMeCommend   = "ask_me"
	addQuestion    = "add_info"

	// Callback data
	backendCallbackData   = "backend"
	frontednCallbackData  = "frontend"
	getAnswerCallbackData = "Get_an_answer"
	// Ask me sub callback
	callbackDataAskMe = []string{
		"1_Go_AskMe", "2_ComputerSince_AskMe", "3_Network_AskMe", "4_DB_AskMe",
		"5_Algorithms_AskMe", "6_Architecture_AskMe", "7_Unix_AskMe", "8_Rust_AskMe",
		"9_DevOps_AskMe", "10_General_AskMe",
	}
	// Ask me sub sub callback
	callbackDataSubSubdirectionAskMe = []string{
		"1_1_Go_Channels_AskMe", "1_2_Go_Interface_AskMe", "1_3_Go_Sheduler_AskMe",
		"1_4_Go_Sync_AskMe", "1_5_Go_Slice_AskMe", "1_6_Go_Context_AskMe",
		"1_7_Go_GC_AskMe", "1_8_Go_Goroutines_AskMe", "1_9_Go_Map_AskMe",
		"1_10_Go_OS_AskMe", "1_11_Go_BuildIn_AskMe", "1_12_Go_Pointers_AskMe",
		"1_13_Go_net/http_AskMe", "1_14_Go_General_AskMe",
	}
	// Add sub question callback
	callbackDataSubdirectionAddQuestion = []string{
		"1_Go_AddQuestion", "2_Computer_Since_AddQuestion", "3_Network_AddQuestion", "4_DB_AddQuestion",
		"5_Algorithm_AddQuestion", "6_Architecture_AddQuestion", "7_Unix_AddQuestion", "8_Rust_AddQuestion",
		"9_DevOps_AddQuestion", "10_General_AddQuestion",
	}
	// Add sub sub quersis callback
	callbackDataSubSubdirectionAddQuestion = []string{
		"1_1_Go_Channels_AddQuestion", "1_2_Go_Interface_AddQuestion", "1_3_Go_Sheduler_AddQuestion",
		"1_4_Go_Sync_AddQuestion", "1_5_Go_Slice_AddQuestion", "1_6_Go_Context_AddQuestion",
		"1_7_Go_GC_AddQuestion", "1_8_Go_Goroutines_AddQuestion", "1_9_Go_Map_AddQuestion",
		"1_10_Go_OS_AddQuestion", "1_11_Go_BuildIn_AddQuestion", "1_12_Go_Pointers_AddQuestion",
		"1_13_Go_net/http_AddQuestion", "1_14_Go_General_AddQuestion",
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
