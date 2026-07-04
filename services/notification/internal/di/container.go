package di

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/apotekly/services/notification/configs"
	"github.com/ritchieridanko/apotekly/services/notification/internal/channels"
	"github.com/ritchieridanko/apotekly/services/notification/internal/dispatcher"
	"github.com/ritchieridanko/apotekly/services/notification/internal/infra"
	"github.com/ritchieridanko/apotekly/services/notification/internal/processor"
	"github.com/ritchieridanko/apotekly/services/notification/internal/repositories"
	"github.com/ritchieridanko/apotekly/services/notification/internal/repositories/database"
	"github.com/ritchieridanko/apotekly/services/notification/internal/transport/event/handlers"

	"github.com/ritchieridanko/apotekly/services/notification/internal/transport/event/middlewares"
	"github.com/ritchieridanko/apotekly/services/notification/internal/usecases"
	infdb "github.com/ritchieridanko/apotekly/services/shared/infra/database"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/infra/mailer"
	"github.com/ritchieridanko/apotekly/services/shared/infra/subscriber"
)

type Container struct {
	config     *configs.Config
	database   *infdb.Database
	transactor *infdb.Transactor
	logger     *logger.Logger
	mailer     *mailer.Mailer

	acs   *subscriber.Subscriber
	aevrs *subscriber.Subscriber

	ec channels.EmailChannel

	eidb database.EventInboxDatabase

	eir repositories.EventInboxRepository

	au  usecases.AuthUsecase
	eiu usecases.EventInboxUsecase

	aeh *handlers.AuthEventHandler

	lem  middlewares.EventMiddleware
	rqem middlewares.EventMiddleware
	rvem middlewares.EventMiddleware
	tem  middlewares.EventMiddleware

	ed *dispatcher.Dispatcher

	ep1 *processor.Processor
	ep2 *processor.Processor
}

func Init(cfg *configs.Config, inf *infra.Infra) (*Container, error) {
	// Infra
	db := infdb.NewDatabase(inf.Database())
	tx := infdb.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())
	m := mailer.NewMailer(inf.Mailer())

	// Subscribers
	acs := subscriber.NewSubscriber(inf.SubscriberAC(), l)
	aevrs := subscriber.NewSubscriber(inf.SubscriberAEVR(), l)

	// Channels
	ec, err := channels.NewEmailChannel(cfg.Client.Addr, cfg.Mailer.From, cfg.App.LogoURL, m)
	if err != nil {
		return nil, err
	}

	// Databases
	eidb := database.NewEventInboxDatabase(db)

	// Repositories
	eir := repositories.NewEventInboxRepository(eidb)

	// Usecases
	au := usecases.NewAuthUsecase(cfg.App.Name, ec)
	eiu := usecases.NewEventInboxUsecase(cfg.App.Name, eir)

	// Event Handlers
	aeh := handlers.NewAuthEventHandler(eiu)

	// Event Middlewares
	lem := middlewares.Logging(l)
	rqem := middlewares.Request()
	rvem := middlewares.Recovery(l)
	tem := middlewares.Tracing()

	// Event Dispatcher
	ed := dispatcher.Init(au, eir, tx)

	// Event Processors
	ep1 := processor.Init("1", ed, l)
	ep2 := processor.Init("2", ed, l)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		mailer:     m,
		acs:        acs,
		aevrs:      aevrs,
		ec:         ec,
		eidb:       eidb,
		eir:        eir,
		au:         au,
		eiu:        eiu,
		aeh:        aeh,
		lem:        lem,
		rqem:       rqem,
		rvem:       rvem,
		tem:        tem,
		ed:         ed,
		ep1:        ep1,
		ep2:        ep2,
	}, nil
}

func (c *Container) RunSubscriberAC(ctx context.Context) error {
	return c.acs.Listen(
		ctx,
		c.rqem(c.rvem(c.tem(c.lem(
			c.aeh.HandleAC,
		)))),
	)
}

func (c *Container) RunSubscriberAEVR(ctx context.Context) error {
	return c.aevrs.Listen(
		ctx,
		c.rqem(c.rvem(c.tem(c.lem(
			c.aeh.HandleAEVR,
		)))),
	)
}

func (c *Container) RunEventProcessor1(ctx context.Context) error {
	return c.ep1.Run(ctx)
}

func (c *Container) RunEventProcessor2(ctx context.Context) error {
	return c.ep2.Run(ctx)
}

func (c *Container) Close() error {
	if err := c.mailer.Close(); err != nil {
		return fmt.Errorf("failed to close mailer: %w", err)
	}
	return nil
}
