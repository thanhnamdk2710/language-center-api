package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/middlewares"
	"github.com/thanhnamdk2710/auth-service/internal/modules"
)

func initAuthRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	auth := router.Group("auth")

	// Dependency injection
	authModule := modules.NewAuthModule()

	// Apis
	auth.POST("/register", middlewares.BindAndValidate[dto.RegisterRequest](), authModule.RegisterHandler.Register)
	auth.POST("/login", middlewares.BindAndValidate[dto.LoginRequest](), authModule.LoginHandler.Login)
	auth.POST("/forgot-password", middlewares.BindAndValidate[dto.ForgotPasswordRequest](), authModule.ForgotPasswordHandler.ForgotPassword)
	auth.POST("/refresh-token", middlewares.BindAndValidate[dto.RefreshTokenRequest](), authModule.RefreshTokenHandler.RefreshToken)

	return auth
}
