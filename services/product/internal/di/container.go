package di

import (
	"github.com/ritchieridanko/apotekly/services/product/configs"
	"github.com/ritchieridanko/apotekly/services/product/internal/clients"
	"github.com/ritchieridanko/apotekly/services/product/internal/infra"
	"github.com/ritchieridanko/apotekly/services/product/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/product/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/product/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/apotekly/services/product/internal/transport/rpc/server"
	"github.com/ritchieridanko/apotekly/services/product/internal/usecases"
	infdb "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/storage"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
)

type Container struct {
	config     *configs.Config
	database   *infdb.Database
	transactor *infdb.Transactor
	logger     *logger.Logger
	storage    *storage.Storage

	pc clients.PharmacyClient

	pdb  database.ProductDatabase
	pidb database.ProductImageDatabase

	pr  repositories.ProductRepository
	pir repositories.ProductImageRepository

	validator *validator.Validator

	pu  usecases.ProductUsecase
	piu usecases.ProductImageUsecase

	ph  *handlers.ProductHandler
	pih *handlers.ProductImageHandler

	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	db := infdb.NewDatabase(inf.Database())
	tx := infdb.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())
	s := storage.NewStorage(inf.Storage())

	// Clients
	pc := clients.NewPharmacyClient(inf.PharmacyService().PharmacyClient())

	// Databases
	pdb := database.NewProductDatabase(db)
	pidb := database.NewProductImageDatabase(db)

	// Repositories
	pr := repositories.NewProductRepository(pdb)
	pir := repositories.NewProductImageRepository(pidb)

	// Utils
	v := validator.Init()

	// Usecases
	piu := usecases.NewProductImageUsecase(cfg.App.Name, int64(cfg.Image.Product.MaxTotal), pr, pir, pc, tx, s, v, l)
	pu := usecases.NewProductUsecase(cfg.App.Name, piu, pr, pc, tx, v, l)

	// Handlers
	ph := handlers.NewProductHandler(pu)
	pih := handlers.NewProductImageHandler(piu)

	// Server
	srv := server.Init(&cfg.Server, l, ph, pih)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		storage:    s,
		pc:         pc,
		pdb:        pdb,
		pidb:       pidb,
		pr:         pr,
		pir:        pir,
		validator:  v,
		pu:         pu,
		piu:        piu,
		ph:         ph,
		pih:        pih,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
