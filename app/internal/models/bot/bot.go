package bot

import (
	"sync"
)

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
	ChatID            int64
	MessageID         int
	CallbackData      string
	SubdirectionID    int
	SubSubdirectionID int
}

type AddInfoParams struct {
	State             int
	QuestionID        int
	SubdirectionID    int
	SubSubdirectionID int
}

type AddInfoSubdirectionParams struct {
	ChatID         int64
	MessageID      int
	CallbackData   string
	SubdirectionID int
}

type AddInfoSubSubdirectionParams struct {
	ChatID            int64
	MessageID         int
	CallbackData      string
	SubdirectionID    int
	SubSubdirectionID int
}

// type AskMeSubdirectionsParams struct {
// 	ChatID         int64
// 	MessageID      int
// 	CallbackData   string
// 	SubdirectionID int
// }
//
// type AskMeSubSubdirectionsParams struct {
// 	ChatID            int64
// 	MessageID         int
// 	CallbackData      string
// 	SubdirectionID    int
// 	SubSubdirectionID int
// }

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

type AskMeCallbackResult struct {
	Question   string
	QuestionID int
}

type SetUpDirection struct {
	ChatID       int64
	MessageID    int
	CallbackData string
	DirectionID  int
}

type DirectionInfo struct {
	DirectionID   int    `db:"direction_id"`
	DirectionName string `db:"direction_name"`
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

type DirectionsData struct {
	DirectionInfo       []DirectionInfo
	SubdirectionInfo    []SubdirectionInfo
	SubSubdirectionInfo []SubSubdirectionInfo
	mutex               sync.Mutex
}

func (d *DirectionsData) Store(
	directionInfo []DirectionInfo,
	subdirectionInfo []SubdirectionInfo,
	subSubdirectionInfo []SubSubdirectionInfo) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.DirectionInfo = directionInfo
	d.SubdirectionInfo = subdirectionInfo
	d.SubSubdirectionInfo = subSubdirectionInfo
}

func (d *DirectionsData) GetSubdirectionsByDirectionID(directionID int) (result []SubdirectionInfo) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, v := range d.SubdirectionInfo {
		if v.DirectionID == directionID {
			result = append(result, v)
		}
	}

	return
}

func (d *DirectionsData) GetSubSubdirectionsBySubdirectionID(subdirectionID int) (result []SubSubdirectionInfo) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for _, v := range d.SubSubdirectionInfo {
		if v.SubdirectionID == subdirectionID {
			result = append(result, v)
		}
	}

	return
}
