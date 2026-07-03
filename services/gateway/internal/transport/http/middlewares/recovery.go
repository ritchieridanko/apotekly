package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

func Recovery(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				l.Error(
					ctx.Request.Context(),
					"PANIC RECOVERED",
					logger.NewField("method", ctx.Request.Method),
					logger.NewField("path", ctx.Request.URL.Path),
					logger.NewField("panic", fmt.Sprintf("%v", r)),
					logger.NewField("stack_trace", string(debug.Stack())),
				)

				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"status":  http.StatusInternalServerError,
						"message": ce.MsgInternalServer,
					},
				)
			}
		}()

		ctx.Next()
	}
}
