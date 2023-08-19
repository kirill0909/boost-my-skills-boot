package repository

import (
	"boost-my-skills-bot/internal/bot"
	models "boost-my-skills-bot/internal/models/bot"
	"context"

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
	if _, err = r.db.ExecContext(ctx, queryUserActivation, params.TgName, params.ChatID, params.UUID); err != nil {
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

func (r *BotPGRepo) GetRandomQuestion(ctx context.Context, params models.AskMeParams) (result models.AskMeResult, err error) {
	rows, err := r.db.QueryContext(ctx, queryGetRandomQuestion, params.ChatID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&result.QuestionID, &result.Question, &result.Code); err != nil {
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
