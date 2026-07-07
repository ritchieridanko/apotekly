package infra

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/apotekly/services/auth/configs"
	"github.com/ritchieridanko/apotekly/services/shared/infra/cache"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/publisher"
	"github.com/ritchieridanko/apotekly/services/shared/infra/tracer"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Infra struct {
	config   *configs.Config
	cache    *redis.Client
	database *pgxpool.Pool
	logger   *zap.Logger
	tracer   *tracer.Tracer
	acp      *kafka.Writer
	aecrp    *kafka.Writer
	aevrp    *kafka.Writer
	aprrp    *kafka.Writer
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(cfg.App.Env)
	if err != nil {
		return nil, err
	}

	cc, err := cache.Init(&cfg.Cache, l)
	if err != nil {
		return nil, err
	}

	db, err := database.Init(&cfg.Database, l)
	if err != nil {
		return nil, err
	}

	t, err := tracer.Init(cfg.App.Env, cfg.App.Name, cfg.Tracer.Addr, l)
	if err != nil {
		return nil, err
	}

	// Publishers
	acp := publisher.Init(&cfg.Broker.AC, cfg.Broker.Brokers, l)
	aecrp := publisher.Init(&cfg.Broker.AECR, cfg.Broker.Brokers, l)
	aevrp := publisher.Init(&cfg.Broker.AEVR, cfg.Broker.Brokers, l)
	aprrp := publisher.Init(&cfg.Broker.APRR, cfg.Broker.Brokers, l)

	return &Infra{
		config:   cfg,
		cache:    cc,
		database: db,
		logger:   l,
		tracer:   t,
		acp:      acp,
		aecrp:    aecrp,
		aevrp:    aevrp,
		aprrp:    aprrp,
	}, nil
}

func (i *Infra) Cache() *redis.Client {
	return i.cache
}

func (i *Infra) Database() *pgxpool.Pool {
	return i.database
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) PublisherAC() *kafka.Writer {
	return i.acp
}

func (i *Infra) PublisherAECR() *kafka.Writer {
	return i.aecrp
}

func (i *Infra) PublisherAEVR() *kafka.Writer {
	return i.aevrp
}

func (i *Infra) PublisherAPRR() *kafka.Writer {
	return i.aprrp
}

func (i *Infra) Close() error {
	if err := i.cache.Close(); err != nil {
		return fmt.Errorf("failed to close cache: %w", err)
	}
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.tracer.Shutdown(); err != nil {
		return fmt.Errorf("failed to close tracer: %w", err)
	}
	if err := i.acp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", i.acp.Topic, err)
	}
	if err := i.aecrp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", i.aecrp.Topic, err)
	}
	if err := i.aevrp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", i.aevrp.Topic, err)
	}
	if err := i.aprrp.Close(); err != nil {
		return fmt.Errorf("failed to close publisher (topic: %s): %w", i.aprrp.Topic, err)
	}

	i.database.Close()
	return nil
}
