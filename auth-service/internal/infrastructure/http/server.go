package http

import (
	"github.com/gin-gonic/gin"
	"github.com/thanhnamdk2710/auth-service/internal/infrastructure/di"
)

func NewRouter(c *di.Container) *gin.Engine {
	r := gin.Default()

	RegisterRoutes(r, c)

	return r
}
