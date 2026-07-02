package middlewares

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/notification/internal/transport/event/handlers"
	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/segmentio/kafka-go"
)

func Request() EventMiddleware {
	return func(next handlers.EventHandler) handlers.EventHandler {
		return func(ctx context.Context, msg kafka.Message) *ce.Error {
			ctx = context.WithValue(ctx, constants.CtxKeyRequestID, utils.EvtMsgRequestID(msg))
			return next(ctx, msg)
		}
	}
}
