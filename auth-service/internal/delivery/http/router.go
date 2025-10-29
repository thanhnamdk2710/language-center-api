package http

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/container"
	"github.com/thanhnamdk2710/auth-service/internal/delivery/http/routes"
)

func NewRouter(c *container.Container) *gin.Engine {
	r := gin.Default()

	routes.InitRoutes(r, c)

	return r
}
