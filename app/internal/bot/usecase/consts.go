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
	directionQuestionMessage       = "In which direction do you want to add a question?"
	readyMessage                   = "Отлично! Как будете готовы жмите \"ask me\" и бот будет присылать вам вопросы"

	// main buttons
	getUUIDButton = "/get_uuid"
	askMeButton   = "/ask_me"
	addInfoButton = "/add_info"
)
