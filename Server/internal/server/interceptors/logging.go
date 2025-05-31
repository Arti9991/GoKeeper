package interceptors

import (
	"context"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// loggingInterceptor перехватчик для логирования вызваного метода
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	logger.Log.Info("Recieved new request", zap.String("Method", info.FullMethod))

	return handler(ctx, req)
}
