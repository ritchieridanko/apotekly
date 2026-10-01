package di

import (
	"context"

	"github.com/ritchieridanko/apotekly/services/auth/configs"
	"github.com/ritchieridanko/apotekly/services/auth/internal/clients"
	"github.com/ritchieridanko/apotekly/services/auth/internal/infra"
	"github.com/ritchieridanko/apotekly/services/auth/internal/infra/oauth"
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

	goa *oauth.OAuth

	acp   *publisher.Publisher
	aecrp *publisher.Publisher
	aevrp *publisher.Publisher
	aprrp *publisher.Publisher

	pc clients.PharmacyClient

	acc  cache.AuthCache
	oacc cache.OAuthCache
	tcc  cache.TokenCache

	adb  database.AuthDatabase
	oadb database.OAuthDatabase
	sdb  database.SessionDatabase

	ar  repositories.AuthRepository
	oar repositories.OAuthRepository
	sr  repositories.SessionRepository
	tr  repositories.TokenRepository

	bcrypt    *bcrypt.BCrypt
	jwt       *jwt.JWT
	validator *validator.Validator

	au  usecases.AuthUsecase
	oau usecases.OAuthUsecase
	su  usecases.SessionUsecase

	ah  *handlers.AuthHandler
	oah *handlers.OAuthHandler

	server *server.Server
}

func Init(ctx context.Context, cfg *configs.Config, inf *infra.Infra) (*Container, error) {
	// Infra
	cc := infcc.NewCache(inf.Cache())
	db := infdb.NewDatabase(inf.Database())
	tx := infdb.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// OAuth
	goa, err := oauth.NewOAuth(ctx, &cfg.OAuth.Google, nil, inf.Logger())
	if err != nil {
		return nil, err
	}

	// Publishers
	acp := publisher.NewPublisher(inf.PublisherAC())
	aecrp := publisher.NewPublisher(inf.PublisherAECR())
	aevrp := publisher.NewPublisher(inf.PublisherAEVR())
	aprrp := publisher.NewPublisher(inf.PublisherAPRR())

	// Clients
	pc := clients.NewPharmacyClient(inf.PharmacyService().PharmacyClient())

	// Caches
	acc := cache.NewAuthCache(cc)
	oacc := cache.NewOAuthCache(cc)
	tcc := cache.NewTokenCache(cc)

	// Databases
	adb := database.NewAuthDatabase(db)
	oadb := database.NewOAuthDatabase(db)
	sdb := database.NewSessionDatabase(db)

	// Repositories
	ar := repositories.NewAuthRepository(acc, adb)
	oar := repositories.NewOAuthRepository(oacc, oadb)
	sr := repositories.NewSessionRepository(sdb)
	tr := repositories.NewTokenRepository(tcc)

	// Utils
	b := bcrypt.Init(cfg.Auth.BCrypt.Cost)
	j := jwt.Init(cfg.Auth.JWT.Issuer, cfg.Auth.JWT.Secret, cfg.Auth.JWT.Duration)
	v := validator.Init()

	// Usecases
	su := usecases.NewSessionUsecase(cfg.App.Name, cfg.Auth.JWT.Duration, cfg.Auth.Duration.Session, sr, pc, tx, j)
	au := usecases.NewAuthUsecase(
		cfg.App.Name,
		cfg.Auth.Duration.EmailChange,
		cfg.Auth.Duration.PasswordReset,
		cfg.Auth.Duration.Verification,
		su, ar, tr, tx,
		acp, aecrp, aevrp, aprrp,
		b, v, l,
	)
	oau := usecases.NewOAuthUsecase(cfg.App.Name, cfg.Auth.Duration.OAuthState, su, ar, oar, tx, acp, goa, l)

	// Handlers
	ah := handlers.NewAuthHandler(au)
	oah := handlers.NewOAuthHandler(oau)

	// Server
	srv := server.Init(&cfg.Server, l, ah, oah)

	return &Container{
		config:     cfg,
		cache:      cc,
		database:   db,
		transactor: tx,
		logger:     l,
		goa:        goa,
		acp:        acp,
		aecrp:      aecrp,
		aevrp:      aevrp,
		aprrp:      aprrp,
		pc:         pc,
		acc:        acc,
		oacc:       oacc,
		tcc:        tcc,
		adb:        adb,
		oadb:       oadb,
		sdb:        sdb,
		ar:         ar,
		oar:        oar,
		sr:         sr,
		tr:         tr,
		bcrypt:     b,
		jwt:        j,
		validator:  v,
		au:         au,
		oau:        oau,
		su:         su,
		ah:         ah,
		oah:        oah,
		server:     srv,
	}, nil
}

func (c *Container) Server() *server.Server {
	return c.server
}
