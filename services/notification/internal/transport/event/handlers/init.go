package handlers

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
	"github.com/segmentio/kafka-go"
)

type EventHandler func(context.Context, kafka.Message) *ce.Error
