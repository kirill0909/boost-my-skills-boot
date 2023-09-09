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
	State             int
	QuestionID        int
	SubdirectionID    int
	SubSubdirectionID int
}

type SaveQuestionParams struct {
	ChatID            int64
	Question          string
	SubdirectionID    int
	SubSubdirectionID int
}

type SaveAnswerParams struct {
	Answer     string
	QuestionID int
}

type GetSubdirectionsParams struct {
	ChatID int64
}

type GetSubSubdirectionsParams struct {
	ChatID         int64
	SubdirectionID int
}

type AksMeCallbackParams struct {
	ChatID            int64
	SubdirectionID    int
	SubSubdirectionID int
}

type SubdirectionsCallbackResult struct {
	Question   string
	QuestionID int
}

type DirectionInfo struct {
	DirectionID   int    `db:"direction_id"`
	DirectionName string `db:"direction_name"`
}

type SubdirectionsData struct {
	SubdirectionInfo    []SubdirectionInfo
	SubSubdirectionInfo []SubSubdirectionInfo
}

type SubdirectionInfo struct {
	DirectionID      int    `db:"direction_id"`
	SubdirectionID   int    `db:"sub_direction_id"`
	SubdirectionName string `db:"sub_direction_name"`
}

type SubSubdirectionInfo struct {
	DirectionID         int    `db:"direction_id"`
	SubdirectionID      int    `db:"sub_direction_id"`
	SubSubdirectionID   int    `db:"sub_sub_direction_id"`
	SubSubdirectionName string `db:"sub_sub_direction_name"`
}
