package dispatcher

import (
	"context"
	"time"

	"github.com/ritchieridanko/apotekly/services/notification/internal/models"
	"github.com/ritchieridanko/apotekly/services/notification/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/notification/internal/usecases"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

const (
	maxRetryBackoff time.Duration = 2 * time.Minute
)

type Dispatcher struct {
	au         usecases.AuthUsecase
	eir        repositories.EventInboxRepository
	transactor *database.Transactor
}

func Init(au usecases.AuthUsecase, eir repositories.EventInboxRepository, tx *database.Transactor) *Dispatcher {
	return &Dispatcher{au: au, eir: eir, transactor: tx}
}

func (d *Dispatcher) Dispatch(ctx context.Context) *ce.Error {
	return d.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		// Event Claiming
		ei, err := d.eir.Claim(ctx)
		if err != nil {
			return err
		}

		evtIDField := logger.NewField("event_id", ei.ID.String())
		evtTopicField := logger.NewField("event_topic", ei.Topic)

		// Event Processing
		switch ei.Topic {
		case "auth.created":
			data, convertErr := utils.FromJSONRawMessage[models.EventAC](ei.Payload)
			if convertErr != nil {
				return ce.NewError(
					ce.CodeJSONUnmarshallingFailed,
					ce.MsgInternalServer,
					convertErr,
					evtIDField,
					evtTopicField,
				)
			}

			err = d.au.ProcessEventAC(ctx, data)
		case "auth.email.change.requested":
			data, convertErr := utils.FromJSONRawMessage[models.EventAECR](ei.Payload)
			if convertErr != nil {
				return ce.NewError(
					ce.CodeJSONUnmarshallingFailed,
					ce.MsgInternalServer,
					convertErr,
					evtIDField,
					evtTopicField,
				)
			}

			err = d.au.ProcessEventAECR(ctx, data)
		case "auth.email.verification.requested":
			data, convertErr := utils.FromJSONRawMessage[models.EventAEVR](ei.Payload)
			if convertErr != nil {
				return ce.NewError(
					ce.CodeJSONUnmarshallingFailed,
					ce.MsgInternalServer,
					convertErr,
					evtIDField,
					evtTopicField,
				)
			}

			err = d.au.ProcessEventAEVR(ctx, data)
		default:
			return ce.NewError(
				ce.CodeEventTopicNotRegistered,
				ce.MsgInternalServer,
				nil,
				evtIDField,
				evtTopicField,
			)
		}

		// Set retrying if Event Processing fails
		if err != nil {
			backoff := min(
				time.Second*time.Duration(1<<ei.RetryCount),
				maxRetryBackoff,
			)

			retryErr := d.eir.Retry(
				ctx,
				ei.ID,
				&models.RetryEventInbox{
					NextRetryAt: time.Now().UTC().Add(backoff),
					Error:       err.Message(),
				},
			)
			if retryErr != nil {
				return retryErr.Append(evtIDField, evtTopicField)
			}
			return err.Append(evtIDField, evtTopicField)
		}

		// Set completed if Event Processing is OK
		if err := d.eir.Complete(ctx, ei.ID); err != nil {
			return err.Append(evtIDField, evtTopicField)
		}
		return nil
	})
}
