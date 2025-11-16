package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/di"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/middleware"
)

func NewRouter(c *di.Container) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")

	// Health check
	api.GET("/health", c.HealthHandler.Check)

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", c.AuthHandler.Register)
		auth.POST("/verify-otp", c.AuthHandler.VerifyOTP)
		auth.POST("/login",
			middleware.NewRateLimiter(5, 15*time.Minute).Middleware(),
			c.AuthHandler.Login,
		)
	}

	return r
}
