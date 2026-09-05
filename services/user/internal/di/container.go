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
	config     *configs.Config
	database   *infdb.Database
	transactor *infdb.Transactor
	logger     *logger.Logger

	udb database.UserDatabase
	adb database.AddressDatabase

	ur repositories.UserRepository
	ar repositories.AddressRepository

	validator *validator.Validator

	uu usecases.UserUsecase
	au usecases.AddressUsecase

	uh *handlers.UserHandler
	ah *handlers.AddressHandler

	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	db := infdb.NewDatabase(inf.Database())
	tx := infdb.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Databases
	udb := database.NewUserDatabase(db)
	adb := database.NewAddressDatabase(db)

	// Repositories
	ur := repositories.NewUserRepository(udb)
	ar := repositories.NewAddressRepository(adb)

	// Utils
	v := validator.Init()

	// Usecases
	uu := usecases.NewUserUsecase(cfg.App.Name, ur, v, l)
	au := usecases.NewAddressUsecase(cfg.App.Name, ar, tx, v, l)

	// Handlers
	uh := handlers.NewUserHandler(uu)
	ah := handlers.NewAddressHandler(au)

	// Server
	srv := server.Init(&cfg.Server, l, uh, ah)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		udb:        udb,
		adb:        adb,
		ur:         ur,
		ar:         ar,
		validator:  v,
		uu:         uu,
		au:         au,
		uh:         uh,
		ah:         ah,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
