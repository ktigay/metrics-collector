package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// WithLogging логирует запрос.
func WithLogging(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()

		logger.Infow(
			"request",
			"request", req,
		)

		resp, err = handler(ctx, req)

		logger.Infow(
			"response",
			"response", resp,
			"duration", time.Since(start),
			"size", proto.Size(resp.(proto.Message)),
			"error", err,
		)
		return
	}
}
