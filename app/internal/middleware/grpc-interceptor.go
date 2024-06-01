package middleware

import (
	"context"
	"log/slog"

	"boost-my-skills-bot/app/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryInterceptor(grpcApiKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		apiKey := md["api-key"]
		if len(apiKey) == 0 || apiKey[0] != grpcApiKey {
			slog.Error("UnaryInterceptor()", "rpc method", info.FullMethod, "error", utils.InvalidApiKey)
			return nil, status.Errorf(codes.Unauthenticated, utils.InvalidApiKey)
		}

		slog.Info("UnaryInterceptor()", "rpc method", info.FullMethod)
		return handler(ctx, req)
	}
}
