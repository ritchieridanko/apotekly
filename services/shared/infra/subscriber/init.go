package subscriber

import (
	"strings"

	"github.com/ritchieridanko/apotekly/services/shared/configs"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Init(cfg *configs.Subscriber, appName, brokers string, l *zap.Logger) *kafka.Reader {
	r := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:        strings.Split(brokers, ","),
			GroupID:        appName,
			Topic:          cfg.Name,
			MaxBytes:       cfg.MaxBytes,
			MaxWait:        cfg.MaxWait,
			CommitInterval: cfg.CommitInterval,
		},
	)

	l.Sugar().Infof("[SUBSCRIBER] initialized (topic=%s, brokers=%s)", cfg.Name, brokers)
	return r
}
