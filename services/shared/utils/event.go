package utils

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/constants"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

type EvtHeader []kafka.Header

func NewEvtHeader(ctx context.Context) EvtHeader {
	headers := EvtHeader{
		kafka.Header{
			Key:   constants.EvtHeaderKeyRequestID,
			Value: []byte(CtxRequestID(ctx)),
		},
		kafka.Header{
			Key:   "content-type",
			Value: []byte("application/x-protobuf"),
		},
	}

	otel.GetTextMapPropagator().Inject(ctx, &headers)
	return headers
}

func (eh EvtHeader) Get(key string) string {
	for _, header := range eh {
		if header.Key == key {
			return string(header.Value)
		}
	}
	return ""
}

func (eh EvtHeader) Set(key, value string) {
	eh = append(
		eh,
		kafka.Header{
			Key:   key,
			Value: []byte(value),
		},
	)
}

func (eh EvtHeader) Keys() []string {
	keys := make([]string, len(eh))
	for i, header := range eh {
		keys[i] = header.Key
	}
	return keys
}

// Get Request ID from Kafka Message
func EvtMsgRequestID(msg kafka.Message) string {
	for _, header := range msg.Headers {
		if header.Key == constants.EvtHeaderKeyRequestID {
			return string(header.Value)
		}
	}
	return ""
}
