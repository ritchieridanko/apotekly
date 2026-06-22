package interceptors

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recovery(l *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				l.Error(
					ctx,
					"PANIC RECOVERED",
					logger.NewField("method", info.FullMethod),
					logger.NewField("panic", fmt.Sprintf("%v", r)),
					logger.NewField("stack_trace", debug.Stack()),
				)
				err = status.Error(codes.Internal, ce.MsgInternalServer)
			}
		}()

		return handler(ctx, req)
	}
}
