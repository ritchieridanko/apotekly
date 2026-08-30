package interceptors

import (
	"context"
	"errors"
	"time"

	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Logging(l *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		// Check if Request is OK
		fields := []logger.Field{
			logger.NewField("method", info.FullMethod),
			logger.NewField("latency", time.Since(start).String()),
		}
		if err == nil {
			fields = append(fields, logger.NewField("status", codes.OK.String()))
			l.Info(ctx, "REQUEST OK", fields...)
			return resp, nil
		}

		// Check if Request fails
		var e *ce.Error
		if errors.As(err, &e) {
			grpcErr := e.ToGRPCErr()
			st, _ := status.FromError(grpcErr)
			code := st.Code()

			fields = append(fields, logger.NewField("status", code.String()))
			fields = append(fields, e.Fields()...)
			fields = append(
				fields,
				logger.NewField("error_code", e.Code()),
				logger.NewField("error", e.Error()),
			)

			switch code {
			case
				codes.Aborted,
				codes.AlreadyExists,
				codes.Canceled,
				codes.FailedPrecondition,
				codes.InvalidArgument,
				codes.NotFound,
				codes.OutOfRange,
				codes.PermissionDenied,
				codes.Unauthenticated:
				l.Warn(ctx, "REQUEST ERROR", fields...)
			case
				codes.DataLoss,
				codes.Internal,
				codes.Unavailable,
				codes.Unknown:
				l.Error(ctx, "SYSTEM ERROR", fields...)
			}
			return nil, grpcErr
		}

		// Unknown Error Fallback
		fields = append(
			fields,
			logger.NewField("status", codes.Unknown.String()),
			logger.NewField("error_code", ce.CodeUnknown),
			logger.NewField("error", err.Error()),
		)

		l.Error(ctx, "UNKNOWN ERROR", fields...)
		return nil, status.Error(codes.Unknown, ce.MsgInternalServer)
	}
}
