package middlewares

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/notification/internal/transport/event/handlers"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

func Tracing() EventMiddleware {
	return func(next handlers.EventHandler) handlers.EventHandler {
		return func(ctx context.Context, msg kafka.Message) *ce.Error {
			ctx = otel.GetTextMapPropagator().Extract(
				ctx,
				utils.EvtHeader(msg.Headers),
			)
			return next(ctx, msg)
		}
	}
}
