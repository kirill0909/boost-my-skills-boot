package usecase

const (
	notAdmin = "Вы не администратор! Для получения инвайт токена обратитесь к @kirillkorunov"

	// stat machine
	idle                    = 1
	awaitingQuestion        = 2
	awaitingAnswer          = 3
	awaitingSubdirection    = 4
	awaitingSubSubdirection = 4

	// mesages
	noOneSubdirectionsFoundMessage = "No subdirections were found for your direction. To add subdirections refer to @kirillkorunov"
	subdirectionAddInfoMessage     = "In which sub direction do you want to add a question?"
	subSubdirectionAddInfoMessage  = "In which sub sub direction do you want to add a question?"
	readyMessage                   = "Отлично! Как будете готовы жмите \"ask me\" и бот будет присылать вам вопросы"
	enterQuestionMessage           = "Alright, Enter yout question"
	subdirectionAskMeMessage       = "Choose sub direction for question"
	notQuestionsMessage            = "Not a single question was found in this direction"

	// main buttons
	getUUIDButton = "/get_uuid"
	askMeButton   = "/ask_me"
	addInfoButton = "/add_info"

	// buttons
	getAnswerButton = "Get an answer"
)
