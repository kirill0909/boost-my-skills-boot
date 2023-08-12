package handler

import (
	"aifory-pay-admin-bot/internal/bot"
	"aifory-pay-admin-bot/pkg/utils"
	"context"

	pb "gitlab.axarea.ru/main/aiforypay/package/admin-bot-proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BotHandler struct {
	pb.UnimplementedAdminBotServer
	tgUC bot.Usecase
}

func NewBotHandler(tgUC bot.Usecase) pb.AdminBotServer {
	return &BotHandler{tgUC: tgUC}
}

func (h *BotHandler) NotifyGeneralAppealPercent(ctx context.Context, req *pb.NotifyGeneralAppealPercentRequest) (
	res *emptypb.Empty, err error) {
	ctx, span := utils.StartGrpcTrace(ctx, "BotHandler.NotifyGeneralAppealPercent")
	defer span.End()

	if err = h.tgUC.NotifyGeneralAppealPercent(ctx, req); err != nil {
		return &emptypb.Empty{}, utils.FormatErr(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *BotHandler) NotifyTraderAppealPercent(ctx context.Context, req *pb.NotifyTraderAppealPercentRequest) (
	res *emptypb.Empty, err error) {
	ctx, span := utils.StartGrpcTrace(ctx, "BotHandler.NotifyTraderAppealPercent")
	defer span.End()

	if err = h.tgUC.NotifyTraderAppealPercent(ctx, req); err != nil {
		return &emptypb.Empty{}, utils.FormatErr(err)
	}

	return &emptypb.Empty{}, nil
}
