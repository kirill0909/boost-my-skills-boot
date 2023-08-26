package repository

import (
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type BotPGRepo struct {
	db *sqlx.DB
}

func NewBotPGRepo(db *sqlx.DB) bot.PgRepository {
	return &BotPGRepo{db: db}
}

func (r *BotPGRepo) GetUUID(ctx context.Context) (result string, err error) {
	if err = r.db.GetContext(ctx, &result, queryGetUUID); err != nil {
		return
	}

	return
}

func (r *BotPGRepo) IsAdmin(ctx context.Context, params models.GetUUID) (result bool, err error) {
	if err = r.db.GetContext(ctx, &result, queryIsAdmin, params.TgName, params.ChatID); err != nil {
		return
	}

	return
}

func (r *BotPGRepo) UserActivation(ctx context.Context, params models.UserActivation) (err error) {
	result, err := r.db.ExecContext(ctx, queryUserActivation, params.TgName, params.ChatID, params.UUID)
	if err != nil {
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if affected != 1 {
		err = fmt.Errorf("Wrong number of rows affected %d != 1", affected)
		return
	}

	return
}

func (r *BotPGRepo) SetUpBackendDirection(ctx context.Context, chatID int64) (err error) {
	if _, err = r.db.ExecContext(ctx, querySetUpBackendDirection, chatID); err != nil {
		return
	}

	return
}

func (r *BotPGRepo) SetUpFrontendDirection(ctx context.Context, chatID int64) (err error) {
	if _, err = r.db.ExecContext(ctx, querySetUpFrontedDirection, chatID); err != nil {
		return
	}

	return
}

func (r *BotPGRepo) GetRandomQuestion(ctx context.Context, params models.SubdirectionsCallbackParams) (
	result models.SubdirectionsCallbackResult, err error) {
	rows, err := r.db.QueryContext(ctx, queryGetRandomQuestion, params.ChatID, params.SubdirectionID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&result.QuestionID, &result.Question); err != nil {
			return
		}
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

func (r *BotPGRepo) GetAnswer(ctx context.Context, questionID int) (result string, err error) {
	if err = r.db.GetContext(ctx, &result, queryGetAnswer, questionID); err != nil {
		return
	}

	return
}

func (r *BotPGRepo) SaveQuestion(ctx context.Context, params models.SaveQuestionParams) (result int, err error) {
	if err = r.db.GetContext(ctx, &result, querySaveQuestion, params.ChatID, params.Question); err != nil {
		return
	}

	return
}

func (r *BotPGRepo) SaveAnswer(ctx context.Context, params models.SaveAnswerParams) (err error) {
	if _, err = r.db.ExecContext(ctx, querySaveAnswer, params.Answer, params.QuestionID); err != nil {
		return
	}

	return
}

func (r *BotPGRepo) GetSubdirections(ctx context.Context, params models.GetSubdirectionsParams) (result []string, err error) {
	rows, err := r.db.QueryContext(ctx, queryGetSubdirectons, params.ChatID)
	if err != nil {
		return
	}
	defer rows.Close()

	var res string
	for rows.Next() {
		if err = rows.Scan(&res); err != nil {
			return
		}

		result = append(result, res)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
