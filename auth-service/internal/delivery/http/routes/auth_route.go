package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/container"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/dto"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/middlewares"
)

func initAuthRoutes(router *gin.RouterGroup, c *container.Container) *gin.RouterGroup {
	auth := router.Group("auth")

	// Apis
	auth.POST("/register", middlewares.BindAndValidate[dto.RegisterRequest](), c.RegisterHandler.Register)
	auth.POST("/login", middlewares.BindAndValidate[dto.LoginRequest](), c.LoginHandler.Login)
	auth.POST("/forgot-password", middlewares.BindAndValidate[dto.ForgotPasswordRequest](), c.ForgotPasswordHandler.ForgotPassword)
	auth.POST("/refresh-token", middlewares.BindAndValidate[dto.RefreshTokenRequest](), c.RefreshTokenHandler.RefreshToken)

	return auth
}
