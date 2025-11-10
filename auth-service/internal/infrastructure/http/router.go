package http

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/di"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/http/route"
)

func NewRouter(c *di.Container) *gin.Engine {
	r := gin.Default()

	route.InitRoutes(r, c)

	return r
}
