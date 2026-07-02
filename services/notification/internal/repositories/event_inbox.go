package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/notification/internal/models"
	"github.com/ritchieridanko/apotekly/services/notification/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type EventInboxRepository interface {
	Create(ctx context.Context, data *models.CreateEventInbox) (err *ce.Error)
	Claim(ctx context.Context) (ei *models.EventInbox, err *ce.Error)
	Retry(ctx context.Context, eventID uuid.UUID, data *models.RetryEventInbox) (err *ce.Error)
	Complete(ctx context.Context, eventID uuid.UUID) (err *ce.Error)
}

type eventInboxRepository struct {
	database database.EventInboxDatabase
}

func NewEventInboxRepository(db database.EventInboxDatabase) EventInboxRepository {
	return &eventInboxRepository{database: db}
}

func (r *eventInboxRepository) Create(ctx context.Context, data *models.CreateEventInbox) *ce.Error {
	return r.database.Create(ctx, data)
}

func (r *eventInboxRepository) Claim(ctx context.Context) (*models.EventInbox, *ce.Error) {
	return r.database.Claim(ctx)
}

func (r *eventInboxRepository) Retry(ctx context.Context, eventID uuid.UUID, data *models.RetryEventInbox) *ce.Error {
	return r.database.Retry(ctx, eventID, data)
}

func (r *eventInboxRepository) Complete(ctx context.Context, eventID uuid.UUID) *ce.Error {
	return r.database.Complete(ctx, eventID)
}
