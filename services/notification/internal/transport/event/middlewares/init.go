package middlewares

import "github.com/ritchieridanko/apotekly/services/notification/internal/transport/event/handlers"

type EventMiddleware func(handlers.EventHandler) handlers.EventHandler
