package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/handlers"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/middlewares"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	router *gin.Engine
}

func Init(appName string, l *logger.Logger, ah *handlers.AuthHandler) *Router {
	r := gin.New()
	r.ContextWithFallback = true

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "OK",
			},
		)
	})

	v1 := r.Group(
		"/v1",
		middlewares.Request(l),
		middlewares.Recovery(l),
		otelgin.Middleware(appName),
		middlewares.Logging(l),
	)

	// AUTH ENDPOINTS
	auth := v1.Group("/auth")
	{
		auth.POST("/signup", ah.SignUp)
		auth.POST("/signin", ah.SignIn)
		auth.POST("/refresh", ah.RotateAuthToken)

		// Emails
		email := auth.Group("/email")
		{
			// Availability
			email.GET("/available", ah.IsEmailAvailable)
		}
	}

	return &Router{router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
