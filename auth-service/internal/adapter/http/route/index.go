package route

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/http/handler"
	di "github.com/thanhnamdk2710/auth-service/internal/app"
)

func InitRoutes(r *gin.Engine, c *di.Container) {
	r.GET("/health", handler.HealthCheck())

	api := r.Group("/api/v1")
	initAuthRoutes(api, c)
}
