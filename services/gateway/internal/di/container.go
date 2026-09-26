package di

import (
	"github.com/ritchieridanko/apotekly/services/gateway/configs"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/clients"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/infra"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/repositories/cache"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/handlers"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/middlewares"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/router"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/server"
	infcc "github.com/ritchieridanko/apotekly/services/shared/infra/cache"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/cookie"
	"github.com/ritchieridanko/apotekly/services/shared/utils/jwt"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
)

type Container struct {
	config *configs.Config
	cache  *infcc.Cache
	logger *logger.Logger

	ac  clients.AuthClient
	phc clients.PharmacyClient
	prc clients.ProductClient
	uc  clients.UserClient
	uac clients.AddressClient

	rlcc cache.RateLimiterCache

	cookie    *cookie.Cookie
	jwt       *jwt.JWT
	validator *validator.Validator

	ah  *handlers.AuthHandler
	phh *handlers.PharmacyHandler
	prh *handlers.ProductHandler
	uh  *handlers.UserHandler
	uah *handlers.AddressHandler

	rlm *middlewares.RateLimiterMiddleware

	router *router.Router
	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	cc := infcc.NewCache(inf.Cache())
	l := logger.NewLogger(inf.Logger())

	// Clients
	ac := clients.NewAuthClient(inf.AuthService().AuthClient())
	phc := clients.NewPharmacyClient(inf.PharmacyService().PharmacyClient())
	prc := clients.NewProductClient(inf.ProductService().ProductClient())
	uc := clients.NewUserClient(inf.UserService().UserClient())
	uac := clients.NewAddressClient(inf.UserService().AddressClient())

	// Caches
	rlcc := cache.NewRateLimiterCache(cc)

	// Utils
	c := cookie.Init(cfg.App.Env, "")
	j := jwt.Init(cfg.JWT.Secret, cfg.JWT.Secret, cfg.JWT.Duration)
	v := validator.Init()

	// Handlers
	ah := handlers.NewAuthHandler(ac, v, c)
	phh := handlers.NewPharmacyHandler(phc, ac, v, c, l)
	prh := handlers.NewProductHandler(prc)
	uh := handlers.NewUserHandler(uc)
	uah := handlers.NewAddressHandler(uac)

	// Middlewares
	rlm := middlewares.NewRateLimiterMiddleware(rlcc, l)

	// Router
	r := router.Init(cfg.App.Name, cfg.Client.Addr, j, l, rlm, ah, phh, prh, uh, uah)

	// Server
	srv := server.Init(&cfg.Server, r, l)

	return &Container{
		config:    cfg,
		cache:     cc,
		logger:    l,
		ac:        ac,
		phc:       phc,
		prc:       prc,
		uc:        uc,
		uac:       uac,
		rlcc:      rlcc,
		cookie:    c,
		jwt:       j,
		validator: v,
		ah:        ah,
		phh:       phh,
		prh:       prh,
		uh:        uh,
		uah:       uah,
		rlm:       rlm,
		router:    r,
		server:    srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
