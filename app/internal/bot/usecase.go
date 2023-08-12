package bot

import (
	"context"

	pb "gitlab.axarea.ru/main/aiforypay/package/admin-bot-proto"
)

type Usecase interface {
	NotifyGeneralAppealPercent(ctx context.Context, req *pb.NotifyGeneralAppealPercentRequest) (err error)
	NotifyTraderAppealPercent(ctx context.Context, req *pb.NotifyTraderAppealPercentRequest) (err error)
}
