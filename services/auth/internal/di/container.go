package di

import (
	"github.com/ritchieridanko/apotekly/services/auth/configs"
	"github.com/ritchieridanko/apotekly/services/auth/internal/infra"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories/cache"
	"github.com/ritchieridanko/apotekly/services/auth/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/auth/internal/transport/rpc/handlers"
	"github.com/ritchieridanko/apotekly/services/auth/internal/transport/rpc/server"
	"github.com/ritchieridanko/apotekly/services/auth/internal/usecases"
	infcc "github.com/ritchieridanko/apotekly/services/shared/infra/cache"
	infdb "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/publisher"
	"github.com/ritchieridanko/apotekly/services/shared/utils/bcrypt"
	"github.com/ritchieridanko/apotekly/services/shared/utils/jwt"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
)

type Container struct {
	config     *configs.Config
	cache      *infcc.Cache
	database   *infdb.Database
	transactor *infdb.Transactor
	logger     *logger.Logger

	acp   *publisher.Publisher
	aevrp *publisher.Publisher

	acc cache.AuthCache
	tcc cache.TokenCache

	adb database.AuthDatabase
	sdb database.SessionDatabase

	ar repositories.AuthRepository
	sr repositories.SessionRepository
	tr repositories.TokenRepository

	bcrypt    *bcrypt.BCrypt
	jwt       *jwt.JWT
	validator *validator.Validator

	au usecases.AuthUsecase
	su usecases.SessionUsecase

	ah *handlers.AuthHandler

	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	cc := infcc.NewCache(inf.Cache())
	db := infdb.NewDatabase(inf.Database())
	tx := infdb.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Publishers
	acp := publisher.NewPublisher(inf.PublisherAC())
	aevrp := publisher.NewPublisher(inf.PublisherAEVR())

	// Caches
	acc := cache.NewAuthCache(cc)
	tcc := cache.NewTokenCache(cc)

	// Databases
	adb := database.NewAuthDatabase(db)
	sdb := database.NewSessionDatabase(db)

	// Repositories
	ar := repositories.NewAuthRepository(acc, adb)
	sr := repositories.NewSessionRepository(sdb)
	tr := repositories.NewTokenRepository(tcc)

	// Utils
	b := bcrypt.Init(cfg.Auth.BCrypt.Cost)
	j := jwt.Init(cfg.Auth.JWT.Issuer, cfg.Auth.JWT.Secret, cfg.Auth.JWT.Duration)
	v := validator.Init()

	// Usecases
	su := usecases.NewSessionUsecase(cfg.App.Name, cfg.Auth.JWT.Duration, cfg.Auth.Duration.Session, sr, tx, j)
	au := usecases.NewAuthUsecase(cfg.App.Name, cfg.Auth.Duration.Verification, su, ar, tr, tx, acp, aevrp, b, v, l)

	// Handlers
	ah := handlers.NewAuthHandler(au)

	// Server
	srv := server.Init(&cfg.Server, l, ah)

	return &Container{
		config:     cfg,
		cache:      cc,
		database:   db,
		transactor: tx,
		logger:     l,
		acp:        acp,
		aevrp:      aevrp,
		acc:        acc,
		tcc:        tcc,
		adb:        adb,
		sdb:        sdb,
		ar:         ar,
		sr:         sr,
		tr:         tr,
		bcrypt:     b,
		jwt:        j,
		validator:  v,
		au:         au,
		su:         su,
		ah:         ah,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
