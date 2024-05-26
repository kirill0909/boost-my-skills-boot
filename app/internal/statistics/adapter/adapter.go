package adapter

import (
	pb "boost-my-skills-bot/app/pkg/proto/github.com/kirill0909/boost-my-skills-boot/app/pkg/proto/boost_bot_proto"
	"context"
)

type Statistics struct {
	pb.UnimplementedStatisticsServer
}

func NewStatistics() *Statistics {
	return &Statistics{}
}

func (a *Statistics) GetStatistics(ctx context.Context, req *pb.GetStatisticsRequest) (*pb.GetStatisticsResponse, error) {
	return &pb.GetStatisticsResponse{InfosAdded: 10}, nil
}
