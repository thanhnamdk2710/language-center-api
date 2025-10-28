package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/handlers"
)

func InitRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/health", handlers.HealthCheck())

	api := r.Group("/api/v1")
	initAuthRoutes(api, db)
}
