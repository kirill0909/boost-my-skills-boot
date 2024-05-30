package usecase

import (
	"boost-my-skills-bot/app/internal/statistics"
	"boost-my-skills-bot/app/internal/statistics/models"
	"context"
)

type StatisticsUseCase struct {
	pgRepo statistics.PgRepository
}

func NewStatisticsUsecase(pgRepo statistics.PgRepository) statistics.UseCase {
	return &StatisticsUseCase{pgRepo: pgRepo}
}

func (u *StatisticsUseCase) GetStatistics(ctx context.Context, params models.GetStatisticsRequest) (models.GetStatisticsResult, error) {
	res, err := u.pgRepo.GetStatistics(ctx, params)
	if err != nil {
		return models.GetStatisticsResult{}, err
	}

	return models.GetStatisticsResult{InfosAdded: res.InfosAdded}, nil
}
