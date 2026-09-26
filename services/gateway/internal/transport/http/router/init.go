package router

import (
	"net/http"
	"time"

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

func Init(
	appName,
	clientAddr string,
	j *jwt.JWT,
	l *logger.Logger,
	rlm *middlewares.RateLimiterMiddleware,
	ah *handlers.AuthHandler,
	phh *handlers.PharmacyHandler,
	prh *handlers.ProductHandler,
	uh *handlers.UserHandler,
	uah *handlers.AddressHandler,
) *Router {
	r := gin.New()
	r.ContextWithFallback = true

	r.Use(middlewares.CORS(clientAddr))

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
		// Sign Up
		auth.POST("/signup", rlm.LimitByIPAddress("auth:su", 1*time.Minute, 10), ah.SignUp)

		// Sign In
		auth.POST("/signin", rlm.LimitByIPAddress("auth:si", 1*time.Minute, 10), ah.SignIn)

		// Sign Out
		auth.POST("/signout", middlewares.Auth(j), rlm.LimitByAuthID("auth:so", 1*time.Minute, 10), ah.SignOut)

		// Rotate Token
		auth.POST("/refresh", rlm.LimitByIPAddress("auth:r", 1*time.Minute, 20), ah.RotateAuthToken)

		// Emails
		email := auth.Group("/email")
		{
			// Check Availability
			email.GET("/available", rlm.LimitByIPAddress("auth:eac", 1*time.Minute, 20), ah.IsEmailAvailable)

			// Verifications
			verification := email.Group("/verification", middlewares.Auth(j), rlm.LimitByAuthID("auth:ev", 1*time.Minute, 10))
			{
				// Resend
				verification.POST("", ah.ResendVerification)

				// Confirm
				verification.POST("/confirm", ah.VerifyEmail)
			}

			// Changes
			change := email.Group("/change", middlewares.Auth(j), rlm.LimitByAuthID("auth:ec", 1*time.Minute, 10))
			{
				// Request
				change.POST("", ah.ChangeEmail)

				// Confirm
				change.POST("/confirm", ah.ConfirmEmailChange)
			}
		}

		// Passwords
		password := auth.Group("/password")
		{
			// Update
			password.PATCH("", middlewares.Auth(j), rlm.LimitByAuthID("auth:pc", 1*time.Minute, 10), ah.ChangePassword)

			// Resets
			reset := password.Group("/reset", rlm.LimitByIPAddress("auth:pr", 1*time.Minute, 20))
			{
				// Request
				reset.POST("", ah.ResetPassword)

				// Confirm
				reset.POST("/confirm", ah.ConfirmPasswordReset)

				// Check Validity
				reset.GET("/valid", ah.IsPasswordResetTokenValid)
			}
		}
	}

	// PHARMACY ENDPOINTS
	pharmacy := v1.Group("/pharmacies")
	{
		// Create
		pharmacy.POST("", middlewares.Auth(j), rlm.LimitByAuthID("phar:pc", 1*time.Minute, 5), phh.CreatePharmacy)

		// Fetch
		pharmacy.GET("/:pharmacy_id", rlm.LimitByIPAddress("phar:pid", 1*time.Minute, 50), phh.GetPharmacyByID)

		// Fetch All
		pharmacy.GET("", rlm.LimitByIPAddress("phar:pfa", 1*time.Minute, 50), middlewares.AuthOptional(j), phh.GetAllPharmacies)

		// Me
		me := pharmacy.Group("/me", middlewares.Auth(j))
		{
			// Fetch
			me.GET("", rlm.LimitByAuthID("phar:me", 1*time.Minute, 50), phh.GetMe)

			// Update
			me.PATCH("", rlm.LimitByAuthID("phar:pu", 1*time.Minute, 20), phh.UpdatePharmacy)

			// Update Profile Picture
			me.PUT("/profile-picture", rlm.LimitByAuthID("phar:ppu", 1*time.Minute, 10), phh.UpdateProfilePicture)

			// Update Profile Banner
			me.PUT("/profile-banner", rlm.LimitByAuthID("phar:pbu", 1*time.Minute, 10), phh.UpdateProfileBanner)
		}
	}

	// PRODUCT ENDPOINTS
	product := v1.Group("/products")
	{
		// Create
		product.POST("", middlewares.Auth(j), rlm.LimitByAuthID("prod:pc", 1*time.Minute, 5), prh.CreateProduct)
	}

	// USER ENDPOINTS
	user := v1.Group("/users")
	{
		// Create
		user.POST("", middlewares.Auth(j), rlm.LimitByAuthID("user:uc", 1*time.Minute, 5), uh.CreateUser)

		// Me
		me := user.Group("/me", middlewares.Auth(j))
		{
			// Fetch
			me.GET("", rlm.LimitByAuthID("user:me", 1*time.Minute, 50), uh.GetMe)

			// Update
			me.PATCH("", rlm.LimitByAuthID("user:uu", 1*time.Minute, 20), uh.UpdateUser)

			// Update Profile Picture
			me.PUT("/profile-picture", rlm.LimitByAuthID("user:ppu", 1*time.Minute, 10), uh.UpdateProfilePicture)

			// Update Profile Banner
			me.PUT("/profile-banner", rlm.LimitByAuthID("user:pbu", 1*time.Minute, 10), uh.UpdateProfileBanner)

			// Addresses
			address := me.Group("/addresses", rlm.LimitByAuthID("user:addr", 1*time.Minute, 30))
			{
				// Create
				address.POST("", uah.CreateAddress)

				// Fetch All
				address.GET("", uah.GetAllAddresses)

				// Update
				address.PATCH("/:address_id", uah.UpdateAddress)

				// Set Primary
				address.PUT("/:address_id/primary", uah.SetPrimaryAddress)

				// Delete
				address.DELETE("/:address_id", uah.DeleteAddress)
			}
		}
	}

	return &Router{router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
