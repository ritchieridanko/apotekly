package di

import (
	infdb "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
	"github.com/ritchieridanko/apotekly/services/user/configs"
	"github.com/ritchieridanko/apotekly/services/user/internal/infra"
	"github.com/ritchieridanko/apotekly/services/user/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/user/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/user/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/apotekly/services/user/internal/transport/rpc/server"
	"github.com/ritchieridanko/apotekly/services/user/internal/usecases"
)

type Container struct {
	config   *configs.Config
	database *infdb.Database
	logger   *logger.Logger

	udb database.UserDatabase

	ur repositories.UserRepository

	validator *validator.Validator

	uu usecases.UserUsecase

	uh *handlers.UserHandler

	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	db := infdb.NewDatabase(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Databases
	udb := database.NewUserDatabase(db)

	// Repositories
	ur := repositories.NewUserRepository(udb)

	// Utils
	v := validator.Init()

	// Usecases
	uu := usecases.NewUserUsecase(cfg.App.Name, ur, v, l)

	// Handlers
	uh := handlers.NewUserHandler(uu)

	// Server
	srv := server.Init(&cfg.Server, l, uh)

	return &Container{
		config:    cfg,
		database:  db,
		logger:    l,
		udb:       udb,
		ur:        ur,
		validator: v,
		uu:        uu,
		uh:        uh,
		server:    srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
