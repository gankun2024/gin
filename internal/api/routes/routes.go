package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"

	"github.com/gankun2024/gin-demo-project/internal/api/handlers"
	// "github.com/gankun2024/gin-demo-project/internal/api/middleware"
)

// Setup configures all routes for the application
func Setup(r *gin.Engine, logger *slog.Logger) {
	// Add middleware
	r.Use(sloggin.New(logger))
	r.Use(gin.Recovery())

	// Public routes
	public := r.Group("/v1")
	{
		// Health check
		public.GET("/health", handlers.HealthCheck)

		// Auth routes
		// auth := public.Group("/auth")
		// {
		// 	auth.POST("/register", handlers.Register)
		// 	auth.POST("/login", handlers.Login)
		// 	auth.POST("/forgot-password", handlers.ForgotPassword)
		// 	auth.POST("/reset-password", handlers.ResetPassword)
		// }
	}

	// // Protected routes
	// protected := r.Group("/api/v1")
	// protected.Use(middleware.Auth())
	// {
	// 	// User routes
	// 	user := protected.Group("/user")
	// 	{
	// 		user.GET("/profile", handlers.GetProfile)
	// 		user.PUT("/profile", handlers.UpdateProfile)
	// 	}

	// 	// Payment routes
	// 	payment := protected.Group("/payments")
	// 	{
	// 		payment.POST("/create-checkout", handlers.CreateCheckoutSession)
	// 		payment.GET("/success", handlers.PaymentSuccess)
	// 		payment.GET("/cancel", handlers.PaymentCancel)
	// 	}
	// }

	// // Webhook routes (Stripe)
	// webhook := r.Group("/webhooks")
	// {
	// 	webhook.POST("/stripe", handlers.StripeWebhook)
	// }
}
