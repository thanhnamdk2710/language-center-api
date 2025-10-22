package http

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/routes"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	routes.InitRoutes(r)

	return r
}
