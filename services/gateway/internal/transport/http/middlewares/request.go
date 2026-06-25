package middlewares

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

func Request(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := strings.TrimSpace(ctx.GetHeader("X-Request-ID"))
		if requestID == "" {
			uuid, err := utils.GenerateUUIDv7()
			if err != nil {
				requestID = fmt.Sprintf("fallback-%d", time.Now().UnixNano())

				l.Warn(
					ctx.Request.Context(),
					"failed to generate request id, using fallback",
					logger.NewField("fallback_id", requestID),
					logger.NewField("error_code", ce.CodeUUIDGenerationFailed),
					logger.NewField("error", err),
				)
			} else {
				requestID = uuid.String()
			}
		}

		ctx.Writer.Header().Set("X-Request-ID", requestID)
		ctx.Request = ctx.Request.WithContext(
			context.WithValue(
				ctx.Request.Context(),
				constants.CtxKeyRequestID,
				requestID,
			),
		)
		ctx.Next()
	}
}
