package repository

import (
	"boost-my-skills-bot/app/internal/statistics"
	"context"

	"boost-my-skills-bot/app/internal/statistics/models"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type StatisticsPgRepo struct {
	db *sqlx.DB
}

func NewStatisticsPgRepo(db *sqlx.DB) statistics.PgRepository {
	return &StatisticsPgRepo{db: db}
}

func (r *StatisticsPgRepo) GetStatistics(ctx context.Context, params models.GetStatisticsRequest) (models.GetStatisticsResult, error) {
	var result int64
	if err := r.db.GetContext(ctx, &result, queryGetStatistics, params.DateFrom, params.DateTo); err != nil {
		err = errors.Wrapf(err, "StatisticsPgRepo.GetStatistics.queryGetStatistics: params(%+v)", params)
		return models.GetStatisticsResult{}, err
	}

	return models.GetStatisticsResult{InfosAdded: result}, nil
}
