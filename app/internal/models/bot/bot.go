package bot

import "github.com/guregu/null"

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

type AskMeResult struct {
	Question   string
	QuestionID int
	Code       null.String
}

type AddQuestionParams struct {
	State      int
	QuestionID int
}

type SaveQuestionParams struct {
	ChatID   int64
	Question string
}

type SaveAnswerParams struct {
	Answer     string
	QuestionID int
}

type GetSubdirectionsParams struct {
	ChatID int64
}
