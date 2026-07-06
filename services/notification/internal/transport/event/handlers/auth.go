package handlers

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/notification/internal/models"
	"github.com/ritchieridanko/apotekly/services/notification/internal/usecases"
	"github.com/ritchieridanko/apotekly/services/shared/contract/events/v1"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"
)

type AuthEventHandler struct {
	eiu usecases.EventInboxUsecase
}

func NewAuthEventHandler(eiu usecases.EventInboxUsecase) *AuthEventHandler {
	return &AuthEventHandler{eiu: eiu}
}

func (eh *AuthEventHandler) HandleAC(ctx context.Context, msg kafka.Message) *ce.Error {
	var evt events.AuthCreated
	if err := proto.Unmarshal(msg.Value, &evt); err != nil {
		return ce.NewError(ce.CodeProtobufParsingFailed, ce.MsgInternalServer, err)
	}

	evtIDField := logger.NewField("event_id", evt.GetId())

	payload := models.EventAC{
		ID:                utils.ToUUID(evt.GetId()),
		AuthID:            evt.GetAuthId(),
		Email:             evt.GetEmail(),
		Role:              evt.GetRole(),
		EmailVerifiedAt:   utils.ToTime(evt.GetEmailVerifiedAt()),
		Session:           evt.Session,
		VerificationToken: evt.VerificationToken,
		CreatedAt:         utils.ToTime(evt.GetCreatedAt()),
	}
	rm, err := utils.ToJSONRawMessage(payload)
	if err != nil {
		return ce.NewError(
			ce.CodeJSONRawEncodingFailed,
			ce.MsgInternalServer,
			err,
			evtIDField,
		)
	}

	storeErr := eh.eiu.StoreEvent(
		ctx,
		&models.Event{
			ID:      payload.ID,
			Topic:   msg.Topic,
			Payload: rm,
		},
	)
	if storeErr != nil {
		return storeErr.Append(evtIDField)
	}
	return nil
}

func (eh *AuthEventHandler) HandleAECR(ctx context.Context, msg kafka.Message) *ce.Error {
	var evt events.AuthEmailChangeRequested
	if err := proto.Unmarshal(msg.Value, &evt); err != nil {
		return ce.NewError(ce.CodeProtobufParsingFailed, ce.MsgInternalServer, err)
	}

	evtIDField := logger.NewField("event_id", evt.GetId())

	payload := models.EventAECR{
		ID:        utils.ToUUID(evt.GetId()),
		AuthID:    evt.GetAuthId(),
		OldEmail:  evt.GetOldEmail(),
		NewEmail:  evt.GetNewEmail(),
		Role:      evt.GetRole(),
		Token:     evt.GetToken(),
		CreatedAt: utils.ToTime(evt.GetCreatedAt()),
	}
	rm, err := utils.ToJSONRawMessage(payload)
	if err != nil {
		return ce.NewError(
			ce.CodeJSONRawEncodingFailed,
			ce.MsgInternalServer,
			err,
			evtIDField,
		)
	}

	storeErr := eh.eiu.StoreEvent(
		ctx,
		&models.Event{
			ID:      payload.ID,
			Topic:   msg.Topic,
			Payload: rm,
		},
	)
	if storeErr != nil {
		return storeErr.Append(evtIDField)
	}
	return nil
}

func (eh *AuthEventHandler) HandleAEVR(ctx context.Context, msg kafka.Message) *ce.Error {
	var evt events.AuthEmailVerificationRequested
	if err := proto.Unmarshal(msg.Value, &evt); err != nil {
		return ce.NewError(ce.CodeProtobufParsingFailed, ce.MsgInternalServer, err)
	}

	evtIDField := logger.NewField("event_id", evt.GetId())

	payload := models.EventAEVR{
		ID:        utils.ToUUID(evt.GetId()),
		AuthID:    evt.GetAuthId(),
		Email:     evt.GetEmail(),
		Role:      evt.GetRole(),
		Token:     evt.GetToken(),
		CreatedAt: utils.ToTime(evt.GetCreatedAt()),
	}
	rm, err := utils.ToJSONRawMessage(payload)
	if err != nil {
		return ce.NewError(
			ce.CodeJSONRawEncodingFailed,
			ce.MsgInternalServer,
			err,
			evtIDField,
		)
	}

	storeErr := eh.eiu.StoreEvent(
		ctx,
		&models.Event{
			ID:      payload.ID,
			Topic:   msg.Topic,
			Payload: rm,
		},
	)
	if storeErr != nil {
		return storeErr.Append(evtIDField)
	}
	return nil
}
