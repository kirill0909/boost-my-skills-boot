package middleware

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	invalidApiKey          = "invalid api key"
	requestFailed          = "request failed"
	unauthenticatedRequest = "unauthenticated request"
	requestSuccess         = "request success"
)

func UnaryInterceptor(grpcApiKey string, log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)

		apiKey := md["api-key"]
		if len(apiKey) == 0 || apiKey[0] != grpcApiKey {
			log.Error(unauthenticatedRequest,
				"errorPatch", "UnaryInterceptor()",
				"method", info.FullMethod,
				"status", codes.Unauthenticated,
				"errorDetails", invalidApiKey)
			return nil, status.Errorf(codes.Unauthenticated, invalidApiKey)
		}

		res, err := handler(ctx, req)
		if err != nil {
			st, _ := status.FromError(err)
			log.Error(requestFailed, "method", info.FullMethod, "status", st.Code(), "errorDetails", err.Error())
			return nil, status.Error(st.Code(), requestFailed)
		}

		log.Info(requestSuccess, "method", info.FullMethod)
		return res, nil
	}
}
