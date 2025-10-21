package http

import (
	"github.com/gin-gonic/gin"
	handler "github.com/thanhnamdk2710/auth-service/internal/delivery/http/handlers"
	usecase "github.com/thanhnamdk2710/auth-service/internal/usecase/register"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger(), gin.Recovery())

	// Health check handler
	r.GET("/health", handler.HealthCheck())

	// Init all route groups
	api := r.Group("/api/v1")
	initAuthRoutes(api)

	return r
}

func initAuthRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("auth")

	// Dependency injection
	registerUsecase := usecase.NewRegisterUsecase()
	registerHandler := handler.NewAuthHandler(registerUsecase)

	// Apis
	authGroup.POST("/register", registerHandler.Register)
}
