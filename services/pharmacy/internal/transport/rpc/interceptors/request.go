package interceptors

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Request(l *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Check if Metadata exists
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			l.Error(
				ctx,
				"REQUEST NOT ACCEPTED",
				logger.NewField("method", info.FullMethod),
				logger.NewField("error_code", ce.CodeMissingMetadata),
			)
			return nil, status.Error(codes.Internal, ce.MsgInternalServer)
		}

		// Request ID Metadata
		if values := md.Get(constants.MDKeyRequestID); len(values) > 0 {
			ctx = context.WithValue(ctx, constants.CtxKeyRequestID, values[0])
		}

		return handler(ctx, req)
	}
}
