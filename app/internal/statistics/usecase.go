package statistics

import (
	"boost-my-skills-bot/app/internal/statistics/models"
	"context"
)

type UseCase interface {
	GetStatistics(context.Context, models.GetStatisticsRequest) (models.GetStatisticsResult, error)
}
