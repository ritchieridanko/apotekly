package middlewares

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/ritchieridanko/apotekly/services/notification/internal/transport/event/handlers"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/segmentio/kafka-go"
)

func Recovery(l *logger.Logger) EventMiddleware {
	return func(next handlers.EventHandler) handlers.EventHandler {
		return func(ctx context.Context, msg kafka.Message) (err *ce.Error) {
			defer func() {
				if r := recover(); r != nil {
					l.Error(
						ctx,
						"PANIC RECOVERED",
						logger.NewField("topic", msg.Topic),
						logger.NewField("partition", msg.Partition),
						logger.NewField("offset", msg.Offset),
						logger.NewField("key", string(msg.Key)),
						logger.NewField("panic", fmt.Sprintf("%v", r)),
						logger.NewField("stack_trace", string(debug.Stack())),
					)
					err = ce.NewError(ce.CodePanicOccurred, ce.MsgInternalServer, fmt.Errorf("%v", r))
				}
			}()

			return next(ctx, msg)
		}
	}
}
