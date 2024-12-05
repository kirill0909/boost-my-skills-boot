package repository

import (
	"boost-my-skills-bot/app/internal/bot"
	"boost-my-skills-bot/app/internal/bot/models"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type botPGRepo struct {
	db *sqlx.DB
}

func NewBotPGRepo(db *sqlx.DB) bot.PgRepository {
	return &botPGRepo{db: db}
}

func (r *botPGRepo) GetMainKeyboard(ctx context.Context) ([]models.GetMainKeyboardResult, error) {
	rows, err := r.db.QueryContext(ctx, queryGetMainKeyboards)
	if err != nil {
		err = errors.Wrap(err, "BotPGRepo.GetMainKeyboard.queryGetMainKeyboards")
		return []models.GetMainKeyboardResult{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			err = errors.Wrap(err, "BotPGRepo.GetMainKeyboard.Close()")
		}
	}()

	var buttons []models.GetMainKeyboardResult
	var button models.GetMainKeyboardResult
	for rows.Next() {
		if err := rows.Scan(&button.ID, &button.Name, &button.OnlyForAdmin, &button.CreatedAt, &button.UpdatedAt); err != nil {
			err = errors.Wrap(err, "BotPGRepo.GetMainKeyboard.Scan")
			return []models.GetMainKeyboardResult{}, err
		}

		buttons = append(buttons, button)
	}

	return buttons, nil
}

// TODO: Remove
func (r *botPGRepo) GetMainButtons(ctx context.Context) ([]models.GetMainButtonsResult, error) {
	rows, err := r.db.QueryContext(ctx, queryGetMainButtons)
	if err != nil {
		err = errors.Wrap(err, "BotPGRepo.GetMainButtons.queryGetMainButtons")
		return []models.GetMainButtonsResult{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			err = errors.Wrap(err, "BotPGRepo.GetMainButtons.Close()")
		}
	}()

	var buttons []models.GetMainButtonsResult
	var button models.GetMainButtonsResult
	for rows.Next() {
		if err := rows.Scan(&button.Name, &button.OnlyForAdmin); err != nil {
			err = errors.Wrap(err, "BotPGRepo.GetMainButtons.Scan")
			return []models.GetMainButtonsResult{}, err
		}

		buttons = append(buttons, button)
	}

	return buttons, nil
}

func (r *botPGRepo) GetActiveUsers(ctx context.Context) ([]models.GetActiveUsersResult, error) {
	rows, err := r.db.QueryContext(ctx, queryGetActiveUsers)
	if err != nil {
		err = errors.Wrap(err, "BotPGRepo.GetActiveUsers.queryGetActiveUsers")
		return []models.GetActiveUsersResult{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			err = errors.Wrap(err, "BotPGRepo.GetActiveUsers.Close()")
		}
	}()

	var users []models.GetActiveUsersResult
	var user models.GetActiveUsersResult
	for rows.Next() {
		if err := rows.Scan(&user.ChatID, &user.IsAdmin); err != nil {
			err = errors.Wrap(err, "BotPGRepo.GetActiveUsers.Scan")
			return []models.GetActiveUsersResult{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *botPGRepo) GetUpdatedButtons(ctx context.Context, param int64) ([]models.GetUpdatedButtonsResult, error) {
	rows, err := r.db.QueryContext(ctx, queryGetUpdatedButtons, param)
	if err != nil {
		err = errors.Wrap(err, "BotPGRepo.GetUpdatedButtons.queryGetUpdatedButtons")
		return []models.GetUpdatedButtonsResult{}, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			err = errors.Wrap(err, "BotPGRepo.GetUpdatedButtons.Close()")
		}
	}()

	var buttons []models.GetUpdatedButtonsResult
	var button models.GetUpdatedButtonsResult
	for rows.Next() {
		if err := rows.Scan(&button.Name, &button.OnlyForAdmin); err != nil {
			err = errors.Wrap(err, "BotPGRepo.GetUpdatedButtons.Scan")
			return []models.GetUpdatedButtonsResult{}, err
		}

		buttons = append(buttons, button)
	}

	return buttons, nil

}

func (r *botPGRepo) SetUserActive(ctx context.Context, params models.SetUserActiveParams) error {
	res, err := r.db.ExecContext(ctx, querySetUserActive, params.ChatID, params.UUID, params.TgName)
	if err != nil {
		return errors.Wrapf(err, "BotPGRepo.SetUserActive.querySetUserActive. params(%+v)", params)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "BotPGRepo.SetUserActive.RowsAffected. params(%+v)", params)

	}

	if rowsAffected != 1 {
		return fmt.Errorf("BotPGRepo.SetUserActive rowsAffected != 1. params(%+v)", params)
	}

	return nil
}

func (r *botPGRepo) GetUserDirection(ctx context.Context, params models.GetUserDirectionParams) ([]models.UserDirection, error) {
	rows, err := r.db.QueryContext(ctx, queryGetUserDirection, params.ChatID, params.ParentDirectionID)
	if err != nil {
		err = errors.Wrapf(err, "BotPGRepo.GetUserDirection.queryGetUserDirection. params(%+v)", params)
		return []models.UserDirection{}, err
	}

	directionList := make([]models.UserDirection, 0, 50)
	var direction models.UserDirection
	for rows.Next() {
		if err := rows.Scan(&direction.ID, &direction.Direction, &direction.ParentDirectionID, &direction.CreatedAt, &direction.UpdatedAt); err != nil {
			err = errors.Wrapf(err, "BotPGRepo.GetUserDirection.queryGetUserDirection. params(%+v)", params)
			return []models.UserDirection{}, err
		}
		directionList = append(directionList, direction)
	}

	return directionList, nil
}

func (r *botPGRepo) CreateDirection(ctx context.Context, params models.CreateDirectionParams) (string, error) {
	var direction string
	if err := r.db.GetContext(ctx, &direction, queryCreateDirection, params.DirectionName, params.ChatID, params.ParentDirectionID); err != nil {
		return "", errors.Wrapf(err, "botPGRepo.CreateDirection.queryCreateDirection. params(%+v)", params)
	}

	return direction, nil
}

func (r *botPGRepo) SaveQuestion(ctx context.Context, params models.SaveQuestionParams) (int, error) {
	var infoID int
	if err := r.db.GetContext(ctx, &infoID, querySaveQuestion, params.Question, params.DirectionID); err != nil {
		return 0, errors.Wrapf(err, "botPGRepo.SaveQuestion.querySaveQuestion. params(%+v)", params)
	}

	return infoID, nil
}

func (r *botPGRepo) SaveAnswer(ctx context.Context, params models.SaveAnswerParams) error {
	if _, err := r.db.ExecContext(ctx, querySaveAnswer, params.Answer, params.InfoID); err != nil {
		return errors.Wrapf(err, "botPGRepo.SaveAnswer.querySaveAnswer. params(%+v)", params)
	}

	return nil
}

func (r *botPGRepo) GetQuestionsByDirectionID(ctx context.Context, directionID int) ([]models.Question, error) {
	rows, err := r.db.QueryContext(ctx, queryGetQuestionsByDirectionID, directionID)
	if err != nil {
		err = errors.Wrapf(err, "botPGRepo.GetQuestionsByDirectionID.queryGetQuestionsByDirectionID. directionID: %d", directionID)
		return []models.Question{}, err
	}

	questions := make([]models.Question, 0, 100)
	var question models.Question
	for rows.Next() {
		if err := rows.Scan(&question.ID, &question.Text); err != nil {
			err = errors.Wrapf(err, "botPGRepo.GetQuestionsByDirectionID.Scan(). directionID: %d", directionID)
			return []models.Question{}, err
		}

		questions = append(questions, question)
	}

	return questions, nil
}

func (r *botPGRepo) GetAnswerByInfoID(ctx context.Context, infoID int) (string, error) {
	var answer string
	if err := r.db.GetContext(ctx, &answer, queryGetAnswerByInfoID, infoID); err != nil {
		return "", errors.Wrapf(err, "botPGRepo.GetAnswerByInfoID.queryGetAnswerByInfoID. infoID: %d", infoID)
	}

	return answer, nil
}

func (r *botPGRepo) CreateInActiveUser(ctx context.Context) (string, error) {
	var result string
	if err := r.db.GetContext(ctx, &result, queryCreateInActiveUser); err != nil {
		return "", errors.Wrap(err, "botPGRepo.CreateInActiveUser.queryCreateInActiveUser")
	}

	return result, nil
}

func (r *botPGRepo) GetUserInfo(ctx context.Context, chatID int64) (models.UserInfo, error) {
	var result models.UserInfo
	if err := r.db.GetContext(ctx, &result, queryGetUserInfo, chatID); err != nil {
		err = errors.Wrapf(err, "botPGRepo.GetUserInfo.queryGetUserInfo chatID: %d", chatID)
		return models.UserInfo{}, err
	}

	return result, nil
}

func (r *botPGRepo) AddNewButtonToMainKeyboard(ctx context.Context, params models.AddNewButtonToMainKeyboardParams) error {
	result, err := r.db.ExecContext(ctx, queryAddNewButtonToMainKeyboard, params.ButtonName, params.OnlyForAdmin)
	if err != nil {
		return errors.Wrapf(err, "botPGRepo.AddNewButtonToMainKeyboard.queryAddNewButtonToMainKeyboard. params(%+v)", params)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "botPGRepo.AddNewButtonToMainKeyboard.RowsAffected. params(%+v)", params)
	}

	if rowsAffected != 1 {
		err := fmt.Errorf("wrong number of rows affected %d != 1", rowsAffected)
		return errors.Wrapf(err, "botPGRepo.AddNewButtonToMainKeyboard. params(%+v)", params)
	}

	return nil
}
