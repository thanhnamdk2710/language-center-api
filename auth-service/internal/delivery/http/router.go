package http

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/routes"
)

func NewRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()

	routes.InitRoutes(r, db)

	return r
}
