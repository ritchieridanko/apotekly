package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ritchieridanko/apotekly/services/notification/internal/models"
	db "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type EventInboxDatabase interface {
	Create(ctx context.Context, data *models.CreateEventInbox) (err *ce.Error)
	Claim(ctx context.Context) (ei *models.EventInbox, err *ce.Error)
	Retry(ctx context.Context, eventID uuid.UUID, data *models.RetryEventInbox) (err *ce.Error)
	Complete(ctx context.Context, eventID uuid.UUID) (err *ce.Error)
}

type eventInboxDatabase struct {
	database *db.Database
}

func NewEventInboxDatabase(db *db.Database) EventInboxDatabase {
	return &eventInboxDatabase{database: db}
}

func (d *eventInboxDatabase) Create(ctx context.Context, data *models.CreateEventInbox) *ce.Error {
	query := `
		INSERT INTO event_inbox (id, topic, payload)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO NOTHING
	`

	err := d.database.Execute(
		ctx, query,
		data.ID,
		data.Topic,
		data.Payload,
	)
	if err != nil {
		return ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to create event inbox: %w", err),
		)
	}

	return nil
}

func (d *eventInboxDatabase) Claim(ctx context.Context) (*models.EventInbox, *ce.Error) {
	query := `
		SELECT
			id, topic, payload, received_at, completed_at, retry_count,
			next_retry_at, last_attempted_at, last_error
		FROM
			event_inbox
		WHERE
			completed_at IS NULL
			AND next_retry_at <= NOW()
		ORDER BY
			next_retry_at,
			received_at
		LIMIT
			1
		FOR UPDATE SKIP LOCKED
	`

	var ei models.EventInbox
	err := d.database.Query(
		ctx, query,
	).Scan(
		&ei.ID,
		&ei.Topic,
		&ei.Payload,
		&ei.ReceivedAt,
		&ei.CompletedAt,
		&ei.RetryCount,
		&ei.NextRetryAt,
		&ei.LastAttemptedAt,
		&ei.LastError,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBQueryNoRows) {
			return nil, ce.NewError(
				ce.CodeNoPendingEventInbox,
				ce.MsgNoPendingEventInbox,
				fmt.Errorf("failed to claim event inbox: %w", err),
			)
		}
		return nil, ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to claim event inbox: %w", err),
		)
	}

	return &ei, nil
}

func (d *eventInboxDatabase) Retry(ctx context.Context, eventID uuid.UUID, data *models.RetryEventInbox) *ce.Error {
	query := `
		UPDATE
			event_inbox
		SET
			retry_count = retry_count + 1,
			next_retry_at = $2,
			last_attempted_at = NOW(),
			last_error = $3
		WHERE
			id = $1
	`

	err := d.database.Execute(
		ctx, query,
		eventID,
		data.NextRetryAt,
		data.Error,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(
				ce.CodeOrphanedEventInbox,
				ce.MsgOrphanedEventInbox,
				fmt.Errorf("failed to retry event inbox: %w", err),
			)
		}
		return ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to retry event inbox: %w", err),
		)
	}

	return nil
}

func (d *eventInboxDatabase) Complete(ctx context.Context, eventID uuid.UUID) *ce.Error {
	query := "UPDATE event_inbox SET completed_at = NOW() WHERE id = $1"

	err := d.database.Execute(
		ctx, query,
		eventID,
	)
	if err != nil {
		if errors.Is(err, ce.ErrDBAffectNoRows) {
			return ce.NewError(
				ce.CodeOrphanedEventInbox,
				ce.MsgOrphanedEventInbox,
				fmt.Errorf("failed to complete event inbox: %w", err),
			)
		}
		return ce.NewError(
			ce.CodeDBQueryExec,
			ce.MsgInternalServer,
			fmt.Errorf("failed to complete event inbox: %w", err),
		)
	}

	return nil
}
