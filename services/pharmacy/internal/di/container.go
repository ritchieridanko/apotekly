package di

import (
	"github.com/ritchieridanko/apotekly/services/pharmacy/configs"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/infra"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/transport/rpc/server"
	"github.com/ritchieridanko/apotekly/services/pharmacy/internal/usecases"
	infdb "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
)

type Container struct {
	config   *configs.Config
	database *infdb.Database
	logger   *logger.Logger

	pdb database.PharmacyDatabase

	pr repositories.PharmacyRepository

	validator *validator.Validator

	pu usecases.PharmacyUsecase

	ph *handlers.PharmacyHandler

	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	db := infdb.NewDatabase(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Databases
	pdb := database.NewPharmacyDatabase(db)

	// Repositories
	pr := repositories.NewPharmacyRepository(pdb)

	// Utils
	v := validator.Init()

	// Usecases
	pu := usecases.NewPharmacyUsecase(cfg.App.Name, pr, v, l)

	// Handlers
	ph := handlers.NewPharmacyHandler(pu)

	// Server
	srv := server.Init(&cfg.Server, l, ph)

	return &Container{
		config:    cfg,
		database:  db,
		logger:    l,
		pdb:       pdb,
		pr:        pr,
		validator: v,
		pu:        pu,
		ph:        ph,
		server:    srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
