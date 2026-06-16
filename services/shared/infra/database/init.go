package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/apotekly/services/shared/configs"
	"go.uber.org/zap"
)

func Init(cfg *configs.Database, l *zap.Logger) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to build database config: %w", err)
	}

	conf.MaxConns = int32(cfg.MaxConns)
	conf.MinConns = int32(cfg.MinConns)
	conf.MaxConnLifetime = cfg.MaxConnLifetime
	conf.MaxConnIdleTime = cfg.MaxConnIdleTime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}
	if err := p.Ping(ctx); err != nil {
		p.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	l.Sugar().Infof("[DATABASE] connected (host=%s, port=%d, name=%s)", cfg.Host, cfg.Port, cfg.Name)
	return p, nil
}
