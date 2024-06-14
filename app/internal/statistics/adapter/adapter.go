package adapter

import (
	"boost-my-skills-bot/app/internal/statistics"
	"boost-my-skills-bot/app/internal/statistics/models"
	pb "boost-my-skills-bot/app/pkg/proto/github.com/kirill0909/boost-my-skills-boot/app/pkg/proto/boost_bot_proto"
	"context"
	"log/slog"
)

type Statistics struct {
	uc  statistics.UseCase
	log *slog.Logger
	pb.UnimplementedStatisticsServer
}

func NewStatistics(uc statistics.UseCase, log *slog.Logger) *Statistics {
	return &Statistics{uc: uc, log: log}
}

func (a *Statistics) GetStatistics(ctx context.Context, req *pb.GetStatisticsRequest) (*pb.GetStatisticsResponse, error) {
	params := models.GetStatisticsRequest{DateFrom: req.DateFrom, DateTo: req.DateTo}

	res, err := a.uc.GetStatistics(ctx, params)
	if err != nil {
		return &pb.GetStatisticsResponse{}, err
	}

	return &pb.GetStatisticsResponse{InfosAdded: res.InfosAdded}, nil
}
