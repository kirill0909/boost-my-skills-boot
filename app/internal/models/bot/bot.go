package bot

type GetUUID struct {
	ChatID int64
	TgName string
}

type UserActivation struct {
	ChatID int64
	TgName string
	UUID   string
}

type AskMeParams struct {
	ChatID int64
}

type AddQuestionParams struct {
	State          int
	QuestionID     int
	SubdirectionID int
}

type SaveQuestionParams struct {
	ChatID         int64
	Question       string
	SubdirectionID int
}

type SaveAnswerParams struct {
	Answer     string
	QuestionID int
}

type GetSubdirectionsParams struct {
	ChatID int64
}

type SubdirectionsCallbackParams struct {
	ChatID         int64
	SubdirectionID int
}

type SubdirectionsCallbackResult struct {
	Question   string
	QuestionID int
}
