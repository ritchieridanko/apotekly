package infra

import (
	"fmt"

	"github.com/ritchieridanko/apotekly/services/gateway/configs"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/services"
	"github.com/ritchieridanko/apotekly/services/shared/infra/tracer"
	"go.uber.org/zap"
)

type Infra struct {
	config *configs.Config
	logger *zap.Logger
	tracer *tracer.Tracer
	as     *services.AuthService
	ps     *services.PharmacyService
	us     *services.UserService
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(cfg.App.Env)
	if err != nil {
		return nil, err
	}

	t, err := tracer.Init(cfg.App.Env, cfg.App.Name, cfg.Tracer.Addr, l)
	if err != nil {
		return nil, err
	}

	// Services
	as, err := services.NewAuthService(&cfg.Service.Auth, l)
	if err != nil {
		return nil, err
	}
	ps, err := services.NewPharmacyService(&cfg.Service.Pharmacy, l)
	if err != nil {
		return nil, err
	}
	us, err := services.NewUserService(&cfg.Service.User, l)
	if err != nil {
		return nil, err
	}

	return &Infra{
		config: cfg,
		logger: l,
		tracer: t,
		as:     as,
		ps:     ps,
		us:     us,
	}, nil
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) AuthService() *services.AuthService {
	return i.as
}

func (i *Infra) PharmacyService() *services.PharmacyService {
	return i.ps
}

func (i *Infra) UserService() *services.UserService {
	return i.us
}

func (i *Infra) Close() error {
	if err := i.logger.Sync(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}
	if err := i.tracer.Shutdown(); err != nil {
		return fmt.Errorf("failed to close tracer: %w", err)
	}
	if err := i.as.Close(); err != nil {
		return fmt.Errorf("failed to close auth service connection: %w", err)
	}
	if err := i.ps.Close(); err != nil {
		return fmt.Errorf("failed to close pharmacy service connection: %w", err)
	}
	if err := i.us.Close(); err != nil {
		return fmt.Errorf("failed to close user service connection: %w", err)
	}
	return nil
}
