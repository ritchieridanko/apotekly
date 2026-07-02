package usecases

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/notification/internal/models"
	"github.com/ritchieridanko/apotekly/services/notification/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"go.opentelemetry.io/otel"
)

type EventInboxUsecase interface {
	StoreEvent(ctx context.Context, data *models.Event) (err *ce.Error)
}

type eventInboxUsecase struct {
	appName string
	eir     repositories.EventInboxRepository
}

func NewEventInboxUsecase(appName string, eir repositories.EventInboxRepository) EventInboxUsecase {
	return &eventInboxUsecase{appName: appName, eir: eir}
}

func (u *eventInboxUsecase) StoreEvent(ctx context.Context, data *models.Event) *ce.Error {
	ctx, span := otel.Tracer(u.appName).Start(ctx, "event_inbox.usecase.StoreEvent")
	defer span.End()

	return u.eir.Create(
		ctx,
		&models.CreateEventInbox{
			ID:      data.ID,
			Topic:   data.Topic,
			Payload: data.Payload,
		},
	)
}
