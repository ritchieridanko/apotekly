package middlewares

import (
	"context"
	"time"

	"github.com/ritchieridanko/apotekly/services/notification/internal/transport/event/handlers"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/segmentio/kafka-go"
)

func Logging(l *logger.Logger) EventMiddleware {
	return func(next handlers.EventHandler) handlers.EventHandler {
		return func(ctx context.Context, msg kafka.Message) *ce.Error {
			start := time.Now()
			err := next(ctx, msg)

			// Check if Event Handling is OK
			fields := []logger.Field{
				logger.NewField("topic", msg.Topic),
				logger.NewField("partition", msg.Partition),
				logger.NewField("offset", msg.Offset),
				logger.NewField("key", string(msg.Key)),
				logger.NewField("latency", time.Since(start).String()),
			}
			if err == nil {
				l.Info(ctx, "EVENT HANDLING OK", fields...)
				return nil
			}

			// Check if Event Handling fails
			fields = append(fields, err.Fields()...)
			fields = append(
				fields,
				logger.NewField("error_code", err.Code()),
				logger.NewField("error", err.Error()),
			)

			l.Error(ctx, "EVENT HANDLING ERROR", fields...)
			return err
		}
	}
}
