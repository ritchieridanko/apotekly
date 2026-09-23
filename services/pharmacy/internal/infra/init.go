package infra

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/apotekly/services/pharmacy/configs"
	"github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/services"
	"github.com/ritchieridanko/apotekly/services/shared/infra/storage"
	"github.com/ritchieridanko/apotekly/services/shared/infra/tracer"
	"go.uber.org/zap"
)

type Infra struct {
	config   *configs.Config
	database *pgxpool.Pool
	logger   *zap.Logger
	storage  *cloudinary.Cloudinary
	tracer   *tracer.Tracer
	as       *services.AuthService
	us       *services.UserService
}

func Init(cfg *configs.Config) (*Infra, error) {
	l, err := logger.Init(cfg.App.Env)
	if err != nil {
		return nil, err
	}

	db, err := database.Init(&cfg.Database, l)
	if err != nil {
		return nil, err
	}

	s, err := storage.Init(&cfg.Storage, l)
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
	us, err := services.NewUserService(&cfg.Service.User, l)
	if err != nil {
		return nil, err
	}

	return &Infra{
		config:   cfg,
		database: db,
		logger:   l,
		storage:  s,
		tracer:   t,
		as:       as,
		us:       us,
	}, nil
}

func (i *Infra) Database() *pgxpool.Pool {
	return i.database
}

func (i *Infra) Logger() *zap.Logger {
	return i.logger
}

func (i *Infra) Storage() *cloudinary.Cloudinary {
	return i.storage
}

func (i *Infra) AuthService() *services.AuthService {
	return i.as
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
	if err := i.us.Close(); err != nil {
		return fmt.Errorf("failed to close user service connection: %w", err)
	}

	i.database.Close()
	return nil
}
