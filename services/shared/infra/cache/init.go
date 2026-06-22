package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly/services/shared/configs"
	"go.uber.org/zap"
)

func Init(cfg *configs.Cache, l *zap.Logger) (*redis.Client, error) {
	if cfg.Pass == "" {
		l.Sugar().Warnln("[CACHE] connecting without password...")
	}

	c := redis.NewClient(
		&redis.Options{
			Addr:                  cfg.Addr,
			Password:              cfg.Pass,
			PoolSize:              cfg.PoolSize,
			MinIdleConns:          cfg.MinIdleConns,
			MaxActiveConns:        cfg.MaxActiveConns,
			ConnMaxLifetime:       cfg.ConnMaxLifetime,
			ConnMaxIdleTime:       cfg.ConnMaxIdleTime,
			DialTimeout:           cfg.Timeout.Dial,
			ReadTimeout:           cfg.Timeout.Read,
			WriteTimeout:          cfg.Timeout.Write,
			ContextTimeoutEnabled: true,
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.Ping(ctx).Err(); err != nil {
		c.Close()
		return nil, fmt.Errorf("failed to ping cache: %w", err)
	}

	l.Sugar().Infof("[CACHE] connected (host=%s, port=%d)", cfg.Host, cfg.Port)
	return c, nil
}
