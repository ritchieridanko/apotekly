package di

import (
	"github.com/ritchieridanko/apotekly/services/gateway/configs"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/clients"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/infra"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/handlers"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/router"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/server"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/cookie"
	"github.com/ritchieridanko/apotekly/services/shared/utils/jwt"
	"github.com/ritchieridanko/apotekly/services/shared/utils/validator"
)

type Container struct {
	config *configs.Config
	logger *logger.Logger

	ac  clients.AuthClient
	phc clients.PharmacyClient
	prc clients.ProductClient
	uc  clients.UserClient
	uac clients.AddressClient

	cookie    *cookie.Cookie
	jwt       *jwt.JWT
	validator *validator.Validator

	ah  *handlers.AuthHandler
	phh *handlers.PharmacyHandler
	prh *handlers.ProductHandler
	uh  *handlers.UserHandler
	uah *handlers.AddressHandler

	router *router.Router
	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	l := logger.NewLogger(inf.Logger())

	// Clients
	ac := clients.NewAuthClient(inf.AuthService().AuthClient())
	phc := clients.NewPharmacyClient(inf.PharmacyService().PharmacyClient())
	prc := clients.NewProductClient(inf.ProductService().ProductClient())
	uc := clients.NewUserClient(inf.UserService().UserClient())
	uac := clients.NewAddressClient(inf.UserService().AddressClient())

	// Utils
	c := cookie.Init(cfg.App.Env, "")
	j := jwt.Init(cfg.JWT.Secret, cfg.JWT.Secret, cfg.JWT.Duration)
	v := validator.Init()

	// Handlers
	ah := handlers.NewAuthHandler(ac, v, c)
	phh := handlers.NewPharmacyHandler(phc)
	prh := handlers.NewProductHandler(prc)
	uh := handlers.NewUserHandler(uc)
	uah := handlers.NewAddressHandler(uac)

	// Router
	r := router.Init(cfg.App.Name, cfg.Client.Addr, j, l, ah, phh, prh, uh, uah)

	// Server
	srv := server.Init(&cfg.Server, r, l)

	return &Container{
		config:    cfg,
		logger:    l,
		ac:        ac,
		phc:       phc,
		prc:       prc,
		uc:        uc,
		uac:       uac,
		cookie:    c,
		jwt:       j,
		validator: v,
		ah:        ah,
		phh:       phh,
		prh:       prh,
		uh:        uh,
		uah:       uah,
		router:    r,
		server:    srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
