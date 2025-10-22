package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/handlers"
)

func InitRoutes(r *gin.Engine) {
	r.GET("/health", handlers.HealthCheck())

	api := r.Group("/api/v1")
	initAuthRoutes(api)
}
