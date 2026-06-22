package publisher

import (
	"strings"

	"github.com/ritchieridanko/apotekly/services/shared/configs"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func Init(cfg *configs.Publisher, brokers string, l *zap.Logger) *kafka.Writer {
	w := kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:      strings.Split(brokers, ","),
			Topic:        cfg.Name,
			Balancer:     balancer(cfg.Balancer),
			BatchSize:    cfg.BatchSize,
			BatchTimeout: cfg.BatchTimeout,
			RequiredAcks: int(kafka.RequireAll),
			Async:        false,
		},
	)

	l.Sugar().Infof("[PUBLISHER] initialized (topic=%s, balancer=%s, brokers=%s)", cfg.Name, cfg.Balancer, brokers)
	return w
}
