package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/http/middleware"
	di "github.com/thanhnamdk2710/auth-service/internal/app"
)

func initAuthRoutes(router *gin.RouterGroup, c *di.Container) *gin.RouterGroup {
	auth := router.Group("auth")

	// Apis
	auth.POST("/register", middleware.BindAndValidate[dto.RegisterRequest](), c.RegisterHandler.Register)
	auth.POST("/verify-otp", middleware.BindAndValidate[dto.VerifyOTPRequest](), c.VerifyOTPHandler.VerifyOTP)

	rateLimiter := middleware.NewRateLimiter(5, 15*time.Minute)
	auth.POST("/login",
		rateLimiter.Middleware(),
		middleware.BindAndValidate[dto.LoginRequest](),
		c.LoginHandler.Login,
	)

	auth.POST("/forgot-password", middleware.BindAndValidate[dto.ForgotPasswordRequest](), c.ForgotPasswordHandler.ForgotPassword)
	auth.POST("/refresh-token", middleware.BindAndValidate[dto.RefreshTokenRequest](), c.RefreshTokenHandler.RefreshToken)

	return auth
}
