package interceptors

import (
	"context"
	"time"

	"github.com/Ord1nI/netStats/internal/logger"
	"google.golang.org/grpc"
)

func LoggerInterceptor(log logger.Logger) grpc.UnaryServerInterceptor{
	return grpc.UnaryServerInterceptor(func (ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any,error) {

		tn := time.Now()


		res, err := handler(ctx, req)

		log.Info(
			"Get Request",
			"Time to procces: ", time.Since(tn),
		)

		return res, err

	})
}
