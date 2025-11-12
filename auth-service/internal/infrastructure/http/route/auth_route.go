package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/di"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/middleware"
)

func initAuthRoutes(router *gin.RouterGroup, c *di.Container) *gin.RouterGroup {
	auth := router.Group("auth")

	// Apis
	auth.POST("/register", c.RegisterHandler.Register)
	auth.POST("/verify-otp", c.VerifyOTPHandler.VerifyOTP)

	rateLimiter := middleware.NewRateLimiter(5, 15*time.Minute)
	auth.POST("/login",
		rateLimiter.Middleware(),
		c.LoginHandler.Login,
	)

	auth.POST("/forgot-password", c.ForgotPasswordHandler.ForgotPassword)
	auth.POST("/refresh-token", c.RefreshTokenHandler.RefreshToken)

	return auth
}
