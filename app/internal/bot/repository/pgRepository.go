package repository

import (
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type BotPGRepo struct {
	db *sqlx.DB
}

func NewBotPGRepo(db *sqlx.DB) bot.PgRepository {
	return &BotPGRepo{db: db}
}

func (r *BotPGRepo) GetMainButtons(ctx context.Context) ([]models.GetMainButtonsResult, error) {
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

func (r *BotPGRepo) GetActiveUsers(ctx context.Context) ([]models.GetActiveUsersResult, error) {
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

func (r *BotPGRepo) GetUpdatedButtons(ctx context.Context, param int64) ([]models.GetUpdatedButtonsResult, error) {
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

func (r *BotPGRepo) SetUserActive(ctx context.Context, params models.SetUserActiveParams) error {
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
