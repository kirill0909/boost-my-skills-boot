package repository

import (
	"boost-my-skills-bot/internal/bot"
	"boost-my-skills-bot/internal/bot/models"
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
