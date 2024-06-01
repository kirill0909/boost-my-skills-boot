package middleware

import (
	"context"
	"fmt"
	"log/slog"

	"boost-my-skills-bot/app/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(grpcApiKey string, log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		apiKey := md["api-key"]
		if len(apiKey) == 0 || apiKey[0] != grpcApiKey {
			log.Error("UnaryInterceptor()", "error", fmt.Sprintf("rpc  method: %s, error: %s", info.FullMethod, utils.InvalidApiKey))
			return nil, status.Errorf(codes.Unauthenticated, utils.InvalidApiKey)
		}

		log.Info("UnaryInterceptor()", "info", info.FullMethod)
		return handler(ctx, req)
	}
}
