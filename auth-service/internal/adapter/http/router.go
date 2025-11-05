package http

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/adapter/http/route"
	"github.com/thanhnamdk2710/auth-service/internal/container"
)

func NewRouter(c *container.Container) *gin.Engine {
	r := gin.Default()

	route.InitRoutes(r, c)

	return r
}
