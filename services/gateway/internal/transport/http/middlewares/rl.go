package middlewares

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/models"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/repositories/cache"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

const ctxTimeout time.Duration = 100 * time.Millisecond

type RateLimiterMiddleware struct {
	cache  cache.RateLimiterCache
	logger *logger.Logger
}

func NewRateLimiterMiddleware(cc cache.RateLimiterCache, l *logger.Logger) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{cache: cc, logger: l}
}

func (m *RateLimiterMiddleware) LimitByIPAddress(namespace string, window time.Duration, limit uint16) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		ip := ctx.ClientIP()
		err := m.cache.LimitByFixedWindow(
			c,
			&models.LimitByFixedWindowRL{
				IPAddress: ip,
				Namespace: namespace,
				Window:    window,
				Limit:     limit,
			},
		)
		if err != nil && err.Code() == ce.CodeCacheScriptExec {
			m.logger.Warn(
				c,
				"failed to execute rate limiter, falling back to allow request",
				err.Append(
					logger.NewField("ip_address", ip),
					logger.NewField("error_code", err.Code()),
					logger.NewField("error", err.Unwrap()),
				).Fields()...,
			)

			ctx.Next()
			return
		}
		if err != nil {
			err.Bind(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// NOTE: Shall be placed AFTER the Auth() middleware
func (m *RateLimiterMiddleware) LimitByAuthID(namespace string, window time.Duration, limit uint16) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), ctxTimeout)
		defer cancel()

		authCtx := utils.CtxAuth(ctx.Request.Context())
		if authCtx == nil {
			ce.NewError(
				ce.CodeMissingContextValue,
				ce.MsgInternalServer,
				errors.New("auth missing from context"),
			).Bind(
				ctx,
			)

			ctx.Abort()
			return
		}

		err := m.cache.LimitBySlidingWindow(
			c,
			&models.LimitBySlidingWindowRL{
				AuthID:    authCtx.AuthID,
				Namespace: namespace,
				Now:       time.Now(),
				Window:    window,
				Limit:     limit,
			},
		)
		if err != nil && err.Code() == ce.CodeCacheScriptExec {
			m.logger.Warn(
				c,
				"failed to execute rate limiter, falling back to allow request",
				err.Append(
					logger.NewField("auth_id", authCtx.AuthID),
					logger.NewField("error_code", err.Code()),
					logger.NewField("error", err.Unwrap()),
				).Fields()...,
			)

			ctx.Next()
			return
		}
		if err != nil {
			err.Bind(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
