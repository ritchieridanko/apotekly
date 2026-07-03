package interceptors

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ritchieridanko/apotekly/services/auth/internal/transport/rpc/policies"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Auth() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		// Check if method has an auth policy
		policy, exists := policies.AuthPolicies[info.FullMethod]
		if !exists {
			return handler(ctx, req)
		}

		// Check if metadata exists
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ce.NewError(ce.CodeMissingMetadata, ce.MsgInternalServer, nil)
		}

		// Check if policy requires authentication
		var authID uint64
		if policy.MustBeAuthenticated() {
			values := md.Get(constants.MDKeyAuthID)
			if len(values) == 0 {
				return nil, ce.NewError(
					ce.CodeUnauthenticated,
					ce.MsgUnauthenticated,
					errors.New("auth_id missing from metadata"),
				)
			}

			id, err := strconv.ParseUint(values[0], 10, 64)
			if err != nil {
				return nil, ce.NewError(
					ce.CodeTypeConversionFailed,
					ce.MsgInternalServer,
					fmt.Errorf("failed to convert auth_id (%v) to uint64: %w", values[0], err),
				)
			}

			authID = id
		}

		// Check if policy requires verification
		var isVerified bool
		if policy.MustBeVerified() {
			values := md.Get(constants.MDKeyIsVerified)
			if len(values) == 0 {
				return nil, ce.NewError(
					ce.CodeEmailNotVerified,
					ce.MsgEmailNotVerified,
					errors.New("is_email_verified missing from metadata"),
				)
			}

			verified, err := strconv.ParseBool(values[0])
			if err != nil {
				return nil, ce.NewError(
					ce.CodeTypeConversionFailed,
					ce.MsgInternalServer,
					fmt.Errorf("failed to convert is_email_verified (%v) to bool: %w", values[0], err),
				)
			}
			if !verified {
				return nil, ce.NewError(ce.CodeEmailNotVerified, ce.MsgEmailNotVerified, nil)
			}

			isVerified = verified
		}

		// Check if policy requires role authorization
		var role string
		if policy.RequireRole() {
			values := md.Get(constants.MDKeyRole)
			if len(values) == 0 {
				return nil, ce.NewError(
					ce.CodeRoleNotAuthorized,
					ce.MsgUnauthorized,
					errors.New("role missing from metadata"),
				)
			}
			if role = values[0]; !policy.IsRoleAuthorized(role) {
				return nil, ce.NewError(
					ce.CodeRoleNotAuthorized,
					ce.MsgUnauthorized,
					errors.New("role unauthorized"),
					logger.NewField("role", role),
				)
			}
		}

		return handler(
			context.WithValue(
				ctx,
				constants.CtxKeyAuth,
				&utils.AuthContext{
					AuthID:          authID,
					Role:            role,
					IsEmailVerified: isVerified,
				},
			),
			req,
		)
	}
}
