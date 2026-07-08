package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/handlers"
	"github.com/ritchieridanko/apotekly/services/gateway/internal/transport/http/middlewares"
	"github.com/ritchieridanko/apotekly/services/shared/infra/logger"
	"github.com/ritchieridanko/apotekly/services/shared/utils/jwt"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	router *gin.Engine
}

func Init(appName string, j *jwt.JWT, l *logger.Logger, ah *handlers.AuthHandler) *Router {
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
		auth.POST("/signout", middlewares.Auth(j), ah.SignOut)
		auth.POST("/refresh", ah.RotateAuthToken)

		// Emails
		email := auth.Group("/email")
		{
			// Availability
			email.GET("/available", ah.IsEmailAvailable)

			// Verifications
			verification := email.Group("/verification")
			{
				// Resend
				verification.POST("", middlewares.Auth(j), ah.ResendVerification)

				// Confirm
				verification.POST("/confirm", middlewares.Auth(j), ah.VerifyEmail)
			}

			// Changes
			change := email.Group("/change")
			{
				// Change
				change.POST("", middlewares.Auth(j), ah.ChangeEmail)

				// Confirm
				change.POST("/confirm", middlewares.Auth(j), ah.ConfirmEmailChange)
			}
		}

		// Passwords
		password := auth.Group("/password")
		{
			// Update
			password.PATCH("", middlewares.Auth(j), ah.ChangePassword)

			// Resets
			reset := password.Group("/reset")
			{
				// Reset
				reset.POST("", ah.ResetPassword)

				// Confirm
				reset.POST("/confirm", ah.ConfirmPasswordReset)

				// Validity
				reset.GET("/valid", ah.IsPasswordResetTokenValid)
			}
		}
	}

	return &Router{router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
